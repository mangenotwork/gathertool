/*
	Description : 加密解码相关封装方法
	Author : ManGe
	Version : v0.1
	Date : 2021-08-21
*/

package gathertool

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

const (
	AES_CBC = "AES CBC"
	AES_ECB = "AES ECB"
	AES_CFB = "AES CFB"
	AES_CTR = "AES CTR"
)


// NewAES :  use NewAES(AES_CBC)
//	- CBC (the default)
//	- CFB
//	- CTR
//	- OFB (not)
//	- ECB
func NewAES(typeName string, arg ...[]byte) AES {
	iv := []byte{1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6}
	if len(arg) != 0 {
		iv = arg[0]
	}
	switch typeName {
	case "cbc", "CBC", "Cbc", AES_CBC:
		return &CBC{
			iv: iv,
		}
	case "ecb", "ECB", "Ecb", AES_ECB:
		return &ECB{}
	case "cfb", "CFB", "Cfb", AES_CFB:
		return &CFB{}
	case "ctr","CTR", "Ctr", AES_CTR:
		return &CTR{
			count: iv,
		}
	default:
		return &CBC{
			iv: iv,
		}
	}
}

type AES interface {
	Encrypt(str, key []byte) ([]byte, error)
	Decrypt(str, key []byte) ([]byte, error)
}

// 密码分组链接模式（Cipher Block Chaining (CBC)） default
type CBC struct {
	iv []byte
}

// AES CBC Encrypt
func (cbc *CBC) Encrypt(str, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		loger("[AES-CBC] ERROR:" +err.Error())
		return []byte(""), err
	}
	blockSize := block.BlockSize()
	originData := cbc.pkcs5Padding(str, blockSize)
	blockMode := cipher.NewCBCEncrypter(block,cbc.iv)
	cryptData := make([]byte,len(originData))
	blockMode.CryptBlocks(cryptData,originData)
	return cryptData, nil
}

// AES CBC Decrypt
func (cbc *CBC) Decrypt(str, key []byte) ([]byte, error) {
	block,err := aes.NewCipher(key)
	if err != nil {
		loger("[AES-CBC] ERROR:" +err.Error())
		return []byte(""), err
	}
	blockMode := cipher.NewCBCDecrypter(block, cbc.iv)
	originStr := make([]byte,len(str))
	blockMode.CryptBlocks(originStr,str)
	return cbc.pkcs5UnPadding(originStr), nil
}

func (cbc *CBC) pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func (cbc *CBC) pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadDing := int(origData[length-1])
	return origData[:(length - unpadDing)]
}

// 电码本模式（Electronic Codebook Book (ECB)）
type ECB struct {
}

// AES ECB Encrypt
func (ecb *ECB) Encrypt(str, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(ecb.generateKey(key))
	if err != nil {
		loger("[AES-ECB] ERROR:" +err.Error())
		return []byte(""), err
	}
	blockSize := block.BlockSize()
	//padding
	paddingCount := aes.BlockSize - len(str)%aes.BlockSize
	if paddingCount != 0 {
		str = append(str, bytes.Repeat([]byte{byte(0)}, paddingCount)...)
	}
	//返回加密结果
	encryptData := make([]byte, len(str))
	//存储每次加密的数据
	tmpData := make([]byte, blockSize)

	//分组分块加密
	for index := 0; index < len(str); index += blockSize {
		block.Encrypt(tmpData, str[index:index+blockSize])
		copy(encryptData, tmpData)
	}
	return encryptData, nil
}

// AES ECB Decrypt
func (ecb *ECB) Decrypt(str, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(ecb.generateKey(key))
	if err != nil {
		loger("[AES-ECB] ERROR:" +err.Error())
		return []byte(""), err
	}
	blockSize := block.BlockSize()
	//返回加密结果
	decryptData := make([]byte, len(str))
	//存储每次加密的数据
	tmpData := make([]byte, blockSize)

	//分组分块加密
	for index := 0; index < len(str); index += blockSize {
		block.Decrypt(tmpData, str[index:index+blockSize])
		copy(decryptData, tmpData)
	}

	return ecb.unPadding(decryptData), nil
}

func (ecb *ECB) generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

func (ecb *ECB) unPadding(src []byte) []byte {
	for i := len(src) - 1; ; i-- {
		if src[i] != 0 {
			return src[:i+1]
		}
	}
	return nil
}


type aesObj struct {

}

// 密码反馈模式（Cipher FeedBack (CFB)）
type CFB struct {
}

// AES CFB Encrypt
func (cfb *CFB) Encrypt(str, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		loger("[AES-CFB] ERROR:" +err.Error())
		return []byte(""), err
	}
	encrypted := make([]byte, aes.BlockSize+len(str))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], str)
	return encrypted, nil
}

// AES CFB Decrypt
func (cfb *CFB) Decrypt(str, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		loger("[AES-CFB] ERROR:" +err.Error())
		return []byte(""), err
	}
	if len(str) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := str[:aes.BlockSize]
	str = str[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(str, str)
	return str, nil
}

// 计算器模式（Counter (CTR)）
type CTR struct {
	count []byte //指定计数器,长度必须等于block的块尺寸
}

// AES CTR Encrypt
func (ctr *CTR) Encrypt(str, key []byte) ([]byte, error) {
	return ctr.crypto(str, key, ctr.count)
}

// AES CTR Decrypt
func (ctr *CTR) Decrypt(str, key []byte) ([]byte, error) {
	return ctr.crypto(str, key, ctr.count)
}

func (ctr *CTR) crypto(str, key, count []byte) ([]byte, error) {
	block,err:=aes.NewCipher(key)
	if err != nil {
		loger("[AES-CTR] ERROR:" +err.Error())
		return []byte(""), err
	}
	//指定分组模式
	blockMode:=cipher.NewCTR(block,count)
	//执行加密、解密操作
	res:=make([]byte,len(str))
	blockMode.XORKeyStream(res,str)
	//返回明文或密文
	return res, nil
}

// 输出反馈模式（Output FeedBack (OFB)）
type OFB struct {
}


// TODO DES

// TODO Rabbit

// TODO RC4

// TODO PBKDF2

// TODO HmacMD5

// TODO HmacSHA1

// TODO HmacSHA256

// TODO HmacSHA512

// TODO MD5

// TODO SHA-1

// TODO SHA-2

// TODO SHA-3

// TODO RIPEMD-160