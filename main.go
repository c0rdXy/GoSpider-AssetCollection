package main

import (
	"flag"
	"fmt"
	"github.com/c0rdXy/GoSpider-AssetCollection/file"
	"github.com/c0rdXy/GoSpider-AssetCollection/mysql"
	"github.com/c0rdXy/GoSpider-AssetCollection/spider"
	"os"
	"sync"
)

func main() {
	var (
		companyName string
		filePath    string
	)

	flag.StringVar(&companyName, "company", "", "指定公司名称")
	flag.StringVar(&filePath, "file", "", "指定Excel文件路径")
	flag.Parse()

	// 初始化数据库
	mysql.InitDB()
	defer mysql.CloseDB() // 程序结束时关闭数据库连接

	if companyName != "" {
		spider.SpiderCompanyInfo(companyName)
	} else if filePath != "" {
		batchSpiderFromExcel(filePath)
	} else {
		fmt.Println("请使用 -company 或 -file 参数指定操作")
		os.Exit(1)
	}
}

func batchSpiderFromExcel(filePath string) {
	companies, err := file.ReadExcel(filePath)
	if err != nil {
		fmt.Println("读取Excel失败:", err)
		return
	}

	var wg sync.WaitGroup
	for _, company := range companies {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			spider.SpiderCompanyInfo(c)
		}(company)
	}

	wg.Wait()
	fmt.Println("批量爬取成功")
}
