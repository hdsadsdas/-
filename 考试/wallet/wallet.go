package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"errors"
	"github.com/boltdb/bolt"
	"公链系统开发/考试/tools"
)

const PRIVATE_BUCKET = "private_bucket"

//用于封装若干个功能的，比如生成地址，地址校验，保存私钥，查看私钥等
type Wallet struct {
	DB *bolt.DB
}

//实例化一个wallet对象
func NewWallet(db *bolt.DB)(*Wallet,error){
	if db ==nil{
		return nil,errors.New("db错误")
	}
	//创建出来一个桶
	err:=db.Update(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(PRIVATE_BUCKET))
		if bk == nil{
			//没有桶
			_,err:= tx.CreateBucket([]byte(PRIVATE_BUCKET))
			if err !=nil{
				return err
			}
		}
		return nil
	})
	if err!=nil{
		return nil, err
	}

	return &Wallet{db},nil
}

//保存数据到桶中
func(w *Wallet) SavePrivateKey(addr string ,pri *ecdsa.PrivateKey)error{
	db:=w.DB
	err:=db.Update(func(tx *bolt.Tx) error {
		bk:=tx.Bucket([]byte(PRIVATE_BUCKET))
		if bk == nil{
			return errors.New("桶不存在")
		}
		//序列化，把私钥转为[]byte
		pribyte, err := tools.Serialize(pri)
		if err!=nil{
			return err
		}
		err = bk.Put([]byte(addr), pribyte)
		if err !=nil{
			return err
		}
		return nil
	})
	return err
}

//生成地址
func (w *Wallet)NewAddress()(string,*ecdsa.PrivateKey,error){
	//获取公钥
	pri, pub, err := NewPubKeys()
	if err !=nil{
		return "",nil,err
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
	return addr,pri,nil
}


//用来获取私钥信息 ，真正的保存比特币私钥，是先多私钥进行AES加密，保存的是加密信息
func(w *Wallet) ShowPrivateKey(address string)(ecdsa.PrivateKey,error){
	db:=w.DB
	var privatekey ecdsa.PrivateKey
	err:=db.View(func(tx *bolt.Tx) error {
		bk:=tx.Bucket([]byte(PRIVATE_BUCKET))
		if bk == nil{
			return nil
		}
		priByts := bk.Get([]byte(address))

		//反序列化
		de := gob.NewDecoder(bytes.NewBuffer(priByts))
		gob.Register(elliptic.P256())
		err := de.Decode(&privatekey)
		if err !=nil{
			return err
		}
		return nil
	})
	return privatekey,err
}


func(w *Wallet)Checkaddr(addr string)bool{
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