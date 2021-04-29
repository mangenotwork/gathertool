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

func NewCsv(fileName string) *Csv {
	//创建文件
	f, err := os.Create(fileName)
	if err == nil {
		// 写入UTF-8 BOM
		_,_ = f.WriteString("\xEF\xBB\xBF")
	}
	if err != nil {
		loger(err.Error())
	}
	defer f.Close()

	csvObj := &Csv{FileName: fileName}

	csvObj.W = csv.NewWriter(f)
	csvObj.R = csv.NewReader(f)

	return csvObj
}

func (c *Csv) Add(data []string) error{
	return c.W.Write(data)
}

func (c *Csv) ReadAll() ([][]string, error){
	return c.R.ReadAll()
}