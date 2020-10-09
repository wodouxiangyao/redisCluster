package com

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

const (
	redisClusterNetwork = "redis"
)

func ConfigAll(app *cli.App) {
	configBaseInfo(app)
	configCommand(app)

	configBefore(app)
	configAfter(app)
}

/*配置基础信息*/
func configBaseInfo(app *cli.App) {
	app.Name = "single-docker-redis-cluster"
	app.Usage = "单机版redis的多docker容器部署cluster集群"
	app.Author = "WangLian"
	app.Version = "0.0.1-ALPHA"
	app.UsageText = "redis-cluster command [options] "
}

func configBefore(app *cli.App) {
	app.Before = func(c *cli.Context) error {

		existNetworkString := fmt.Sprintf(" docker network ls |grep %s|wc -l", redisClusterNetwork)
		output, err := exec.Command("/bin/sh", "-c", existNetworkString).Output()

		if err != nil {
			log.Fatal("查看是否存在redis桥接网络报错  ", err)
		}
		sprintf := fmt.Sprintf("%s", output)

		if strings.TrimSpace(sprintf) == strconv.Itoa(0) {
			log.Println("创建容器桥接网络...")
			/*创建redis桥接网络*/
			createNetworkString := fmt.Sprintf("docker network create --subnet=172.30.188.0/24 --gateway=172.30.188.1 %s", redisClusterNetwork)
			output, err = exec.Command("/bin/sh", "-c", createNetworkString).Output()
			if err != nil {
				log.Fatal("创建redis桥接网络报错  ", err)
			}
		}
		return nil
	}
}


func configAfter(app *cli.App) {
	app.After = func(c *cli.Context) error {
		command := exec.Command("/bin/sh", "-c", "rm -rf /tmp/cluster/*")
		command.Run()
		return nil
	}
}

/*配置命令*/
func configCommand(app *cli.App) {
	app.Commands = []cli.Command{
		createCommand(),
		rmCommand(),

		/*{
			Name:     "db",
			Usage:    "database operations",
			Category: "database",
			Subcommands: []cli.Command{
				{
					Name:  "insert",
					Usage: "insert data",
					Action: func(c *cli.Context) error {
						fmt.Println("insert subcommand")
						return nil
					},
				},
				{
					Name:  "delete",
					Usage: "delete data",
					Action: func(c *cli.Context) error {
						fmt.Println("delete subcommand")
						return nil
					},
				},
			},
		},*/
	}
}
