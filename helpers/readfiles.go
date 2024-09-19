package helpers

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func ReadFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("ERROR > Error opening file: ", err)
		return nil, err
	}
	defer file.Close()

	linesMap := make(map[int]string)
	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		linesMap[lineNumber] = line
		lineNumber++
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("ERROR > Error reading file: ", err)
		return nil, err
	}

	// Extract keys and sort them
	var keys []int
	for k := range linesMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// Create a sorted list of lines
	var sortedLines []string
	for _, k := range keys {
		sortedLines = append(sortedLines, linesMap[k])
	}

	return sortedLines, nil
}
