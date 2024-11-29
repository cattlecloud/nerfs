package main

import (
	"os"

	"cattlecloud.net/nerfs/cmds/nerfs-builds/commands"
	"noxide.lol/go/babycli"
)

func main() {
	args := babycli.Arguments()
	rc := commands.Invoke(args)
	os.Exit(rc)
}
