package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/cupcake/rdb"
	"github.com/cupcake/rdb/nopdecoder"
)

// KeyInfo 保存键的信息
type KeyInfo struct {
	Key     string
	Size    int64
	Percent float64
}

// main 函数
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path_to_rdb_file>")
		os.Exit(1)
	}

	rdbFile := os.Args[1]
	file, err := os.Open(rdbFile)
	if err != nil {
		fmt.Printf("Error opening RDB file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// 使用自定义解码器
	bigKeys := make(map[string]int64)
	totalSize := int64(0)

	decoder := &CustomDecoder{
		bigKeys:   bigKeys,
		totalSize: &totalSize,
	}
	err = rdb.Decode(file, decoder)
	if err != nil {
		fmt.Printf("Error decoding RDB file: %v\n", err)
		os.Exit(1)
	}

	// 计算和排序
	var keyInfoList []KeyInfo
	for key, size := range bigKeys {
		percent := (float64(size) / float64(*decoder.totalSize)) * 100
		keyInfoList = append(keyInfoList, KeyInfo{Key: key, Size: size, Percent: percent})
	}

	sort.Slice(keyInfoList, func(i, j int) bool {
		return keyInfoList[i].Size > keyInfoList[j].Size
	})

	// 输出结果
	fmt.Printf("Total size: %d bytes\n", *decoder.totalSize)
	fmt.Println("BigKeys analysis:")
	fmt.Println("Key\tSize(Bytes)\tPercent(%)")
	for _, info := range keyInfoList {
		fmt.Printf("%s\t%d\t%.2f%%\n", info.Key, info.Size, info.Percent)
	}
}

// CustomDecoder 自定义解码器
type CustomDecoder struct {
	nopdecoder.NopDecoder // 嵌入空实现的解码器
	bigKeys               map[string]int64
	totalSize             *int64
}

// Set 处理每个键
func (d *CustomDecoder) Set(key, value []byte, expiry int64) {
	size := int64(len(value))
	d.bigKeys[string(key)] = size
	*d.totalSize += size
}

// StartHash 处理哈希键
func (d *CustomDecoder) StartHash(key []byte, length, expiry int64) {
	d.bigKeys[string(key)] = 0 // 初始大小
}

// HSet 哈希字段
func (d *CustomDecoder) HSet(key, field, value []byte) {
	size := int64(len(field) + len(value))
	d.bigKeys[string(key)] += size
	*d.totalSize += size
}

// StartSet 处理集合键
func (d *CustomDecoder) StartSet(key []byte, cardinality, expiry int64) {
	d.bigKeys[string(key)] = 0
}

// Sadd 集合成员
func (d *CustomDecoder) Sadd(key, member []byte) {
	size := int64(len(member))
	d.bigKeys[string(key)] += size
	*d.totalSize += size
}

// StartList 处理列表键
func (d *CustomDecoder) StartList(key []byte, length, expiry int64) {
	d.bigKeys[string(key)] = 0
}

// RPush 列表成员
func (d *CustomDecoder) RPush(key, value []byte) {
	size := int64(len(value))
	d.bigKeys[string(key)] += size
	*d.totalSize += size
}

// StartZSet 处理有序集合键
func (d *CustomDecoder) StartZSet(key []byte, cardinality, expiry int64) {
	d.bigKeys[string(key)] = 0
}

// ZAdd 有序集合成员
func (d *CustomDecoder) ZAdd(key []byte, score float64, member []byte) {
	size := int64(len(member))
	d.bigKeys[string(key)] += size
	*d.totalSize += size
}
