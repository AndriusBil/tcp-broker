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
	sw.Value = append(sw.Value, strings.Trim(string(p), "\n"))
	sw.mu.Unlock()
	return len(p), nil
}
