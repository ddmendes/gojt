package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	pipedData, err := getPipedData()
	if err != nil {
		panic(err)
	}
	fmt.Println(pipedData)
}

func getPipedData() (string, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	if info.Mode()&os.ModeNamedPipe != 0 && info.Size() > 0 {
		var output []rune
		reader := bufio.NewReader(os.Stdin)
		for {
			input, _, err := reader.ReadRune()
			if err == io.EOF {
				break
			}
			output = append(output, input)
		}
		return string(output), nil
	} else {
		return "", errors.New("No piped data")
	}
}
