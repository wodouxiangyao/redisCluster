package com

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	clusterCommand = "redis-cli --cluster create --cluster-replicas %d %s"
)

/*创建命令的配置*/
func createCommand() cli.Command {

	return cli.Command{
		Name:                   "create",
		Usage:                  "创建集群",
		Category:               "cluster",
		Flags:                  configCreateOption(),
		UseShortOptionHandling: true,
		Action:                 handleAction,
	}
}

/*处理创建命令下的响应*/
func handleAction(c *cli.Context) error {

	/*params := &CreateParams{c.String("host"), c.String("ports"),c.Int("replicas"), c.Int("master")}*/
	params := &CreateParams{c.String("host"), c.String("ports"), 1, 3}

	verifyParams(params)
	createImage()
	runContainer(params)
	return nil
}

/*验证所有的参数*/
func verifyParams(params *CreateParams) {
	if ok, errMsg := params.verifyHost(); !ok {
		log.Fatal(errMsg)
		os.Exit(1)
	}

	if ok, errMsg := params.verifyMaster(); !ok {
		log.Fatal(errMsg)
		os.Exit(1)
	}

	if ok, errMsg := params.verifyReplicas(); !ok {
		log.Fatal(errMsg)
		os.Exit(1)
	}

	if ok, errMsg := params.verifyPorts(); !ok {
		log.Fatal(errMsg)
		os.Exit(1)
	}
}

func createImage() {

	/*判断是否存在该镜像,存在则不删除*/
	output, err := exec.Command("/bin/sh", "-c", IsExistImage).Output()
	if err != nil {
		log.Fatal("		查看镜像发生错误  ", err)
	}
	sprintf := fmt.Sprintf("%s", output)

	if strings.TrimSpace(sprintf) == strconv.Itoa(0) {
		log.Println("即将生成镜像...")
		err := exec.Command("/bin/sh", "-c", buildImageString()).Run()
		if err != nil {
			log.Fatal("		创建镜像发生错误  ", err)
		}
		log.Println("生成镜像 redis:cluster")
	} else {
		log.Println("镜像已存在,无需创建")
	}
	exec.Command("/bin/sh", "-c", "docker rmi alpine").Run()
}

func runContainer(params *CreateParams) {
	/*判断是否存在容器，存在就删除*/
	output, err := exec.Command("/bin/sh", "-c", IsExistContainer).Output()
	if err != nil {
		log.Fatal("		查看容器列表发生错误  ", err)
	}

	/*存在则全部删除*/
	if strings.TrimSpace(string(output)) != strconv.Itoa(0) {
		output, _ := exec.Command("/bin/sh", "-c", ClusterContainerName).Output()

		/*将命令的结果转换为以空格分隔的一行字符串*/
		allContainerWithSpace := strings.ReplaceAll(string(output), "\n", " ")
		_, err := exec.Command("/bin/sh", "-c", rmContainerString(allContainerWithSpace)).Output()
		if err != nil {
			log.Fatal("删除容器发生错误  ", err)
		}

		log.Printf("删除容器:   %s\n", strings.ReplaceAll(string(output), "\n", "\n\t\t\t\t"))
	}
	/*创建容器*/

	createContainer, allHostStr := getCreateContainerString(*params)

	log.Println(createContainer.String())
	command := exec.Command("/bin/sh", "-c", createContainer.String())
	command.Wait()

	output, err = command.Output()
	if err != nil {
		log.Fatal("创建容器发生错误  ", err)
	}

	log.Printf(`创建容器:   %s`, strings.ReplaceAll(string(output), "\n", "\n\t\t\t\t"))

	createCluster(params.replicas, strings.Join(allHostStr, " "))

	createContainer.Reset()
	split := strings.Split(params.ports, ",")
	for i := 0; i < params.master; i++ {
		createContainer.WriteString(params.host)
		createContainer.WriteString(":")
		createContainer.WriteString(split[i])
		createContainer.WriteString("\t")
	}

	log.Println("redis的集群创建成功...")
	log.Printf("集群访问地址为: %s \n", createContainer.String())
}

/*创建集群*/
func createCluster(replicas int, allHost string) {

	redisCliComm := fmt.Sprintf(clusterCommand, replicas, allHost)
	interactionComm := fmt.Sprintf(`expect -c "
			   spawn %s
			   expect {
					\"(type 'yes' to accept):\" {send \"yes\r\";exp_continue;}
			}"`, redisCliComm)
	commStr := fmt.Sprintf("docker exec redis1 %s", interactionComm)
	log.Println(commStr)
	err := exec.Command("/bin/sh", "-c", commStr).Run()

	if err != nil {
		log.Fatal("创建集群发生错误  ", err)
	}
}

/*创建容器的配置选项*/
func configCreateOption() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:     "host, H",
			Usage:    "该主机的IP",
			Required: true,
			Hidden:   false,
			Value:    "",
		},

		cli.IntFlag{
			Name:  "master, m",
			Usage: "创建的主机数,默认为3,应该大于3且为奇数",
			Value: 3,
		},

		cli.IntFlag{
			Name:  "replicas, r",
			Usage: "创建的副本数,默认为1",
			Value: 1,
		},

		cli.StringFlag{
			Name:  "ports, p",
			Value: "",
			Usage: `集群多容器的端口,默认可使用系统生成(7001开始,逐步加1),
					也可自己指定,个数为master与replicas加1的乘积(用逗号分隔)
					确保端口没有被占用`,
		},
	}
}
