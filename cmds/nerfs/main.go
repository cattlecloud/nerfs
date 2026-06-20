package main

import (
	"os"

	"cattlecloud.net/go/babycli"
	"cattlecloud.net/go/nerfs/cmds/nerfs/commands"
)

func main() {
	args := babycli.Arguments()
	rc := commands.Invoke(args)
	os.Exit(rc)
}
