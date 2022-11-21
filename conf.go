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
	"os"

	yaml "gopkg.in/yaml.v3"
)

var Config *conf

type conf struct {
	Path string
	Data map[string]interface{}
}

func NewConf(appConfigPath string) error {
	Config = &conf{
		Path: appConfigPath,
		Data: make(map[string]interface{}),
	}
	err := Config.Init()
	return err
}

func (c *conf) Init() error {
	if !fileExists(c.Path) {
		return fmt.Errorf("未找到配置文件 [%v] !", c.Path)
	}
	Info("读取配置文件:", c.Path)
	//读取yaml文件到缓存中
	config, err := ioutil.ReadFile(c.Path)
	if err != nil {
		Errorf("读取配置文件[%v]失败", c.Path)
		return err
	}
	return yaml.Unmarshal(config, c.Data)
}

func (c *conf) GetInt(key string) int {
	if c.Data == nil {
		_ = c.Init()
	}
	return Any2Int(c.Data[key])
}

func (c *conf) Get(key string) interface{} {
	if c.Data == nil {
		_ = c.Init()
	}
	return c.Data[key]
}

func (c *conf) GetStr(key string) string {
	if c.Data == nil {
		_ = c.Init()
	}
	return Any2String(c.Data[key])
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
