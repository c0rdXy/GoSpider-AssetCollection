package spider

import (
	"encoding/json"
	"fmt"
	"github.com/c0rdXy/GoSpider-AssetCollection/http"
	"log"
	"strconv"
	"strings"
)

// CompanyInfo 结构体定义
type CompanyInfo struct {
	Company string
	Source  string
	Product string
}

// SearchResult 结构体定义
type SearchResult struct {
	State   string            `json:"state"`
	Message string            `json:"message"`
	Data    []CompanyListInfo `json:"data"`
	// 其他字段...
}

// CompanyListInfo 结构体定义
type CompanyListInfo struct {
	ID      int64  `json:"id"`
	ComName string `json:"comName"`
}

// ExternalCompanyResult 结构体定义
type ExternalCompanyResult struct {
	State   string                  `json:"state"`
	Message string                  `json:"message"`
	Data    ExternalCompanyDataInfo `json:"data"`
	// 其他字段...
}

// ExternalCompanyDataInfo 结构体定义
type ExternalCompanyDataInfo struct {
	Result []ExternalCompanyInfo `json:"result"`
	// 其他字段...
}

// ExternalCompanyInfo 结构体定义
type ExternalCompanyInfo struct {
	Name        string `json:"name"`
	Percent     string `json:"percent"`
	RegStatus   string `json:"regStatus"`
	ProductName string `json:"productName"`
	// 其他字段...
}

// SpiderCompanyInfo 爬取公司信息
func SpiderCompanyInfo(companyName string) []ExternalCompanyInfo {
	// 发送 POST 请求获取公司列表
	companyList, err := searchCompanyList(companyName)
	if err != nil {
		log.Println("获取公司列表失败:", err)
		return nil
	}

	for _, company := range companyList {

		// 获取公司的 ID
		if company.ComName == companyName {
			companyID := company.ID

			externalCompanies, err := getExternalCompanies(companyID)
			if err != nil {
				log.Println("获取对外资产公司列表失败:", err)
				return nil
			}

			// 过滤投资比例大于50%且经营状态为存续的公司
			filteredCompanies := filterCompanies(externalCompanies)

			return filteredCompanies
		}
	}

	return nil
}

// searchCompanyList 发送 POST 请求获取公司列表
func searchCompanyList(companyName string) ([]CompanyListInfo, error) {
	url := "https://capi.tianyancha.com/cloud-tempest/search/suggest/v3"
	payload := map[string]string{"keyword": companyName}

	resp, err := http.MyHTTPPost(url, payload)
	if err != nil {
		return nil, err
	}

	//fmt.Println(string(resp))

	var response SearchResult
	if err := json.Unmarshal(resp, &response); err != nil {
		return nil, err
	}

	if response.State != "ok" {
		return nil, fmt.Errorf("获取公司列表失败，错误信息：%s", response.Message)
	}
	return response.Data, nil
}

// getExternalCompanies 获取对外资产公司列表
func getExternalCompanies(companyID int64) ([]ExternalCompanyInfo, error) {
	url := "https://capi.tianyancha.com/cloud-company-background/company/investListV2"
	payload := map[string]interface{}{
		"gid":          strconv.FormatInt(companyID, 10),
		"pageSize":     10,
		"pageNum":      1,
		"province":     "-100",
		"percentLevel": "-100",
		"category":     "-100",
	}

	body, err := http.MyHTTPPost(url, payload)
	if err != nil {
		return nil, err
	}

	var response ExternalCompanyResult
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if response.State != "ok" {
		return nil, fmt.Errorf("获取对外资产公司列表失败，错误信息：%s", response.Message)
	}

	return response.Data.Result, nil
}

// filterCompanies 过滤投资比例大于50%且经营状态为存续的公司
func filterCompanies(companies []ExternalCompanyInfo) []ExternalCompanyInfo {
	var filteredCompanies []ExternalCompanyInfo
	for _, c := range companies {
		// 去除百分号
		str := strings.TrimRight(c.Percent, "%")

		percent, err := strconv.Atoi(str)
		if err != nil {
			//log.Printf("无法解析投资比例：%s\n", c.Percent)
			continue
		}

		// 过滤投资比例大于50%且经营状态为存续的公司
		if percent > 50 && (strings.Contains(c.RegStatus, "存续") || strings.Contains(c.RegStatus, "在营") || strings.Contains(c.RegStatus, "开业") || strings.Contains(c.RegStatus, "在册")) {
			filteredCompanies = append(filteredCompanies, c)
		}
	}
	return filteredCompanies
}
