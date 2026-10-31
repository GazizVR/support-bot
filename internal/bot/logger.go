package bot

import (
	"fmt"
	"os"
)

type FileLogger struct {
	file *os.File
}

func NewFileLogger(filepath string) (*FileLogger, error) {
	file, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	return &FileLogger{file: file}, nil
}

func (l *FileLogger) Debugf(format string, args ...any) {
	fmt.Fprintf(l.file, format, args)
	fmt.Fprintln(l.file, "\n")
	l.file.Sync()
}

func (l *FileLogger) Errorf(format string, args ...any) {
	fmt.Fprintf(l.file, format, args)
	fmt.Fprintln(l.file, "\n")
	l.file.Sync()
}
