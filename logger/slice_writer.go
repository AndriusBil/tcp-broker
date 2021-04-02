package logger

import (
	"strings"
	"sync"
)

type SliceWriter struct {
	mu    sync.Mutex
	Value []string
}

func (sw *SliceWriter) Write(p []byte) (int, error) {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	sw.Value = append(sw.Value, strings.Trim(string(p), "\n"))
	return len(p), nil
}

func (sw *SliceWriter) Get() []string {
	sw.mu.Lock()
	defer sw.mu.Unlock()
	return sw.Value
}
