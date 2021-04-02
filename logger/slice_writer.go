package logger

import "strings"

type SliceWriter struct {
	Slice []string
}

func (sw *SliceWriter) Write(p []byte) (int, error) {
	sw.Slice = append(sw.Slice, strings.Trim(string(p), "\n"))
	return len(p), nil
}
