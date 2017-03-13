package main

import (
	"time"
	"github.com/Andrew512512/nginx_log_analy/tools"
	"net"
	"errors"
	"fmt"
)

var (
	filePath string = "/usr/local/openresty/nginx/logs/access.log"
	endPoint string
	prefix string = "test_v5_"
)

func main() {
	var err error
	tools.InitCache()
	tools.InitFileChecker(filePath)

	endPoint, err = init_ip()
	if err != nil {
		panic(err)
	}
	fmt.Println("当前服务器endPoint为", endPoint)

	go fileTimer()
	go uploadTimer()

	//酱油主线程
	for true {
		time.Sleep(60 * time.Second)
	}
}

//日志文件变动扫描定时器
func fileTimer() {
	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		for {
			<-ticker.C
			tools.CheckFileOnce(filePath)
		}
	}()
}

//上传定时器
func uploadTimer() {
	ticker := time.NewTicker(60 * time.Second)
	go func() {
		for {
			<-ticker.C
			tools.UploadOnce(endPoint, prefix)
		}
	}()
}

//根据本机局域网ip匹配endPoint
func init_ip() (string, error) {
	ip_list := make(map[string]string)
	ip_list["172.16.1.11"] = "www"
	ip_list["172.16.1.12"] = "www.chat"
	ip_list["172.16.1.13"] = "www.static"
	ip_list["172.16.1.31"] = "www.exam"
	ip_list["172.16.1.34"] = "dev"
	ip_list["172.16.1.35"] = "dev.chat"
	ip_list["172.16.1.37"] = "dev.exam"
	ip_list["172.16.1.67"] = "www.device"

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				endPoint = ip_list[ipnet.IP.String()]
				if len(endPoint) > 0 {
					return endPoint, nil
				}
				break
			}
		}
	}
	return "", errors.New("未正确匹配ip")
}