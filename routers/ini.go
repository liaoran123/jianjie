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
	//创建数据库连接
	Xb = xb
	//创建表信息结构
	dbinfo := xbdb.NewTableInfo(Xb)
	//删除后添加=修改
	dbinfo.Del("u")
	dbinfo.Del("j")
	dbinfo.Del("d")
	dbinfo.Del("qz")
	dbinfo.Del("wz")
	dbinfo.Del("hf")
	dbinfo.Del("sc")
	dbinfo.Del("admin")
	dbinfo.Del("test")

	//创建表
	createtbs(dbinfo)
	//打开表操作结构
	Table = make(map[string]*xbdb.Table)
	Table["u"] = xbdb.NewTable(Xb, "u")
	Table["j"] = xbdb.NewTable(Xb, "j")
	Table["d"] = xbdb.NewTable(Xb, "d")
	Table["qz"] = xbdb.NewTable(Xb, "qz")
	Table["wz"] = xbdb.NewTable(Xb, "wz")
	Table["hf"] = xbdb.NewTable(Xb, "hf")
	Table["sc"] = xbdb.NewTable(Xb, "sc")
	Table["admin"] = xbdb.NewTable(Xb, "admin")

	//测试代码
	Table["test"] = xbdb.NewTable(Xb, "test")
	/*
		params := map[string]string{"zd0": "12", "zd1": "阿弥陀佛"}
		r := Table["test"].Ins(params)
		fmt.Printf("r: %v\n", r)

			params1 := map[string]string{"id": "1", "zd1": "众生颠倒"} //, "zd1": "众生颠倒",, "zd0": "1000"
			r1 := Table["test"].Upd(params1)
			fmt.Printf("r1: %v\n", r1)
	*/
	//打印数据库//用于测试代码
	//Table["u"].Select.For(Pr0)
	Table["test"].Select.ForTb(Pr)

}

func Pr(k, v []byte) bool {
	fmt.Println(string(k), string(v))
	return true
}

