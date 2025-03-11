package test

import (
	gt "github.com/mangenotwork/gathertool"
	"testing"
)

func TestApplicationTerminalOutClose(t *testing.T) {
	gt.ApplicationTerminalOutClose()
}

func TestNewHostScanUrl(t *testing.T) {
	scan := gt.NewHostScanUrl("http://127.0.0.1:8881/home", 3)
	gt.Info(scan.Run())
}

func TestNewHostScanExtLinks(t *testing.T) {
	scan := gt.NewHostScanExtLinks("http://127.0.0.1:8881/home")
	gt.Info(scan.Run())
}

func TestNewHostScanBadLink(t *testing.T) {
	scan := gt.NewHostScanBadLink("http://127.0.0.1:8881/home", 3)
	gt.Info(scan.Run())
}

func TestNewHostPageSpeedCheck(t *testing.T) {
	scan := gt.NewHostPageSpeedCheck("http://127.0.0.1:8881/home", 3)
	gt.Info(scan.Run())
}
