package main

import (
	"os"
	"path"
)

var pprofPath = guessPprofPath()

func exe() string {
	exe, err := os.Readlink("/proc/self/exe")
	if err != nil {
		return ""
	}
	return exe
}

func exeDir() string {
	exe, err := os.Readlink("/proc/self/exe")
	if err != nil {
		return ""
	}
	dir, _ := path.Split(exe)
	return dir
}

func guessPprofPath() string {
	path := os.Getenv("GOROOT") + "/misc/pprof"
	_, err := os.Stat(path)
	if err == nil {
		return path
	}
	path = exeDir() + "pprof"
	_, err = os.Stat(path)
	if err != nil {
		panic("could not find pprof")
	}
	return path
}
