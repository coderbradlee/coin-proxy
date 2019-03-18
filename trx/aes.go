package trx

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	// "encoding/base64"
	// "fmt"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	// "log"
)

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	// log.Println("leng:", length)
	// log.Println("unpadding:", length)
	if length <= unpadding {
		return origData[:]
	}
	return origData[:(length - unpadding)]
}

func encrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}
func AesEncrypt(origData, keys []byte) ([]byte, error) {
	// length := len(keys)
	// var key []byte
	// if length < 32 {

	// } else if length > 32 {
	// 	key = keys[:32]
	// }
	hash := sha3.NewKeccak256()
	var buf []byte
	hash.Write(keys)
	buf = hash.Sum(buf)
	return encrypt(origData, buf)
}
func AesDecrypt(crypted, keys []byte) ([]byte, error) {
	// log.Println("52")
	hash := sha3.NewKeccak256()
	var buf []byte
	hash.Write(keys)
	buf = hash.Sum(buf)
	return decrypt(crypted, buf)
}
func decrypt(crypted, key []byte) ([]byte, error) {
	// log.Println("60")
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	// log.Println("65")
	blockSize := block.BlockSize()
	// log.Println("72")
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	// log.Println("73")
	origData := make([]byte, len(crypted))
	// log.Println("74")
	blockMode.CryptBlocks(origData, crypted)
	// log.Println("74")
	origData = PKCS7UnPadding(origData)
	// log.Println("75")
	return origData, nil
}

// func main() {
//     key := []byte("0123456789abcdef")
//     result, err := AesEncrypt([]byte("hello world"), key)
//     if err != nil {
//         panic(err)
//     }
//     fmt.Println(base64.StdEncoding.EncodeToString(result))
//     origData, err := AesDecrypt(result, key)
//     if err != nil {
//         panic(err)
//     }
//     fmt.Println(string(origData))
// }
