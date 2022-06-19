package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"公链系统开发/考试/tools"
	"errors"
)
//版本号
const VERSION = 0X00
//生成公钥和私钥
func NewPubKeys()(*ecdsa.PrivateKey,*ecdsa.PublicKey, error){
	//生成一个p256的一个曲线
	curve := elliptic.P256()
	//创建私钥
	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err !=nil{
		return nil,nil,err
	}
	//生成公钥
	publicKey := privateKey.PublicKey
	return privateKey,&publicKey,nil
}

//生成比特币地址
func generateAddr()(string,error){

	//获取公钥
	_, pub, err := NewPubKeys()
	if err !=nil{
		return "",err
	}
	//对公钥进行类型转换，序列化
	pub_byte := elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	//对公钥进行sha256计算
	pub_sha256 := tools.GetHash(pub_byte)
	//进行ripemd160计算  长度 160位，20个字节
	pubhash := tools.Ripemd160(pub_sha256)
	//将版本号和pubhash进行拼接
	ver_pubhash:=append([]byte{VERSION},pubhash...)
	//进行双sha256计算
	first_hash:=tools.GetHash(ver_pubhash)
	second_hash:=tools.GetHash(first_hash)
	check:=second_hash[:4]
	//将check和ver_pubhash进行拼接
	ver_pubhash_check:=append(ver_pubhash,check...)
	//将ver_pubhash_check进行base58转化
	addr := tools.Encode(ver_pubhash_check)
	return addr,nil
}

func checkAddr(addr string)bool{

	//进行解码
	ver_pubhash_check := tools.Decode(addr)
	//截取check
	check:=ver_pubhash_check[len(ver_pubhash_check)-4:]
	//截取ver_pubhash
	ver_pubhash:=ver_pubhash_check[:len(ver_pubhash_check)-4]
	//进行双sha256计算
	first:=tools.GetHash(ver_pubhash)
	second:=tools.GetHash(first)
	//对second进行截取前四位
	check2:=second[:4]
	//check和check2进行比较
	return bytes.Compare(check2,check)==0

}

//根据地址获取公钥hash
func GetPubHash(address string)([]byte,error){
	istrue := checkAddr(address)
	if !istrue{
		return nil,errors.New("地址不合法")
	}
	//解码
	ver_pubhash_check := tools.Decode(address)

	//得到ver_pubhash
	ver_pubhash:=ver_pubhash_check[:len(ver_pubhash_check)-4]

	pubhash:=ver_pubhash[1:]
	return pubhash,nil
}