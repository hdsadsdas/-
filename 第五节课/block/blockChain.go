package block

import (
	"errors"
	"github.com/boltdb/bolt"

)

//要保存的文件地址
const CHAIN_DB_PATH = "./chain.db"

//存区块的桶的名字
const BUCKET_BLOCK = "chain_blocks"

//保存最后区块hash值的桶的名字
const BUCKET_STATUS = "chain_status"

//用来存最后一个区块的hash值
const LAST_HASH = "last_hash"

type BlockChain struct {
	//Blcocks  []*Block
	DB *bolt.DB

	LastHash []byte
}

func NewChain(data []byte)(*BlockChain,error){
	//打开数据库
	db, err := bolt.Open(CHAIN_DB_PATH, 0600, nil)
	if err !=nil{
		return nil,err
	}
	var lastHash []byte
	//向数据库中添加数据
	//同一个时间内，只能有一个人来进行写操作
	err = db.Update(func(tx *bolt.Tx) error {

		bk := tx.Bucket([]byte(BUCKET_BLOCK))
		//判断是否存在桶
		if bk == nil{
			//获取创世区块
			genesis := GenesisBlock(data)
			//创建桶
			bk, err := tx.CreateBucket([]byte(BUCKET_BLOCK))
			if err !=nil{
				return err
			}
			//将创世区块序列化
			serialize, err := genesis.Serialize()
			if err !=nil{
				return err
			}
			//将创世区块放入到桶1中
			bk.Put(genesis.Hash,serialize)
			//创建桶2
			bk2, err := tx.CreateBucket([]byte(BUCKET_STATUS))
            //将创世区块的hash值放到桶2中
			bk2.Put([]byte("LAST_HASH"),genesis.Hash)
			//得到最后一位hash值
			lastHash = genesis.Hash
		}else{
			//当桶存在时
			bk2 := tx.Bucket([]byte(BUCKET_STATUS))
			//直接从桶中得到最后一位hash值
			lastHash = bk2.Get([]byte(LAST_HASH))
		}
		return nil
	})

	bc:=BlockChain{
	DB: db,
	LastHash: lastHash,
	}
	
	return &bc,err
}

func (bc *BlockChain)AddBlock(data []byte)error{

	new:=NewBlock(bc.LastHash,data)

	err :=bc.DB.Update(func(tx *bolt.Tx) error {

		bk := tx.Bucket([]byte(BUCKET_BLOCK))

		if bk == nil{

			return errors.New("没有创建桶")

		}

		serialize, _ := new.Serialize()

		bk.Put(new.Hash,serialize)

		bk2 := tx.Bucket([]byte(BUCKET_STATUS))

		if bk2 == nil{
			return errors.New("没有创建桶2")
		}

		bk2.Put([]byte(LAST_HASH),new.Hash)

		bc.LastHash = new.Hash

		return nil
	})
	return err
}

//创建一个迭代器对象,迭代器只能在有区块链的情况下才可以使用迭代器
func (bc *BlockChain)Iterator()*ChainIterator{

	iterator:=ChainIterator{
		DB: bc.DB  ,
		currentHash: bc.LastHash,
	}

	return &iterator
}
