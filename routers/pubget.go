package routers

import (
	"jianjie/pubgo"
	"jianjie/xbdb"
	"net/http"
	"strconv"
)

//根据索引查询，返回json
func PubIDXGet(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows("/PubIDXGet/" + req.Method)
	params := getparas(req)
	idxfield := params["idxfield"]
	idxvalue := params[idxfield]
	key := Table[params["tbname"]].Ifo.FieldChByte(idxfield, idxvalue)
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

//根据打开表的Top，顺序或倒序，返回json
func PubGetTB(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows("/PubGetTB/" + req.Method)
	params := getparas(req)

	b, _ := strconv.Atoi(params["b"])
	count, _ := strconv.Atoi(params["count"])
	tbname := params["tbname"]
	tbd := Table[tbname].Select.FindPrefix([]byte(tbname+xbdb.Split), params["asc"] == "1", b, count)
	//tbd := Table[params["tbname"]].Select.WhereIdx([]byte(params["idxfield"]), key, params["asc"] == "1", b, count)
	if tbd != nil {
		json := Table[tbname].DataToJson(tbd)
		w.Write(json.Bytes())
		json.Reset()
		xbdb.Bufpool.Put(json)
	} else {
		w.Write([]byte(""))
	}
}

//根据打开表的Top，顺序或倒序，返回json
func PubGetTBOne(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*") //同源策略，不加客户端调用不了。
	w.Header().Set("Content-Type", "application/json")

	pubgo.Tj.Brows("/PubGetTBOne/" + req.Method)
	params := getparas(req)

	tbname := params["tbname"]
	id := params["id"]
	bid := Table[tbname].Ifo.FieldChByte(Table[tbname].Ifo.Fields[0], id) //默认第一个必须是主键。Table[tbname].Ifo.Fields[0]
	tbd := Table[tbname].Select.Record(bid)
	if tbd != nil {
		json := Table[tbname].DataToJson(tbd)
		w.Write(json.Bytes())
		json.Reset()
		xbdb.Bufpool.Put(json)
	} else {
		w.Write([]byte(""))
	}
}
