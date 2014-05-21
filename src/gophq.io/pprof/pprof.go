package cfpprof

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"runtime"
	"strconv"
)

func memStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	b, err := json.MarshalIndent(m, "", "  ")
	if err == nil {
		w.Write(b)
		w.Write([]byte{'\n'})
	}
}

func rateHandler(name string, setRate func(int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rates := r.URL.Query()["rate"]
		var rate int
		var err error
		if len(rates) > 0 {
			rate, err = strconv.Atoi(rates[0])
		} else {
			err = errors.New("no rate parameter")
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(fmt.Sprintf("%s(%d)\n", name, rate)))
		setRate(rate)
	}
}

type Handler struct {
	Name string
	Func http.HandlerFunc
}

func listenPprofHandlers(handlers ...Handler) {
	path := "@/go/pprof/" + strconv.Itoa(os.Getpid())
	l, err := net.Listen("unix", path)
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	for _, handler := range handlers {
		mux.Handle(handler.Name, handler.Func)
	}
	go http.Serve(l, mux)
}

var defaultHandlers = []Handler{
	{"/debug/memstats", http.HandlerFunc(memStats)},
	{"/debug/pprof/", http.HandlerFunc(pprof.Index)},
	{"/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline)},
	{"/debug/pprof/profile", http.HandlerFunc(pprof.Profile)},
	{"/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol)},
	{"/debug/blockrate", rateHandler("SetBlockProfileRate", runtime.SetBlockProfileRate)},
	{"/debug/cpurate", rateHandler("SetCPUProfileRate", runtime.SetCPUProfileRate)},
}

func init() {
	listenPprofHandlers(defaultHandlers...)
}
