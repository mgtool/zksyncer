package tools

import (
	"fmt"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

var (
	TotalNodes        = 0
	BlackNodes        = 0
	CreateNodes       = 0
	CreateFailedNodes = 0
	ModifyNodes       = 0
	ModifyFailedNodes = 0
	SyncedNodes       = 0
)

func Start(srcHosts []string, dstHosts []string, whiteList []string, blackList []string) {

	src, _, err := zk.Connect(srcHosts, time.Second*5)
	dst, _, err := zk.Connect(dstHosts, time.Second*5)

	defer src.Close()
	defer dst.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	// 白名单，只同步白名单内节点
	for _, v := range whiteList {
		if strings.TrimSpace(v) != "/" {
			SyncData(strings.TrimSpace(v), blackList, src, dst)
		}
		CheckPath(strings.TrimSpace(v), blackList, src, dst)
	}

}

func CheckPath(basePath string, blackList []string, src *zk.Conn, dst *zk.Conn) {

	// 1.获取子节点
	c, _, err := src.Children(basePath)
	if err != nil {
		fmt.Println("节点: " + basePath + ", 不存在，检测失败！！！")
		panic(err)
	}

	// 2.判断是否存在子节点
	if len(c) > 0 { // 存在
		for _, child := range c {

			// 2.1 判断父节点是否是"/"
			tempNode := basePath + "/" + child
			if basePath == "/" {
				tempNode = basePath + child
			}

			// 2.2 同步数据
			SyncData(tempNode, blackList, src, dst)

			// 2.3 递归调用
			CheckPath(tempNode, blackList, src, dst)
		}
	}

}

func SyncData(node string, blackList []string, src *zk.Conn, dst *zk.Conn) {

	TotalNodes++

	fmt.Println("-------------------------------------------------------")
	fmt.Println("节点: " + node + ", 同步开始！！！")

	// 黑名单中的节点不同步
	for _, v := range blackList {

		blackNode := strings.TrimSpace(v)

		if strings.Contains(blackNode, "/*") { // 黑名单，通配符匹配
			s := (strings.Split(blackNode, "*"))[0]
			if strings.Contains(node, s) || strings.Contains(node+"/", s) {
				fmt.Println("节点: " + node + ", 通配黑名单，同步终止！！！")
				BlackNodes++
				return
			}
		} else {
			if node == blackNode { // 黑名单，普通字符匹配
				fmt.Println("节点: " + node + ", 普通黑名单，同步终止！！！")
				BlackNodes++
				return
			}
		}

	}

	srcValue, _, err := src.Get(node)
	if err != nil {
		fmt.Println("节点: " + node + ", 源节点查询异常，同步失败！！！")

		panic(err)
	}

	dstValue, dstStat, err := dst.Get(node)

	if err != nil { // 插入数据
		_, err := dst.Create(node, srcValue, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			fmt.Println("节点: " + node + ", 新增数据，同步失败！！！")
			CreateFailedNodes++
			panic(err)
		}
		CreateNodes++
		fmt.Println("节点: " + node + ", 新增数据，同步成功！！！")
	} else if string(dstValue) != string(srcValue) { // 修改数据
		_, err = dst.Set(node, srcValue, dstStat.Version)
		if err != nil {
			fmt.Println("节点: " + node + ", 修改数据，同步失败！！！")
			ModifyFailedNodes++
			panic(err)
		}
		ModifyNodes++
		fmt.Println("节点: " + node + ", 修改数据，同步成功！！！")
	} else {
		SyncedNodes++
		fmt.Println("节点: " + node + ", 数据一致，同步终止！！！")
	}

}
