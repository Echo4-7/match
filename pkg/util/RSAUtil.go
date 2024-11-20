package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

// Encrypt 加密
func Encrypt(msg string, pub *rsa.PublicKey) (string, error) {
	ciphertext, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		pub,
		[]byte(msg),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %v", err)
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 使用私钥解密，并返回解密后的明文字符串
func Decrypt(ciphertext string, priv *rsa.PrivateKey) (string, error) {
	// 将 Base64 编码的密文解码为字节切片
	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed: %v", err)
	}
	plaintext, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		priv,
		decodedCiphertext,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %v", err)
	}
	// 将解密后的字节切片转换为字符串
	return string(plaintext), nil
}

func GetPrivateKey() *rsa.PrivateKey {
	// 从文件中读取PEM数据
	pemBytes, err := os.ReadFile("private.pem")
	if err != nil {
		log.Fatalf("Error reading private key file: %v", err)
	}

	// 解析PEM数据
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		log.Fatal("Failed to parse PEM block containing the key")
	}

	// 解析私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// 如果你使用的是PKCS8格式的私钥，可以尝试以下代码
		// privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		log.Fatalf("Error parsing private key: %v", err)
	}
	return privateKey
}

// generateRSAKeyPair 生成公钥和私钥
//func generateRSAKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
//	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
//	if err != nil {
//		fmt.Println("Failed to generate private key:", err)
//		return nil, nil
//	}
//	publicKey := &privateKey.PublicKey
//	savePEMKey("private.pem", privateKey)
//	savePublicPEMKey("public.pem", publicKey)
//	return privateKey, publicKey
//}

// savePEMKey 保存私钥
//func savePEMKey(fileName string, key *rsa.PrivateKey) {
//	file, err := os.Create(fileName)
//	if err != nil {
//		fmt.Println("Failed to create key file:", err)
//		return
//	}
//	defer file.Close()
//	privateKeyPEM := pem.EncodeToMemory(
//		&pem.Block{
//			Type:  "RSA PRIVATE KEY",
//			Bytes: x509.MarshalPKCS1PrivateKey(key),
//		},
//	)
//	file.Write(privateKeyPEM)
//}

// savePublicPEMKey保存公钥
//func savePublicPEMKey(fileName string, pubkey *rsa.PublicKey) {
//	file, err := os.Create(fileName)
//	if err != nil {
//		fmt.Println("Failed to create key file:", err)
//		return
//	}
//	defer file.Close()
//	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubkey)
//	if err != nil {
//		fmt.Println("Failed to marshal public key:", err)
//		return
//	}
//	publicKeyPEM := pem.EncodeToMemory(
//		&pem.Block{
//			Type:  "PUBLIC KEY",
//			Bytes: pubKeyBytes,
//		},
//	)
//	file.Write(publicKeyPEM)
//}
