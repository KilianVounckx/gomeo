package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for range c {
			fmt.Println()
			fmt.Println()
			main()
		}
	}()

	builtinFunctions["run"] = RUN

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("gomeo> ")
		scanner.Scan()
		text := scanner.Text()

		if strings.Trim(text, " \t") == "" {
			fmt.Println()
			continue
		}

		value, err := run("<stdin>", text)
		if err != nil {
			fmt.Println(err.AsString())
			continue
		} else if value != nil {
			fmt.Println(value)
		}

		fmt.Println()
	}
}
