package main

import (
	"TgToys/express"
	"TgToys/weather"
	"flag"
	"fmt"
)

func usage() {
	fmt.Println("Usage: tgtoys command [arguments] ...")
	fmt.Println("command is one of:")
	fmt.Println("\t express\t Express tracking service.\n\t weather\t Weather forecast query.")

}

func main() {
	flag.Parse()
	args := flag.Args()
	if args == nil || len(args) < 1 {
		usage()
		return
	}
	switch args[0] {
	case "express":
		if len(args) != 2 {
			fmt.Println("Usage: tgtoys express <expressnumber>")
			return
		}
		express.GetExpressInfo(args[1])
	case "weather":
		if len(args) > 2 {
			fmt.Println("Usage: tgtoys weather -<address>")
			return
		}
		if len(args) == 2 {
			weather.GetWeather(args[1])
		} else {
			weather.GetWeather()
		}

	default:
		usage()
	}

}
