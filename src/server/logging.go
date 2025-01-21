package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	fp "path/filepath"
	"time"
)

/*
 * This file will be where the logging specific code is being held.
 * One of these things will be the log writer so not all log data will go into the same file after
 *	specified amount of time or file size.
 */

type RotateWriter struct {
	Name    string        // Name file file. should contain the path as well if needed
	MaxAge  time.Duration // how long the object can have this file open before it needs ot roate out a new one
	MaxSize int64         // Max size a file can be before a new one is needed to be rolled out

	age          time.Time   // time the file was opened by this object
	file         *os.File    // pointer to the actual files ot write and read from
	fileinfo     os.FileInfo // System information on the file
	infopulltime time.Time   // When os.Stat last updated fileinfo
}

func (r *RotateWriter) Write(b []byte) (int, error) {
	var err error
	if r.file == nil {
		if r.fileinfo, err = os.Stat(r.Name); errors.Is(err, os.ErrNotExist) {
			if r.file, err = os.OpenFile(r.Name, os.O_APPEND, 777); err != nil {
				return 0, err
			}
			r.age = time.Now()
			return r.file.Write(b)
		} else if err != nil {
			return 0, err
		}
		r.infopulltime = time.Now()
		newname := fp.Join(fp.Dir(r.Name), string(time.Now().Unix())+"_"+r.fileinfo.Name())
		if err = os.Rename(r.Name, newname); err != nil {
			return 0, err
		}
		if r.file, err = os.OpenFile(r.Name, os.O_APPEND, 777); err != nil {
			return 0, err
		}
		r.age = time.Now()
		return r.file.Write(b)
	}

	if time.Since(r.infopulltime) >= time.Second*5 {
		r.fileinfo, err = os.Stat(r.Name)
		r.infopulltime = time.Now()
	}

	if time.Since(r.age) <= r.MaxAge && r.fileinfo.Size() < r.MaxSize {
		return r.file.Write(b)
	}
	newname := fp.Join(fp.Dir(r.Name), string(time.Now().Unix())+"_"+r.fileinfo.Name())
	if err = os.Rename(r.Name, newname); err != nil {
		return 0, err
	}
	return r.file.Write(b)
}
