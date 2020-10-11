package com

import (
	"fmt"
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
	IsExistContainer     = "docker ps -a | awk '$NF ~ \"redis.+\" {print $NF}'|wc -l"
	ClusterContainerName = "docker ps -a | awk '$NF ~ \"redis.+\" {print $NF}'"
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
func getCreateContainerString(params CreateParams) (createContainerBuilder strings.Builder, allHostStr []string) {
	count := params.master * (params.replicas + 1)

	split := strings.Split(params.ports, ",")

	/*创建集群副本的所有主机IP端口*/
	allHostStr = make([]string, count)

	for i := 0; i < count; i++ {
		sprintf := fmt.Sprintf("docker run -itd --name redis%d -h redis%d  --network %s --ip 172.30.188.%d -p %s:6379 -p 1%s:16379 -e redisPort=%s  -e redisHost=%s %s:%s",
			i+1, i+1, NetName, 101+i, split[i], split[i], split[i], params.host, ImageName, ImageTag)

		createContainerBuilder.WriteString(sprintf)

		if i < count-1 {
			createContainerBuilder.WriteString("   && ")
		}
		allHostStr[i] = "172.30.188." + strconv.Itoa(101+i) + ":6379"
	}
	return
}
