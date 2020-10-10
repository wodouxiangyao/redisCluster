package com

import (
	"strconv"
	"strings"
)

/*集群的桥接网络相关*/
const (
	NetName    = "redis"
	IsExistNet = "docker network ls |grep " + NetName + "|wc -l "
	CreateNet  = "docker network create --subnet=172.30.188.0/24 --gateway=172.30.188.1 " + NetName
)

/*镜像相关*/
const (
	ImageName    = "redis"
	ImageTag     = "cluster"
	IsExistImage = "docker images|awk '$1==\"" + ImageName + "\" && $2==\"" + ImageTag + "\" {print " + ImageName + "}'|wc -l"
)

const (
	IsExistContainer = "docker ps -a | awk '$NF ~ \"redis.+\" {print $NF}'|wc -l"
)

/*获取构建镜像的命令*/
func buildImageString() string {
	return "docker build --force-rm -q --no-cache -t " + ImageName + ":" + ImageTag + " -f " + tempDir + "resources/Dockerfile " + tempDir + "resources/"
}

/*删除容器的命令*/
func rmContainerString(containers string) string {
	var builder strings.Builder
	builder.WriteString("docker rm -f ")
	builder.WriteString(containers)
	return builder.String()
}

/**

 */
func getCreateContainerString(params CreateParams) (createString strings.Builder, allHostStr []string) {
	count := params.master * (params.replicas + 1)

	var createContainer strings.Builder
	split := strings.Split(params.ports, ",")

	/*创建集群副本的所有主机IP端口*/
	allHostStr = make([]string, count)

	for i := 0; i < count; i++ {
		createContainer.WriteString("docker run -itd --name redis")
		createContainer.WriteString(strconv.Itoa(i + 1))
		createContainer.WriteString(" -h redis")
		createContainer.WriteString(strconv.Itoa(i + 1))
		createContainer.WriteString(" --network ")
		createContainer.WriteString(NetName)
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
		createContainer.WriteString(ImageName)
		createContainer.WriteString(":")
		createContainer.WriteString(ImageTag)
		if i < count-1 {
			createContainer.WriteString("   && ")
		}
	}
	return createContainer, allHostStr
}
