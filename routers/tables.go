package routers

import (
	"jianjie/xbdb"
)

func InsOrUpd(name string, params map[string]string, insorupd string) (r xbdb.ReInfo) {
	iou := map[string]func(name string, v [][]byte) (r xbdb.ReInfo){
		"ins": ins,
		"upd": upd,
	}
	if insorupd == "" {
		insorupd = "ins"
	}
	/*
		var vals [][]byte
		f, v := "", ""
		for i := 0; i < len(Table[name].Ifo.FieldType); i++ {
			f = Table[name].Ifo.FieldType[i]
			v = Table[name].Ifo.Fields[i]
			vals = append(vals, Table[name].Ifo.TypeChByte(f, params[v]))
		}*/
	vals := Table[name].StrToByte(params)
	r = iou[insorupd](name, vals)
	return
}
func ins(name string, vals [][]byte) (r xbdb.ReInfo) {
	r = Table[name].Insert(vals)
	return
}
func upd(name string, vals [][]byte) (r xbdb.ReInfo) {
	r = Table[name].Updata(vals)
	return
}

/*
//删除一条表记录
func deleterd(tbname, id string) (r xbdb.ReInfo) {
	//bid := Table[tbname].Ifo.TypeChByte(Table[tbname].Ifo.FieldType[0], id)
	r = Table[tbname].Del(id)
	return
}
*/
