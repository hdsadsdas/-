package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"公链系统开发/第三节课/block"
)

func main() {

	bc, err := block.NewChain([]byte("创始区块"))
	if err != nil {
		fmt.Println("hahha",err.Error())
		return
	}

	err = bc.AddBlock([]byte("哈哈哈"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bc.DB.View(func(tx *bolt.Tx) error {

		bucket := tx.Bucket([]byte(block.BUCKET_BLOCK))

		get := bucket.Get(bc.LastHash)

		serialize, err2 := block.DeSerialize(get)

		if err2 != nil{
			return err2
		}

		fmt.Println(string(serialize.Data))

		return nil

	})

}
