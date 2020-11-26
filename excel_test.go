package main

import (
	"testing"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

// 简单的结构化文件推荐csv文件，体积小很多，记事本可打开，但不支持excel高级功能，不支持sheet
func TestExcel(t *testing.T) {
	f := excelize.NewFile()
	// 创建一个工作表
	index := f.NewSheet("Sheet2")
	// 设置单元格的值
	f.SetCellValue("Sheet2", "A2", "Hello world.")
	f.SetCellValue("Sheet1", "B2", 100)
	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	s := []string{"A", "B"}
	f.SetSheetRow("Sheet1", "A1", &s)
	// 根据指定路径保存文件
	if err := f.SaveAs("Book1.xlsx"); err != nil {
		t.Error(err)
	}
}

func TestExcel2(t *testing.T) {
	f := excelize.NewFile()
	streamWriter, _ := f.NewStreamWriter("Sheet1")
	streamWriter.SetRow("A1", []interface{}{"用户内部ID", "外部ID"})
	streamWriter.SetRow("A2", []interface{}{"A", "B"})

	streamWriter.Flush()
	err := f.SaveAs("Book2.xlsx")
	if err != nil {
		t.Error(err)
	}
}
