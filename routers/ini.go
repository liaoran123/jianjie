package routers

import (
	"fmt"
	"log"
	"jianjie/xbdb"

	"github.com/syndtr/goleveldb/leveldb"
)

var Xb *leveldb.DB
var Table map[string]*xbdb.Table

func Ini() {
	//打开或创建数据库
	dbpath := ConfigMap["dbpath"].(string)
	xb, err := leveldb.OpenFile(dbpath+"db", nil)
	if err != nil {
		log.Fatal(err)
	}
	//建表
	Xb = xb
	dbinfo := xbdb.NewTableInfo(Xb)
	if dbinfo.GetInfo("user").FieldType == nil {
		createuser(dbinfo)
	}
	if dbinfo.GetInfo("jianjie").FieldType == nil {
		createjianjie(dbinfo)
	}
	if dbinfo.GetInfo("dzan").FieldType == nil {
		createdzan(dbinfo)
	}
	//打开表操作结构
	Table = make(map[string]*xbdb.Table)
	Table["user"] = xbdb.NewTable(Xb, "user")
	Table["jianjie"] = xbdb.NewTable(Xb, "jianjie")
	Table["dzan"] = xbdb.NewTable(Xb, "dzan")
	//目录入加载内存

	//Table["jianjie"].Select.For(Pr)
}

func Pr(rd []byte) bool {
	fmt.Println(string(rd))
	return true
}

//创建用户表
func createuser(tbifo *xbdb.TableInfo) {
	name := "user"                                                                 //目录表
	fields := []string{"id", "email", "psw", "fahao", "jianjie", "sj"}             //字段，编码，邮箱，密码，法号。简介，注册时间。
	fieldType := []string{"int", "string", "string", "string", "string", "string"} //字段
	idxs := []string{"1", "3"}                                                     //索引字段,fields的下标对应的字段。支持组合查询，字段之间用,分隔
	fullText := []string{}                                                         //考据级全文搜索索引字段的下标。
	ftlen := "7"                                                                   //全文搜索的长度，中文默认是7
	r := tbifo.Create(name, ftlen, fields, fieldType, idxs, fullText)
	fmt.Printf("r: %v\n", r)
}

//创建见解表
func createjianjie(tbifo *xbdb.TableInfo) {
	name := "jianjie"                                                                        //目录表
	fields := []string{"id", "userid", "fahao", "secid", "sectext", "text", "sj"}            //字段，编码，对应的用户id编码，经文，内容，发布时间。将jingwen记录下来，以免将来改变大藏经结构导致数据不一致。
	fieldType := []string{"int", "string", "string", "string", "string", "string", "string"} //字段
	idxs := []string{"1", "3"}                                                               ////索引字段,fields的下标对应的字段。支持组合查询，字段之间用,分隔
	fullText := []string{}                                                                   //考据级全文搜索索引字段的下标。
	ftlen := "7"                                                                             //全文搜索的长度，中文默认是7
	r := tbifo.Create(name, ftlen, fields, fieldType, idxs, fullText)
	fmt.Printf("r: %v\n", r)
}

//点赞数表
func createdzan(tbifo *xbdb.TableInfo) {
	name := "dzan"                      //目录表，
	fields := []string{"id", "shu"}     //字段 见解表的编码，点赞数
	fieldType := []string{"int", "int"} //字段
	idxs := []string{"1"}               //索引字段,fields的下标对应的字段。支持组合查询，字段之间用,分隔
	fullText := []string{}              //考据级全文搜索索引字段的下标。
	ftlen := "7"                        //全文搜索的长度，中文默认是7
	r := tbifo.Create(name, ftlen, fields, fieldType, idxs, fullText)
	fmt.Printf("r: %v\n", r)
}
