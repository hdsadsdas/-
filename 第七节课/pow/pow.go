package pow

import (
	"bytes"
	"math/big"
	"strconv"
	"公链系统开发/第七节课/tools"
	"公链系统开发/第七节课/transaction"
)

const BITS = 1 //难度系数，前面有多少个0

/**
 * 区块的hash值 < 系统给定的hash值
 */

type ProofOfWork struct {
	//Block  *block.Block
	//TimeStamp int64
	//PrevHash  []byte
	//Data      []byteB

	Block BlockInterface

	Target *big.Int //系统给定的值
	//hash值转成2进制的 256位 不能用int64
}

type BlockInterface interface {
	GetTimeStamp() int64

	GetPrevHash() []byte

	GetTxs() []transaction.Transaction
}

/**
实例化一个pow结构体，并且返回
*/
func NewPow(block BlockInterface) *ProofOfWork {

	target := big.NewInt(1) //声明一个大整数类型的1
	//hash值 256 - 1 - bits  255 - bits
	target = target.Lsh(target, 255-BITS)
	pow := ProofOfWork{
		Block:  block,
		Target: target,
	}
	return &pow
}

//用来寻找随机数
func (pow *ProofOfWork) Run() ([]byte, int64) {
	var nonce int64 //随机数
	nonce = 0
	//block := pow.Block
	//block.Nonce = nonce

	timeBytes := []byte(strconv.FormatInt(pow.Block.GetTimeStamp(), 10))

	num := big.NewInt(0)

	//转型  []byte转成大整数
	for {
		nonceBytes := []byte(strconv.FormatInt(nonce, 10))

		txsBytes := []byte{}

		for _, v := range pow.Block.GetTxs() {

			txsByte, _ := v.Serialize()

			txsBytes = append(txsByte, txsByte...)

		}

		hashByets := bytes.Join([][]byte{txsBytes, pow.Block.GetPrevHash(), timeBytes, nonceBytes}, []byte{})
		hash := tools.GetHash(hashByets) //当前区块的hash值
		//fmt.Println("正在寻找nonce,当前的nonce为",nonce)
		num = num.SetBytes(hash) //用来转换成大整数的
		/*
			 	*hash  []byte  -》 大整数类型
				pow.Target  大整数类型
				A.cmp(B)
		*/
		//if(hash < pow.Target){
		//
		//}
		if num.Cmp(pow.Target) == -1 {
			return hash, nonce
		}
		nonce++
	}
	return nil, 0
}
