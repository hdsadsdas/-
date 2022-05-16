package transaction

/**
* @author : 哈哈
* @email : 598421227@qq.com
* @phone : 18816473550
* @DateTime : 2022/5/2 10:00
**/

//未消费的交易输出
type UTXO struct {

	Txid  []byte
	Index int
	*Output
}

//创建UTXO结构体
func NewUTXO(txid []byte,index int,output *Output)*UTXO  {

	return &UTXO{txid,index,output}

}
