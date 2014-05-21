package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	signals := make(chan os.Signal, 10)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case sig := <-signals:
			log.Printf("signal: %v", sig)
			switch sig {
			case syscall.SIGTERM:
				log.Printf("terminating")
				return
			case syscall.SIGINT:
				log.Printf("terminating due to interrupt")
				return
			}
		case now := <-ticker.C:
			log.Printf("tick at %v", now)
		}
	}
}
