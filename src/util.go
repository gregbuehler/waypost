package main

import "log"

func warnIfErr(err error) {
	if err != nil {
		log.Print(err)
	}
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
