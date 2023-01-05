package routers

import (
	"jianjie/pubgo"
	"jianjie/xbdb"
	"net/http"
	"strconv"
)

//根据索引查询，返回json
func Pubget(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows("/pubget/" + req.Method)
	params := getparas(req)
	key := Table[params["tbname"]].Ifo.FieldChByte(params["idxfield"], params["idxvalue"])
	b, _ := strconv.Atoi(params["b"])
	count, _ := strconv.Atoi(params["count"])
	tbd := Table[params["tbname"]].Select.WhereIdx([]byte(params["idxfield"]), key, params["asc"] == "1", b, count)
	if tbd != nil {
		json := Table[params["tbname"]].DataToJson(tbd)
		w.Write(json.Bytes())
		json.Reset()
		xbdb.Bufpool.Put(json)
	} else {
		w.Write([]byte(""))
	}
}

//打开表，返回json
func Pubgettb(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows("/Pubgettb/" + req.Method)
	params := getparas(req)

	b, _ := strconv.Atoi(params["b"])
	count, _ := strconv.Atoi(params["count"])
	tbd := Table[params["tbname"]].Select.FindPrefix([]byte(params["tbname"]+xbdb.Split), params["asc"] == "1", b, count)
	//tbd := Table[params["tbname"]].Select.WhereIdx([]byte(params["idxfield"]), key, params["asc"] == "1", b, count)
	if tbd != nil {
		json := Table[params["tbname"]].DataToJson(tbd)
		w.Write(json.Bytes())
		json.Reset()
		xbdb.Bufpool.Put(json)
	} else {
		w.Write([]byte(""))
	}
}
