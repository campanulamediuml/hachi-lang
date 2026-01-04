package interpreter

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type ComputerModel struct {
	DataTape [65536]byte
	Ptr      int
	Command  []string
}

type CommandCode struct {
	Name string
	Tag  int
}

var Tokenize = []string{
	"哈基米",
	"南北绿豆",
	"阿西卡",
	"奶龙",
	"曼波",
	"叮咚鸡",
}

func Parser(AllCode string) ([]string, error) {
	var result []string

	// 为了提高效率，将token按长度降序排序
	// 这样可以优先匹配最长的token
	stringLines := strings.Split(AllCode, "\n")
	for line, codeLine := range stringLines {
		codeList := strings.Split(codeLine, "//")
		code := codeList[0]
		code = strings.TrimSpace(code)
		code = strings.ReplaceAll(code, " ", "")
		if code == "" {
			continue
		}
		sortedTokens := make([]string, len(Tokenize))
		copy(sortedTokens, Tokenize)
		sort.Slice(sortedTokens, func(i, j int) bool {
			return len(sortedTokens[i]) > len(sortedTokens[j])
		})

		i := 0
		for i < len(code) {
			matched := false
			// 尝试匹配所有token
			for _, token := range sortedTokens {
				tokenLen := len(token)
				if i+tokenLen <= len(code) && code[i:i+tokenLen] == token {
					result = append(result, token)
					i += tokenLen
					matched = true
					break
				}
			}

			if !matched {
				// 无法匹配任何token，打印错误信息
				fmt.Printf("解析错误：位置 %d:%d 处遇到未定义的词汇 '%c'\n", line+1, i, code[i])
				// 跳过当前字符继续解析
				return nil, errors.New("解析错误：遇到未定义的词汇")
			}
		}
	}
	return result, nil

}

func InitMachine(code string) (*ComputerModel, error) {
	command, err := Parser(code)
	if err != nil {
		return nil, err
	}
	res := &ComputerModel{
		DataTape: [65536]byte{},
		Ptr:      0,
		Command:  command,
	}
	return res, nil
}

func getJmpIDX(currentIdx int, commandList []string) int {
	// 寻找对应的跳转位置
	depth := 0
	for {
		currentIdx--
		if commandList[currentIdx] == "曼波" {
			depth++
		}
		if commandList[currentIdx] == "奶龙" {
			if depth == 0 {
				return currentIdx
			} else {
				depth--
			}
		}
		if currentIdx <= 0 {
			return -1
		}
	}
}

func (cm *ComputerModel) Compile() ([]CommandCode, error) {
	commandList := cm.Command
	//fmt.Println(commandList)
	machineCode := make([]CommandCode, 0)
	for idx, cmd := range commandList {
		if cmd != "曼波" {
			machineCode = append(machineCode, CommandCode{
				Name: cmd,
				Tag:  0,
			})
		} else {
			//遇到跳转指令，寻找对应的跳转位置
			jmpIdx := getJmpIDX(idx, commandList)
			if jmpIdx == -1 {
				return nil, errors.New("跳转指令匹配错误")
			}
			machineCode = append(machineCode, CommandCode{
				Name: cmd,
				Tag:  jmpIdx,
			})
		}
	}
	return machineCode, nil
}

func (cm *ComputerModel) Run() error {
	machineCode, err := cm.Compile()
	if err != nil {
		return err
	}
	pc := 0
	for pc < len(machineCode) {
		cmd := machineCode[pc]
		switch cmd.Name {
		case "哈基米":
			err = cm.OperatorR()
		case "南北绿豆":
			err = cm.OperatorInbox()
		case "阿西卡":
			err = cm.OperatorLambda()
		case "叮咚鸡":
			fmt.Println("O<", cm.DataTape[cm.Ptr])
		case "奶龙":
			// 不管
		case "曼波":
			//无条件跳转
			if cm.DataTape[cm.Ptr] != 0 {
				pc = cmd.Tag
				continue
			}
		default:
			return errors.New("未知指令")
		}
		if err != nil {
			return err
		}
		pc++
	}

	return nil
}

func (cm *ComputerModel) OperatorR() error {
	if cm.Ptr >= len(cm.DataTape)-1 {
		return errors.New("数据指针越界")
	}
	cm.Ptr++
	return nil
}

func (cm *ComputerModel) OperatorLambda() error {
	if cm.Ptr <= 0 {
		return errors.New("数据指针越界")
	}
	cm.DataTape[cm.Ptr]++
	cm.Ptr--
	return nil
}

func (cm *ComputerModel) OperatorInbox() error {
	//读取控制台输入
	var input string
	fmt.Printf("I> ")
	_, err := fmt.Scanln(&input)
	if err != nil {
		return err
	}
	digit, err := strconv.Atoi(input)
	if err != nil {
		return err
	}
	cm.DataTape[cm.Ptr] = byte(digit)
	return nil
}
