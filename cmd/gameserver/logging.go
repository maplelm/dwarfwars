package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type RotationWriter struct {
	file       *os.File
	name       string
	dir        string
	lastRotate time.Time
	maxSize    int64
}

func NewRotationWriter(path, name string, maxsize int64) *RotationWriter {
	return &RotationWriter{
		name:       name,
		dir:        path,
		lastRotate: time.Unix(0, 0),
		maxSize:    maxsize,
	}
}

func (rw *RotationWriter) Write(b []byte) (n int, err error) {
	tobig := false
	info, err := os.Stat(filepath.Join(rw.dir, rw.name))
	if err == nil && info.Size()/1_000_000 >= rw.maxSize {
		tobig = true
	}
	if rw.file == nil || tobig {
		err = rw.Rotate()
		if err != nil {
			fmt.Printf("RotationWriter Failed: %s", err)
			os.Exit(1)
			return
		}

		rw.lastRotate = time.Now()
	}
	return rw.file.Write(b)
}

func (rw *RotationWriter) Rotate() (err error) {
	file := filepath.Join(rw.dir, rw.name)
	if rw.file != nil {
		rw.file.Close()
		rw.file = nil
	}

	info, err := os.Stat(file)
	if err == nil && info.Size()/1_000_000 > rw.maxSize {
		err = os.Rename(file, fmt.Sprintf("%s_%s_%s.log", file, time.Now().Format(time.DateOnly), time.Now().Format(time.TimeOnly)))
	}

	rw.file, err = os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("RotationWriter.Rotate(): %s\n", err)
	}
	return
}

func LogError(err error) {
	log.Printf("Error: %s", err)
}

func LogWarning(msg string) {
	log.Printf("Warning: %s", msg)
}

func LogInfo(msg string) {
	log.Printf("Info: %s", msg)
}
