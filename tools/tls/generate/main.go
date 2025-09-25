// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a MIT License
// that can be found in the LICENSE file.

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

func main() {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exePath = filepath.Dir(exePath)
	fmt.Println("Executable dir:", exePath)

	// 生成RSA私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// 设置证书模板
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   "localhost", // 证书通用名称，通常是域名
			Organization: []string{"My Organization"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour), // 有效期1年
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,                                // 不是CA证书
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},   // 允许的IP地址
		DNSNames:              []string{"localhost", "example.com"}, // 允许的域名
	}

	// 生成证书
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		panic(err)
	}

	// 保存证书
	certOut, err := os.Create(filepath.Join(exePath, "server.crt"))
	if err != nil {
		panic(err)
	}
	defer certOut.Close()
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		panic(err)
	}

	// 保存私钥
	keyOut, err := os.Create(filepath.Join(exePath, "server.key"))
	if err != nil {
		panic(err)
	}
	defer keyOut.Close()
	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		panic(err)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		panic(err)
	}

	fmt.Println("自签名证书生成成功:")
	fmt.Println("证书文件: server.crt")
	fmt.Println("私钥文件: server.key")

	// read the generated files and print their contents
	certData, err := os.ReadFile(filepath.Join(exePath, "server.crt"))
	if err != nil {
		panic(err)
	}

	pemBlock, _ := pem.Decode(certData)

	if pemBlock == nil || pemBlock.Type != "CERTIFICATE" {
		panic("failed to decode PEM block containing certificate")
	}

	cert, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n证书详情:\n")
	fmt.Printf("  版本: %d\n", cert.Version)
	fmt.Printf("  序列号: %s\n", cert.SerialNumber)
	fmt.Printf("  签发者: %s\n", cert.Issuer)
	fmt.Printf("  有效期: %s 至 %s\n", cert.NotBefore.Format(time.RFC3339), cert.NotAfter.Format(time.RFC3339))
	fmt.Printf("  主体: %s\n", cert.Subject)
	fmt.Printf("  公钥算法: %s\n", cert.PublicKeyAlgorithm)
	fmt.Printf("  签名算法: %s\n", cert.SignatureAlgorithm)
	fmt.Printf("  DNS 名称: %v\n", cert.DNSNames)
	fmt.Printf("  IP 地址: %v\n", cert.IPAddresses)
	fmt.Printf("  是否为CA: %v\n", cert.IsCA)
	fmt.Printf("  密钥用法: %v\n", cert.KeyUsage)
	fmt.Printf("  扩展密钥用法: %v\n", cert.ExtKeyUsage)

	// 判断证书是否生效和是否过期、剩余天数
	now := time.Now()
	if now.Before(cert.NotBefore) {
		fmt.Printf("⌚️证书尚未生效, 生效时间: %s\n", cert.NotBefore.Format(time.RFC3339))
	} else if now.After(cert.NotAfter) {
		fmt.Printf("⚠️证书已过期, 过期时间: %s\n", cert.NotAfter.Format(time.RFC3339))
	} else {
		daysLeft := int(cert.NotAfter.Sub(now).Hours() / 24)
		fmt.Printf("⌚️证书有效, 剩余天数: %d 天\n", daysLeft)
	}
}
