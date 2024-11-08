package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fileName := "input.txt" // 这里假设文件名为 input.txt，可以根据需求修改

	// 打开文件
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("无法打开文件:", err)
		return
	}
	defer file.Close()

	sum := 0
	scanner := bufio.NewScanner(file)

	// 逐行读取文件内容
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line) // 分割字段

		// 确保字段数至少为2，否则跳过
		if len(fields) >= 2 {
			number, err := strconv.Atoi(fields[1]) // 将第二个字段转换为整数
			if err == nil {
				sum += number // 累加数值
			} else {
				fmt.Println("跳过无法转换的值:", fields[1])
			}
		}
	}

	// 检查读取过程中的错误
	if err := scanner.Err(); err != nil {
		fmt.Println("读取文件时出错:", err)
		return
	}

	fmt.Println("第二列数值之和:", sum)
}
