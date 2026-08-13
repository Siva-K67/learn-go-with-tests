package main

import (
	"io"
	"os"
)

// tape wraps a ReadWriteSeeker so writes always start from the beginning of the file
type tape struct {
	file *os.File
}

// Write empties the file, rewinds, then writes fresh data — prevents leftover bytes from old writes
func (t *tape) Write(p []byte) (n int, err error) {
	t.file.Truncate(0)           // ADDED: empty the file completely
	t.file.Seek(0, io.SeekStart) // rewind cursor to start
	return t.file.Write(p)       // write fresh data
}
