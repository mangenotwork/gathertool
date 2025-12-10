package test

import (
	gt "github.com/mangenotwork/gathertool"
	"testing"
)

func TestStressTesting(t *testing.T) {
	test := gt.NewTestUrl("https://baidu.com", "Get", 20, 10)
	test.Run()
}

// TODO
func TestScanNewHostScanUrl(t *testing.T) {

}
