package tools

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type FileHandler struct {
	reader *bytes.Buffer
	isSkip bool
}

func NewFileHandlerWithReader(reader *bytes.Buffer) *FileHandler {
	return &FileHandler{
		reader: reader,
		isSkip: false,
	}
}

func NewFileHandler() *FileHandler {
	return &FileHandler{
		isSkip: false,
	}
}

func (handler *FileHandler) Result() string {
	if handler.reader == nil {
		return ""
	}
	return handler.reader.String()
}

func (handler *FileHandler) ResultLines() []string {
	if handler.reader == nil {
		return []string{}
	}
	arr := strings.Split(handler.reader.String(), "\n")
	retArr := []string{}
	for _, line := range arr {
		l := strings.TrimSpace(line)
		if l != "" {
			retArr = append(retArr, l)
		}
	}
	return retArr
}

func (handler *FileHandler) GreapWithFile(patten string, filePath string) (bool, error) {
	if handler.isSkip {
		return handler.isSkip, nil
	}
	cmd := exec.Command("grep", patten, filePath)
	var out bytes.Buffer
	handler.reader = &bytes.Buffer{}
	cmd.Stdout = &out
	err := cmd.Run()
	handler.reader = &out
	if err != nil {
		if existErr, ok := err.(*exec.ExitError); ok {
			if existErr.ExitCode() == 1 {
				fmt.Println("未匹配到")
				handler.isSkip = true
				return handler.isSkip, nil
			} else {
				fmt.Printf("执行 grep 错误：%v\n", err)
			}
		} else {
			fmt.Printf("执行 grep 命令出错：%v\n", err)
		}
		return false, err
	}
	return false, nil
}

func (handler *FileHandler) Greap(patten string) (bool, error) {
	if handler.isSkip {
		return handler.isSkip, nil
	}
	if handler.reader != nil {
		cmd := exec.Command("grep", patten)
		var out bytes.Buffer
		cmd.Stdin = handler.reader
		cmd.Stdout = &out
		err := cmd.Run()
		handler.reader = &out
		if err != nil {
			if existErr, ok := err.(*exec.ExitError); ok {
				if existErr.ExitCode() == 1 {
					fmt.Println("未匹配到")
					handler.isSkip = true
					return handler.isSkip, nil
				} else {
					fmt.Printf("执行 grep 错误：%v\n", err)
				}
			} else {
				fmt.Printf("执行 grep 命令出错：%v\n", err)
			}
			return false, err
		}
	}
	return false, nil
}

func (handler *FileHandler) Awk(expression string) (bool, error) {
	if handler.isSkip {
		return handler.isSkip, nil
	}
	if handler.reader != nil {
		cmd := exec.Command("awk", expression)
		cmd.Stdin = handler.reader
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		handler.reader = &out
		if err != nil {
			if existErr, ok := err.(*exec.ExitError); ok {
				if existErr.ExitCode() == 1 {
					fmt.Println("未匹配到")
					handler.isSkip = true
					return handler.isSkip, nil
				} else {
					fmt.Printf("执行 grep 错误：%v\n", err)
				}
			} else {
				fmt.Printf("执行 grep 命令出错：%v\n", err)
			}
			return false, err
		}
	}
	return false, nil
}

func (handler *FileHandler) Head(count int) (bool, error) {
	if handler.isSkip {
		return handler.isSkip, nil
	}
	if handler.reader != nil {
		cmd := exec.Command("head", fmt.Sprintf("-n %d", count))
		cmd.Stdin = handler.reader
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		handler.reader = &out
		if err != nil {
			if existErr, ok := err.(*exec.ExitError); ok {
				if existErr.ExitCode() == 1 {
					fmt.Println("未匹配到")
					handler.isSkip = true
					return handler.isSkip, nil
				} else {
					fmt.Printf("执行 grep 错误：%v\n", err)
				}
			} else {
				fmt.Printf("执行 grep 命令出错：%v\n", err)
			}
			return false, err
		}
	}
	return false, nil
}

func ListDir(dir string) []string {
	cmd := exec.Command("ls", dir)
	var out bytes.Buffer

	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		if existErr, ok := err.(*exec.ExitError); ok {
			if existErr.ExitCode() == 1 {
				fmt.Println("未匹配到")
			} else {
				fmt.Printf("执行 grep 错误：%v\n", err)
			}
		} else {
			fmt.Printf("执行 grep 命令出错：%v\n", err)
		}
		return []string{}
	}
	files := out.String()
	arr := strings.Split(files, "\n")
	return arr
}
