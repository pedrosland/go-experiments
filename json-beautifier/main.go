package main

import (
	"bufio"
	"encoding/json"
	"os"

	"log"
)

func main() {
	s := bufio.NewScanner(os.Stdin)

	logger := log.New(os.Stderr, "[json]", 0)

	for s.Scan() {
		if err := s.Err(); err != nil {
			logger.Fatalln(err)
			break
		}

		text := s.Bytes()
		if len(text) == 0 {
			// Skip blank lines
			continue
		}

		var a interface{}

		err := json.Unmarshal(text, &a)
		if err != nil {
			log.Println(err)
			continue
		}

		out, err := json.MarshalIndent(a, "", "  ")
		if err != nil {
			logger.Print(err)
			continue
		}

		os.Stdout.Write(out)
		os.Stdout.WriteString("\n")
	}
}
