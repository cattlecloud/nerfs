package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"cattlecloud.net/go/babycli"
	"cattlecloud.net/nerfs"
	"cattlecloud.net/nerfs/cmds/nerfs-builds/domains"
	"cattlecloud.net/nerfs/cmds/nerfs-builds/wordlist"
)

func Invoke(args []string) babycli.Code {
	return invoke(args)
}

func invoke(args []string) babycli.Code {
	return babycli.New(&babycli.Configuration{
		Arguments: args,
		Version:   "v0.0.0-dev",
		Top: &babycli.Component{
			Name:        "nerfs-builds",
			Help:        "run the artifact builder(s)",
			Description: "Builds an artifact",
			Flags:       babycli.Flags{},
			Components: babycli.Components{
				{
					Name: "wordlist",
					Help: "generate the wordlist.txt artifact",
					Flags: babycli.Flags{
						{
							Type:    babycli.StringFlag,
							Long:    "output",
							Require: true,
							Short:   "o",
							Help:    "specify output FILE",
							Default: &babycli.Default{
								Value: filepath.Join(os.TempDir(), nerfs.WordsFile),
								Show:  true,
							},
						},
					},
					Function: func(c *babycli.Component) babycli.Code {
						output := c.GetString("output")
						b := wordlist.NewBuilder()
						if err := b.Build(output); err != nil {
							fmt.Println("build failure:", err)
							return babycli.Failure
						}
						return babycli.Success
					},
				},
				{
					Name: "domains",
					Help: "generate the domains.txt artifact",
					Flags: babycli.Flags{
						{
							Type:    babycli.StringFlag,
							Long:    "output",
							Require: true,
							Short:   "o",
							Help:    "specify output FILE",
							Default: &babycli.Default{
								Value: filepath.Join(os.TempDir(), nerfs.DomainsFile),
								Show:  true,
							},
						},
					},
					Function: func(c *babycli.Component) babycli.Code {
						output := c.GetString("output")
						b := domains.NewBuilder()
						if err := b.Build(output); err != nil {
							fmt.Println("build failure:", err)
							return babycli.Failure
						}
						return babycli.Success
					},
				},
			},
		},
	}).Run()
}
