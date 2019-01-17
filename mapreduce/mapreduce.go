package mapreduce

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// service for RPC
type Word int

func (t *Word) Map(chunk string, result *map[int][]string) error {

	fmt.Println(os.Getpid(), ":  mapping..")

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
	fmt.Println(os.Getpid(), ":  map complete!")
	return nil
}


func (t *Word) Reduce(reduceMap map[int][]string, result *map[int]int) error {

	fmt.Println(os.Getpid(), ":  reducing..")
	for k, words := range reduceMap {
		// could it be 	(*result)[k] = len(words)
		for _ = range words {
			(*result)[k]++
		}
	}

	fmt.Println(os.Getpid(), ":  reduce complete!")
	return nil
}
