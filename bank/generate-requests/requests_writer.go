package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type RequestsWriter struct {
	file   *os.File
	writer *bufio.Writer
	buf    []byte
}

func NewRequestsWriter(path string) *RequestsWriter {
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("Creating the output file failed: %v", err)
	}
	return &RequestsWriter{
		file:   file,
		writer: bufio.NewWriterSize(file, 64*1024),
	}
}

func (rw *RequestsWriter) WriteHeader(config Config) {
	_, err := rw.writer.WriteString(fmt.Sprintf("%d %d\n", config.accounts, config.initialBalance))
	if err != nil {
		log.Fatalf("Writing the header failed: %v", err)
	}
}

func (rw *RequestsWriter) WriteLine(i uint64, isSuccessful bool, fromID uint64, toID uint64, amount uint64) {
	rw.buf = rw.buf[:0]
	rw.buf = strconv.AppendUint(rw.buf, i, 10)
	rw.buf = append(rw.buf, ' ')
	if isSuccessful {
		rw.buf = append(rw.buf, "1 "...)
	} else {
		rw.buf = append(rw.buf, "0 "...)
	}
	rw.buf = strconv.AppendUint(rw.buf, fromID, 10)
	rw.buf = append(rw.buf, ' ')
	rw.buf = strconv.AppendUint(rw.buf, toID, 10)
	rw.buf = append(rw.buf, ' ')
	rw.buf = strconv.AppendUint(rw.buf, amount, 10)
	rw.buf = append(rw.buf, '\n')
	rw.writer.Write(rw.buf)
}

func (rw *RequestsWriter) Close() {
	if err := rw.writer.Flush(); err != nil {
		log.Fatalf("Flushing the output failed: %v", err)
	}
	if err := rw.file.Close(); err != nil {
		log.Fatalf("Closing the output file failed: %v", err)
	}
}
