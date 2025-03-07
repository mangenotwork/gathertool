package test

import (
	gt "github.com/mangenotwork/gathertool"
	"testing"
)

func TestPost(t *testing.T) {
	ctx, err := gt.Post("http://127.0.0.1/post/test", []byte("xx"), "application/json")
	if err != nil {
		t.Error(err)
	}
	gt.Info(ctx.Json)
}

func TestPostForm(t *testing.T) {
	ctx, err := gt.PostForm("http://127.0.0.1/postFrom/test", map[string]string{"xx": "xx"})
	if err != nil {
		t.Error(err)
	}
	gt.Info(ctx.RespBody)
}

func TestPostFile(t *testing.T) {
	ctx, err := gt.PostFile("http://127.0.0.1/postFile/test", "file", "./test.txt")
	if err != nil {
		t.Error(err)
	}
	gt.Info(ctx.RespBody)
}

func TestPut(t *testing.T) {
	ctx, err := gt.Put("http://127.0.0.1/put/test", []byte("xx"), "application/json")
	if err != nil {
		t.Error(err)
	}
	gt.Info(ctx.RespBody)
}

func TestDelete(t *testing.T) {
	ctx, err := gt.Delete("http://127.0.0.1/delete/test")
	if err != nil {
		t.Error(err)
	}
	gt.Info(ctx.RespBody)
}

func TestOptions(t *testing.T) {
	ctx, err := gt.Options("http://127.0.0.1/options/test")
	if err != nil {
		t.Error(err)
	}
	gt.Info(ctx.RespBody)
}

func TestHead(t *testing.T) {
	ctx, err := gt.Head("http://127.0.0.1/head/test")
	if err != nil {
		t.Error(err)
	}
	gt.Info(ctx.RespBody)
}

func TestPatch(t *testing.T) {
	ctx, err := gt.Patch("http://127.0.0.1/patch/test")
	if err != nil {
		t.Error(err)
	}
	gt.Info(ctx.RespBody)
}

func TestNewProxyIP(t *testing.T) {
	gt.NewProxyIP("127.0.0.1", 8888, "", "", false)
}

func TestNewProxyPool(t *testing.T) {
	gt.NewProxyPool()
}

func TestNewCookiePool(t *testing.T) {
	gt.NewCookiePool()
}
