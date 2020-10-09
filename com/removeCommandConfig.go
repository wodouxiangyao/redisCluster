package com

import (
	"fmt"
	"github.com/urfave/cli"
)

/*创建命令的配置*/
func rmCommand() cli.Command {

	return cli.Command{
		Name: "rm",
		//Aliases:  []string{"s"},
		Usage:    "删除",
		Category: "remove",
		Flags:    configRmOption(),
		Action: func(c *cli.Context) error {
			fmt.Println("5 - 3 = ", 5-3)
			return nil
		},
	}
}

func configRmOption() []cli.Flag {

	return []cli.Flag{
		cli.BoolTFlag{
			Name:  "force, f",
			Usage: "强制删除",
		},
	}

}
