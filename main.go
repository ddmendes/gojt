package main

import (
	"bufio"
	"encoding/json"
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

	jsonDoc, err := unmarshalJSON(pipedData)
	if err != nil {
		panic(err)
	}

	output, err := marshalJSON(jsonDoc, true)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(output))
}

func getPipedData() ([]byte, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return []byte{}, err
	}

	if info.Mode()&os.ModeNamedPipe != 0 && info.Size() > 0 {
		var output []byte
		reader := bufio.NewReader(os.Stdin)
		for {
			input, err := reader.ReadByte()
			if err == io.EOF {
				break
			}
			output = append(output, input)
		}
		return []byte(output), nil
	}
	return []byte{}, errors.New("No piped data")
}

func unmarshalJSON(input []byte) (interface{}, error) {
	var document interface{}
	err := json.Unmarshal(input, &document)
	return document, err
}

func marshalJSON(document interface{}, beautify bool) ([]byte, error) {
	if beautify {
		return json.MarshalIndent(document, "", "    ")
	}
	return json.Marshal(document)
}
