package main

import (
	"fmt"
	"io/ioutil"
	"mgt-zookeeper/pkg/tools"
	"strconv"

	"github.com/olebedev/config"
)

func main() {

	file, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}
	yamlString := string(file)

	cfg, err := config.ParseYaml(yamlString)
	if err != nil {
		panic(err)
	}

	source, err := cfg.List("zookeeper.server.source")
	if err != nil {
		panic(err)
	}

	destination, err := cfg.List("zookeeper.server.destination")
	if err != nil {
		panic(err)
	}

	whitelist, err := cfg.List("zookeeper.url.whitelist")
	if err != nil {
		panic(err)
	}

	blacklist, err := cfg.List("zookeeper.url.blacklist")
	if err != nil {
		panic(err)
	}

	tools.Start(TransformData(source...), TransformData(destination...), TransformData(whitelist...), TransformData(blacklist...))

	fmt.Println("#########################统计结果#########################")
	fmt.Println("一致节点：" + strconv.Itoa(tools.SyncedNodes))
	fmt.Println("创建节点：" + strconv.Itoa(tools.CreateNodes))
	fmt.Println("创建失败节点：" + strconv.Itoa(tools.CreateFailedNodes))
	fmt.Println("修改节点：" + strconv.Itoa(tools.ModifyNodes))
	fmt.Println("修改失败节点：" + strconv.Itoa(tools.ModifyFailedNodes))
	fmt.Println("黑名单节点：" + strconv.Itoa(tools.BlackNodes))
	fmt.Println("节点总数：" + strconv.Itoa(tools.TotalNodes))

}

func TransformData(params ...interface{}) []string {

	strArray := make([]string, len(params))
	for i, arg := range params {
		strArray[i] = arg.(string)
	}
	return strArray
}
