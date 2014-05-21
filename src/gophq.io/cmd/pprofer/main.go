package main

import (
	"bufio"
	"bytes"
	"flag"
	"github.com/gorilla/mux"
	_ "gophq.io/pprof"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

type unixDialer struct {
	path string
}

func (this *unixDialer) Dial(network, addr string) (net.Conn, error) {
	return net.Dial("unix", this.path)
}

const pprofSockets = "@/go/pprof/"

func get(pid int, path, rawQuery string) ([]byte, error) {
	pprofSocket := pprofSockets + strconv.Itoa(pid)
	dialer := &unixDialer{pprofSocket}
	transport := &http.Transport{
		Dial: dialer.Dial,
	}
	client := &http.Client{
		Transport: transport,
	}

	urlStr := "http://unix" + path + "?" + rawQuery
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	// Don't leave connections between pprofer and
	// other processes lying around.
	// Setting MaxIdleConnsPerHost in the transport
	// did not accomplish this.
	req.Close = true

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

const rootText = `<html>
<head>
<title>pprofer</title>
<body>
<table>
{{range .}}
<tr>
<td>{{.Cmdline}}</td>
<td>{{.Pid}}</td>
<td><a href="/{{.Pid}}/memstats">memstats</a></td>
<td><a href="/{{.Pid}}/blockrate?rate=1">blockrate=1</a></td>
<td><a href="/{{.Pid}}/pprof/profile?seconds=30">cpu 30s</a></td>
<td><a href="/{{.Pid}}/pprof/block">block</a></td>
<td><a href="/{{.Pid}}/pprof/goroutine">goroutine</a></td>
<td><a href="/{{.Pid}}/pprof/goroutine?debug=2">stacktrace</a></td>
<td><a href="/{{.Pid}}/pprof/heap">heap</a></td>
</tr>
{{end}}
</table>
</body>
</html>
`

var rootTmpl = template.Must(template.New("/").Parse(rootText))

type process struct {
	Pid     int
	Cmdline string
}

type processesByPid []process

func (s processesByPid) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s processesByPid) Len() int           { return len(s) }
func (s processesByPid) Less(i, j int) bool { return s[i].Pid < s[j].Pid }

func handleSlash(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		httpErrCode(w, http.StatusBadRequest)
		return
	}

	f, err := os.Open("/proc/net/unix")
	if err != nil {
		httpErr(w, err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var procs []process
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		if len(parts) != 8 {
			continue
		}

		flags, err := strconv.ParseInt(parts[3], 16, 64)
		if err != nil {
			continue
		}
		// constant from linux/net.h
		const soAcceptCon = 1 << 16
		if flags != soAcceptCon {
			continue
		}

		addr := parts[len(parts)-1]
		if !strings.HasPrefix(addr, pprofSockets) {
			continue
		}
		pidStr := addr[len(pprofSockets):]
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			continue
		}

		f, err := os.Open("/proc/" + pidStr + "/cmdline")
		if err != nil {
			continue
		}
		b, err := ioutil.ReadAll(f)
		f.Close()
		if err != nil {
			continue
		}
		b = bytes.Replace(b, []byte{0}, []byte{' '}, -1)
		cmdline := string(b)
		cmdline = strings.TrimSpace(cmdline)
		procs = append(procs, process{pid, cmdline})
	}

	sort.Sort(processesByPid(procs))

	buf := new(bytes.Buffer)
	if err := rootTmpl.Execute(buf, procs); err != nil {
		httpErr(w, err)
		return
	}
	b := buf.Bytes()

	header := w.Header()
	header.Set("Content-Length", strconv.Itoa(len(b)))
	header.Set("Content-Type", "text/html; charset=utf-8")
	header.Set("Content-Language", "en")
	w.Write(b)
}

type pidHandler struct {
	path string
}

func (this *pidHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// TODO use gorilla's MethodHandler for this
	if req.Method != "GET" {
		httpErrCode(w, http.StatusBadRequest)
		return
	}

	vars := mux.Vars(req)
	pid, err := strconv.Atoi(vars["pid"])
	if err != nil {
		httpErr(w, err)
		return
	}
	b, err := get(pid, this.path, req.URL.RawQuery)
	if err != nil {
		httpErr(w, err)
		return
	}

	header := w.Header()
	header.Set("Content-Length", strconv.Itoa(len(b)))
	header.Set("Content-Type", "text/plain; charset=utf-8")
	header.Set("Content-Language", "en")
	w.Write(b)
}

type profileHandler struct {
	profile string
}

func (this *profileHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// TODO use gorilla's MethodHandler for this
	if req.Method != "GET" {
		httpErrCode(w, http.StatusBadRequest)
		return
	}

	vars := mux.Vars(req)
	pid, err := strconv.Atoi(vars["pid"])
	if err != nil {
		httpErr(w, err)
		return
	}

	url := "/debug/pprof/" + this.profile
	profile, err := get(pid, url, req.URL.RawQuery)
	if err != nil {
		httpErr(w, err)
		return
	}

	tmpfile, err := ioutil.TempFile("", "pprof")
	if err != nil {
		httpErr(w, err)
		return
	}
	defer tmpfile.Close()
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write(profile); err != nil {
		httpErr(w, err)
		return
	}

	exe := "/proc/" + strconv.Itoa(pid) + "/exe"
	cmd := exec.Command(pprofPath, "--svg", exe, tmpfile.Name())
	svg, err := cmd.Output()
	if err != nil {
		httpErr(w, err)
		return
	}

	header := w.Header()
	header.Set("Content-Length", strconv.Itoa(len(svg)))
	header.Set("Content-Type", "image/svg+xml")
	w.Write(svg)
}

func handlePprof(w http.ResponseWriter, req *http.Request) {
	// TODO use gorilla's MethodHandler for this
	if req.Method != "GET" {
		httpErrCode(w, http.StatusBadRequest)
		return
	}

	vars := mux.Vars(req)
	pid, err := strconv.Atoi(vars["pid"])
	if err != nil {
		httpErr(w, err)
		return
	}

	b, err := get(pid, "/debug/pprof/"+vars["profile"], req.URL.RawQuery)
	if err != nil {
		httpErr(w, err)
		return
	}

	header := w.Header()
	header.Set("Content-Length", strconv.Itoa(len(b)))
	header.Set("Content-Type", "text/plain; charset=utf-8")
	header.Set("Content-Language", "en")
	w.Write(b)
}

func main() {
	listenAddr := flag.String("addr", "127.0.0.1:9111", "listen address")
	flag.Parse()

	router := mux.NewRouter()
	router.Handle("/{pid}/pprof/profile", &profileHandler{"profile"})
	router.Handle("/{pid}/pprof/heap", &profileHandler{"heap"})
	router.Handle("/{pid}/pprof/block", &profileHandler{"block"})
	router.HandleFunc("/{pid}/pprof/{profile}", handlePprof)
	router.Handle("/{pid}/memstats", &pidHandler{"/debug/memstats"})
	router.Handle("/{pid}/blockrate", &pidHandler{"/debug/blockrate"})
	router.Handle("/{pid}/cpurate", &pidHandler{"/debug/cpurate"})
	router.HandleFunc("/", handleSlash)

	serveMux := http.NewServeMux()
	serveMux.Handle("/", router)
	srv := &http.Server{Handler: serveMux}

	l, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Panicf("Listen: %v", err)
	}
	log.Panicf("Serve: %v", srv.Serve(l))
}
