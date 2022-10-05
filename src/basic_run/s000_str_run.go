package basic_run

import (
	"fmt"
	"strconv"
	"strings"
)

/** 迭代测试 begin */
func StringIterRun() {
	fmt.Printf("Begin: =================StringIterRun()=========================\n")
	var s string
	s = "abc123@#$我你"
	fmt.Printf("=======直接迭代字符串 %s======\n", s)
	// len(s) 的总长度等于 len([]byte(s)) 相当于是字符串占用空间的长度
	// 但是range会按照字符返回 但是索引是按照 bytes来的 所以需要特别注意一下
	fmt.Printf("=======转换为char迭代 总长度%d======\n", len(s))
	for i, val := range s {
		fmt.Printf("%d:char=%c\n", i, val)
	}
	// len([]rune(s)) 的长度试unicode字符的个数 中文字符 一个字符占用一个unicode字符
	fmt.Printf("=======转换为rune迭代 总长度%d======\n", len([]rune(s)))
	for i, val := range []rune(s) {
		fmt.Printf("%d:rune=%d\n", i, val)
	}
	// 按照底层utf-8存储的字节来迭代
	fmt.Printf("=======转换为bytes迭代 总长度 %d======\n", len([]byte(s)))
	for i, val := range []byte(s) {
		fmt.Printf("%d:byte=%d\n", i, val)
	}
	// 字符串 截取 截取索引需要按照字节索引来 而不是字符索引
	start_idx := 8
	end_idx := 15
	fmt.Printf("截取字符串%s 索引[%d:%d] 截取后 %s\n", s, start_idx, end_idx, s[start_idx:end_idx])
	fmt.Printf("End: =================StringIterRun()=========================\n")
}

// 字符串数字互转 以及split join
func StrConvRun() {
	fmt.Printf("Begin: =================StrConvRun()=========================\n")

	// 字符串转换为数字
	conved_int, _ := strconv.Atoi("4599")
	fmt.Printf("conved int=%d\n", conved_int)

	// 数字转字符串
	conved_string := strconv.Itoa(4399)
	fmt.Printf("conved string=%s\n", conved_string)

	// split %v 输出所有对象
	raw_string := "aaa|bbb|ccc"
	splited_strings := strings.Split(raw_string, "|")
	fmt.Printf("splited_strings=%v\n", splited_strings)

	// join
	splited_strings = []string{"192", "168", "0", "1"}
	joined_string := strings.Join(splited_strings, ".")
	fmt.Printf("joined_string=%s\n", joined_string)

	// 查找子串
	total_str := "wanghongkai112 graduated from bupt and xdu"
	target_str := "bupt"
	contain_bool := strings.Contains(total_str, target_str)
	contain_idx := strings.Index(total_str, target_str)
	fmt.Printf("total_str=%s || target_str=%s contain_bool=%t, containe_idx=%d\n", total_str, target_str, contain_bool, contain_idx)
	fmt.Printf("End: =================StrConvRun()=========================\n")
}
