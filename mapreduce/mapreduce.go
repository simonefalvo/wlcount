package mapreduce

import (
	"bufio"
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
		word = strings.Trim(word, ".,:;()[]{}!?'\"\"")
		wlen = len(word)
		res[wlen] = append(res[wlen], word)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	*result = res;
	return nil
}


func (t *Word) Reduce(reduceMap map[int][]string, result *map[int]int) error {

	for k, words := range reduceMap {
		// could it be 	(*result)[k] = len(words)
		for _ = range words {
			(*result)[k]++
		}
	}

	return nil
}
