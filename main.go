package main

import (
	"bufio"
	"fmt"
	"github.com/jessevdk/go-flags"
	"io"
	"math/rand"
	"os"
	"time"
)

type cmdOptions struct {
	Rate float64 `short:"r" long:"rate" description:"Sampling rate" default:"0.1"`
}

func main() {
	opts := &cmdOptions{}
	p := flags.NewParser(opts, flags.HelpFlag)
	args, err := p.Parse()
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	rate := opts.Rate

	rand.Seed(time.Now().Unix())

	lineCh := make(chan []byte)
	exitCh := make(chan bool)

	go readInputs(os.Stdin, args, lineCh, exitCh)
	printSampledLines(lineCh, exitCh, rate)
}

func readInputs(stdin io.Reader, args []string, lineCh chan []byte, exitCh chan bool) {
	if len(args) > 0 {
		for _, filepath := range args {
			file, err := os.Open(filepath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to open file: %s\n", err.Error())
				os.Exit(1)
			}

			readLines(file, lineCh)

			if err = file.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "failed to close file: %s\n", err.Error())
				os.Exit(1)
			}
		}
	} else {
		readLines(stdin, lineCh)
	}

	exitCh <- true
}

func readLines(file io.Reader, lineCh chan []byte) {
	reader := bufio.NewReader(file)
	var line []byte
	var err error
	for ; err == nil; line, err = reader.ReadBytes('\n') {
		lineCh <- line
	}
}

func printSampledLines(lineCh chan []byte, exitCh chan bool, rate float64) {
	for {
		select {
		case line := <-lineCh:
			if rand.Float64() < rate {
				os.Stdout.Write(line)
			}

		case <-exitCh:
			return
		}
	}
}
