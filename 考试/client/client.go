package client

import (
	"flag"
	"fmt"
	"os"
	"公链系统开发/考试/block"
	"公链系统开发/考试/tools"
	"公链系统开发/考试/transaction"
)

/**
用户的程序交互接口
该模块只负责读取用户传递的命令和参数，并进行解析
再去调用对用的功能
*/
type Cli struct {
}

func (cl *Cli) Run() {
	//获取用户的所有的输入

	//确定有哪些功能，这些功能需不需要参数
	/**
		1.创建带有创世区块的区块链  参数：1  创世区块的交易信息
		main.exe createchain --data "交易"
		2.发起一笔交易  参数：3
		mian.exe send --from "zhang" --to "liu" --amount 50
		3.获取区块链中所有区块的个数  参数：无
		main.exe getblockcount
	   	4.获取所有区块的信息  参数：无
		main.exe  allblock
		5.获取地址的余额  参数：1  哪一个地址的余额
		main.exe getbalance --address "地址"
		6.输出当前系统的使用说明  参数：无
		main.exe help
		7.获取对应地址的私钥
		main.exe showPirkey --address ""
	*/
	//把createchain变成功能
	switch os.Args[1] {
	//生成区块链
	case "createchain":
		cl.createchain()
		//查询私钥
	case "showPirkey":
		cl.showPirkey()
		//发起交易
	case "send":
		cl.send()
		//获取区块总数
	case "getblockcount":
		cl.getblockcount()
		//获取所有区块
	case "allblock":
		cl.allblock()
       //获取余额
	case "getbalance":
		cl.getbalance()
		//生成地址
	case "generateaddr":
		cl.generateaddr()
		//验证地址
	case "checkaddr":
		cl.checkaddr()
	case "help":
		cl.help()
	default:
		fmt.Println("没有对应的功能")
		os.Exit(1)

	}

}

func (cl *Cli) showPirkey() {
	showPirkey := flag.NewFlagSet("showPirkey", flag.ExitOnError)
	addr := showPirkey.String("addr", "", "要查询的地址")
	showPirkey.Parse(os.Args[2:])
	bc, err := block.CreatChain("")
	defer bc.DB.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	privateKey, err := bc.Wallet.ShowPrivateKey(*addr)
	if err != nil {
		fmt.Println(err.Error())
		return

	}
	fmt.Println("私钥为：", privateKey)
}

func (cl *Cli) checkaddr() {
	checkaddr := flag.NewFlagSet("checkaddr", flag.ExitOnError)
	addr := checkaddr.String("addr", "", "要判断的地址")
	checkaddr.Parse(os.Args[2:])
	bc, err := block.CreatChain("")
	defer bc.DB.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	istrue := bc.Wallet.Checkaddr(*addr)
	if istrue {
		fmt.Println("地址合法")
	} else {
		fmt.Println("地址不合法")
	}
}

func (cl *Cli) generateaddr() {

	bc, err := block.CreatChain("")
	defer bc.DB.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	addr, pri, err := bc.Wallet.NewAddress()
	if err != nil {
		fmt.Println("创建地址失败",err.Error())
	}

	err = bc.Wallet.SavePrivateKey(addr, pri)
	if err != nil {
		fmt.Println("保存地址错误",err.Error())
		return
	}
	fmt.Println("生成的地址为",addr)
}

func (cl *Cli) getbalance() {
	exist := tools.FileExist("./chain.db")
	if !exist {
		fmt.Println("区块链不存在")
		return
	}
	bc, err := block.CreatChain("")
	defer bc.DB.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	getbalance := flag.NewFlagSet("getbalance", flag.ExitOnError)
	address := getbalance.String("address", "", "矿工的地址")
	getbalance.Parse(os.Args[2:])
	istrue := bc.Wallet.Checkaddr(*address)
	if !istrue {
		fmt.Println("地址不合法")
		return
	}
	balance := bc.GetBalance(*address)
	fmt.Printf("地址%s的余额为：%d\n", *address, balance)
}

