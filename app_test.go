package main

import (
	"jianjie/xbdb"
	"testing"
)

func TestA(t *testing.T) { //不能使用Testxbdb 类似名称，邪门
	xbdb.OpenDb("F:/dababase/fojingjianjie/")
	xbdb.OpenTableStructs()
}
