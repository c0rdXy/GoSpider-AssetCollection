package spider

import (
	"fmt"
	"github.com/c0rdXy/GoSpider-AssetCollection/mysql"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type CompanyInfo struct {
	Company string
	Source  string
	Product string
}

var wg sync.WaitGroup

func SpiderCompanyInfo(companyName string) {
	defer wg.Done()

	// 实现爬虫逻辑
	companyInfo, err := crawlCompanyInfo(companyName)
	if err != nil {
		log.Println("爬取公司信息失败:", err)
		return
	}

	// 插入数据库
	sourceID, err := mysql.InsertSource(companyInfo.Company)
	if err != nil {
		log.Println("插入公司信息失败:", err)
		return
	}

	if companyInfo.Source != "" {
		err = mysql.InsertInvest(sourceID, companyInfo.Product)
		if err != nil {
			log.Println("插入对外资产信息失败:", err)
			return
		}
	}

	fmt.Printf("公司 %s 爬取成功\n", companyName)
}

func crawlCompanyInfo(companyName string) (*CompanyInfo, error) {
	// 这里根据实际情况填写爬虫逻辑
	// 以下是示例代码，您需要根据目标网站的实际结构和数据获取方式进行修改

	url := fmt.Sprintf("https://www.tianyancha.com/search?key=%s", companyName)

	// 发送 HTTP 请求获取网页内容
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP请求失败，状态码：%d", resp.StatusCode)
	}

	// 使用 goquery 解析 HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// 提取公司名称
	company := doc.Find(".search-item .header a").Text()

	// 提取对外资产公司名称（投资占比50%以上）
	var source, product string
	doc.Find(".out-investment-list .header").Each(func(i int, s *goquery.Selection) {
		investmentText := s.Text()
		percentageText := s.Parent().Find(".percentage").Text()

		// 正则表达式匹配百分比
		percentageRegex := regexp.MustCompile(`(\d+)%`)
		matches := percentageRegex.FindStringSubmatch(percentageText)

		if len(matches) > 1 {
			percentage, _ := strconv.Atoi(matches[1])
			if percentage >= 50 {
				// 提取对外资产公司名称和关联产品名
				source = investmentText
				product = s.Parent().Find(".product").Text()
			}
		}
	})

	return &CompanyInfo{
		Company: company,
		Source:  source,
		Product: product,
	}, nil
}