func (cl *Cli) help() {
	fmt.Println("本系统有一下功能")
	fmt.Println("1.创建带有创世区块的区块链  参数：1  创世区块的交易信息")
	fmt.Println("2.发送交易  参数：1  新区块的交易信息")
	fmt.Println("3.获取区块链中所有区块的个数  参数：无")
	fmt.Println("4.获取所有区块的信息  参数：无")
	fmt.Println("5.获取地址的余额  参数：1  哪一个地址的余额")
	fmt.Println("6.创建地址")
	fmt.Println("7.判断地址是否符合要求 参数 1 要判断的地址")
	fmt.Println("8.输出当前系统的使用说明  参数：无")
}
func (cl *Cli) allblock() {
	//区块hash值和交易信息
	exist := tools.FileExist("./chain.db")
	if !exist {
		fmt.Println("区块链不存在")
		return
	}
	bc, err := block.CreatChain("")
	defer bc.DB.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	blocks, err := bc.GetAllBlock()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//遍历获取每一个区块
	for _, value := range blocks {
		fmt.Printf("区块hash:%x,交易个数:%d\n", value.Hash, len(value.Txs))
		//遍历交易集合
		for _, tx := range value.Txs {
			fmt.Printf("\t交易hash:%x\n", tx.TXHash)
			fmt.Printf("\t\t有%d个交易输入\n", len(tx.Input))
			for index, input := range tx.Input {
				fmt.Printf("\t\t\t消费%d,来自%x,下标%d\n", index, input.Txid, input.Vout)
			}
			fmt.Printf("\t\t有%d个交易输出\n", len(tx.Output))
			for index, output := range tx.Output {
				fmt.Printf("\t\t\t收入下标%d,金额%d,属于%x\n", index, output.Value, output.ScriptPubkey)
			}
		}
	}

}

func (cl *Cli) getblockcount() {
	exist := tools.FileExist("./chain.db")
	if !exist {
		fmt.Println("区块链不存在")
		return
	}
	bc, err := block.CreatChain("")
	defer bc.DB.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	blocks, err := bc.GetAllBlock()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("一共有%d个区块", len(blocks))
}

//发起一笔交易，把区块添加到区块链中
func (cl *Cli) send() {
	send := flag.NewFlagSet("send", flag.ExitOnError)
	from := send.String("from", "", "交易发起者的地址")
	to := send.String("to", "", "交易接收者的地址")
	//正整数
	amount := send.Uint("amount", 0, "交易的数量")
	//参数的解析
	err := send.Parse(os.Args[2:])
	if err != nil {
		fmt.Println("解析失败", err.Error())
		return
	}
	//1.创建一个普通的交易
	//将构建新交易作为区块链的一个功能提供出来
	//2.把这个交易放到区块中，然后在把区块存储到区块链中
	//在2这个过程中，产生新区块的过程中，回产生一个coinbase交易
	bc, err := block.CreatChain("")
	defer bc.DB.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	from_istrue := bc.Wallet.Checkaddr(*from)
	if !from_istrue {
		fmt.Println("from地址不合法")
		return
	}
	to_istrue := bc.Wallet.Checkaddr(*to)
	if !to_istrue {
		fmt.Println("to的地址不合法")
		return
	}
	tx, err := bc.NewTransaction(*from, *to, *amount)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//创建一个coinbase交易
	cb, err := bc.NewCoinBase(*from)
	err = bc.AddBlock([]transaction.Transaction{*tx, *cb})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("交易成功")
}

func (cl *Cli) createchain() {

	//判断区块链是否存在
	exist := tools.FileExist("./chain.db")
	if exist {
		fmt.Println("区块链已存在")
		return
	}

	//调用创建区块链的方法  生成一个区块链
	bc, err := block.CreatChain("")
	defer bc.DB.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("创建成功")
}
