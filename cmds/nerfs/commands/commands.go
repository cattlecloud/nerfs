package commands

import (
	"os"

	"cattlecloud.net/go/babycli"
	"cattlecloud.net/go/ulog"
	"cattlecloud.net/nerfs/cmds/nerfs/domains"
	"cattlecloud.net/nerfs/cmds/nerfs/wordlist"
)

func Invoke(args []string) babycli.Code {
	return invoke(args)
}

func invoke(args []string) babycli.Code {
	return babycli.New(&babycli.Configuration{
		Arguments: args,
		Version:   "v0.0.0",
		Top: &babycli.Component{
			Name:        "nerfs",
			Help:        "run the artifact builder(s)",
			Description: "Builds an artifact",
			Flags:       babycli.Flags{},
			Components: babycli.Components{
				{
					Name: "build",
					Help: "generate the artifact files",
					Flags: babycli.Flags{
						{
							Type:    babycli.StringFlag,
							Long:    "output",
							Require: true,
							Short:   "o",
							Help:    "specify output DIR",
							Default: &babycli.Default{
								Value: os.TempDir(),
								Show:  true,
							},
						},
					},
					Function: func(c *babycli.Component) babycli.Code {
						log := ulog.New("nerfs")

						output := c.GetString("output")

						buildDomains := domains.NewBuilder()
						if err := buildDomains.Build(output); err != nil {
							log.E.Fmt("unable to build domains list: %v", err)
							return babycli.Failure
						}

						buildWords := wordlist.NewBuilder()
						if err := buildWords.Build(output); err != nil {
							log.E.Fmt("unable to build words list: %v", err)
							return babycli.Failure
						}

						return babycli.Success
					},
				},
			},
		},
	}).Run()
}
