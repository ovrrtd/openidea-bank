package main

import "github.com/ovrrtd/openidea-bank/cmd"

func main() {
	if err := cmd.Server(); err != nil {
		panic(err)
	}
}
