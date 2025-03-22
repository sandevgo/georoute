package main

import (
	"log"

	"github.com/sandevgo/georoute/internal/options"
	"github.com/sandevgo/georoute/internal/ripencc"
)

func main() {
	opts, err := options.NewOptions()
	if err != nil {
		log.Fatalln(err)
	}

	registry := ripencc.NewRegistry(opts.Country, opts.Format)

	err = registry.GetDelegated()
	if err != nil {
		log.Fatalln(err)
	}
	defer registry.Close()

	err = registry.Process()
	if err != nil {
		log.Fatalln(err)
	}
}
