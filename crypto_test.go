package gathertool

import (
	"log"
	"testing"
)

func TestAES(t *testing.T){
	aes_bcb := NewAES(CBC, []byte("2222222222222222"))
	data1, err := aes_bcb.Encrypt([]byte("aaaaaa"),[]byte("1111111111111111"))
	loger(string(data1), err)
	data2, err := aes_bcb.Decrypt(data1, []byte("1111111111111111"))
	loger(string(data2), err)

	des_bcb := NewDES(CBC, []byte("22222222"))
	data11, err := des_bcb.Encrypt([]byte("aaaaaa"),[]byte("11111111"))
	loger(string(data11), err)
	data22, err := des_bcb.Decrypt(data11, []byte("11111111"))
	loger(string(data22), err)

	aes_ecb := NewAES(ECB)
	data3, err := aes_ecb.Encrypt([]byte("aaaaaa"), []byte("1111111111111111"))
	loger(string(data3), err)
	data4, err := aes_ecb.Decrypt(data3, []byte("1111111111111111"))
	loger(string(data4), err)

	des_ecb := NewDES(ECB)
	data33, err := des_ecb.Encrypt([]byte("aaaaaa"), []byte("11111111"))
	loger(string(data33), err)
	data44, err := des_ecb.Decrypt(data33, []byte("11111111"))
	loger(string(data44), err)

	aes_cfb := NewAES(CFB)
	data5, err := aes_cfb.Encrypt([]byte("aaaaaa"), []byte("1111111111111111"))
	loger(string(data5), err)
	data6, err := aes_cfb.Decrypt(data5, []byte("1111111111111111"))
	loger(string(data6), err)

	des_cfb := NewDES(CFB)
	data55, err := des_cfb.Encrypt([]byte("aaaaaa"), []byte("11111111"))
	loger(string(data55), err)
	data66, err := des_cfb.Decrypt(data55, []byte("11111111"))
	loger(string(data66), err)

	aes_ctr := NewAES(CTR, []byte("2222222222222222"))
	data7, err := aes_ctr.Encrypt([]byte("aaaaaa"), []byte("1111111111111111"))
	loger(string(data7), err)
	data8, err := aes_ctr.Decrypt(data7, []byte("1111111111111111"))
	loger(string(data8), err)

	des_ctr := NewDES(CTR, []byte("2222222222222222"))
	data77, err := des_ctr.Encrypt([]byte("aaaaaa"), []byte("11111111"))
	loger(string(data77), err)
	data88, err := des_ctr.Decrypt(data77, []byte("11111111"))
	loger(string(data88), err)

	txt := "api.eol.cn/gh5/api?local_province_id=11&local_type_id=1&page=1&school_id=3602&size=30&uri=apidata/api/gk/score/province&year=2015"
	data9 := HmacSHA1(txt, "D23ABC@#56")
	log.Println(data9)

	data10 := PBKDF2([]byte("D23ABC@#56"), []byte("secret"))
	log.Println(string(data10))

}