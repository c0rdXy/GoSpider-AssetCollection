package output

import (
	"github.com/c0rdXy/GoSpider-AssetCollection/spider"
	"log"
)

func FromNameOutputToMysql(companyName string) {
	spider.SpiderCompanyInfo(companyName)
}

func FromNameOutputToExcel(companyName string) {

}

func FromNameOutputToTerminal(companyName string) {
	companys := spider.SpiderCompanyInfo(companyName)
	// 打印结果或存储到数据库/文件
	for _, c := range companys {
		log.Printf("公司名称: %s, 投资比例: %s, 经营状态: %s, 产品名称: %s\n", c.Name, c.Percent, c.RegStatus, c.ProductName)
	}
}

func FromExcelOutputToExcel() {

}

func FromExcelOutputToMysql() {

}

func FromExcelOutputToTerminal() {

}
