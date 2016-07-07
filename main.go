// get-tengine-upstream project main.go

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type NginxUpstream struct {
	//	定义upstream状态struct
	Servers struct {
		Total      int `json:"total"`
		Generation int `json:"generation"`
		Server     []struct {
			Index    int    `json:"index"`
			Upstream string `json:"upstream"`
			Name     string `json:"name"`
			Status   string `json:"status"`
			Rise     int    `json:"rise"`
			Fall     int    `json:"fall"`
			Type     string `json:"type"`
			Port     int    `json:"port"`
		} `json:"server"`
	} `json:"servers"`
}

func PathExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func getUrl() string {

	path := filepath.Dir(os.Args[0])
	conf := path + "/url.conf"
	if !PathExist(conf) {
		//	检测配置文件是否存在
		fmt.Printf("您必须将[url.conf]，放到[%s]目录下!\r\n", path)
		os.Exit(1)
	}
	URL, err := ioutil.ReadFile(conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	return strings.Replace(strings.Replace(string(URL), "\n", "", -1), "\r", "", -1)
}

func getUpstreamStatus(url string) []byte {
	//	输入一个url，返回请求该http得到的json字节集。
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body
}

func Fatal(errMsg string) {
	log.Println(errMsg)
	log.Fatalln("如果您想获取Index为0的Server的Status，您可以这么输入：get-tengine-upstream 0 status")
}

func main() {

	//  获取upstream状态
	statusByte := getUpstreamStatus(getUrl())
	var status NginxUpstream
	json.Unmarshal(statusByte, &status)
	if status.Servers.Total < 1 {
		//	判断负载池中服务器的数量，如果数量小于1，就报错。
		log.Fatal("负载均衡池中服务器数量为0！")
	}

	//	获取命令行参数
	flag.Parse()

	//	根据命令行参数开始判断
	if flag.NArg() < 1 {
		//	如果没有定义命令行参数，就直接输出请求到的json文本，然后退出。
		fmt.Printf("%s", statusByte)
		os.Exit(0)
	}

	//	获取第一个命令行参数
	First, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		//	如果第一个参数输入错误，就报错。
		Fatal("第一个参数必须输入一个大于等于0的数字！")
	}
	if First >= status.Servers.Total {
		//	如果第一个参数的数值超过了负载池中的数量，就报错。
		Fatal(fmt.Sprintf("注意：Index从0开始计数，您输入的第一个参数已经超过了Server的总数量[%d]，允许输入的数值为[0-%d]!", status.Servers.Total, status.Servers.Total-1))
	}

	//	获取第二个命令行参数
	Second := flag.Arg(1)
	//	fmt.Println(status.Servers.Server[First])
	switch strings.ToLower(Second) {
	case "index":
		fmt.Print(status.Servers.Server[First].Index)
	case "upstream":
		fmt.Print(status.Servers.Server[First].Upstream)
	case "name":
		fmt.Print(status.Servers.Server[First].Name)
	case "status":
		//	处理一下status的返回值，将up和down状态转换成更易于处理的数字，up=1,down=0。
		if status.Servers.Server[First].Status == "up" {
			fmt.Print(1)
		} else if status.Servers.Server[First].Status == "down" {
			fmt.Print(0)
		} else {
			fmt.Print(-1)
		}
	case "rise":
		fmt.Print(status.Servers.Server[First].Rise)
	case "fall":
		fmt.Print(status.Servers.Server[First].Fall)
	case "type":
		fmt.Print(status.Servers.Server[First].Type)
	case "port":
		fmt.Print(status.Servers.Server[First].Port)
	default:
		fmt.Println(status.Servers.Server[First])
		fmt.Println()
		Fatal("您输入的第二个参数必须是这几个中的一个[index|upstream|name|status|rise|fall|type|port]!")
	}
}
