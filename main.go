package main

import (
	"flag"
	"fmt"
	"github.com/c0rdXy/GoSpider-AssetCollection/mysql"
	"github.com/c0rdXy/GoSpider-AssetCollection/output"
	"os"
)

var (
	companyName  string
	filePath     string
	outputMethod string
)

func init() {

	flag.StringVar(&companyName, "company", "", "指定公司名称")
	flag.StringVar(&filePath, "file", "", "指定Excel文件路径")
	flag.StringVar(&outputMethod, "om", "",
		`指定输出方式:
				-om excel 存入到Excel文件
				-om mysql 存入到MySQL数据库`)
	flag.Parse()

	// 初始化数据库
	mysql.InitDB()
	defer mysql.CloseDB() // 程序结束时关闭数据库连接
}

func main() {
	//spider.SpiderCompanyInfo("深圳市腾讯计算机系统有限公司")

	if companyName != "" && filePath == "" {
		if outputMethod == "excel" {

		} else if outputMethod == "mysql" {

		} else {
			output.FromNameOutputToTerminal(companyName)
		}
	} else if filePath != "" && companyName == "" {
		if outputMethod == "excel" {

		} else if outputMethod == "mysql" {

		} else {

		}
	} else {
		fmt.Println("请使用 -company 或 -file 参数指定操作，使用 -om 参数指定输出方式（不指定默认命令行输出）")
		os.Exit(1)
	}
}
