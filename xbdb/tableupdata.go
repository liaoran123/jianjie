//小白数据库
//表信息
package xbdb

//修改整条记录。等于先删除后添加。
func (t *Table) Updata(vals [][]byte) (r ReInfo) {
	r = t.Delete(vals[0])
	if !r.Succ {
		return
	}
	r = t.Insert(vals)
	return
}

//修改某个某些字段数据。不需要修改的字段索引不会删除和重复添加，提高性能。
func (t *Table) Upd(params map[string]string) (r ReInfo) {
	var updatefield []bool
	vals := t.StrToByte(params)
	for _, v := range vals {
		if len(v) == 0 {
			updatefield = append(updatefield, false)
		} else {
			updatefield = append(updatefield, true)
		}
	}
	key := JoinBytes(t.Select.GetTbKey(), vals[0])
	data, err := t.Db.Get(key, nil) //获取旧数据
	if err != nil {
		r.Info = err.Error()
		return
	}
	//组织新数据
	var newvals [][]byte
	newvals = append(newvals, vals[0])
	newvals = append(newvals, SplitRd(data)...)
	r = t.Acts(newvals, "delete", updatefield) //删除旧数据
	if !r.Succ {
		return
	}
	for i, v := range vals {
		if len(v) != 0 { //即是要修改的字段
			newvals[i] = v //更改要更改的字段值
		}
	}
	r = t.Acts(newvals, "insert", updatefield) //添加新数据
	return
}
