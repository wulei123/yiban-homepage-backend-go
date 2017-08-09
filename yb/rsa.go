package yb

import (
	"encoding/pem"
	"crypto/x509"
	"crypto/rsa"
	"crypto/rand"
	"encoding/base64"
)

func GetBlockFromPem(pemKey []byte) []byte {
	block, _ := pem.Decode(pemKey)
	if block == nil {
		panic(" key error!")
	}
	return block.Bytes
}
func RsaEncrypt(origData []byte, pemKey []byte) string {
	Bytes := GetBlockFromPem(pemKey) //获取公钥pem的block
	pubInterface, err := x509.ParsePKIXPublicKey(Bytes) //解析公钥
	if err != nil {
		panic(err)
	}
	pub := pubInterface.(*rsa.PublicKey)
	encypt, err := rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
	if err != nil {
		panic(err)
	}
	return string(base64.StdEncoding.EncodeToString(encypt)) //由于加密后是字节流，直接输出查看会乱码 用base64加密
}