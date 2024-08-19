package main

import (
	"os"
	"path/filepath"
	"time"
)

type RotationWriter struct {
	file        *os.File
	name        string
	dir         string
	rotateCheck time.Duration
	lastRotate  time.Time
}

func (rw *RotationWriter) Write(b []byte) (n int, err error) {
	if rw.file == nil || time.Since(rw.lastRotate) >= rw.rotateCheck {
		if err = rw.Rotate(); err != nil {
			return
		}
		rw.lastRotate = time.Now()
	}
	return rw.file.Write(b)
}

func (rw *RotationWriter) Rotate() (err error) {
	return
}
