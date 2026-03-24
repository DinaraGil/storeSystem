package logic

import (
	"bufio"
	"fmt"
	"io"
	"mime/multipart"
)

func FileHandler(file *multipart.File) error {
	reader := bufio.NewReader(*file)
	lineNum := 0

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		lineNum++
		fmt.Printf("Строка %d: %s", lineNum, line)
	}
	return nil
}
