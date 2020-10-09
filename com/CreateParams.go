package com

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type CreateParams struct {
	host     string
	ports    string
	replicas int
	master   int
}

/*验证主机是否输入*/
func (p *CreateParams) verifyHost() (bool, interface{}) {
	if p.host == "" {
		return false, "主机IP不能为空"
	} else if address := net.ParseIP(p.host); address == nil {
		return false, "主机IP不合法"
	}
	return true, nil
}

func (p *CreateParams) verifyMaster() (bool, interface{}) {
	if p.master < 3 {
		return false, "主机数不能小于3"
	} else if p.master%2 == 0 {
		return false, "主机数必须为奇数"
	}

	return true, nil
}

func (p *CreateParams) verifyReplicas() (bool, interface{}) {
	if p.replicas < 1 {
		return false, "副本应该大于等于1"
	}
	return true, nil
}

func (p *CreateParams) verifyPorts() (bool, interface{}) {
	count := p.master * (p.replicas + 1)
	/*使用系统默认生成的端口*/
	if p.ports == "" {
		var ports strings.Builder
		start := 7001
		for i := 0; i < count; i++ {
			ports.WriteString(strconv.Itoa(start))
			if i < count-1 {
				ports.WriteString(",")
				start++
			}
		}
		p.ports = ports.String()
	} else { /*用户自己输入的端口列表*/
		split := strings.Split(p.ports, ",")
		if len(split) != count {
			return false, "		端口的个数不对"
		}
		/*判断是否有重复的*/
		distinctMap := make(map[int64]interface{}, count)
		for _, element := range split {
			portInt, err := strconv.ParseInt(element, 10, 10)
			if err != nil {
				return false, "		端口不合法,请输入整数"
			}
			if _, ok := distinctMap[portInt]; ok {
				return false, fmt.Sprintf("		存在相同端口: %v", portInt)
			}
			distinctMap[portInt] = true
		}
	}
	log.Printf("集群将使用如下端口:%s 以及容器间通信使用的端口 %s \n", p.ports, "1"+strings.ReplaceAll(p.ports, ",", ",1"))

	return true, nil
}
