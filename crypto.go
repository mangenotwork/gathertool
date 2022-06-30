/*
	Description : 加密解码相关封装方法
	Author : ManGe
			2912882908@qq.com
			https://github.com/mangenotwork/gathertool
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
	"fmt"
	"github.com/dgrijalva/jwt-go"
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

// AES  AES interface
type AES interface {
	Encrypt(str, key []byte) ([]byte, error)
	Decrypt(str, key []byte) ([]byte, error)
}

// DES DES interface
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

// NewDES :  use NewDES(DES_CBC)
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
		Error("["+cbc.cryptoType+"-CBC] ERROR:" +err.Error())
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
		Error("["+cbc.cryptoType+"-CBC] ERROR:" +err.Error())
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
		Error("["+ecb.cryptoType+"-ECB] ERROR:" +err.Error())
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
		Error("["+ecb.cryptoType+"-ECB] ERROR:" +err.Error())
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
}

func (ecb *ecbObj) pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
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
		Error("["+cfb.cryptoType+"-CFB] ERROR:" +err.Error())
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
		Error("["+cfb.cryptoType+"-CFB] ERROR:" +err.Error())
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
		Error("[AES-CTR] ERROR:" +err.Error())
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
type ofbObj struct {}

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

// jwtEncrypt
func jwtEncrypt(token *jwt.Token, data map[string]interface{}, secret string) (string, error) {
	claims := make(jwt.MapClaims)
	for k, v := range data {
		claims[k] = v
	}
	token.Claims = claims
	return token.SignedString([]byte(secret))
}

// JwtEncrypt
func JwtEncrypt(data map[string]interface{}, secret, method string) (string, error){
	switch method {
	case "256":
		return jwtEncrypt(jwt.New(jwt.SigningMethodHS256), data, secret)
	case "384":
		return jwtEncrypt(jwt.New(jwt.SigningMethodHS384), data, secret)
	case "512":
		return jwtEncrypt(jwt.New(jwt.SigningMethodHS512), data, secret)
	}
	return "",fmt.Errorf("未知method; method= 256 or 384 or 512 ")
}

// JwtEncrypt256
func JwtEncrypt256(data map[string]interface{}, secret string) (string, error){
	token := jwt.New(jwt.SigningMethodHS256)
	return jwtEncrypt(token, data, secret)
}

// JwtEncrypt384
func JwtEncrypt384(data map[string]interface{}, secret string) (string, error){
	token := jwt.New(jwt.SigningMethodHS384)
	return jwtEncrypt(token, data, secret)
}

// JwtEncrypt512
func JwtEncrypt512(data map[string]interface{}, secret string) (string, error){
	token := jwt.New(jwt.SigningMethodHS512)
	return jwtEncrypt(token, data, secret)
}

// JwtDecrypt
func JwtDecrypt(tokenString, secret string) (data map[string]interface{}, err error) {
	data = make(map[string]interface{})
	var secretFunc = func() jwt.Keyfunc { //按照这样的规则解析
		return func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		}
	}
	token, err := jwt.Parse(tokenString, secretFunc())
	if err != nil {
		err = fmt.Errorf("未知Token")
		return
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return
	}
	if !token.Valid {
		// 令牌错误
		return
	}
	for k, v := range claim {
		data[k] =v
	}
	return
}


// TODO Rabbit

// TODO RC4

// TODO RIPEMD-160

/*
0: “sha1”
1: “sha224”
2: “sha256”
3: “sha384”
4: “sha512”
5: “md5”
6: “rmd160”
7: “sha224WithRSAEncryption”
8: “RSA-SHA224”
9: “sha256WithRSAEncryption”
10: “RSA-SHA256”
11: “sha384WithRSAEncryption”
12: “RSA-SHA384”
13: “sha512WithRSAEncryption”
14: “RSA-SHA512”
15: “RSA-SHA1”
16: “ecdsa-with-SHA1”
17: “sha256”
18: “sha224”
19: “sha384”
20: “sha512”
21: “DSA-SHA”
22: “DSA-SHA1”
23: “DSA”
24: “DSA-WITH-SHA224”
25: “DSA-SHA224”
26: “DSA-WITH-SHA256”
27: “DSA-SHA256”
28: “DSA-WITH-SHA384”
29: “DSA-SHA384”
30: “DSA-WITH-SHA512”
31: “DSA-SHA512”
32: “DSA-RIPEMD160”
33: “ripemd160WithRSA”
34: “RSA-RIPEMD160”
35: “md5WithRSAEncryption”
36: “RSA-MD5”
 */

