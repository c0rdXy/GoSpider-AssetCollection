package file

import (
	"fmt"
	"github.com/c0rdXy/GoSpider-AssetCollection/spider"
	"github.com/xuri/excelize/v2"
	"time"
)

// ReadExcel 读取Excel文件
func ReadExcel(filePath string) ([]string, error) {
	var companies []string

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	rows, err := f.GetRows(`Sheet1`)
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		companies = append(companies, row[0])
	}

	return companies, nil
}

// WriteToExcel 写入Excel文件
func WriteToExcel(companies []spider.CompanyInfo) (string, error) {
	f := excelize.NewFile()
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		return "", err
	}

	f.SetCellValue("Sheet1", "A1", "公司名称")
	f.SetCellValue("Sheet1", "B1", "对外资产公司名称")
	f.SetCellValue("Sheet1", "C1", "对外资产公司对应的关联产品名")

	f.SetActiveSheet(index)

	for i, info := range companies {
		row := i + 2
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), info.Company)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), info.Source)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), info.Product)
	}

	fileName, err := CloseExcelFile(f)
	if err != nil {
		return "", err
	}

	return fileName, err
}

// CloseExcelFile 关闭Excel文件
func CloseExcelFile(f *excelize.File) (string, error) {
	// 使用当前时间作为文件名的一部分
	timeFormat := time.Now().Format("20060102_150405") // 格式化时间
	fileName := fmt.Sprintf("result_%s.xlsx", timeFormat)

	err := f.SaveAs(fileName)
	if err != nil {
		return "", err
	}
	return fileName, err
}

// AppendToExcel 追加数据到Excel文件
func AppendToExcel(companies []spider.CompanyInfo, filePath string) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}

	// 获取 Sheet1 的最后一行
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return err
	}

	lastRow := len(rows) + 1

	// 追加新数据到 Sheet1
	for i, info := range companies {
		row := lastRow + i
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), info.Company)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), info.Source)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), info.Product)
	}

	CloseExcelFile(f)

	return err
}
