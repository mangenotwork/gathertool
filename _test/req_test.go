package test

import (
	gt "github.com/mangenotwork/gathertool"
	"testing"
)

func TestReq(t *testing.T) {
	ctx, err := gt.Request("https://www.doubao.com/", "get", []byte(""), "")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ctx.Html)

	ctx = gt.NewRequest("https://www.doubao.com/", "get", []byte(""), "")
	if ctx.Err != nil {
		t.Fatal(err)
	}
	t.Log(ctx.Html)

	t.Log(gt.IsUrl("https://www.doubao.com/"))
	t.Log(gt.UrlStr("www.doubao.com/"))

	gt.SetSleepMs(1, 3)

	ck := gt.NewCookie(map[string]string{
		"case1": "1",
		"case2": "2",
	})
	t.Log(ck.Delete("case1"))
	t.Log(ck.Get("case1"))

	gt.SearchDomain("180.97.248.210")

	gt.SearchPort("10.0.40.3")

	gt.ClosePingTerminalPrint()

	gt.GetCertificateInfo("https://www.doubao.com/")

	gt.SetStatusCodeSuccessEvent(203)
	gt.SetStatusCodeRetryEvent(401)
	gt.SetStatusCodeFailEvent(403)
	gt.SetAgent(1, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")
}
