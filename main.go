package main

import (
	"bufio"
	"fmt"
	"github.com/jessevdk/go-flags"
	"math/rand"
	"os"
	"time"
)

type cmdOptions struct {
	Rate float64 `short:"r" long:"rate" description:"Sampling rate" default:"0.1"`
}

func main() {
	opts := &cmdOptions{}
	p := flags.NewParser(opts, flags.HelpFlag+flags.PrintErrors)
	_, err := p.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
	}

	rate := opts.Rate

	rand.Seed(time.Now().Unix())

	reader := bufio.NewReader(os.Stdin)

	var line []byte

	for ; err == nil; line, err = reader.ReadBytes('\n') {
		if rand.Float64() < rate {
			os.Stdout.Write(line)
		}
	}
}
