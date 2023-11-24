package test

import (
	"github.com/mangenotwork/gathertool"
	"testing"
)

func TestStressTesting(t *testing.T) {
	test := gathertool.NewTestUrl("https://baidu.com", "Get", 20, 10)
	test.Run()
}

// TODO
func TestScanNewHostScanUrl(t *testing.T) {

}
