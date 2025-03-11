/*
*	Description : 输出数据到文本相关的操作； 其中含csv的相关方法， excel，文本操作等  TODO 测试
*	Author 		: ManGe
*	Mail 		: 2912882908@qq.com
**/

package gathertool

import (
	"archive/tar"
	"archive/zip"
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
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

func (c *Csv) Add(data []string) error {
	err := c.W.Write(data)
	if err != nil {
		return err
	}
	c.W.Flush()
	return nil
}

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

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// OutJsonFile 将data输出到json文件
func OutJsonFile(data any, fileName string) error {
	var (
		f   *os.File
		err error
	)
	if Exists(fileName) { //如果文件存在
		f, err = os.OpenFile(fileName, os.O_APPEND, 0666) //打开文件
	} else {
		f, err = os.Create(fileName) //创建文件
	}
	if err != nil {
		return err
	}
	str, err := Any2Json(data)
	if err != nil {
		return err
	}
	_, err = io.WriteString(f, str)
	if err != nil {
		return err
	}
	return nil
}

// GetAllFile 获取目录下的所有文件
func GetAllFile(pathname string) ([]string, error) {
	s := make([]string, 0)
	rd, err := os.ReadDir(pathname)
	if err != nil {
		Error("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		if !fi.IsDir() {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}

func subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)
	if start < 0 || start > length {
		Error("start is wrong")
		return ""
	}
	if end < start || end > length {
		Error("end is wrong")
		return ""
	}
	return string(rs[start:end])
}

// DeCompressZIP zip解压文件
func DeCompressZIP(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = reader.Close()
	}()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		filename := dest + file.Name
		err = os.MkdirAll(subString(filename, 0, strings.LastIndex(filename, "/")), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		_ = w.Close()
		_ = rc.Close()
	}
	return nil
}

// DeCompressTAR tar 解压文件
func DeCompressTAR(tarFile, dest string) error {
	file, err := os.Open(tarFile)
	if err != nil {
		Error(err)
		return err
	}
	defer func() {
		_ = file.Close()
	}()
	tr := tar.NewReader(file)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		filename := dest + hdr.Name
		err = os.MkdirAll(subString(filename, 0, strings.LastIndex(filename, "/")), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, tr)
		if err != nil {
			return err
		}
		_ = w.Close()
	}
	return nil
}

// DecompressionZipFile zip压缩文件
func DecompressionZipFile(src, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		_ = reader.Close()
	}()
	for _, file := range reader.File {
		filePath := path.Join(dest, file.Name)
		if file.FileInfo().IsDir() {
			_ = os.MkdirAll(filePath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				return err
			}
			inFile, err := file.Open()
			if err != nil {
				return err
			}
			outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
			_ = inFile.Close()
			_ = outFile.Close()
		}
	}
	return nil
}

// CompressFiles 压缩很多文件
// files 文件数组，可以是不同dir下的文件或者文件夹
// dest 压缩文件存放地址
func CompressFiles(files []string, dest string) error {
	d, _ := os.Create(dest)
	defer func() {
		_ = d.Close()
	}()
	w := zip.NewWriter(d)
	defer func() {
		_ = w.Close()
	}()
	for _, file := range files {
		err := compressFiles(file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

func compressFiles(filePath string, prefix string, zw *zip.Writer) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f := file.Name() + "/" + fi.Name()
			err = compressFiles(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = prefix + "/" + header.Name
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		_ = file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// CompressDirZip 压缩目录
func CompressDirZip(src, outFile string) error {
	_ = os.RemoveAll(outFile)
	zipFile, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = zipFile.Close()
	}()
	archive := zip.NewWriter(zipFile)
	defer func() {
		_ = archive.Close()
	}()
	return filepath.Walk(src, func(path string, info os.FileInfo, _ error) error {
		if path == src {
			return nil
		}
		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(path, src+`/`)
		if info.IsDir() {
			header.Name += `/`
		} else {
			header.Method = zip.Deflate
		}
		writer, _ := archive.CreateHeader(header)
		if !info.IsDir() {
			file, _ := os.Open(path)
			defer func() {
				_ = file.Close()
			}()
			_, _ = io.Copy(writer, file)
		}
		return nil
	})
}

func Exists(path string) bool {
	if stat, err := os.Stat(path); stat != nil && !os.IsNotExist(err) {
		return true
	}
	return false
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

// GetNowPath 获取当前运行路径
func GetNowPath() string {
	pathData, err := os.Getwd()
	if err != nil {
		Error(err)
	}
	if runtime.GOOS == "windows" {
		return strings.Replace(pathData, "\\", "/", -1)
	}
	return pathData
}

// FileMd5sum 文件 Md5
func FileMd5sum(fileName string) string {
	fin, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if err != nil {
		Info(fileName, err)
		return ""
	}
	defer func() {
		_ = fin.Close()
	}()
	Buf, bufErr := os.ReadFile(fileName)
	if bufErr != nil {
		Info(fileName, bufErr)
		return ""
	}
	m := md5.Sum(Buf)
	return hex.EncodeToString(m[:16])
}

// PathExists 目录不存在则创建
func PathExists(path string) {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll(path, 0777)
	}
}

// FileSizeFormat 字节的单位转换 保留两位小数
func FileSizeFormat(fileSize int64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

// AbPathByCaller 获取当前执行文件绝对路径（go run）
func AbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return path.Join(abPath, "../../../")
}
