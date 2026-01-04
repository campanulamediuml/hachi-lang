package main

import (
	"fmt"
	"hachimi-lang/interpreter"
)

func main() {
	code := "哈基米阿西卡哈基米阿西卡哈基米阿西卡哈基米叮咚鸡"
	vm, err := interpreter.InitMachine(code)
	if err != nil {
		fmt.Println("初始化错误:", err)
		return
	}
	err = vm.Run()
	if err != nil {
		fmt.Println("运行错误:", err)
		return
	}
}
