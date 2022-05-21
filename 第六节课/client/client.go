package client

import (
	"flag"
	"fmt"
	"os"
	"公链系统开发/第六节课/block"
	"公链系统开发/第六节课/tools"
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
		b.添加新的区块到区块链中  参数：有  1个  新区区块的交易信息  string
		c.打印所有区块信息       参数：无
		d.获取当前区块链中区块的个数  参数：无
		e.输出当前系统的使用说明    参数：无
	*/

	switch args[1] {

	case "createchain":

		cl.createchain()

	case "addblock":

		cl.addblock()

	case "printblock":

		cl.printblock()

	case "getblockcount":

		cl.getblockcount()

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
	s := createchain.String("data", "", "创世区块的交易信息")
	//解析参数
	createchain.Parse(os.Args[2:])

	//判断文件是否存在（是否存在链）
	if tools.FileExist("./chain.db") {

		fmt.Println("文件已经存在")
		return

	}

	//生成区块链
	bc, err := block.NewChain([]byte(*s))
	defer bc.DB.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("创建区块链成功")

}

func (cl *Cli) addblock() {
	//判断文件是否存在（是否存在链）
	if !tools.FileExist("./chain.db") {

		fmt.Println("区块链不存在")
		return

	}

	addblock := flag.NewFlagSet("addblock", flag.ExitOnError)

	s := addblock.String("data", "", "添加新区块的信息")

	addblock.Parse(os.Args[2:])

	bc, _ := block.NewChain(nil)
	defer bc.DB.Close()

	err := bc.AddBlock([]byte(*s))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("添加区块成功")
}

func (cl *Cli) printblock() {

	//判断文件是否存在（是否存在链）
	if !tools.FileExist("./chain.db") {

		fmt.Println("区块链不存在")
		return

	}

	bc, _ := block.NewChain(nil)
	defer bc.DB.Close()

	blocks, _ := bc.GetAllBlock()

	for _, v := range blocks {

		fmt.Println(string(v.Data))

	}

}

func (cl *Cli) getblockcount() {

	//判断文件是否存在（是否存在链）
	if !tools.FileExist("./chain.db") {

		fmt.Println("区块链不存在")
		return

	}

	bc, _ := block.NewChain(nil)
	defer bc.DB.Close()

	blocks, _ := bc.GetAllBlock()

	fmt.Println(len(blocks))

}
