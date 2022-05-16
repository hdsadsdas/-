package client

import (
	"flag"
	"fmt"
	"os"
	"公链系统开发/第七节课/block"
	"公链系统开发/第七节课/tools"
	"公链系统开发/第七节课/transaction"
)

/**
* @author : 哈哈
* @email : 598421227@qq.com
* @phone : 18816473550
* @DateTime : 2022/3/28 10:54
**/

type Cli struct {
}

func (cl *Cli) Run() {

	//获取用户输入的参数
	args := os.Args

	/*
		确定系统需要哪些功能，需要哪些参数
		a.创建带有创世区块的区块链 参数：有 1个  创世区块的交易信息  string
		b.发起一笔交易 参数：有  3个
		c.打印所有区块信息       参数：无
		d.获取当前区块链中区块的个数  参数：无
		e.输出当前系统的使用说明    参数：无
	*/

	switch args[1] {

	case "createchain":

		cl.createchain()

	case "printblock":

		cl.printblock()

	case "getblockcount":

		cl.getblockcount()
	case "send":

		cl.send()
	case "allmoney":
		cl.allmoney()

	case "help":

		fmt.Println("a.创建带有创世区块的区块链，参数：有 1个  创世区块的交易信息  string\n\tmain.exe createchain --data \"123\"\n\tb.添加新的区块到区块链中， 参数：有  1个  新区区块的交易信息  string\n\tmain.exe addblock --data \"456\"\n\tc.打印所有区块信息   ，参数：无\n\tmain.exe printblock\n\td.获取当前区块链中区块的个数 ,参数：无\n\tmain.exe getblockcount\n\te.输出当前系统的使用说明 ，参数：无")

	default:
		fmt.Println("没有对应的功能")
		os.Exit(1)

	}

}

func (cl *Cli) createchain() {

	//创建命令
	createchain := flag.NewFlagSet("createchain", flag.ExitOnError)
	//命令包含的参数
	address := createchain.String("address", "", "账户名称")
	//解析参数
	createchain.Parse(os.Args[2:])

	//判断文件是否存在（是否存在链）
	if tools.FileExist("./chain.db") {

		fmt.Println("文件已经存在")
		return

	}

	//生成区块链
	bc, err := block.NewChain(*address)
	defer bc.DB.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("创建区块链成功")

}

//func (cl *Cli) addblock() {
//	//判断文件是否存在（是否存在链）
//	if !tools.FileExist("./chain.db") {
//
//		fmt.Println("区块链不存在")
//		return
//
//	}
//
//	addblock := flag.NewFlagSet("addblock", flag.ExitOnError)
//
//	s := addblock.String("data", "", "添加新区块的信息")
//
//	addblock.Parse(os.Args[2:])
//
//	bc, _ := block.NewChain(nil)
//	defer bc.DB.Close()
//
//	err := bc.AddBlock([]byte(*s))
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//
//	fmt.Println("添加区块成功")
//}

func (cl *Cli) printblock() {

	//判断文件是否存在（是否存在链）
	if !tools.FileExist("./chain.db") {

		fmt.Println("区块链不存在")
		return

	}

	bc, _ := block.NewChain("")
	defer bc.DB.Close()

	blocks, _ := bc.GetAllBlock()

	for _, v := range blocks {

		fmt.Println("区块hash:", v.Hash, ",交易的个数:", len(v.Txs))

		for _, tx := range v.Txs {

			for inputKey, inputV := range tx.Input {
				fmt.Println("第", inputKey, "个交易输入：")
				fmt.Println("		消费", inputV.Vout, "来自", inputV.TXid, "下标为", inputKey)
			}

			for outKey, outV := range tx.Output {
				fmt.Println("第", outKey, "个交易输出：")
				fmt.Println("		收入", outV.Value, "属于", string(outV.ScriptPutKey))
			}
		}

	}

}

func (cl *Cli) getblockcount() {

	//判断文件是否存在（是否存在链）
	if !tools.FileExist("./chain.db") {

		fmt.Println("区块链不存在")
		return

	}

	bc, _ := block.NewChain("")
	defer bc.DB.Close()

	blocks, _ := bc.GetAllBlock()

	fmt.Println(len(blocks))

}

//发起一笔交易,再添加到区块链中
func (cl *Cli) send() {
	//创建命令
	send := flag.NewFlagSet("send", flag.ExitOnError)

	from := send.String("from", "", "交易发起者的地址")

	to := send.String("to", "", "交易接收者的地址")

	amount := send.Uint("amount", 0, "交易金额")

	err := send.Parse(os.Args[2:])
	if err != nil {
		fmt.Println("解析失败", err.Error())
		return
	}

	//2.把这个交易放到区块中，再把区块添加到区块链中
	bc, err := block.NewChain("")
	defer bc.DB.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	tx, err := bc.NewTransaction(*from, *to, *amount)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	cb, err := bc.ChainCoinbase(*from)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = bc.AddBlock([]transaction.Transaction{*tx, *cb})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func (cl *Cli) allmoney() {

	//创建命令
	createchain := flag.NewFlagSet("allmoney", flag.ExitOnError)
	//命令包含的参数
	address := createchain.String("address", "", "账户名称")
	//解析参数
	createchain.Parse(os.Args[2:])

	bc, _ := block.NewChain("")
	defer bc.DB.Close()

	money := bc.FindAllMoney(*address)

	fmt.Println("可用金额为", money)

}
