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
	imageName      = "redis"
	imageTag       = "cluster"
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
	existImages := fmt.Sprintf("docker images|awk '$1==\"%s\" && $2==\"%s\" {print $1 $2}'|wc -l", imageName, imageTag)
	output, err := exec.Command("/bin/sh", "-c", existImages).Output()
	if err != nil {
		log.Fatal("		查看镜像发生错误  ", err)
	}
	sprintf := fmt.Sprintf("%s", output)

	if strings.TrimSpace(sprintf) == strconv.Itoa(0) {
		log.Println("即将生成镜像...")
		buildCommand := fmt.Sprintf("docker build --force-rm -q --no-cache -t %s:%s -f %s %s", imageName, imageTag, tempDir+"resources/Dockerfile", tempDir+"resources/")
		err := exec.Command("/bin/sh", "-c", buildCommand).Run()
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
	existContainer := "docker ps -a | awk '$NF ~ \"redis.+\" {print $NF}'|wc -l"
	output, err := exec.Command("/bin/sh", "-c", existContainer).Output()
	if err != nil {
		log.Fatal("		查看容器列表发生错误  ", err)
	}
	sprintf := fmt.Sprintf("%s", output)
	/*存在则全部删除*/
	if strings.TrimSpace(sprintf) != strconv.Itoa(0) {
		output, _ := exec.Command("/bin/sh", "-c", "docker ps -a | awk '$NF ~ \"redis.+\" {print $NF}'").Output()
		sprintf = fmt.Sprintf("%s", output)
		split := strings.Split(sprintf, "\n")
		var containers strings.Builder
		containers.WriteString("docker rm -f ")
		containers.WriteString(strings.Join(split, " "))
		output, err := exec.Command("/bin/sh", "-c", containers.String()).Output()
		if err != nil {
			log.Fatal("删除容器发生错误  ", err)
		}
		log.Printf("删除容器:   %s\n", strings.Join(strings.Split(fmt.Sprintf("%s", output), "\n"), "\n\t\t\t\t"))
	}
	/*创建容器*/
	containerCount := params.master * (params.replicas + 1)
	var createContainer strings.Builder

	split := strings.Split(params.ports, ",")

	/*创建集群副本的所有主机IP端口*/
	allHostStr := make([]string, containerCount)

	for i := 0; i < containerCount; i++ {
		createContainer.WriteString("docker run -itd --name redis")
		createContainer.WriteString(strconv.Itoa(i + 1))
		createContainer.WriteString(" -h redis")
		createContainer.WriteString(strconv.Itoa(i + 1))
		createContainer.WriteString(" --network ")
		createContainer.WriteString(redisClusterNetwork)
		createContainer.WriteString(" --ip 172.30.188.")

		allHostStr[i] = "172.30.188." + strconv.Itoa(100+i+1) + ":6379"
		createContainer.WriteString(strconv.Itoa(100 + i + 1)) //容器的IP從100開始
		createContainer.WriteString(" -p ")
		createContainer.WriteString(split[i])
		createContainer.WriteString(":6379")
		createContainer.WriteString(" -p 1")
		createContainer.WriteString(split[i])
		createContainer.WriteString(":16379")
		createContainer.WriteString(" -e redisPort=")
		createContainer.WriteString(split[i])
		createContainer.WriteString(" -e redisHost=")
		createContainer.WriteString(params.host)
		createContainer.WriteString(" ")
		createContainer.WriteString(imageName)
		createContainer.WriteString(":")
		createContainer.WriteString(imageTag)
		if i < containerCount-1 {
			createContainer.WriteString("   && ")
		}
	}

	command := exec.Command("/bin/sh", "-c", createContainer.String())
	command.Wait()

	output, err = command.Output()
	if err != nil {
		log.Fatal("创建容器发生错误  ", err)
	}

	log.Printf(`创建容器:   %s`, strings.Join(strings.Split(fmt.Sprintf("%s", output), "\n"), "\n\t\t\t\t"))

	createCluster(params.replicas, strings.Join(allHostStr, " "))

	createContainer.Reset()
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

	//_, err := exec.Command("/bin/sh", "-c", commStr).Output()
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
