package transaction

import (
	"bytes"
	"公链系统开发/考试/tools"
	"公链系统开发/考试/wallet"
)

//创建一个交易输入的结构体
type Input struct {
	//交易哈希
	Txid  []byte
	//交易输出索引位置
	Vout  int
	//解锁脚本
	//ScriptSig []byte
	Sig []byte //签名
	Pubk []byte //公钥

}

func NewInput(txid []byte,vout int,sig []byte ,pubk []byte)Input{
	return Input{txid,vout,sig,pubk}
}

//判断input是某个人的消费
func(input *Input)IsLocked(name string)bool{
	if name ==""{
		return false
	}
	pubHash, err := wallet.GetPubHash(name)
	if err !=nil{
		return false
	}

	pub_sha256 := tools.GetHash(input.Pubk)
	pubhash2 := tools.Ripemd160(pub_sha256)
	return bytes.Compare(pubhash2,pubHash) == 0
}