/*
func Pr0(k []byte) bool {
	fmt.Println(string(k))
	return true
	rdmap := Table["u"].RDtoMap(k)
	if rdmap["sj"] == "" {
		params1 := map[string]string{"id": rdmap["id"], "sj": "2023-01-10 17:12:52"}
		r1 := Table["u"].Upd(params1)
		fmt.Printf("r1: %v\n", r1)
	} else {
		if strings.Contains(rdmap["sj"], "CST") {
			sj := rdmap["sj"]
			sj = strings.Split(sj, ".")[0]
			params1 := map[string]string{"id": rdmap["id"], "sj": sj}
			r1 := Table["u"].Upd(params1)
			fmt.Printf("r1: %v\n", r1)
		}
	}

	return true
}*/
func createtbs(dbinfo *xbdb.TableInfo) {
	if dbinfo.GetInfo("u").FieldType == nil { //创建用户表
		name := "u"                                                                  //user                                                              //目录表
		fields := []string{"id", "email", "psw", "fahao", "jianjie", "sj"}           //字段，编码，邮箱，密码，法号。简介，通过（默认。0，则不通过。删除用户设置为0即可）,注册时间。
		fieldType := []string{"int", "string", "string", "string", "string", "time"} //字段类型
		idxs := []string{"1", "3"}                                                   //索引字段,fields的下标对应的字段。支持组合查询，字段之间用,分隔
		fullText := []string{}                                                       //考据级全文搜索索引字段的下标。
		ftlen := "7"                                                                 //全文搜索的长度，中文默认是7
		r := dbinfo.Create(name, ftlen, fields, fields, fieldType, idxs, fullText)
		fmt.Printf("r: %v\n", r)

	}
	if dbinfo.GetInfo("j").FieldType == nil { //创建见解表
		name := "j"                                                                            //jianjie                                                                   //目录表
		fields := []string{"id", "userid", "fahao", "secid", "sectext", "text", "sj"}          //字段，编码，对应的用户id编码，经文，内容，发布时间。将jingwen记录下来，以免将来改变大藏经结构导致数据不一致。
		fieldType := []string{"int", "string", "string", "string", "string", "string", "time"} //字段类型
		idxs := []string{"1", "3"}                                                             ////索引字段,fields的下标对应的字段。支持组合查询，字段之间用,分隔
		fullText := []string{}                                                                 //考据级全文搜索索引字段的下标。
		ftlen := "7"                                                                           //全文搜索的长度，中文默认是7
		r := dbinfo.Create(name, ftlen, fields, fields, fieldType, idxs, fullText)
		fmt.Printf("r: %v\n", r)

	}
	if dbinfo.GetInfo("d").FieldType == nil { //点赞数表
		name := "d"                     //dzan               //点赞表，
		fields := []string{"id"}        //字段 id是见解表的编码和userid组成，每个用户的点赞都记录下来
		fieldType := []string{"string"} //字段类型
		idxs := []string{}              //索引字段,fields的下标对应的字段。支持组合查询，字段之间用,分隔
		fullText := []string{}          //考据级全文搜索索引字段的下标。
		ftlen := "7"                    //全文搜索的长度，中文默认是7
		r := dbinfo.Create(name, ftlen, fields, fields, fieldType, idxs, fullText)
		fmt.Printf("r: %v\n", r)

	}
	if dbinfo.GetInfo("qz").FieldType == nil { //创建群组表
		name := "qz"                                                                      //目录表
		fields := []string{"id", "mc", "userid", "fahao", "jianjie", "zlan", "sj"}        //字段，编码，名称，用户id，简介，是否专栏（1，是；0，否。），创建时间。
		fieldType := []string{"int", "string", "int", "string", "string", "bool", "time"} //字段类型
		idxs := []string{"1", "2"}                                                        //索引字段,fields的下标对应的字段。支持组合查询，字段之间用,分隔
		fullText := []string{"1"}                                                         //考据级全文搜索索引字段的下标。
		ftlen := "7"                                                                      //全文搜索的长度，中文默认是7
		r := dbinfo.Create(name, ftlen, fields, fields, fieldType, idxs, fullText)
		fmt.Printf("r: %v\n", r)

	}
	if dbinfo.GetInfo("wz").FieldType == nil { //创建文章表
		name := "wz"                                                                                         //目录表
		fields := []string{"id", "type", "typename", "userid", "fahao", "title", "cont", "pass", "pl", "sj"} //字段，编码，type,类型，类型名称,用户编码，标题，内容，是否通过审核,可评论（0，关闭pl）,创建时间。
		//type,类型。当为群组时，即是群组编码；当是见解时，即是见解secid。type是string，兼容int或字符串的id。
		//当是新闻资讯时，type="xw"等等。除见解外特殊外，所有文章都用这个表，以type区分。
		//加入类型名称，是为了空间换时间。
		//pass,1是通过。
		fieldType := []string{"int", "string", "string", "int", "string", "string", "string", "bool", "bool", "time"} //字段类型
		idxs := []string{"1", "3"}                                                                                    //索引字段,fields的下标对应的字段。支持组合查询，字段之间用,分隔。仅提供title搜索。
		fullText := []string{"5"}                                                                                     //考据级全文搜索索引字段的下标。
		ftlen := "7"                                                                                                  //全文搜索的长度，中文默认是7
		r := dbinfo.Create(name, ftlen, fields, fields, fieldType, idxs, fullText)
		fmt.Printf("r: %v\n", r)

	}
	if dbinfo.GetInfo("hf").FieldType == nil { //创建回复表，包括见解、和群组文章或其他
		name := "hf"                                                                       //目录表
		fields := []string{"id", "wid", "userid", "fahao", "wtitle", "cont", "pass", "sj"} //字段，编码，文章编码,用户编码,文章标题，内容，创建时间。
		//wid为string，兼容int和string的id
		fieldType := []string{"int", "string", "int", "string", "string", "string", "string", "time"} //字段类型
		idxs := []string{"1", "2"}                                                                    //索引字段,fields的下标对应的字段。支持组合查询，字段之间用,分隔。组合查询，就是为了避开 where ...and...的情况，直接用组合索引代替解决。
		fullText := []string{}                                                                        //考据级全文搜索索引字段的下标。
		ftlen := "7"                                                                                  //全文搜索的长度，中文默认是7
		r := dbinfo.Create(name, ftlen, fields, fields, fieldType, idxs, fullText)
		fmt.Printf("r: %v\n", r)

	}
	/*
		if dbinfo.GetInfo("gbhf").FieldType == nil { //创建创建关闭回复表
			name := "gbhf"
			//如果userid不为空，则是该用户全部文章关闭评论。
			fields := []string{"id", "wid"}        //字段，编码，文章编码（1,见解；2，群组文章），所属（1,见解；2，群组文章），创建时间。
			fieldType := []string{"int", "string"} //字段类型
			idxs := []string{"1"}                  //索引字段,fields的下标对应的字段。支持组合查询，字段之间用,分隔。组合查询，就是为了避开 where ...and...的情况，直接用组合索引代替解决。
			fullText := []string{}                 //考据级全文搜索索引字段的下标。
			ftlen := "7"                           //全文搜索的长度，中文默认是7
			r := dbinfo.Create(name, ftlen, fields,fields, fieldType, idxs, fullText)
			fmt.Printf("r: %v\n", r)

		}*/
	if dbinfo.GetInfo("sc").FieldType == nil { //创建收藏表
		name := "sc"                                                              //收藏表
		fields := []string{"id", "userid", "url", "title", "cata", "sj"}          //字段，编码，用户编码，地址，标题,分类，时间。
		fieldType := []string{"int", "int", "string", "string", "string", "time"} //字段类型
		idxs := []string{"1"}                                                     //索引字段,fields的下标对应的字段。支持组合查询，字段之间用,分隔。组合查询，就是为了避开 where ...and...的情况，直接用组合索引代替解决。
		fullText := []string{}                                                    //考据级全文搜索索引字段的下标。
		ftlen := "7"                                                              //全文搜索的长度，中文默认是7
		r := dbinfo.Create(name, ftlen, fields, fields, fieldType, idxs, fullText)
		fmt.Printf("r: %v\n", r)
	}
	if dbinfo.GetInfo("admin").FieldType == nil { //管理员表
		name := "admin"
		//type,类型。当为群组时，即是群组编码；type是string，兼容int或字符串的id。
		//当是新闻资讯时，type="xw", 当是见解时，type="jj" 等等。                                  //收藏表
		fields := []string{"id", "userid", "fahao", "type", "pass", "sj"}       //字段，编码，用户编码，类型，是否通过，时间。
		fieldType := []string{"int", "int", "string", "string", "bool", "time"} //字段类型
		idxs := []string{"1"}                                                   //索引字段,fields的下标对应的字段。支持组合查询，字段之间用,分隔。组合查询，就是为了避开 where ...and...的情况，直接用组合索引代替解决。
		fullText := []string{}                                                  //考据级全文搜索索引字段的下标。
		ftlen := "7"                                                            //全文搜索的长度，中文默认是7
		r := dbinfo.Create(name, ftlen, fields, fields, fieldType, idxs, fullText)
		fmt.Printf("r: %v\n", r)
	}
	if dbinfo.GetInfo("test").FieldType == nil { //创建专用测试表
		name := "test"                                   //收藏表
		fields := []string{"id", "zd0", "zd1"}           //字段，编码，用户编码，地址，标题，时间。
		fieldType := []string{"int", "string", "string"} //字段类型
		idxs := []string{"1"}                            //索引字段,fields的下标对应的字段。支持组合查询，字段之间用,分隔。组合查询，就是为了避开 where ...and...的情况，直接用组合索引代替解决。
		fullText := []string{"2"}                        //考据级全文搜索索引字段的下标。
		ftlen := "7"                                     //全文搜索的长度，中文默认是7
		r := dbinfo.Create(name, ftlen, fields, fields, fieldType, idxs, fullText)
		fmt.Printf("r: %v\n", r)
	}
}
