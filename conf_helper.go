/*
*	Description : 配置文件 yaml  TODO 扩展到所有能支持的配置文件
*	Author 		: ManGe
*	Mail 		: 2912882908@qq.com
**/

package gathertool

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var Config *conf

type conf struct {
	Path string
	Data map[string]any
}

// NewConf 读取配置，只支持yaml
func NewConf(appConfigPath string) error {
	Config = &conf{
		Path: appConfigPath,
		Data: make(map[string]any),
	}
	err := Config.Init()
	return err
}

func (c *conf) Init() error {
	if !FileExists(c.Path) {
		return fmt.Errorf("未找到配置文件: %v", c.Path)
	}
	Info("读取配置文件:", c.Path)
	//读取yaml文件到缓存中
	config, err := os.ReadFile(c.Path)
	if err != nil {
		ErrorF("读取配置文件 %v 失败", c.Path)
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

func (c *conf) get(key string) (interface{}, bool) {
	var (
		d  interface{}
		ok bool
	)
	keyList := strings.Split(key, "::")
	temp := make(map[string]interface{})
	temp = c.Data
	for _, v := range keyList {
		d, ok = temp[v]
		if !ok {
			break
		}
		temp = Any2Map(d)
	}
	return d, ok
}

// Get  获取配置,多级使用::，例如：user::name
func (c *conf) Get(key string) (any, bool) {
	if c.Data == nil {
		_ = c.Init()
	}
	return c.get(key)
}

// GetString  获取配置,配置不存在返回 "", false
// ex: conf.GetString("user::name")
func (c *conf) GetString(key string) (string, bool) {
	if c.Data == nil {
		_ = c.Init()
	}
	val, has := c.get(key)
	if !has {
		return "", has
	}
	return Any2String(val), has
}

func (c *conf) GetAll() map[string]any {
	return c.Data
}
