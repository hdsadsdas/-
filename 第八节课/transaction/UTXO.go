package transaction
/**
	用来描述用户的可以消费的交易输出的信息
	1.该交易输出在哪个交易中：txid
    2.该交易输出在交易中是第几个：index
	3.该交易输出属于谁
	4.该交易输出的金额
 */
type UTXO struct {
	Txid []byte
	Index int
	Output //匿名字段  utxo这个结构体就会默认包含了output中的两个字段
}

//创建utxo结构体的,
func NewUTXO(txid []byte,index int,output Output)UTXO{
	return UTXO{txid,index,output}
}