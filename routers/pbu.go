package routers

import (
	"crypto/md5"
	"encoding/hex"
	"sync"
)

var (
	ConfigMap map[string]interface{} //配置文件
	mu        sync.RWMutex
	/*
		bufpool   = sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		}*/
)

// 返回一个32位md5加密后的字符串
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

/*
func GetUrlData(url string) (r []byte) {
	//url := "http://127.0.0.1:9007/hello?age=20&id=1&name=lisi"
	req, _ := http.Get(url)
	data, _ := ioutil.ReadAll(req.Body)
	fmt.Printf("%v", string(data))
	return data
}
func GetData(wurl string) (r []byte) {
	req := ConfigMap["sev"].(string) + wurl
	var res *http.Response
	if strings.Contains(req, " ") {
		u, _ := url.Parse(req)
		q := u.Query()
		u.RawQuery = q.Encode() //urlencode//url需要转义
		res, _ = http.Get(u.String())
	} else {
		res, _ = http.Get(req) //有些 不需要转义
	}
	//req, _ := http.Get(rurl)
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		data = []byte(err.Error())
	}
	defer res.Body.Close()
	return data
}

*/
