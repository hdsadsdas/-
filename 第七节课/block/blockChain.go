package block

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"公链系统开发/第七节课/transaction"
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

//func GetBlockChain() (*BlockChain,error) {
//
//	blockChain := &BlockChain{}
//
//	db, err := bolt.Open(CHAIN_DB_PATH, 0600, nil)
//
//	blockChain.DB = db
//
//	if err !=nil{
//		return nil,err
//	}
//	db.Update(func(tx *bolt.Tx) error {
//
//		bucket := tx.Bucket([]byte(BUCKET_STATUS))
//
//		LastHash := bucket.Get([]byte(LAST_HASH))
//
//		blockChain.LastHash = LastHash
//
//		return nil
//	})
//
//	return blockChain,nil
//
//}

func NewChain(address string) (*BlockChain, error) {
	//打开数据库
	db, err := bolt.Open(CHAIN_DB_PATH, 0600, nil)
	if err != nil {
		return nil, err
	}
	var lastHash []byte
	//向数据库中添加数据
	//同一个时间内，只能有一个人来进行写操作
	err = db.Update(func(tx *bolt.Tx) error {

		bk := tx.Bucket([]byte(BUCKET_BLOCK))
		//判断是否存在桶
		if bk == nil {

			coinbase, err := transaction.NewCoinbase(address)
			if err != nil {
				return err
			}

			//获取创世区块
			genesis := GenesisBlock(*coinbase)
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
			//将创世区块放入到桶1中
			bk.Put(genesis.Hash, serialize)
			//创建桶2
			bk2, err := tx.CreateBucket([]byte(BUCKET_STATUS))
			//将创世区块的hash值放到桶2中
			bk2.Put([]byte(LAST_HASH), genesis.Hash)
			//得到最后一位hash值
			lastHash = genesis.Hash
		} else {
			//当桶存在时
			bk2 := tx.Bucket([]byte(BUCKET_STATUS))
			//直接从桶中得到最后一位hash值
			lastHash = bk2.Get([]byte(LAST_HASH))
		}
		return nil
	})

	bc := BlockChain{
		DB:       db,
		LastHash: lastHash,
	}

	return &bc, err
}

func (bc *BlockChain) AddBlock(tx []transaction.Transaction) error {

	new := NewBlock(bc.LastHash, tx)

	err := bc.DB.Update(func(tx *bolt.Tx) error {

		bk := tx.Bucket([]byte(BUCKET_BLOCK))

		if bk == nil {

			return errors.New("没有创建桶")

		}

		serialize, _ := new.Serialize()

		bk.Put(new.Hash, serialize)

		bk2 := tx.Bucket([]byte(BUCKET_STATUS))

		if bk2 == nil {
			return errors.New("没有创建桶2")
		}

		bk2.Put([]byte(LAST_HASH), new.Hash)

		bc.LastHash = new.Hash

		return nil
	})
	return err
}

//创建一个迭代器对象,迭代器只能在有区块链的情况下才可以使用迭代器
func (bc *BlockChain) Iterator() *ChainIterator {

	iterator := ChainIterator{
		DB:          bc.DB,
		currentHash: bc.LastHash,
	}

	return &iterator
}

func (bc *BlockChain) GetAllBlock() ([]*Block, error) {

	blocks := []*Block{}

	iterator := bc.Iterator()
	for {
		if iterator.HasNext() {
			bk, err := iterator.Next()
			if err != nil {
				return nil, err
			}
			blocks = append(blocks, bk)

		} else {

			break
		}
	}

	return blocks, nil

}

//找某个地址的所有的收入（output）
func (bc *BlockChain) FindAllOutput(address string) []transaction.UTXO {

	allUTXO := make([]transaction.UTXO, 0)

	//先找所有的区块，再获取每个区块中的所有的交易
	blocks, err := bc.GetAllBlock()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	//获取每一个区块
	for _, block := range blocks {

		//获取每一笔交易
		for _, tx := range block.Txs {
			//获取每一笔交易中的交易输出
			for outIndex, output := range tx.Output {

				if output.IsUnlock(address) {

					utxo := transaction.NewUTXO(tx.TXid, outIndex, &output)
					allUTXO = append(allUTXO, *utxo)

				}

			}

		}

	}

	return allUTXO

}

//获取某个人的交易输入（Input）

func (bc *BlockChain) FindAllInput(address string) []transaction.Input {

	allInput := make([]transaction.Input, 0)

	blocks, err := bc.GetAllBlock()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	//得到每个区块
	for _, block := range blocks {

		for _, tx := range block.Txs {

			for _, input := range tx.Input {

				if input.IsLocked(address) {

					allInput = append(allInput, input)

				}

			}

		}

	}

	return allInput

}

//通过某个人的收入和消费，并返回可用的金额 (可用的output)
func (bc *BlockChain) FindSpendOutPut(address string,amount uint) ([]transaction.UTXO,uint){
	allOutput := bc.FindAllOutput(address)
	allInput := bc.FindAllInput(address)

	for _, input := range allInput {

		for index, outputs := range allOutput {

			if bytes.Compare(outputs.Txid, input.TXid) == 0 || outputs.Index == input.Vout {

				allOutput = append(allOutput[:index],allOutput[index+1:]...)
				break
			}

		}

	}

	var totalAmount uint = 0
	utxo := make([]transaction.UTXO,0)

	for _,outputs := range allOutput {
		totalAmount += outputs.Value
		utxo = append(utxo,outputs)
		if totalAmount >= amount {
			break
		}

	}

	return utxo,totalAmount

}

//找到所有可用的金额
func (bc *BlockChain) FindAllMoney(address string) uint {
	allOutput := bc.FindAllOutput(address)
	allInput := bc.FindAllInput(address)

	var Money uint = 0

	for _, input := range allInput {

		for index, output := range allOutput {

			if bytes.Compare(output.Txid, input.TXid) == 0 || output.Index == input.Vout {

				allOutput = append(allOutput[:index],allOutput[index+1:]...)
				break

			}

		}

	}

	for _,k :=range allOutput{

		Money+=k.Value

	}

	return Money

}

func (bc *BlockChain)NewTransaction(from string,to string,amount uint)(*transaction.Transaction,error)  {
	//1.构建input交易输入
	//在输入者中寻找所有交易的输出
	//从所有的可用的交易输出中先判断 是否 够  再取出

	utxos,totalAmount := bc.FindSpendOutPut(from,amount)
	if utxos == nil {
		return nil,errors.New("余额为空")
	}
	if totalAmount < amount {
		return nil, errors.New("余额不足")
	}

	newTransaction, err := transaction.NewTransaction(from, to, amount, utxos,totalAmount)
	if err != nil {
		return nil,errors.New("")
	}

	return newTransaction,nil

}


//创建普通的交易
func (bc *BlockChain) ChainCoinbase(address string) (*transaction.Transaction,error) {

	if address == "" || len(address)==0 {
		return nil,errors.New("地址错误")
	}

	return transaction.NewCoinbase(address)
}
