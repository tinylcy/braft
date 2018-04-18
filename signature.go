package raft

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	// "fmt"
	"log"
)

var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEApzo9DwzkZ4TsYxo1DBpw
HZt6SvRghMwEpzYlJz8OvC/wbpQWlnIcr3zBS13WQSUUk8LSa0DSqTU4V8yLD1rR
luW5lukqO5MaBCAln2/LiCupxDRTFkrWmkf515GcmpWHJDnBJNfKmD7qQh++8O90
V42rsV8mZhdiMKQrp4YKkx0tD9gTQ+sTNi/g2qZrHKPYHafHN0Ads515XsdHg/MF
DnzDdjIsvXNAU3J8eAOK6O7AWmL7jV10EXIJpVDGka9li/1uCtE2XlusJMx+WxYz
rZErCjC/hEpqHQ2VfO8URHqb84FD36o+Y1hHhxF6AOP6G5WgAy9Tblbm9cT1HFx9
jQIDAQAB
-----END PUBLIC KEY-----
`)

var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEApzo9DwzkZ4TsYxo1DBpwHZt6SvRghMwEpzYlJz8OvC/wbpQW
lnIcr3zBS13WQSUUk8LSa0DSqTU4V8yLD1rRluW5lukqO5MaBCAln2/LiCupxDRT
FkrWmkf515GcmpWHJDnBJNfKmD7qQh++8O90V42rsV8mZhdiMKQrp4YKkx0tD9gT
Q+sTNi/g2qZrHKPYHafHN0Ads515XsdHg/MFDnzDdjIsvXNAU3J8eAOK6O7AWmL7
jV10EXIJpVDGka9li/1uCtE2XlusJMx+WxYzrZErCjC/hEpqHQ2VfO8URHqb84FD
36o+Y1hHhxF6AOP6G5WgAy9Tblbm9cT1HFx9jQIDAQABAoIBAECDbDjZLWhuVE+r
NZnUvTn+2EAAZRf2KTlk3xJz9jhNekD+qnQh08UzqNJtghGhv319pHWyDVMv7+uX
QnKLA95mA6Ifk6ZmCpxa1ojatTd0OMszsHYiKwZcDBvI1hSg6QDlswiGo2b2pqMZ
4izLBCQeyITmA0dRcBT50MmRIZU9BXULU5iXA0KwJCaAtcG8YvbiT2N+ZS8sY32Q
QY8hSPH4vThO4bpPYlMcv8nPc50mfEJ1PI7cH9dYq8N7tMZD8tn6rio3S59Ll401
SoOQuPOZK0OBZgxcHSlDImWAN/kuhZLzGneHKcPpqrXj3Is3tADcNUOoRf7J6UBJ
QlywL2ECgYEAxpTDAtUMatTXObpiVZ2r68iav4KbDiEj1h7RpezNIx9MNe1D8ybU
68eeqka0EkWxoWLL59OHG0wzDcJ8se0uIqDsgX5FMiL/NFdegfp9dpRflXkzcCee
XvAP4xZfNHSrGvUIIrY148ECnVo6UhTynK6IN9GKVCFgHk6MtazHwuUCgYEA15Sk
FCZi2JmwokeS99q+Uqqtr88nu2yr2lUX/tn08iHhaTRcw2xYCI9xAivPVyMFIjKd
uBmyjnsaUZ23zD8U8FtKWo1cyIQ7PO9+NaNHrgS6pI4jNr1RMYbhknBBglcfAsbB
p0a9qKjXQ2gzJACYmobme9lBL03v3iW7LsKIXYkCgYBu5zv2A/gYXeAJfH9Yo2MV
noZWOGHSRU3XUoTxbsuuNteAMo9FZ8V4HJcPL8d3gPbQU/Xe9fK5mxfUMm8ji3u9
mTQcqeGJO6RdngHJA5U4OWscdoD0vRukl9u3jpIDILlCp+AwSqTUGsIUEQULPGm2
eX7X9a2UiMM+ic3p1KIHxQKBgCH0JUuPKC5ZNnq4ryseZq96dlSkWeupGAAROvBG
v8+LCoeZWarl24+tl+zxnXxp5ZsXQcQHOBo8xU5petNOdCvPFQziCuUB/pqAVe54
wwdjc0oLoPw0IR+d0NVRnN+8fQPg7gs8lw6DWTQiqztWZLKh4JdNBnk+2zKv2qVo
ujWRAoGBAIPiseHRQwYS6sAwnJUk6Gl04MgFIzWtjr5zHCJU/O+mp+XVfmS2NxwI
63aElv0lbxX/OHGIC7hJ/5ipwTLh8cQUVqdul5A0fIbtdkNlLa7tfII5FRDhTCZT
vAohlGrLXWPCwzSPY7GIWdsni7wDb9Sfvef/k5pKsoAn8CgU2olQ
-----END RSA PRIVATE KEY-----
`)

var origindata = []byte("I am tinylcy.")

func signature(origindata []byte) []byte {
	// 生成私钥
	block, _ := pem.Decode(privateKey)
	if block == nil {
		log.Println("private key error")
		return nil
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Println(err)
		return nil
	}

	// Message - Signature
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto
	PSSmessage := origindata
	newhash := crypto.SHA256 // Hash function
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)

	// 使用私钥对哈希值进行加密
	signature, err := rsa.SignPSS(rand.Reader, priv, newhash, hashed, &opts)
	if err != nil {
		log.Println(err)
		return nil
	}

	return signature
}

func verifySignature(origindata []byte, sig []byte) bool {
	// 生成公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		log.Println("public key error")
		return false
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Println(err)
		return false
	}
	pub := pubInterface.(*rsa.PublicKey)

	// Message - Signature
	var opts rsa.PSSOptions
	opts.SaltLength = rsa.PSSSaltLengthAuto
	PSSmessage := origindata
	newhash := crypto.SHA256 // Hash function
	pssh := newhash.New()
	pssh.Write(PSSmessage)
	hashed := pssh.Sum(nil)

	// 使用公钥验证数字签名
	err = rsa.VerifyPSS(pub, newhash, hashed, sig, &opts)
	if err != nil {
		log.Println(err)
		return false
	}

	// log.Println("Verify signature successful")
	return true
}

//
// 假设数据并不需要加密传输，只需要通过数字签名来防止被篡改。
// 步骤：
// 1. 使用私钥对消息摘要进行加密，形成数字签名；
// 2. 使用公钥对签名进行解密，得到摘要原文；
// 3. 计算消息摘要，并与解密得到的摘要原文进行比对。
//
// func main() {
// 	sig := signature(origindata)
// 	result := verifySignature(origindata, sig)
// 	fmt.Println(result)
// }
