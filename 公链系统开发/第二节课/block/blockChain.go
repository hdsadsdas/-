package block

/**
* @author : 哈哈
* @email : 598421227@qq.com
* @phone : 18816473550
* @DateTime : 2022/2/21 10:33
**/

//该结构体  区块链的作用 ：用来存储区块
type BlockChain struct {
	Blocks []*Block
}

/**
创建区块链  拥有创世区块
 */
func NewChain(data []byte)*BlockChain {

	bc := BlockChain{}

	//创建一个创世区块 ， 并存到区块链中
	genesis := GenesisBlock(data)

	bc.Blocks = []*Block{genesis}

	return &bc

}

/**
把区块添加到区块链中
 */
func (bc *BlockChain)AddBlock(data []byte){

	//创建区块 添加到区块链中
	newBlock := NewBlock(data, bc.Blocks[len(bc.Blocks)-1].Hash)
	bc.Blocks = append(bc.Blocks,newBlock)
}