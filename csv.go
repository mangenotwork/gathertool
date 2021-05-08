/*
	Description : csv的相关方法
	Author : ManGe
	Version : v0.1
	Date : 2021-04-27
*/

package gathertool

import (
	"encoding/csv"
	"os"
)

type Csv struct {
	FileName string
	W *csv.Writer
	R *csv.Reader
}

// 新创建一个csv对象
func NewCsv(fileName string) (*Csv,error) {
	//创建文件
	f, err := os.Create(fileName)
	defer f.Close()
	if err != nil {
		loger(err.Error())
		return nil,err
	}
	_,_ = f.WriteString("\xEF\xBB\xBF")
	csvObj := &Csv{FileName: fileName}
	csvObj.W = csv.NewWriter(f)
	csvObj.R = csv.NewReader(f)

	return csvObj,nil
}

// 写入csv
func (c *Csv) Add(data []string) error{
	return c.W.Write(data)
}

// 读取所有
func (c *Csv) ReadAll() ([][]string, error){
	return c.R.ReadAll()
}

// csv file -> [][]string 行列
func ReadCsvFile(filename string) [][]string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.Comma = ','

	allRecords, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	return allRecords
}

