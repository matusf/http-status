package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	usage := "usage: http-status <status code>"
	logger := log.New(os.Stderr, "", 0)

	if len(os.Args) <= 1 {
		logger.Fatalln(usage)
	}

	code, err := strconv.Atoi(os.Args[1])
	if err != nil {
		logger.Fatalf("Invalid argument: '%s'\n%s\n", os.Args[1], usage)
	}

	text := http.StatusText(code)
	if text == "" {
		logger.Fatalf("Invalid code: '%d'\n%s\n", code, usage)
	}

	fmt.Println(text)
}
