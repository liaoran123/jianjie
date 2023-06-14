package routers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

// 备份数据库
func Dbback(w http.ResponseWriter, req *http.Request) {
	back := newbackdb()
	Table["u"].Select.ForDbase(back.do)
	back.xb.Close()
}

type backdb struct {
	xb  *leveldb.DB
	err error
}

func newbackdb() *backdb {
	dbpath := ConfigMap["dbpath"].(string)
	xb, err := leveldb.OpenFile(dbpath+"db"+time.Now().Format("20060102"), nil)
	if err != nil {
		log.Fatal(err)
	}
	return &backdb{
		xb: xb,
	}
}
func (b *backdb) do(k, v []byte) bool {
	b.err = b.xb.Put(k, v, nil)
	if b.err != nil {
		fmt.Println(b.err.Error())
		return false
	}
	return true
}
