package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"errors"
	"github.com/boltdb/bolt"
	"公链系统开发/第八节课/tools"
)

/**
* @author : 哈哈
* @email : 598421227@qq.com
* @phone : 18816473550
* @DateTime : 2022/5/30 9:48
**/

const PRIVATE_BUCKET = "private_bucket"

type Wallet struct {
	DB *bolt.DB
}

func NewWallet(db *bolt.DB) (*Wallet, error) {
	if db == nil {
		return nil, errors.New("db错误")
	}

	err := db.Update(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(PRIVATE_BUCKET))
		if bk == nil {
			_, err := tx.CreateBucket([]byte(PRIVATE_BUCKET))
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &Wallet{db}, nil

}

func (w *Wallet) NewAddress() (string, error) {

	pri, pub, err := NewPubKeys()
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

	db := w.DB

	db.Update(func(tx *bolt.Tx) error {

		bk := tx.Bucket([]byte(PRIVATE_BUCKET))

		if bk == nil {
			return errors.New("桶不存在")
		}

		//序列化把私钥转为[]byte

		pribyte, err := tools.Serialize(pri)
		if err != nil {
			return err
		}

		err = bk.Put([]byte(addr), pribyte)

		if err != nil {
			return err
		}

		return nil
	})

	return addr, nil
}

//用来获取私钥信息，真实的保存比特币私钥，是先对私钥继续AES加密，保存的是加密信息
func (w *Wallet) ShowPrivateKey(address string) (*ecdsa.PrivateKey, error) {

	db := w.DB
	var privatekey ecdsa.PrivateKey

	err := db.View(func(tx *bolt.Tx) error {

		bk := tx.Bucket([]byte(PRIVATE_BUCKET))

		if bk == nil {
			return errors.New("桶为空")
		}

		priBytes := bk.Get([]byte(address))

		de := gob.NewDecoder(bytes.NewBuffer(priBytes))
		err := de.Decode(&privatekey)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &privatekey, nil

}

func (w *Wallet) Checkaddr(addr string) bool {

	ver_pubhash_check := tools.Decode(addr)

	check := ver_pubhash_check[len(ver_pubhash_check)-4:]
	ver_pubhash := ver_pubhash_check[:len(ver_pubhash_check)-4]

	first_hash := tools.GetHash(ver_pubhash)
	second_hash := tools.GetHash(first_hash)

	check2 := second_hash[:4]

	return bytes.Compare(check, check2) == 0

}
