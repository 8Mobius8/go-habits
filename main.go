package main

import "github.com/8Mobius8/go-habits/cmd"

var version = "[Not provided]"

func main() {
	cmd.Version = version
	cmd.Execute()
}
