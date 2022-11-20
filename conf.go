/*
	Description : 配置文件 yaml
	Author : ManGe
			2912882908@qq.com
			https://github.com/mangenotwork/gathertool
*/

package gathertool

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Conf struct {
	Path string
	Data map[string]interface{}
}

func NewConf(appConfigPath string) (*Conf, error) {
	conf := &Conf{
		Path: appConfigPath,
		Data: make(map[string]interface{}),
	}
	err := conf.Init()
	return conf, err
}

func (c *Conf) Init() error {
	if !fileExists(c.Path) {
		return fmt.Errorf("未找到配置文件!")
	}
	log.Println("读取配置文件:", c.Path)
	//读取yaml文件到缓存中
	config, err := ioutil.ReadFile(c.Path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(config, c.Data)
}

func (c *Conf) GetInt(key string) int {
	return c.Data[key].(int)
}

func (c *Conf) Get(key string) interface{} {
	return c.Data[key]
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
