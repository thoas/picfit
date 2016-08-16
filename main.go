package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/thoas/picfit/application"
	"github.com/thoas/picfit/server"
	"github.com/thoas/picfit/signature"
)

func main() {
	app := cli.NewApp()
	app.Name = "picfit"
	app.Author = "thoas"
	app.Email = "florent.messa@gmail.com"
	app.Usage = "Display, manipulate, transform and cache your images"
	app.Version = application.Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Config file path",
			EnvVar: "PICFIT_CONFIG_PATH",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "version",
			ShortName: "v",
			Usage:     "Retrieve the version number",
			Action: func(c *cli.Context) {
				fmt.Printf("picfit %s\n", application.Version)
			},
		},
		{
			Name:      "signature",
			ShortName: "s",
			Usage:     "Verify that your client application is generating correct signatures",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "key",
					Usage: "The signing key",
				},
			},
			Action: func(c *cli.Context) {
				key := c.String("key")
				if key == "" {
					fmt.Fprintf(os.Stderr, "You must provide a key\n")
					os.Exit(1)
				}

				if len(c.Args()) < 1 {
					fmt.Fprintf(os.Stderr, "You must provide a Query String\n")
					os.Exit(1)
				}

				queryString := c.Args()[0]

				sig, err := signature.SignRaw(key, queryString)

				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}

				appended := signature.AppendSign(key, queryString)

				fmt.Fprintf(os.Stdout, "Query String: %s\n", queryString)
				fmt.Fprintf(os.Stdout, "Signature: %s\n", sig)
				fmt.Fprintf(os.Stdout, "Signed Query String: %s\n", appended)
			},
		},
	}
	app.Action = func(c *cli.Context) {
		config := c.String("config")

		if config != "" {
			if _, err := os.Stat(config); err != nil {
				fmt.Fprintf(os.Stderr, "Can't find config file `%s`\n", config)
				os.Exit(1)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Can't find config file\n")
			os.Exit(1)
		}

		err := server.Load(config)

		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	}

	app.Run(os.Args)
}
