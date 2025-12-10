package test

import (
	gt "github.com/mangenotwork/gathertool"
	"testing"
)

func TestSocketProxy(t *testing.T) {
	gt.SocketProxy("0.0.0.0:16666")
}
