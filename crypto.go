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
	"crypto/des"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/pbkdf2"
	"hash"
	"io"
)

const (
	CBC = "CBC"
	ECB = "ECB"
	CFB = "CFB"
	CTR = "CTR"
)

type AES interface {
	Encrypt(str, key []byte) ([]byte, error)
	Decrypt(str, key []byte) ([]byte, error)
}

type DES interface {
	Encrypt(str, key []byte) ([]byte, error)
	Decrypt(str, key []byte) ([]byte, error)
}

// NewAES :  use NewAES(AES_CBC)
func NewAES(typeName string, arg ...[]byte) AES {
	iv := []byte{1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6}
	if len(arg) != 0 {
		iv = arg[0]
	}
	switch typeName {
	case "cbc", "Cbc", CBC:
		return &cbcObj{
			cryptoType: "aes",
			iv: iv,
		}
	case "ecb", "Ecb", ECB:
		return &ecbObj{
			cryptoType: "aes",
		}
	case "cfb", "Cfb", CFB:
		return &cfbObj{
			cryptoType: "aes",
		}
	case "ctr", "Ctr", CTR:
		return &ctrObj{
			cryptoType: "aes",
			count: iv,
		}
	default:
		return &cbcObj{
			iv: iv,
		}
	}
}

// NewAES
func NewDES(typeName string, arg ...[]byte) DES {
	iv := []byte{1,2,3,4,5,6,7,8,9,0,1,2,3,4,5,6}
	if len(arg) != 0 {
		iv = arg[0]
	}
	switch typeName {
	case "cbc", "Cbc", CBC:
		return &cbcObj{
			cryptoType: "des",
			iv: iv,
		}
	case "ecb", "Ecb", ECB:
		return &ecbObj{
			cryptoType: "des",
		}
	case "cfb", "Cfb", CFB:
		return &cfbObj{
			cryptoType: "des",
		}
	case "ctr", "Ctr", CTR:
		return &ctrObj{
			count: iv,
			cryptoType: "des",
		}
	default:
		return &cbcObj{
			iv: iv,
			cryptoType: "des",
		}
	}
}

// CBC : 密码分组链接模式（Cipher Block Chaining (CBC)） default
type cbcObj struct {
	cryptoType string
	iv []byte
}

func (cbc *cbcObj) getBlock(key []byte) (block cipher.Block, err error) {
	if cbc.cryptoType == "aes" {
		block, err = aes.NewCipher(key)
	}
	if cbc.cryptoType == "des" {
		block, err = des.NewCipher(key)
	}
	return
}

// AES CBC Encrypt
func (cbc *cbcObj) Encrypt(str, key []byte) ([]byte, error) {
	block, err := cbc.getBlock(key)
	if err != nil {
		loger("["+cbc.cryptoType+"-CBC] ERROR:" +err.Error())
		return []byte(""), err
	}
	blockSize := block.BlockSize()
	originData := cbc.pkcs5Padding(str, blockSize)
	blockMode := cipher.NewCBCEncrypter(block,cbc.iv)
	cryptData := make([]byte,len(originData))
	blockMode.CryptBlocks(cryptData,originData)
	P2E()
	return cryptData, nil
}

// AES CBC Decrypt
func (cbc *cbcObj) Decrypt(str, key []byte) ([]byte, error) {
	block, err := cbc.getBlock(key)
	if err != nil {
		loger("["+cbc.cryptoType+"-CBC] ERROR:" +err.Error())
		return []byte(""), err
	}
	blockMode := cipher.NewCBCDecrypter(block, cbc.iv)
	originStr := make([]byte,len(str))
	blockMode.CryptBlocks(originStr,str)
	P2E()
	return cbc.pkcs5UnPadding(originStr), nil
}

func (cbc *cbcObj) pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func (cbc *cbcObj) pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadDing := int(origData[length-1])
	return origData[:(length - unpadDing)]
}

// ECB : 电码本模式（Electronic Codebook Book (ECB)）
type ecbObj struct {
	cryptoType string
}

func (ecb *ecbObj) getBlock(key []byte) (block cipher.Block, err error) {
	if ecb.cryptoType == "aes" {
		block, err = aes.NewCipher(ecb.generateKey(key))
	}
	if ecb.cryptoType == "des" {
		block, err = des.NewCipher(key)
	}
	return
}

