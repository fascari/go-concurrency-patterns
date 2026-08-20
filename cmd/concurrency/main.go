package main

import (
	"context"
	"flag"
	"log"
)

func main() {
	raw := flag.String("pattern", "", "scenario to run")
	flag.Parse()

	pattern, err := ParsePattern(*raw)
	if err != nil {
		log.Fatal(err)
	}

	if err := pattern.run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
