package helpers

import (
	"bufio"
	"fmt"
	"os"
)

func WriteMapToFile(filepath string, data []string) error {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, value := range data {
		_, err := writer.WriteString(fmt.Sprintf("%s\n", value))
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}
