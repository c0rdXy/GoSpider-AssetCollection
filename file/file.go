package file

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/c0rdXy/GoSpider-AssetCollection/spider"
)

func ReadExcel(filePath string) ([]string, error) {
	var companies []string

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	rows := f.GetRows(`Sheet1`)

	for _, row := range rows {
		companies = append(companies, row[0])
	}

	return companies, nil
}

func WriteToExcel(companies []spider.CompanyInfo) error {
	f := excelize.NewFile()
	index := f.NewSheet("Sheet1")

	f.SetCellValue("Sheet1", "A1", "公司名称")
	f.SetCellValue("Sheet1", "B1", "对外资产公司名称")
	f.SetCellValue("Sheet1", "C1", "对外资产公司对应的关联产品名")

	for i, info := range companies {
		row := i + 2
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), info.Company)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), info.Source)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), info.Product)
	}

	f.SetActiveSheet(index)

	err := f.SaveAs("result.xlsx")
	return err
}
