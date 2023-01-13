package routers

import (
	"fmt"
	"jianjie/xbdb"
	"log"

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
	Table = make(map[string]*xbdb.Table)
	/*
		//删除后添加=修改
		dbinfo.Del("u")
		dbinfo.Del("j")
		dbinfo.Del("d")
		dbinfo.Del("qz")
		dbinfo.Del("wz")
		dbinfo.Del("hf")
		dbinfo.Del("gbhf")
	*/
	createtbs(dbinfo)
	//打开表操作结构
	Table["u"] = xbdb.NewTable(Xb, "u")       //打开表操作结构
	Table["j"] = xbdb.NewTable(Xb, "j")       //打开表操作结构
	Table["d"] = xbdb.NewTable(Xb, "d")       //打开表操作结构
	Table["qz"] = xbdb.NewTable(Xb, "qz")     //打开表操作结构
	Table["wz"] = xbdb.NewTable(Xb, "wz")     //打开表操作结构
	Table["hf"] = xbdb.NewTable(Xb, "hf")     //打开表操作结构
	Table["gbhf"] = xbdb.NewTable(Xb, "gbhf") //打开表操作结构

	//目录入加载内存

	//Table["j"].Select.ForDb(Pr)
}

func Pr(rd []byte) bool {
	fmt.Println(string(rd))
	return true
}
func createtbs(dbinfo *xbdb.TableInfo) {
	if dbinfo.GetInfo("u").FieldType == nil { //创建用户表
		name := "u"                                                                  //user                                                              //目录表
		fields := []string{"id", "email", "psw", "fahao", "jianjie", "sj"}           //字段，编码，邮箱，密码，法号。简介，注册时间。
		fieldType := []string{"int", "string", "string", "string", "string", "time"} //字段类型
		idxs := []string{"1", "3"}                                                   //索引字段,fields的下标对应的字段。支持组合查询，字段之间用,分隔
		fullText := []string{}                                                       //考据级全文搜索索引字段的下标。
		ftlen := "7"                                                                 //全文搜索的长度，中文默认是7
		r := dbinfo.Create(name, ftlen, fields, fieldType, idxs, fullText)
		fmt.Printf("r: %v\n", r)

	}
	if dbinfo.GetInfo("j").FieldType == nil { //创建见解表
		name := "j"                                                                            //jianjie                                                                   //目录表
		fields := []string{"id", "userid", "fahao", "secid", "sectext", "text", "sj"}          //字段，编码，对应的用户id编码，经文，内容，发布时间。将jingwen记录下来，以免将来改变大藏经结构导致数据不一致。
		fieldType := []string{"int", "string", "string", "string", "string", "string", "time"} //字段类型
		idxs := []string{"1", "3"}                                                             ////索引字段,fields的下标对应的字段。支持组合查询，字段之间用,分隔
		fullText := []string{}                                                                 //考据级全文搜索索引字段的下标。
		ftlen := "7"                                                                           //全文搜索的长度，中文默认是7
		r := dbinfo.Create(name, ftlen, fields, fieldType, idxs, fullText)
		fmt.Printf("r: %v\n", r)

	}
	if dbinfo.GetInfo("d").FieldType == nil { //点赞数表
		name := "d"                     //dzan               //点赞表，
		fields := []string{"id"}        //字段 id是见解表的编码和userid组成，每个用户的点赞都记录下来
		fieldType := []string{"string"} //字段类型
		idxs := []string{}              //索引字段,fields的下标对应的字段。支持组合查询，字段之间用,分隔
		fullText := []string{}          //考据级全文搜索索引字段的下标。
		ftlen := "7"                    //全文搜索的长度，中文默认是7
		r := dbinfo.Create(name, ftlen, fields, fieldType, idxs, fullText)
		fmt.Printf("r: %v\n", r)

	}
	if dbinfo.GetInfo("qz").FieldType == nil { //创建群组表
		name := "qz"                                          //目录表
		fields := []string{"id", "mc", "userid", "sj"}        //字段，编码，名称，用户id，查看密码（空则不用），创建时间。
		fieldType := []string{"int", "string", "int", "time"} //字段类型
		idxs := []string{"2"}                                 //索引字段,fields的下标对应的字段。支持组合查询，字段之间用,分隔
		fullText := []string{"1"}                             //考据级全文搜索索引字段的下标。
		ftlen := "7"                                          //全文搜索的长度，中文默认是7
		r := dbinfo.Create(name, ftlen, fields, fieldType, idxs, fullText)
		fmt.Printf("r: %v\n", r)

	}
	if dbinfo.GetInfo("wz").FieldType == nil { //创建文章表
		name := "wz"                                                           //目录表
		fields := []string{"id", "qzid", "userid", "title", "cont", "sj"}      //字段，编码，群组编码，用户编码，标题，内容，创建时间。
		fieldType := []string{"int", "int", "int", "string", "string", "time"} //字段类型
		idxs := []string{"1", "2"}                                             //索引字段,fields的下标对应的字段。支持组合查询，字段之间用,分隔
		fullText := []string{"3"}                                              //考据级全文搜索索引字段的下标。
		ftlen := "7"                                                           //全文搜索的长度，中文默认是7
		r := dbinfo.Create(name, ftlen, fields, fieldType, idxs, fullText)
		fmt.Printf("r: %v\n", r)

	}
	if dbinfo.GetInfo("hf").FieldType == nil { //创建回复表，包括见解、和群组文章或其他
		name := "hf"                                                                     //目录表
		fields := []string{"id", "wid", "type", "userid", "title", "cont", "sj"}         //字段，编码，文章编码（1,见解；2，群组文章），所属（1,见解；2，群组文章），标题，内容，创建时间。
		fieldType := []string{"int", "int", "string", "int", "string", "string", "time"} //字段类型
		idxs := []string{"1,2", "3"}                                                     //索引字段,fields的下标对应的字段。支持组合查询，字段之间用,分隔。组合查询，就是为了避开 where ...and...的情况，直接用组合索引代替解决。
		fullText := []string{"4"}                                                        //考据级全文搜索索引字段的下标。
		ftlen := "7"                                                                     //全文搜索的长度，中文默认是7
		r := dbinfo.Create(name, ftlen, fields, fieldType, idxs, fullText)
		fmt.Printf("r: %v\n", r)

	}
	if dbinfo.GetInfo("gbhf").FieldType == nil { //创建创建关闭回复表
		name := "gbhf" //目录表
		//如果userid不为空，则是该用户全部文章关闭评论。
		fields := []string{"id", "wid", "type"}       //字段，编码，文章编码（1,见解；2，群组文章），所属（1,见解；2，群组文章），创建时间。
		fieldType := []string{"int", "int", "string"} //字段类型
		idxs := []string{"1,2"}                       //索引字段,fields的下标对应的字段。支持组合查询，字段之间用,分隔。组合查询，就是为了避开 where ...and...的情况，直接用组合索引代替解决。
		fullText := []string{}                        //考据级全文搜索索引字段的下标。
		ftlen := "7"                                  //全文搜索的长度，中文默认是7
		r := dbinfo.Create(name, ftlen, fields, fieldType, idxs, fullText)
		fmt.Printf("r: %v\n", r)

	}
}
