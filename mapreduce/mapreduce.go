package mapreduce

import (
	"bufio"
	"errors"
	"strings"
)

// service for RPC
type Word int

func (t *Word) Map(chunk string, result *map[int][]string) error {

	var word string
	var wlen int
	var res = make(map[int][]string)

	scanner := bufio.NewScanner(strings.NewReader(chunk))
	// Set the split function for the scanning operation.
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word = scanner.Text()
		word = strings.Trim(word, ".,!?'\"\"")
		wlen = len(word)
		res[wlen] = append(res[wlen], word)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	*result = res;
	return nil
}

func (t *Word) Reduce(args map[int]string, result *map[int]string) error {
	if args == nil {
		return errors.New("Divide by zero")
	}
	return nil
}
