package file

import (
	"bufio"
	"errors"
	"math/rand"
	"os"
)

type Buffer []string

func NewBuffer(path string) *Buffer {
	b := new(Buffer)
	var lines, err = readLines(path)
	if err != nil {
		panic(err)
	}
	b = &lines
	return b
}

func readLines(path string) (Buffer, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines Buffer
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func (b Buffer) Get(i ...int) (metric string, err error) {

	switch len(i) {
	case 0:
		metric = b[rand.Intn(len(b))]
	case 1:
		idx := i[0]
		if idx < 0 || idx > b.Len() {
			err = errors.New("Oops. Index out of range")
		} else {
			metric = b[idx]
		}
	default:
		err = errors.New("Get should only have 0 or 1 parameter.")
	}
	return metric, err
}

func (b Buffer) Len() int {
	return len(b)
}
