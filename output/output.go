package output

import (
	"fmt"
	"github.com/c0rdXy/GoSpider-AssetCollection/file"
	"github.com/c0rdXy/GoSpider-AssetCollection/mysql"
	"github.com/c0rdXy/GoSpider-AssetCollection/spider"
	"log"
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
	fmt.Println("作者: xxx	版本: 1.0.1")
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

	err := file.WriteToExcel(companies)
	if err != nil {
		log.Println("写入Excel文件失败:", err)
	}

}

// FromNameOutputToTerminal 从公司名称输出到命令行
func FromNameOutputToTerminal(companyName string) {
	companies := spider.SpiderCompanyInfo(companyName)
	// 打印结果或存储到数据库/文件
	for _, c := range companies {
		log.Printf("公司名称: %s, 投资比例: %s, 经营状态: %s, 产品名称: %s\n", c.Name, c.Percent, c.RegStatus, c.ProductName)
	}
}

// FromExcelOutputToExcel 从Excel文件输出到Excel文件
func FromExcelOutputToExcel(filePath string) {
	excelCompanies, err := file.ReadExcel(filePath)
	if err != nil {
		log.Println("读取Excel文件失败:", err)
	}

	var company = spider.CompanyInfo{}
	var companies []spider.CompanyInfo

	for _, excelCompany := range excelCompanies {

		spiderCompanies := spider.SpiderCompanyInfo(excelCompany)
		for _, info := range spiderCompanies {
			company.Company = excelCompany
			company.Source = info.Name
			company.Product = info.ProductName
			companies = append(companies, company)
		}

		err := file.WriteToExcel(companies)
		if err != nil {
			log.Println("写入Excel文件失败:", err)
		}
	}

}

// FromExcelOutputToMysql 从Excel文件输出到MySQL数据库
func FromExcelOutputToMysql(filePath string) {
	// 初始化数据库
	mysql.InitDB()
	defer mysql.CloseDB() // 程序结束时关闭数据库连接

	excelCompany, err := file.ReadExcel(filePath)
	if err != nil {
		log.Println("读取Excel文件失败:", err)
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
	}

	for _, companyName := range excelCompany {
		companies := spider.SpiderCompanyInfo(companyName)
		// 打印结果或存储到数据库/文件
		for _, c := range companies {
			log.Printf("公司名称: %s, 投资比例: %s, 经营状态: %s, 产品名称: %s\n", c.Name, c.Percent, c.RegStatus, c.ProductName)
		}
	}
}
