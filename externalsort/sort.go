//go:build !solution

package externalsort

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"os"
	"sort"
)

type HeapItem struct {
	line      string
	readerIdx int
	reader    LineReader
}

// Куча реализует интерфейс heap.Interface
type MinHeap []HeapItem

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].line < h[j].line }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Push добавляет элемент в кучу
func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(HeapItem))
}

// Pop удаляет и возвращает минимальный элемент
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

func NewReader(r io.Reader) LineReader {
	return &MyLineReader{bufio.NewReader(r)}
}

func NewWriter(w io.Writer) LineWriter {
	return &MyLineWriter{bufio.NewWriter(w)}
}

func Merge(w LineWriter, readers ...LineReader) error {
	h := &MinHeap{}
	heap.Init(h)

	for i, reader := range readers {
		line, err := reader.ReadLine()
		if line == "" && err == io.EOF {
			continue
		} else if err != nil && err != io.EOF {
			return err
		} else {
			heap.Push(h, HeapItem{line, i, reader})
		}
	}

	for h.Len() > 0 {
		el := heap.Pop(h)
		w.Write(el.(HeapItem).line)

		newLine, err := el.(HeapItem).reader.ReadLine()
		if err == io.EOF && newLine == "" {
			continue
		} else if err != nil && err != io.EOF {
			return err
		} else {
			heap.Push(h, HeapItem{newLine, el.(HeapItem).readerIdx, el.(HeapItem).reader})
		}
	}

	return nil
}

func Sort(w io.Writer, in ...string) error {
	for _, path := range in {
		lines, err := readLines(path)

		if err != nil {
			return err
		}

		sort.Strings(lines)

		err = writeLines(path, lines)

		if err != nil {
			return err
		}
	}

	var lineReaders []LineReader

	for _, path := range in {
		file, err := os.Open(path)

		if err != nil {
			return fmt.Errorf("can't open file %s: %v", path, err)
		}

		var reader io.Reader = file
		lineReaders = append(lineReaders, NewReader(reader))
	}

	err := Merge(NewWriter(w), lineReaders...)

	if err != nil {
		return err
	}

	return nil
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("can't open file %s: %v", path, err)
	}

	defer file.Close()

	var reader io.Reader = file
	var lineReader = NewReader(reader)

	var result []string

	for {
		line, err := lineReader.ReadLine()
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("error while read line in file %s: %v", path, err)
		}
		if err == io.EOF {
			if line != "" {
				result = append(result, line)
			}
			break
		}
		result = append(result, line)
	}

	return result, nil
}

func writeLines(path string, lines []string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("can't open file %s: %v", path, err)
	}

	defer file.Close()

	var writer io.Writer = file
	var lineWriter = NewWriter(writer)

	for _, line := range lines {
		err = lineWriter.Write(line)
		if err != nil {
			return fmt.Errorf("can't write line to file %s: %v", path, err)
		}
	}

	return nil
}
