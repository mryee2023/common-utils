package exEncrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
)

//加密过程：
//  1、处理数据，对数据进行填充，采用PKCS7（当密钥长度不够时，缺几位补几个几）的方式。
//  2、对数据进行加密，采用AES加密方法中CBC加密模式
//  3、对得到的加密数据，进行base64加密，得到字符串
//  解密过程相反

// key 16,24,32位字符串的话，分别对应AES-128，AES-192，AES-256 加密方法

type Aes struct {
	key []byte
	iv  []byte
}

func New(key string) *Aes {
	return &Aes{key: []byte(key)}
}

func NewWithIv(key, iv string) *Aes {
	return &Aes{key: []byte(key), iv: []byte(iv)}
}

// pkcs7Padding 填充
func (a *Aes) pkcs7Padding(data []byte, blockSize int) []byte {
	//判断缺少几位长度。最少1，最多 blockSize
	padding := blockSize - len(data)%blockSize
	//补足位数。把切片[]byte{byte(padding)}复制padding个
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// pkcs7UnPadding 填充的反向操作
func (a *Aes) pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	}
	//获取填充的个数
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

// Encrypt 加密
func (a *Aes) Encrypt(data []byte) ([]byte, error) {
	//创建加密实例
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	//判断加密快的大小
	blockSize := block.BlockSize()

	//填充
	encryptBytes := a.pkcs7Padding(data, blockSize)
	//初始化加密数据接收切片
	crypted := make([]byte, len(encryptBytes))
	//使用cbc加密模式
	blockMode := cipher.NewCBCEncrypter(block, a.getIv(blockSize))

	//执行加密
	blockMode.CryptBlocks(crypted, encryptBytes)
	return crypted, nil
}

func (a *Aes) getIv(blockSize int) []byte {
	if len(a.iv) > 0 {
		return a.iv
	}
	return a.key[:blockSize]
}

// Decrypt 解密
func (a *Aes) Decrypt(data []byte) (res []byte, err error) {
	defer func() {
		if p := recover(); p != nil {
			errMsg := fmt.Sprintf("%v", p)
			err = errors.New(errMsg)
		}
	}()

	//创建实例
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//使用cbc
	blockMode := cipher.NewCBCDecrypter(block, a.getIv(blockSize))
	//初始化解密数据接收切片
	crypted := make([]byte, len(data))
	//执行解密
	blockMode.CryptBlocks(crypted, data)
	//去除填充
	crypted, err = a.pkcs7UnPadding(crypted)
	if err != nil {
		return nil, err
	}
	return crypted, nil
}

func (a *Aes) EncryptHex(data string) (string, error) {
	res, err := a.Encrypt([]byte(data))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(res), nil
}

func (a *Aes) DecryptHex(data string) (string, error) {
	dataByte, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}

	res, err := a.Decrypt(dataByte)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func (a *Aes) EncryptStr(data string) (string, error) {
	res, err := a.Encrypt([]byte(data))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(res), nil
}

func (a *Aes) DecryptStr(data string) ([]byte, error) {
	dataByte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	return a.Decrypt(dataByte)
}
