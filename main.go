package main

import (
	"os"

	"github.com/bilberrry/quarterback/common"
	"gopkg.in/urfave/cli.v1"
)

const (
	appName = "quarterback"
	appDesc = "Manage your backups from one common file"
)

func main() {
	var target string

	app := cli.NewApp()
	app.Name = appName
	app.Usage = appDesc

	app.Commands = []cli.Command{
		{
			Name: "run",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "target",
					Usage:       "Run only specified target",
					Destination: &target,
				},
			},
			Action: func(c *cli.Context) error {
				if len(target) == 0 {
					runAll()
				} else {
					runTarget(target)
				}

				return nil
			},
		},
	}

	app.Run(os.Args)
}

func runAll() {
	for _, targetConfig := range common.Targets {

		target := Target{
			Config: targetConfig,
		}

		target.run()
	}
}

func runTarget(targetName string) {
	for _, targetConfig := range common.Targets {
		if targetConfig.Name == targetName {

			target := Target{
				Config: targetConfig,
			}

			target.run()
			return
		}
	}
}
