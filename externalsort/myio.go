package externalsort

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type (
	MyLineReader struct {
		reader *bufio.Reader
	}
	MyLineWriter struct {
		writer *bufio.Writer
	}
)

func (r *MyLineReader) ReadLine() (string, error) {
	line, err := r.reader.ReadString('\n')

	if err != nil && err != io.EOF {
		return line, err
	}

	return strings.TrimSuffix(line, "\n"), err
}

func (w *MyLineWriter) Write(l string) error {
	_, err := fmt.Fprintln(w.writer, l)
	if err != nil {
		return err
	}
	return w.writer.Flush()
}
