package main

import (
	"fmt"
	"testing"
	"unicode"
)

func TestUnicode(t *testing.T) {
	a := '的'
	unicodePoint := unicode.To(unicode.LowerCase, a) // unicode码点：30340
	println("unicode integer:", unicodePoint)
	fmt.Printf("unicode=>中文：%c\n", unicodePoint) // %c 可直接转为对应字符

	//IsSymbolTest()
	//IsSpaceTest()
	//IsPunctTest()
	//IsLetterTest()
	IsOneOfSet()
}

// IsSymbol 判断 r 是否为一个符号字符
func IsSymbolTest() {
	s := "Hello (<世=界>)"
	for _, r := range s {
		if unicode.IsSymbol(r) {
			fmt.Printf("%c = %v\n", r, unicode.IsSymbol(r))
		}
	} // <=> = true
}

// IsSpace 判断 r 是否为一个空白字符
// 在 Latin-1 字符集中，空白字符为：\t, \n, \v, \f, \r,
// 空格, U+0085 (NEL), U+00A0 (NBSP)
// 其它空白字符的定义有“类别 Z”和“Pattern_White_Space 属性”
func IsSpaceTest() {
	s := "Hello \t世　界！\n\u0000\u0020\u3000\u00A0"
	for i, r := range s {
		if unicode.IsSpace(r) {
			fmt.Printf("%d %c = %v\n", i, r, unicode.IsSpace(r))
		}
	}
}

// 判断标点符号
func IsPunctTest() {
	s := "Hello ,./;'\"[]=-+()世界！!?？" // +`= 不是标点符号
	for _, r := range s {
		if unicode.IsPunct(r) {
			fmt.Printf("%c IsPunct\n", r)
		}
	}
}

// 判断标点符号
func IsLetterTest() {
	s := "Hello ,.~！？=-+" // +`= 不是标点符号
	for _, r := range s {
		if unicode.IsLetter(r) {
			fmt.Printf("%c IsLetter\n", r) // Hello
		}
	}
}

//判断r是否在 set 表范围内
func IsOneOfSet() {
	s := "Hello 世界！"
	// set表设置为“汉字、标点符号, 空白字符”
	set := []*unicode.RangeTable{unicode.Han, unicode.Punct, unicode.Space}
	for _, r := range s {
		fmt.Printf("%c = %v\n", r, unicode.IsOneOf(set, r))
	} // 世界！ = true
}
