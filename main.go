package main

import (
	"./conf"
	"./helper"
	"./routers"
	"encoding/json"
	"flag"
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
	"io/ioutil"
)

func main() {
	err := conf.InitConfig()()
	if err != nil {
		panic(err)
	}

	routers.Init()
	beego.Run()
}

func templateFunc() {
	beego.AddFuncMap("userspaceurl", helper.UserSpaceUrl) // {{userurl .input}} -> input/string 生成新的url
	beego.AddFuncMap("urlfor", helper.UrlFor)             // {{urlfor "User.Signup"}} -> user/signup?... 提交表单
	beego.AddFuncMap("i18n", i18n.Tr)                     // {{i18n .Lang "default.home"}} -> 在加载的语言配置文件中,动态选择语言.
}
