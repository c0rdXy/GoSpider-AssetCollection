package output

import (
	"fmt"
	"github.com/c0rdXy/GoSpider-AssetCollection/file"
	"github.com/c0rdXy/GoSpider-AssetCollection/mysql"
	"github.com/c0rdXy/GoSpider-AssetCollection/spider"
	"github.com/xuri/excelize/v2"
	"log"
	"sync"
)

// GoSpiderInit GoSpider初始化
func GoSpiderInit() {
	text := `
 ________           _________        .__     .___
/  _____/   ____   /   _____/______  |__|  __| _/  ____  _______
/   \  ___  /  _ \  \_____  \ \____ \ |  | / __ | _/ __ \ \_  __ \
\    \_\  \(  <_> ) /        \|  |_> >|  |/ /_/ | \  ___/  |  | \/
 \______  / \____/ /_______  /|   __/ |__|\____ |  \___  > |__|
        \/                 \/ |__|             \/      \/
`

	fmt.Println(text)

	fmt.Println("GoSpider - 一个简单的Go语言爬虫工具")
	fmt.Println("作者: c0rdXy	版本: v2.0.1")
	fmt.Println()
}

// FromNameOutputToMysql 从公司名称输出到MySQL数据库
func FromNameOutputToMysql(companyName string) {
	// 初始化数据库
	mysql.InitDB()
	defer mysql.CloseDB() // 程序结束时关闭数据库连接

	id, err := mysql.InsertSource(companyName)
	if err != nil {
		log.Println("插入数据库失败:", err)
	}

	spiderCompanies := spider.SpiderCompanyInfo(companyName)
	for _, info := range spiderCompanies {
		err := mysql.InsertInvest(id, info.ProductName)
		if err != nil {
			log.Println("插入数据库失败:", err)
		}
	}

}

// FromNameOutputToExcel 从公司名称输出到Excel文件
func FromNameOutputToExcel(companyName string) {
	var company = spider.CompanyInfo{}
	var companies []spider.CompanyInfo

	spiderCompanies := spider.SpiderCompanyInfo(companyName)
	for _, info := range spiderCompanies {
		company.Company = companyName
		company.Source = info.Name
		company.Product = info.ProductName
		companies = append(companies, company)
	}

	fileName, err := file.WriteToExcel(companies)
	if err != nil {
		log.Println("写入Excel文件失败:", err)
	} else {
		log.Println("写入Excel文件成功：", fileName)
	}

}

// FromNameOutputToTerminal 从公司名称输出到命令行
func FromNameOutputToTerminal(companyName string) {
	companies := spider.SpiderCompanyInfo(companyName)
	// 打印结果或存储到数据库/文件
	for _, c := range companies {
		log.Printf("公司名称: %s, 投资公司名称: %s, 投资比例: %s, 经营状态: %s, 产品名称: %s\n", companyName, c.Name, c.Percent, c.RegStatus, c.ProductName)
	}
}

// FromExcelOutputToExcel 从Excel文件输出到Excel文件
func FromExcelOutputToExcel(filePath string) {
	excelCompanies, err := file.ReadExcel(filePath)
	if err != nil {
		log.Fatal("读取Excel文件失败:", err)
	}

	// 创建通道用于传递爬取结果
	resultChannel := make(chan []spider.CompanyInfo)
	// 使用 sync.WaitGroup 等待所有协程完成
	var wg sync.WaitGroup
	var lock sync.Mutex

	// 增加 WaitGroup 计数
	wg.Add(2)

	// 启动协程进行爬取并将结果发送到通道
	go func() {
		defer wg.Done()
		defer close(resultChannel) // 在爬取协程完成后关闭通道
		for _, companyName := range excelCompanies {
			var companyInfos []spider.CompanyInfo
			spiderCompanies := spider.SpiderCompanyInfo(companyName)

			log.Println("爬取公司名称: ", companyName)

			if len(spiderCompanies) == 0 {
				companyInfos = append(companyInfos, spider.CompanyInfo{
					Company: companyName,
					Source:  "",
					Product: "",
				})
			} else {
				for _, externalCompany := range spiderCompanies {
					companyInfos = append(companyInfos, spider.CompanyInfo{
						Company: companyName,
						Source:  externalCompany.Name,
						Product: externalCompany.ProductName,
					})
				}
			}

			//fmt.Println(companyInfos)
			resultChannel <- companyInfos
		}
	}()

	// 启动协程进行数据写入
	go func() {
		defer wg.Done()
		var f *excelize.File

		for companies := range resultChannel {
			lock.Lock()
			if f == nil {
				// 创建 Excel 文件
				fileName, err := file.CreateExcelFile()
				if err != nil {
					log.Fatal("创建Excel文件失败:", err)
				}

				log.Println("创建Excel文件成功：", fileName)

				// 打开 Excel 文件
				f, err = excelize.OpenFile(fileName)
				if err != nil {
					log.Fatal("打开Excel文件失败:", err)
				}

				fmt.Println(fileName)
			}

			// 写入数据到 Excel 文件
			f, err = file.AppendToExcel(companies, f)
			lock.Unlock()

			if err != nil {
				log.Fatal("追加数据到Excel文件失败:", err)
			}
		}

		// 关闭 Excel 文件
		err := f.Save()
		if err != nil {
			log.Fatal("保存Excel文件失败:", err)
		}
	}()

	// 主协程等待所有协程完成
	wg.Wait()

}

// FromExcelOutputToMysql 从Excel文件输出到MySQL数据库
func FromExcelOutputToMysql(filePath string) {
	// 初始化数据库
	mysql.InitDB()
	defer mysql.CloseDB() // 程序结束时关闭数据库连接

	excelCompany, err := file.ReadExcel(filePath)
	if err != nil {
		log.Println("读取Excel文件失败:", err)
	} else {
		log.Println("读取Excel文件成功")
	}

	for _, companyName := range excelCompany {

		id, err := mysql.InsertSource(companyName)
		if err != nil {
			log.Println("插入数据库失败:", err)
		}

		spiderCompanies := spider.SpiderCompanyInfo(companyName)
		for _, info := range spiderCompanies {
			err := mysql.InsertInvest(id, info.ProductName)
			if err != nil {
				log.Println("插入数据库失败:", err)
			}
		}
	}
}

// FromExcelOutputToTerminal 从Excel文件输出到命令行
func FromExcelOutputToTerminal(filePath string) {
	excelCompany, err := file.ReadExcel(filePath)
	if err != nil {
		log.Println("读取Excel文件失败:", err)
	} else {
		log.Println("读取Excel文件成功")
	}

	for _, companyName := range excelCompany {
		companies := spider.SpiderCompanyInfo(companyName)
		// 打印结果或存储到数据库/文件
		for _, c := range companies {
			log.Printf("公司名称: %s, 投资公司名称: %s, 投资比例: %s, 经营状态: %s, 产品名称: %s\n", companyName, c.Name, c.Percent, c.RegStatus, c.ProductName)
		}
	}
}
