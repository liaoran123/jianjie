package routers

import (
	"net/http"

	captcha "github.com/mojocn/base64Captcha"
)

var store = captcha.DefaultMemStore

func NewDriver() *captcha.DriverString {
	driver := new(captcha.DriverString)
	driver.Height = 49
	driver.Width = 108
	driver.NoiseCount = 5
	driver.ShowLineOptions = 2 //captcha.OptionShowSineLine | captcha.OptionShowSlimeLine | captcha.OptionShowHollowLine
	driver.Length = 4
	driver.Source = "1234567890qwertyuipkjhgfdsazxcvbnm"
	driver.Fonts = []string{"wqy-microhei.ttc"}
	return driver
}

// 生成图形验证码
func GenerateCaptchaHandler(w http.ResponseWriter, r *http.Request) {
	var driver = NewDriver().ConvertFonts()
	c := captcha.NewCaptcha(driver, store)
	_, content, answer := c.Driver.GenerateIdQuestionAnswer()
	id := r.URL.Query().Get("id") //"captcha:yufei"
	item, _ := c.Driver.DrawCaptcha(content)
	c.Store.Set(id, answer)
	item.WriteTo(w)
}

/*
// 验证
func CaptchaVerifyHandle(w http.ResponseWriter, r *http.Request) {

	id := "captcha:yufei"
	code := r.FormValue("code")
	body := map[string]interface{}{"code": 1000, "msg": "failed"}
	if store.Verify(id, code, true) {
		body = map[string]interface{}{"code": 1001, "msg": "ok"}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(body)
}
*/
