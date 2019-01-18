package mapreduce

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// service for RPC
type MapReduce int


func (t *MapReduce) Map(chunk string, result *map[int][]string) error {

	if len(chunk) == 0 {
		return errors.New("empty argument string")
	}

	fmt.Println(os.Getpid(), ":  mapping..")

	res := make(map[int][]string)
	scanner := bufio.NewScanner(strings.NewReader(chunk))
	// Set the split function for the scanning operation.
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		word = strings.Trim(word, ".,:;()[]{}!?'\"\"")
		wlen := len(word)
		if wlen != 0 {
			res[wlen] = append(res[wlen], word)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	*result = res
	fmt.Println(os.Getpid(), ":  map complete!")
	return nil
}


func (t *MapReduce) Reduce(reduceMap map[int][]string, result *map[int]int) error {

	if reduceMap == nil {
		return errors.New("nil argument")
	}

	fmt.Println(os.Getpid(), ":  reducing..")
	for k, words := range reduceMap {
		// could it be 	(*result)[k] = len(words)
		for range words {
			(*result)[k]++
		}
	}
	fmt.Println(os.Getpid(), ":  reduce complete!")
	return nil
}
