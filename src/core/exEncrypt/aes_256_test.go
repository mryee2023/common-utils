package exEncrypt

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptStrYQB(t *testing.T) {
	var (
		key = "12345678901234567890123456789012"
		iv  = "mIszfkB0kwKla6pR"
	)
	realKey, err := base64.StdEncoding.DecodeString(key)
	assert.NoError(t, err)

	var aes = NewWithIv(string(realKey), iv)

	planinText := `{"suppId":"8800000016"}`

	fmt.Println("原文：", planinText)

	encrypt, err := aes.EncryptHex(planinText)
	//f17ee300b71bd1b6719d5ec11faa7a3e1001a5ca36be5680612da479c3793303
	fmt.Println("密文:", encrypt, err)

	ss, e := AES256Encrypt(planinText, key, iv)
	fmt.Println(e)
	fmt.Println(ss)

	//assert.Equal(t, "bcade40b70af4c6520d65e12452b8f7f51c9cb99ace944653abfa803612451e4", encrypt)
	//
	//decrypt, err := aes.DecryptHex(encrypt)
	//
	//fmt.Println("jie密文:", decrypt, err)
	//assert.Equal(t, planinText, decrypt)

}
