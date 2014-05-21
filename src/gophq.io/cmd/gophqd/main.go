package main

import (
	"flag"
	"gophq.io/gophqd"
	_ "gophq.io/pprof"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	tcpListenFlag  = flag.String("tcp.listen", "127.0.0.1:9092", "TCP listener address")
	unixListenFlag = flag.String("unix.listen", "", "Unix listener address")
	dataDirFlag    = flag.String("datadir", "", "data directory (empty: disable persistence)")
)

var revision string

func main() {
	flag.Parse()

	log.Printf("gophqd %s", revision)

	var tcpl net.Listener
	if *tcpListenFlag != "" {
		var err error
		tcpl, err = net.Listen("tcp", *tcpListenFlag)
		if err != nil {
			panic(err)
		}
	}

	var unixl net.Listener
	if *unixListenFlag != "" {
		var err error
		unixl, err = net.Listen("unix", *unixListenFlag)
		if err != nil {
			panic(err)
		}
	}

	s := &gophqd.Server{
		DataDir: *dataDirFlag,
	}

	if tcpl != nil {
		go serve(s, tcpl)
	}

	if unixl != nil {
		go serve(s, unixl)
	}

	// Package signal will not block sending to c.
	// The caller must ensure that c has sufficient buffer space.
	signals := make(chan os.Signal, 10)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	for {
		switch sig := <-signals; sig {
		case syscall.SIGTERM, syscall.SIGINT:
			log.Printf("received signal %v", sig)
			return
		default:
			continue
		}
	}

	// TODO shut down the serve goroutines

	// TODO wait for them to be done before exiting
}

func serve(s *gophqd.Server, l net.Listener) {
	err := s.Serve(l)
	if err != nil {
		panic("Serve: " + err.Error())
	}
}