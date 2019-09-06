package main

import (
	"flag"
	app "github.com/randallmlough/goezyrouting"
)

func main() {
	cfg := flag.String("cfg", "./cmd/config.json", "Path to config file")

	a, err := app.Initialize(*cfg)
	if err != nil {
		panic(err)
	}
	a.Start()
}
