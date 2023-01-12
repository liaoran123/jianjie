package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"jianjie/pubgo"
	"jianjie/routers"
	"net/http"
)

func main() {
	//--读取初始化设置数据
	path := pubgo.GetCurrentAbPath() //守护程序读取文件时需要绝对路径
	text, _ := ioutil.ReadFile(path + "config.json")
	routers.ConfigMap = make(map[string]interface{})
	json.Unmarshal(text, &routers.ConfigMap)
	port := routers.ConfigMap["port"].(string) //从配置文件获取port
	routers.Ini()

	pubgo.Tj = pubgo.Newtongji()

	http.HandleFunc("/static/", routers.Static) //静态文件服务器
	http.HandleFunc("/pubtb/", routers.Pubtb)
	http.HandleFunc("/pubget/", routers.Pubget)
	http.HandleFunc("/pubgettb/", routers.Pubgettb)
	http.HandleFunc("/user/", routers.User)
	http.HandleFunc("/jianjie/", routers.FjJianjie)

	http.HandleFunc("/Captcha/", routers.GenerateCaptchaHandler)
	http.HandleFunc("/redir/", routers.Redir)
	//http.HandleFunc("/updb/", routers.Updb)
	http.HandleFunc("/test/", routers.Test)

	http.HandleFunc("/dbback/", routers.Dbback)

	fmt.Println("佛经见解1.0版本")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println(err)
		//log.Fatal(err)
	}
}
