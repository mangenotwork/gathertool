/*
	Description : csv的相关方法
	Author : ManGe
	Mail : 2912882908@qq.com
	Github : https://github.com/mangenotwork/gathertool
*/

package gathertool

import (
	"encoding/csv"
	"log"
	"os"
)

// Csv Csv格式文件
type Csv struct {
	FileName string
	W        *csv.Writer
	R        *csv.Reader
}

// NewCSV 新创建一个csv对象
func NewCSV(fileName string) (*Csv, error) {
	f, err := os.Create(fileName)
	if err != nil {
		Error(err.Error())
		return nil, err
	}
	_, _ = f.WriteString("\xEF\xBB\xBF")
	csvObj := &Csv{FileName: fileName}
	csvObj.W = csv.NewWriter(f)
	csvObj.R = csv.NewReader(f)
	return csvObj, nil
}

func (c *Csv) Close() {
	c.Close()
}

// Add 写入csv
func (c *Csv) Add(data []string) error {
	log.Println("写入csv = ", data)
	err := c.W.Write(data)
	if err != nil {
		return err
	}
	c.W.Flush()
	return nil
}

// ReadAll 读取所有
func (c *Csv) ReadAll() ([][]string, error) {
	return c.R.ReadAll()
}

// ReadCsvFile csv file -> [][]string 行列
func ReadCsvFile(filename string) [][]string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = file.Close()
	}()
	reader := csv.NewReader(file)
	reader.Comma = ','

	allRecords, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	return allRecords
}
