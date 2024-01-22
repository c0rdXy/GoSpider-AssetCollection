package main

import (
	"flag"
	"fmt"
	"github.com/c0rdXy/GoSpider-AssetCollection/output"
	"os"
)

var (
	companyName  string
	filePath     string
	outputMethod string
)

func init() {
	output.GoSpiderInit()

	flag.StringVar(&companyName, "company", "", "指定公司名称")
	flag.StringVar(&filePath, "file", "", "指定Excel文件路径")
	flag.StringVar(&outputMethod, "om", "",
		`指定输出方式（不指定默认命令行输出）:
				-om excel 存入到Excel文件
				-om mysql 存入到MySQL数据库`)
	flag.Parse()

}

func main() {
	if companyName != "" && filePath == "" {
		if outputMethod == "excel" {
			output.FromNameOutputToExcel(companyName)
		} else if outputMethod == "mysql" {
			output.FromNameOutputToMysql(companyName)
		} else {
			output.FromNameOutputToTerminal(companyName)
		}
	} else if filePath != "" && companyName == "" {
		if outputMethod == "excel" {
			output.FromExcelOutputToExcel(filePath)
		} else if outputMethod == "mysql" {
			output.FromExcelOutputToMysql(filePath)
		} else {
			output.FromExcelOutputToTerminal(filePath)
		}
	} else {
		fmt.Println("请使用 -company 或 -file 参数指定操作，使用 -om 参数指定输出方式（不指定默认命令行输出）")
		os.Exit(1)
	}
}
