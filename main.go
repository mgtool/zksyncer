package main

import (
	"io/ioutil"
	"mgt-zookeeper/pkg/tools"

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

}

func TransformData(params ...interface{}) []string {

	strArray := make([]string, len(params))
	for i, arg := range params {
		strArray[i] = arg.(string)
	}
	return strArray
}
