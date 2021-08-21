package gathertool

import (
	"testing"
)

func TestAES(t *testing.T){
	aes_bcb := NewAES(AES_CBC, []byte("2222222222222222"))
	data1, err := aes_bcb.Encrypt([]byte("aaaaaa"),[]byte("1111111111111111"))
	loger(string(data1), err)
	data2, err := aes_bcb.Decrypt(data1, []byte("1111111111111111"))
	loger(string(data2), err)

	aes_ecb := NewAES(AES_ECB)
	data3, err := aes_ecb.Encrypt([]byte("aaaaaa"), []byte("1111111111111111"))
	loger(string(data3), err)
	data4, err := aes_ecb.Decrypt(data3, []byte("1111111111111111"))
	loger(string(data4), err)

	aes_cfb := NewAES(AES_CFB)
	data5, err := aes_cfb.Encrypt([]byte("aaaaaa"), []byte("1111111111111111"))
	loger(string(data5), err)
	data6, err := aes_cfb.Decrypt(data5, []byte("1111111111111111"))
	loger(string(data6), err)

	aes_ctr := NewAES(AES_CTR, []byte("2222222222222222"))
	data7, err := aes_ctr.Encrypt([]byte("aaaaaa"), []byte("1111111111111111"))
	loger(string(data7), err)
	data8, err := aes_ctr.Decrypt(data7, []byte("1111111111111111"))
	loger(string(data8), err)


}