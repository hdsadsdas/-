package main

import (

	"fmt"
	"公链系统开发/考试/client"
)

func main() {

	cl:=client.Cli{}
	cl.Run()
	//args := os.Args
	////把today变成一个功能
	//sum := flag.NewFlagSet("sum", flag.ExitOnError)
	//sub := flag.NewFlagSet("sub", flag.ExitOnError)
	//switch args[1] {
	//case "sum":
	//	num1:= sum.Int("num1", 0, "can参数1")
	//	num2 := sum.Int("num2", 0, "参数2")
	//	sum.Parse(args[2:])
	//	fmt.Println(*num1 + *num2)
	//	break
	//case "sub":
	//	num1:= sub.Int("num1", 0, "can参数1")
	//	num2 := sub.Int("num2", 0, "参数2")
	//	sub.Parse(args[2:])
	//	fmt.Println(*num1 - *num2)
	//	break
	//default:
	//	fmt.Println("没有对用的功能")
	//	os.Exit(1)
	//
	//
	//}
	//
	//判断用户输入的第二个位置上输入的内容
	//if args[1] =="today"{
	//
	//	day := today.String("day", "", "星期一")
	//
	//	//today.Parse()参数的意思：要解析的参数的范围
	//	today.Parse(args[2:])
	//	Today(*day)
	//	fmt.Println("wei")
	//}


	//args := os.Args
	////1.遍历  2.通过下标
	//fmt.Println(args[1])
	//
	//age := flag.Int("age", 16, "年龄")
	//
	//var name string
	//flag.StringVar(&name,"name","张三","年龄")
	//flag.Parse()
	//fmt.Println(*age)
	//fmt.Println(name)


}

func Today(day string){
	fmt.Println(day)
}
