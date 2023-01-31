package main

import (
	"fmt"
	"testing"
)

func TestA(t *testing.T) { //不能使用Testxbdb 类似名称，邪门
	//rb := JoinBytes([]byte("ddd"), []byte("fff"))
	b := []byte{}
	blen := len(b)
	fmt.Printf("blen: %v\n", blen)
	if b == nil {
		fmt.Println("nil")
	} else {
		fmt.Printf("b: %v\n", b)
	}
}
