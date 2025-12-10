package test

import (
	gt "github.com/mangenotwork/gathertool"
	"testing"
)

// go test -v -run=TestConf
func TestConf(t *testing.T) {
	_ = gt.NewConf("./conf.yaml")
	t.Log(gt.Config.Get("test"))
	t.Log(gt.Config.GetString("test::case"))
	t.Log(gt.Config.GetAll())
}
