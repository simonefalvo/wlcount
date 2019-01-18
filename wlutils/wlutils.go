// Package wlutils provides utility functions to the word length count problem

package wlutils

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

// Scan the input reader r according to the split function split
// returning an array of strings on success
func ScanStrings(r io.Reader, split bufio.SplitFunc) ([]string, error) {
	var s []string
	scanner := bufio.NewScanner(r)
	if split != nil {
		scanner.Split(split)
	}
	for scanner.Scan() {
		s = append(s, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return s, nil
}

// Split the file in order to obtain n chunks of approximately the same size in bytes
func SplitFile(file *os.File, n int) []string {

	stat, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fileSize := stat.Size()
	//fmt.Println("file size:", fileSize, "bytes")
	chunkSize := fileSize / int64(n)
	//fmt.Println("chunk size:", chunkSize, "bytes")
	chunks := make([]string, n)

	buffer := make([]byte, chunkSize)
	for i := 0; i < n; i++ {
		if i == n-1 {
			chunkSize = fileSize
			buffer = make([]byte, chunkSize)
		}
		read, err := io.ReadAtLeast(file, buffer, int(chunkSize))
		if err != nil {
			log.Fatal(err)
		}
		chunks[i] = string(buffer)
		//fmt.Printf("-------READ CHUNK %d (%d bytes)-------\n", i, read)
		//fmt.Println(chunks[i])
		if i != n-1 {
			pos := strings.LastIndexByte(chunks[i], ' ')
			chunks[i] = chunks[i][:pos]
			//fmt.Printf("-------CUT CHUNK %d (%d bytes)------\n", i, pos)
			//fmt.Println(chunks[i])
			// Position at the beginning of the cut word
			_, err := file.Seek(-int64(read - pos), io.SeekCurrent)
			if err != nil {
				log.Fatal(err)
			}
			fileSize -= int64(pos)
		}
	}
	return chunks
}
