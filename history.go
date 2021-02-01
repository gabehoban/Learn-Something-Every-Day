package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func timespecToTime(ts syscall.Timespec) string {
	raw := time.Unix(int64(ts.Sec), int64(ts.Nsec))
	return raw.Format("2006-01-02")
}

func updateLog(file string) string {
	fileIO, e := os.OpenFile(file, os.O_RDWR, 0600)
	check(e)
	data, e := ioutil.ReadAll(fileIO)
	check(e)
	lines := strings.Split(string(data), "\n")
	desc := lines[0]
	fileIO.Close()
	return desc
}

func main() {
	searchDir := "."
	fileList := make([]string, 0)
	e := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() && f.Name() == (".git") {
			return filepath.SkipDir
		}
		if filepath.HasPrefix(f.Name(), "README") == false {
			if filepath.Ext(f.Name()) == (".md") {
				fileList = append(fileList, path)
			}
		}
		return err
	})
	check(e)
	currentTime := time.Now()
	readme, e := os.OpenFile("README.md", os.O_APPEND|os.O_WRONLY, 0644)
	check(e)
	beginning := "\n### " + currentTime.Format("2006-01-02")
	if _, e = readme.WriteString(beginning); e != nil {
		panic(e)
	}
	readme.Close()
	for _, file := range fileList {
		finfo, _ := os.Stat(file)
		statT := finfo.Sys().(*syscall.Stat_t)
		readme, e := os.OpenFile("README.md", os.O_APPEND|os.O_WRONLY, 0644)
		check(e)
		if timespecToTime(statT.Ctim) == currentTime.Format("2006-01-02") {
			point := updateLog(file)
			if _, e = readme.WriteString("\n  - (" + file + "): " + point); e != nil {
				panic(e)
			}
			readme.Close()
		}

	}
}
