package interpreter

import (
	"errors"
	"fmt"
	"hachimi-lang/template"
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
	template.MOVER,
	template.INCMOVEL,
	template.JMP0,
	template.JMPTAG,
	template.INPUT,
	template.OUTPUT,
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
				fmt.Printf("哈气：位置 %d:%d 被哈气了！ '%c'\n", line+1, i, code[i])
				// 跳过当前字符继续解析
				return nil, errors.New("哈基耄耋不开心了！")
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
		if commandList[currentIdx] == template.JMP0 {
			depth++
		}
		if commandList[currentIdx] == template.JMPTAG {
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
		if cmd != template.JMP0 {
			machineCode = append(machineCode, CommandCode{
				Name: cmd,
				Tag:  0,
			})
		} else {
			//遇到跳转指令，寻找对应的跳转位置
			jmpIdx := getJmpIDX(idx, commandList)
			if jmpIdx == -1 {
				return nil, errors.New(fmt.Sprintf("%v 找不到 %v 的位置！", template.JMP0, template.JMPTAG))
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
	head := 0
	for head < len(machineCode) {
		cmd := machineCode[head]
		switch cmd.Name {
		case template.MOVER:
			err = cm.OperatorR()
		case template.INPUT:
			err = cm.OperatorInbox()
		case template.INCMOVEL:
			err = cm.OperatorLambda()
		case template.OUTPUT:
			fmt.Println("O<", cm.DataTape[cm.Ptr])
		case template.JMPTAG:
			// 不管
		case template.JMP0:
			//无条件跳转
			if cm.DataTape[cm.Ptr] != 0 {
				head = cmd.Tag
				continue
			}
		default:
			return errors.New("哈基米？哈基米？哈基米？？？？？")
		}
		if err != nil {
			return err
		}
		head++
	}

	return nil
}

func (cm *ComputerModel) OperatorR() error {
	if cm.Ptr == len(cm.DataTape)-1 {
		cm.Ptr = 0
		return nil
	}
	cm.Ptr++
	return nil
}

func (cm *ComputerModel) OperatorLambda() error {
	cm.DataTape[cm.Ptr]++
	if cm.Ptr == 0 {
		cm.Ptr = len(cm.DataTape) - 1
		return nil
	}
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