// AES ECB Encrypt
func (ecb *ecbObj) Encrypt(str, key []byte) ([]byte, error) {
	block, err := ecb.getBlock(key)
	if err != nil {
		loger("["+ecb.cryptoType+"-ECB] ERROR:" +err.Error())
		return []byte(""), err
	}
	blockSize := block.BlockSize()
	if ecb.cryptoType == "aes" {
		str = ecb.pkcs5PaddingAes(str, blockSize)
	}
	if ecb.cryptoType == "des" {
		str = ecb.pkcs5PaddingDes(str, blockSize)
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
	P2E()
	return encryptData, nil
}

func (ecb *ecbObj) pkcs5PaddingDes(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (ecb *ecbObj) pkcs5PaddingAes(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	if padding != 0 {
		ciphertext = append(ciphertext, bytes.Repeat([]byte{byte(0)}, padding)...)
	}
	return ciphertext
}

// AES ECB Decrypt
func (ecb *ecbObj) Decrypt(str, key []byte) ([]byte, error) {
	block, err := ecb.getBlock(key)
	if err != nil {
		loger("["+ecb.cryptoType+"-ECB] ERROR:" +err.Error())
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

	if ecb.cryptoType == "des" {
		return ecb.pkcs5UnPadding(decryptData), nil
	}
	P2E()
	return ecb.unPadding(decryptData), nil
}

func (ecb *ecbObj) generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

func (ecb *ecbObj) unPadding(src []byte) []byte {
	for i := len(src) - 1; ; i-- {
		if src[i] != 0 {
			return src[:i+1]
		}
	}
	return nil
}

func (ecb *ecbObj) pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// 密码反馈模式（Cipher FeedBack (CFB)）
type cfbObj struct {
	cryptoType string
}

func (cfb *cfbObj) getBlock(key []byte) (block cipher.Block, err error) {
	if cfb.cryptoType == "aes" {
		block, err = aes.NewCipher(key)
	}
	if cfb.cryptoType == "des" {
		block, err = des.NewCipher(key)
	}
	return
}

// AES CFB Encrypt
func (cfb *cfbObj) Encrypt(str, key []byte) ([]byte, error) {
	P2E()
	block, err := cfb.getBlock(key)
	if err != nil {
		loger("["+cfb.cryptoType+"-CFB] ERROR:" +err.Error())
		return nil, err
	}

	if cfb.cryptoType == "aes" {
		encrypted := make([]byte, aes.BlockSize+len(str))
		iv := encrypted[:aes.BlockSize]
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return nil, err
		}
		stream := cipher.NewCFBEncrypter(block, iv)
		stream.XORKeyStream(encrypted[aes.BlockSize:], str)
		return encrypted, nil
	}

	if cfb.cryptoType == "des" {
		encrypted := make([]byte, des.BlockSize+len(str))
		iv := encrypted[:des.BlockSize]
		if _, err := io.ReadFull(rand.Reader, iv); err != nil {
			return nil, err
		}
		stream := cipher.NewCFBEncrypter(block, iv)
		stream.XORKeyStream(encrypted[des.BlockSize:], str)
		return encrypted, nil
	}
	return nil, nil
}

// AES CFB Decrypt
func (cfb *cfbObj) Decrypt(str, key []byte) ([]byte, error) {
	P2E()
	block, err := cfb.getBlock(key)
	if err != nil {
		loger("["+cfb.cryptoType+"-CFB] ERROR:" +err.Error())
		return nil, err
	}

	iv := []byte{}
	if cfb.cryptoType == "aes" {
		if len(str) < aes.BlockSize {
			return nil, errors.New("ciphertext too short")
		}
		iv = str[:aes.BlockSize]
		str = str[aes.BlockSize:]
	}

	if cfb.cryptoType == "des" {
		if len(str) < des.BlockSize {
			return nil, errors.New("ciphertext too short")
		}
		iv = str[:des.BlockSize]
		str = str[des.BlockSize:]
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(str, str)
	return str, nil
}

// 计算器模式（Counter (CTR)）
type ctrObj struct {
	count []byte //指定计数器,长度必须等于block的块尺寸
	cryptoType string
}

func (ctr *ctrObj) getBlock(key []byte) (block cipher.Block, err error) {
	if ctr.cryptoType == "aes" {
		block, err = aes.NewCipher(key)
	}
	if ctr.cryptoType == "des" {
		block, err = des.NewCipher(key)
		if len(ctr.count) > des.BlockSize {
			ctr.count = ctr.count[0:des.BlockSize]
		}
	}
	return
}

// AES CTR Encrypt
func (ctr *ctrObj) Encrypt(str, key []byte) ([]byte, error) {
	return ctr.crypto(str, key)
}

// AES CTR Decrypt
func (ctr *ctrObj) Decrypt(str, key []byte) ([]byte, error) {
	return ctr.crypto(str, key)
}

func (ctr *ctrObj) crypto(str, key []byte) ([]byte, error) {
	P2E()
	block,err:=ctr.getBlock(key)
	if err != nil {
		loger("[AES-CTR] ERROR:" +err.Error())
		return []byte(""), err
	}
	//指定分组模式
	blockMode:=cipher.NewCTR(block, ctr.count)
	//执行加密、解密操作
	res:=make([]byte,len(str))
	blockMode.XORKeyStream(res,str)
	//返回明文或密文
	return res, nil
}

// 输出反馈模式（Output FeedBack (OFB)）
type ofbObj struct {
}

func hmacFunc(h func() hash.Hash, str, key []byte) string {
	mac := hmac.New(h, key)
	mac.Write(str)
	res := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return res
}

// HmacMD5
func HmacMD5(str, key string) string {
	return hmacFunc(md5.New, []byte(str), []byte(key))
}

// HmacSHA1
func HmacSHA1(str, key string) string {
	return hmacFunc(sha1.New, []byte(str), []byte(key))
}

// HmacSHA256
func HmacSHA256(str, key string) string {
	return hmacFunc(sha256.New, []byte(str), []byte(key))
}

// HmacSHA512
func HmacSHA512(str, key string) string {
	return hmacFunc(sha512.New, []byte(str), []byte(key))
}

func pbkdf2Func(h func() hash.Hash, str, salt []byte, iterations, keySize int) []byte {
	return pbkdf2.Key(str, salt, iterations, keySize, h)
}

// PBKDF2
func PBKDF2(str, salt []byte, iterations, keySize int) ([]byte) {
	return pbkdf2Func(sha256.New, str, salt, iterations, keySize)
}

// TODO Rabbit

// TODO RC4

// TODO RIPEMD-160
