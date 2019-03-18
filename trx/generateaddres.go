package trx

import (
	// "bytes"
	"encoding/hex"
	// "crypto/rand"
	// "encoding/binary"
	// "encoding/json"
	// "errors"
	"fmt"
	"github.com/sasaxie/go-client-api/common/base58"
	// "io/ioutil"
	"log"
	// "net"
	// "net/http"
	// "strconv"
	"github.com/sasaxie/go-client-api/common/crypto"
	// "github.com/sasaxie/go-client-api/core"
	// "github.com/sasaxie/go-client-api/util"
	// "crypto/ecdsa"
	// "crypto/sha256"
	// "sync"
	// "time"
	"github.com/ethereum/go-ethereum/crypto/sha3"
)

// 6.1 算法描述
// 首先产生密钥对，取公钥，仅包含x，y坐标的64字节的byte数组。
// 对公钥做sha3-256的hash运算。取其最后20字节。测试网在前面填充A0。
// 主网地址在前面补41，得到地址的原始格式。长度为21字节。
// 做两次sha256计算，取其前4字节得到校验码。
// 将校验码附加在地址的原始格式后面，做base58编码，得到base58check格式的地址。测试网地址编码后以27开头，长度为35字节。
// 主网地址编码后以T开头，长度34字节。
// 注意：sha3协议我们使用的是KECCAK-256。

// 6.2 Mainnet地址，以41为前缀
// address = 41||sha3[12,32): 415a523b449890854c8fc460ab602df9f31fe4293f
// sha256_0 = sha256(address): 06672d677b33045c16d53dbfb1abda1902125cb3a7519dc2a6c202e3d38d3322
// sha256_1 = sha256(sha256_0): 9b07d5619882ac91dbe59910499b6948eb3019fafc4f5d05d9ed589bb932a1b4
// checkSum = sha256_1[0, 4): 9b07d561
// addchecksum = address || checkSum: 415a523b449890854c8fc460ab602df9f31fe4293f9b07d561
// base58Address = Base58(addchecksum): TJCnKsPa7y5okkXvQAidZBzqx3QyQ6sxMW

// public static String encode58Check(byte[] input) {
//     byte[] hash0 = Sha256Hash.hash(input);
//     byte[] hash1 = Sha256Hash.hash(hash0);
//     byte[] inputCheck = new byte[input.length + 4];
//     System.arraycopy(input, 0, inputCheck, 0, input.length);
//     System.arraycopy(hash1, 0, inputCheck, input.length, 4);
//     return Base58.encode(inputCheck);
//   }
func (r *RPCClient) generateAddress(net string) (pri, pub, addrori, addr string, err error) {
	// 5d7047e3da9ae069f514874a160c57fe4c819dd22cd0b818148942540dfa41b8
	// 41e9650561309da99ab0423efe090c183f592f6414
	//GetAccount TXFHTpS4EwsuxeuEihfBabJ4GkccDGSfhw
	// priKey, err := crypto.GetPrivateKeyByHexString("5d7047e3da9ae069f514874a160c57fe4c819dd22cd0b818148942540dfa41b8")
	priKey, err := crypto.GenerateKey()
	if err != nil {
		return
	}
	// priKey, err := crypto.GenerateKey()
	if err != nil {
		return
	}
	// fmt.Printf("%x", key.PrivateKey.D.Bytes())
	// fmt.Printf("%x", crypto.FromECDSA(key.PrivateKey))
	pri = fmt.Sprintf("%x", priKey.D.Bytes())
	// pub = fmt.Sprintf("%x", priKey.PublicKey.D.Bytes())
	x := priKey.PublicKey.X.Bytes()
	y := priKey.PublicKey.Y.Bytes()
	byteArray := append(x, y...)
	log.Println("len:", len(byteArray))
	hash := sha3.NewKeccak256()
	var buf []byte
	hash.Write(byteArray)
	buf = hash.Sum(buf)
	log.Println("buf len:", len(buf))
	str := hex.EncodeToString(buf)[24:]
	// str := "47839beb93a34f77adc0b89a8ce0e37509317468"
	log.Println("str:", str)
	if net == "mainnet" {
		str = "41" + str
	} else {
		str = "a0" + str
	}
	addrByte, err := hex.DecodeString(str)
	if err != nil {
		return
	}
	// str := "415a523b449890854c8fc460ab602df9f31fe4293f"
	log.Println("address original:", str)
	addrori = str
	addr = base58.EncodeCheck(addrByte)
	// h := sha256.New()
	// h.Write([]byte(str))
	// firstSha := h.Sum(nil)
	// fmt.Println("firstSha:", hex.EncodeToString(firstSha))
	// // fmt.Printf("firstSha:%x\n", firstSha) //equal with hex.EncodeToString
	// h.Write(firstSha)
	// secondeSha := h.Sum(nil)
	// fmt.Println("secondeSha:", hex.EncodeToString(secondeSha))
	// final := append([]byte(str), secondeSha[:8]...)
	// final := str + hex.EncodeToString(secondeSha[:4])

	// check:=
	log.Println("address:", addr)
	return
}
func (r *RPCClient) TestGenerateAddress(net string) (pri, pub, addrori, addr string, err error) {
	return r.generateAddress("mainnet")

}
