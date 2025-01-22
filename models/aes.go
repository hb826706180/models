package models

import (
	"encoding/base64"
	"fmt"
	"github.com/forgoer/openssl"
	"github.com/go-ini/ini"
)

var Genket string

func init() {
	cfg, err := ini.Load("./conf/app.ini")
	if err != nil {
		fmt.Println("错误 Load() " + err.Error())
	}
	str := cfg.Section("Genkey").Key("genkey").String()
	Genket = str
}

type AES struct{}

// AesCBCEncrypt 加密
func AesCBCEncrypt(data interface{}, key string, iv string, padding string) (string, error) {
	//fmt.Println(">>>>AesCBCEncrypt() 加密<<<<")
	//转换成[]byte
	by := InterfaceToByte(data)
	k := []byte(key)
	i := []byte(iv)
	//fmt.Printf("加密前的数据 =%s key =%s iv =%s padding =%s \n", by, key, iv,padding)
	encrypt, err := openssl.AesCBCEncrypt(by, k, i, padding)
	//fmt.Printf("加密后的数据= %s \n", encrypt)
	if err != nil {
		fmt.Println("错误 AesCBCEncrypt() " + err.Error())
		return "", err
	}
	Encode := base64.StdEncoding.EncodeToString(encrypt)
	//fmt.Println("base64的数据 = " + Encode)
	return Encode, nil
}

// AesCBCDecrypt 解密
func AesCBCDecrypt(data string, key string, iv string, padding string) ([]byte, error) {
	//fmt.Println("解密前的base64数据: " + data)
	by, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		fmt.Println("AesCBCDecrypt解密中,base64解码失败,错误信息: base64 " + err.Error())
		return nil, err
	}
	k := []byte(key)
	i := []byte(iv)
	//fmt.Printf("base64解码后的数据 = %s \n", by)
	//fmt.Printf("解密前的数据 = %s key = %s iv = %s \n", by, key, iv)
	ret, err := openssl.AesCBCDecrypt(by, k, i, padding)
	if err != nil {
		fmt.Println("AesCBCDecrypt解密失败,错误信息:" + err.Error())
		return nil, err
	}
	//fmt.Printf("解密后的数据= %s \n", ret)
	return ret, nil
}

// AesECBEncrypt 加密
func AesECBEncrypt(data interface{}, key string, padding string) string {
	by := InterfaceToByte(data)
	encrypt, err := openssl.AesECBEncrypt(by, []byte(key), padding)
	if err != nil {
		fmt.Println("加密失败 " + err.Error())
		return ""
	}
	return base64.StdEncoding.EncodeToString(encrypt)
}
