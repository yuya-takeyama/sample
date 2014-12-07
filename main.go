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
	p := flags.NewParser(opts, flags.Default)
	args, err := p.Parse()
	if err != nil {
		os.Exit(1)
	}

	rate := opts.Rate

	rand.Seed(time.Now().Unix())

	readInput(os.Stdin, args, rate)
}

func readInput(stdin io.Reader, args []string, rate float64) {
	if len(args) > 0 {
		for _, filepath := range args {
			file, err := os.Open(filepath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to open file: %s\n", err.Error())
				os.Exit(1)
			}

			sample(file, rate)

			if err = file.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "failed to close file: %s\n", err.Error())
				os.Exit(1)
			}
		}
	} else {
		sample(stdin, rate)
	}
}

func sample(file io.Reader, rate float64) {
	reader := bufio.NewReader(file)
	var line []byte
	var err error
	for ; err == nil; line, err = reader.ReadBytes('\n') {
		if rand.Float64() < rate {
			os.Stdout.Write(line)
		}
	}
}
