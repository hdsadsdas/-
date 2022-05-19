package block

import (
	"errors"
	"github.com/boltdb/bolt"
)

/**
* @author : 哈哈
* @email : 598421227@qq.com
* @phone : 18816473550
* @DateTime : 2022/2/21 10:33
**/
//要保存的文件地址
const CHAIN_DB_PATH = "./chain.db"

//存区块的桶的名字
const BUCKET_BLOCK = "chain_blocks"

//保存最后区块hash值的桶的名字
const BUCKET_STATUS = "chain_status"

//用来存最后一个区块的hash值
const LAST_HASH = "last_hash"

//该结构体  区块链的作用 ：用来存储区块
type BlockChain struct {
	//Blocks []*Block
	//bolt数据库   非关系性数据库
	DB *bolt.DB

	LastHash []byte
}

/**
创建区块链  拥有创世区块
*/
func NewChain(data []byte) (*BlockChain, error) {

	//声明区块链
	bc := BlockChain{}
	//创建一个创世区块 ， 并存到区块链中
	genesis := GenesisBlock(data)

	//获得db 数据库连接
	db, err := bolt.Open(CHAIN_DB_PATH, 0600, nil)
	if err != nil {
		return nil, err
	}

	//添加数据
	err = db.Update(func(tx *bolt.Tx) error {

		//创建桶
		bk, err := tx.CreateBucket([]byte(BUCKET_BLOCK))
		if err != nil {
			return err
		}

		//将创世区块序列化
		serialize, err := genesis.Serialize()
		if err != nil {
			return err
		}

		//将创世区块存入桶中
		err = bk.Put(genesis.Hash, serialize)
		if err != nil {
			return err
		}

		//把最后一个区块的hash值放到另一个桶
		bk2, err := tx.CreateBucket([]byte(BUCKET_STATUS))
		if err != nil {
			return err
		}

		err = bk2.Put([]byte(LAST_HASH), genesis.Hash)
		if err != nil {
			return err
		}

		return nil

	})

	//bc.Blocks = []*Block{genesis}

	bc.DB = db
	bc.LastHash = genesis.Hash

	return &bc, nil

}

/**
把区块添加到区块链中
*/
func (bc *BlockChain) AddBlock(data []byte) error {

	//创建区块 添加到区块链中
	newBlock := NewBlock(data,bc.LastHash)

	err := bc.DB.Update(func(tx *bolt.Tx) error {

		bk := tx.Bucket([]byte(BUCKET_BLOCK))
		if bk == nil {
			return errors.New("桶中没有区块，还未创建区块链")
		}

		serialize, err := newBlock.Serialize()
		if err != nil {
			return err
		}

		bk.Put(newBlock.Hash, serialize)

		bk2 := tx.Bucket([]byte(BUCKET_STATUS))
		if bk == nil {
			return errors.New("桶中没有区块，还未创建区块链")
		}

		bk2.Put([]byte(LAST_HASH), newBlock.Hash)

		bc.LastHash = newBlock.Hash

		return nil

	})
	if err != nil {
		return err
	}

	return nil

}

//创建一个迭代器对象,迭代器只能在有区块链的情况下才能使用迭代器
func (bc *BlockChain)Iterator()*ChainIterator{

	iterator := ChainIterator{
		DB: bc.DB,
		currentHash: bc.LastHash,
	}

	return &iterator

}


