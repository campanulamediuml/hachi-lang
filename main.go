package main

import (
	"fmt"
	"hachimi-lang/interpreter"
	"io"
	"os"
	"strings"
)

func ShowHelp() {
	desc := `哈基语虚拟机 (hachimi-lang)
	版本: v1.0.0
	作者: 哈基米
	哈基语是一种基于lambda算子的图灵完备语言，虚拟机支持最长65536字节的内存地址进行操作。
	`
	fmt.Println(desc)
	helpInfo := `哈基语虚拟机使用说明:
	运行命令: hachimi [filename.hachi]
		filename.hachi: 需要运行的.hachi后缀的hachimi语言代码文件
		-h: 显示帮助信息
	`
	fmt.Println(helpInfo)
	grammarInfo := `哈基语语法说明:
	哈基米：向右移动一位地址
	南北绿豆：读取控制台输入的数字
	阿西卡：当前地址的数值加一，并向左移动一位地址
	奶龙 / 曼波：成对出现，曼波处判断地址的数值是否为零，若为零则跳转到对应的【奶龙】处继续执行
	叮咚鸡：输出当前地址的数值到控制台
    
	哈基语支持 // 单行注释，注释内容会被忽略
	`
	fmt.Println(grammarInfo)
}

func main() {
	sysargs := os.Args
	fileName := sysargs[len(sysargs)-1]
	if fileName == "-h" {
		ShowHelp()
		return
	}
	fileTail := strings.Split(fileName, ".")
	if fileTail[len(fileTail)-1] != "hachi" {
		fmt.Println("只能运行.hachi后缀的文件")
		return
	}
	fd, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fd.Close()
	code, err := io.ReadAll(fd)
	//fmt.Println("开始运行虚拟机...")
	vm, err := interpreter.InitMachine(string(code))
	if err != nil {
		fmt.Println("初始化虚拟机错误:", err)
		return
	}
	//fmt.Println("虚拟机初始化成功，开始运行代码...")
	err = vm.Run()
	if err != nil {
		fmt.Println("运行错误:", err)
		return
	}
}
