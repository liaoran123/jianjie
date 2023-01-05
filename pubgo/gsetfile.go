package pubgo

//--读取配置文件
import (
	"io/ioutil"
	"os"
	"path/filepath"
)

//---公共参数

type setfile struct {
	filetext string
}

func Newsetfile(n string) setfile {
	set, _ := ioutil.ReadFile(n)
	return setfile{
		filetext: string(set[:]),
	}
}

//获取程序绝对路径目录
func GetCurrentAbPath() string {
	exePath, err := os.Executable()
	if err != nil {
		//log.Fatal(err)
		return ""
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res + "\\"

}
