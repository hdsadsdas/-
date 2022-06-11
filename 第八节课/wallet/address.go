package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"公链系统开发/第八节课/tools"
)

/**
* @author : 哈哈
* @email : 598421227@qq.com
* @phone : 18816473550
* @DateTime : 2022/5/23 9:15
**/

const VERSION = 0X00

func NewPubKeys() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {

	//生成一个p256曲线
	curve := elliptic.P256()

	//生成reader
	reader := rand.Reader

	//创建私钥
	privatekey, err := ecdsa.GenerateKey(curve, reader)
	if err != nil {
		return nil, nil, err
	}

	//生成公钥
	publickey := privatekey.PublicKey

	return privatekey, &publickey, nil

}

//生成比特币地址
func generateAddr() (string, error) {

	_, pub, err := NewPubKeys()
	if err != nil {
		return "", err
	}

	//将公钥序列化成[]byte
	pub_byte := elliptic.Marshal(pub.Curve, pub.X, pub.Y)

	pub_sha256 := tools.GetHash(pub_byte)

	pubhash := tools.Ripemd160(pub_sha256)

	ver_pubhash := append([]byte{VERSION}, pubhash...)

	first_hash := tools.GetHash(ver_pubhash)
	second_hash := tools.GetHash(first_hash)

	check := second_hash[:4]

	ver_pubhash_check := append(ver_pubhash, check...)
	
	addr := tools.Encode(ver_pubhash_check)

	return addr, nil

}

func checkaddr(addr string) bool {

	ver_pubhash_check := tools.Decode(addr)

	check := ver_pubhash_check[len(ver_pubhash_check)-4:]
	ver_pubhash := ver_pubhash_check[:len(ver_pubhash_check)-4]

	first_hash:=tools.GetHash(ver_pubhash)
	second_hash :=tools.GetHash(first_hash)

	check2 := second_hash[:4]

	return bytes.Compare(check, check2) == 0

}

//根据地址获取公钥hash
func GetPubHash(address string)([]byte,error){

	istrue := checkaddr(address)

	if !istrue {
		return nil,errors.New("地址不合法")
	}

	//解码
	ver_pubhash_check := tools.Decode(address)

	ver_pubhash := ver_pubhash_check[:len(ver_pubhash_check)-4]

	pubhash := ver_pubhash[1:]

	return pubhash,nil

}
