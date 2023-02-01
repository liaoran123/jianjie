package routers

import (
	"encoding/json"
	"jianjie/pubgo"
	"jianjie/xbdb"
	"net/http"
	"strconv"
)

var (
	qzmethod map[string]func(w http.ResponseWriter, req *http.Request) //查询添加修改删除操作处理。

)

//群组
func Qz(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows("/qz/" + req.Method)

	//req.qzmethod
	if qzmethod == nil {
		qzmethod = make(map[string]func(w http.ResponseWriter, req *http.Request), 4)
		qzmethod["POST"] = qzpost //添加
		//qzmethod["GET"] = qzget     //查询
		qzmethod["DELETE"] = qzpost //删除
		qzmethod["PUT"] = qzpost    //qzput       //修改
	}
	if f, ok := qzmethod[req.Method]; ok {
		f(w, req)
	}
}
func qzpost(w http.ResponseWriter, req *http.Request) {
	var r xbdb.ReInfo
	params := postparas(req) //  postparas(req)
	if !store.Verify(params["capid"], params["code"], true) {
		r.Info = "验证码不正确！"
		json.NewEncoder(w).Encode(r)
		return
	}
	tbname := "qz"
	//存在id参数，即是删除。（浏览器不支持delete提交。）
	if id, ok := params["id"]; ok {
		r = delqz(id, params["userid"])
		json.NewEncoder(w).Encode(r)
		return
	}
	buserid := Table[tbname].Ifo.FieldChByte("userid", params["userid"])
	if params["pass"] == "" {
		params["pass"] = "1" //默认是通过。
	}
	tbcount := Table[tbname].Select.WhereIdxCount([]byte("userid"), buserid)
	if tbcount < 3 { //最多能建立3个群组
		if params["sj"] == "" {
			params["sj"] = "now()"
		}
		r = Table[tbname].Ins(params)               //添加群组
		lastid := Table[tbname].Ac.GetidNoInc() - 1 //根据自动增值-1得到最后的一条记录的id值
		aparams := make(map[string]string)
		aparams["userid"] = params["userid"]
		aparams["fahao"] = params["fahao"]
		aparams["type"] = strconv.Itoa(lastid)
		aparams["pass"] = "1"
		Table["admin"].Ins(aparams) //添加群组创建者为管理员

	} else {
		r.Info = "最多能建立3个群组"
	}
	json.NewEncoder(w).Encode(r)
}

//删除群组
func delqz(id, userid string) (r xbdb.ReInfo) {
	tbname := "qz"

	field := Table[tbname].Ifo.Fields[0]
	pkval := Table[tbname].Ifo.FieldChByte(field, id)
	/*
		tbd := Table["wz"].Select.WhereIdx([]byte("type"), pkval, true, 0, 1)
		if tbd != nil {
			r.Info = "群组存在文章，不能删除。"
			tbd.Release()
			return
		}*/
	ct := Table["wz"].Select.WhereIdxExist([]byte("type"), pkval)
	if ct {
		r.Info = "群组存在文章，不能删除。"
		return
	}
	tbd := Table[tbname].Select.Record(pkval)
	rdmap := Table[tbname].RDtoMap(tbd.Rd[0])
	tbd.Release()
	if rdmap["userid"] == userid { //删除的id和用户id对应才能删除，以防数据错乱和攻击。
		r = Table[tbname].Del(id)
	}
	tbd.Release()

	return
}

/*
//打开群组
func qzget(w http.ResponseWriter, req *http.Request) {
	params := getparas(req) //  postparas(req)
	tbname := "qz"
	buserid := Table[tbname].Ifo.FieldChByte("userid", params["userid"])
	tdb := Table[tbname].Select.WhereIdx([]byte("userid"), buserid, true, 0, -1)
	r := Table[tbname].DataToJson(tdb)
	w.Write(r.Bytes())
	r.Reset()
	xbdb.Bufpool.Put(r)
}


func qzdelete(w http.ResponseWriter, req *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	params := postparas(req)
	Table["qz"].Del(params["id"])
}
func qzput(w http.ResponseWriter, req *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	params := postparas(req)
	Table["qz"].Upd(params)
}*/
