package test

import (
	gt "github.com/mangenotwork/gathertool"
	"testing"
)

func TestCrypto(t *testing.T) {
	key := []byte("0123456789012345")
	aesObj := gt.NewAES("cbc")
	t1, err := aesObj.Encrypt([]byte("你好"), key)
	if err != nil {
		t.Fatal(err)
	}
	t1base64 := gt.Byte2Base64(t1)
	t.Log(t1base64)
	t2, err := aesObj.Decrypt(t1, key)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(t2))

	t1b, err := gt.Base642Byte(t1base64)
	if err != nil {
		t.Fatal(err)
	}
	t3, err := aesObj.Decrypt(t1b, key)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(t3))

}
