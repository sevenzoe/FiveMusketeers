package controllers

import (
	"../helper"
	"../models"
	//"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/astaxie/beego/session/redis"
	"github.com/beego/i18n"
	"github.com/dchest/captcha"
	"html/template"
	"regexp"
	"strings"
)

const (
	HTTP_STATUS_CODE_302 = 302

	CACHE_ENGINE_FILE     = "file"
	CACHE_ENGINE_REDIS    = "redis"
	CACHE_ENGINE_MEMORY   = "memory"
	CACHE_ENGINE_MEMCACHE = "memcache"

	CACHE_KEY_LOGIN_NUM           = "login_counter"
	CACHE_KEY_SIGNUP_NUM          = "signup_counter"
	CACHE_KEY_FIND_PASSWORD_NUM   = "findpassword_counter"
	CACHE_KEY_FIND_PASSWORD_TOKEN = "findpassword_token"

	SESSION_ENGINE_FILE     = "file"
	SESSION_ENGINE_REDIS    = "redis"
	SESSION_ENGINE_MYSQL    = "mysql"
	SESSION_ENGINE_MEMCACHE = "memcache"

	SESSION_KEY_UID             = "uid"
	SESSION_KEY_USERNAME        = "username"
	SESSION_KEY_AUTH_USERNAME   = "auth_username"
	SESSION_KEY_AUTH_EMAIL      = "auth_email"
	SESSION_KEY_RESET_USERNAME  = "reset_username"
	SESSION_KEY_RESET_EMAIL     = "reset_email"
	SESSION_KEY_ENCODE_EMAIL    = "encode_email"
	SESSION_KEY_ERROR_MESSAGE   = "error_message"
	SESSION_KEY_SUCCESS_MESSAGE = "success_message"

	VIEW_INPUT_FIELD_USERNAME         = "username"
	VIEW_INPUT_FIELD_PASSWORD         = "password"
	VIEW_INPUT_FIELD_EMAIL            = "email"
	VIEW_INPUT_FIELD_TOKEN            = "token"
	VIEW_INPUT_FIELD_OLD_PASSWORD     = "oldpassword"
	VIEW_INPUT_FIELD_NEW_PASSWORD     = "newpassword"
	VIEW_INPUT_FIELD_CONFIRM_PASSWORD = "confirmpassword"
	VIEW_INPUT_FIELD_INTRODUCTION     = "introduction"
	VIEW_INPUT_FIELD_OCCUPATION       = "occupation"
	VIEW_INPUT_FIELD_PHONENUMBER      = "phoneNumber"
	VIEW_INPUT_FIELD_UNIVERSITY       = "university"
)

type langType struct {
	Lang, Name string
}

type baseController struct {
	beego.Controller
	i18n.Locale
	Uid                int
	User               *models.User
	IsLogin            bool
	IsAuthFindPassword bool
	EnableSignup       bool
	Username           string
}

type baseAdminController struct {
	baseController
}

type baseUserController struct {
	baseController
}

type baseAPIController struct {
	baseController
}

type ResultInfo struct {
	Code    int
	Message string
	Data    interface{}
}

var (
	langTypes []*langType
)

func init() {
	// 设置自定义存储验证码,替换默认的内存存储.必须调用此函数之前生成验证码.
	captcha.SetCustomStore(NewCaptchaStore(models.RedisAddress, 0, 0, -1))

	initSession()

	err := initCache()
	if err != nil {
		beego.Error("init cache error :", err)
	}

	err = initLanguage()
	if err != nil {
		beego.Error("init language error :", err)
	}
}

// 初始化Session
func initSession() {
	// 开启session,配置文件已经开启
	// 过期时间在配置文件中
	// beego.SessionOn = true

	// Session 第三方引擎
	if models.RedisOn {
		beego.SessionProvider = SESSION_ENGINE_REDIS // session 引擎
		beego.SessionSavePath = models.RedisAddress  // 引擎保存路径
	}
}

// 初始化Cache
func initCache() error {
	var err error

	// Cache 第三方引擎
	helper.GlobalCache, err = cache.NewCache(CACHE_ENGINE_REDIS, "{\"conn\": \""+models.RedisAddress+"\"}")
	if err != nil {
		return err
	}

	return nil
}

// 注册本地化文件
func initLanguage() error {
	langs := strings.Split(beego.AppConfig.String("lang::types"), "|")
	names := strings.Split(beego.AppConfig.String("lang::names"), "|")

	langTypes = make([]*langType, 0, len(langs))

	for i, v := range langs {
		langTypes = append(langTypes, &langType{
			Name: names[i],
			Lang: v,
		})
	}

	for _, lang := range langs {
		err := i18n.SetMessage(lang, "languages/"+lang+".ini")
		if err != nil {
			return err
		}
	}

	return nil
}

// 准备函数
func (this *baseController) Prepare() {

	this.TplExt = "html"

	this.setLangVer()

	this.IsLogin = this.CheckLogin()
	this.EnableSignup = true

	//this.GetNavigationBarData()
}

// 初始化控制器语言
func (this *baseController) setLangVer() bool {
	isNeedRedir := false
	hasCookie := false

	// 1. Check URL arguments.
	lang := this.Input().Get("lang")

	// 2. Get language information from cookies.
	if len(lang) == 0 {
		lang = this.Ctx.GetCookie("lang")
		hasCookie = true
	} else {
		isNeedRedir = true
	}

	// Check again in case someone modify by purpose.
	if !i18n.IsExist(lang) {
		lang = ""
		isNeedRedir = false
		hasCookie = false
	}

	// 3. Get language information from 'Accept-Language'.
	if len(lang) == 0 {
		al := this.Ctx.Request.Header.Get("Accept-Language")
		if len(al) > 4 {
			al = al[:5] // Only compare first 5 letters.
			if i18n.IsExist(al) {
				lang = al
			}
		}
	}

	// 4. Default language is English.
	if len(lang) == 0 {
		lang = "zh-CN"
		isNeedRedir = false
	}

	curLang := langType{
		Lang: lang,
	}

	// Save language information in cookies.
	if !hasCookie {
		this.Ctx.SetCookie("lang", curLang.Lang, 1<<31-1, "/")
	}

	restLangs := make([]*langType, 0, len(langTypes)-1)
	for _, v := range langTypes {
		if lang != v.Lang {
			restLangs = append(restLangs, v)
		} else {
			curLang.Name = v.Name
		}
	}

	for _, v := range restLangs {
		fmt.Println("****************", v.Name)
	}

	// Set language properties.
	this.Lang = lang
	this.Data["Lang"] = curLang.Lang
	this.Data["CurLang"] = curLang
	this.Data["RestLangs"] = restLangs

	return isNeedRedir
}

func (this *baseController) SetResultData(result *ResultInfo) {
	this.Data["json"] = result
	this.ServeJson()
}

func (this *baseController) GetNavigationBarData() {
	this.Data["NavList"] = models.NavBarMap[this.Lang]
	//this.Data["FooterList"]
}

// 只获取一次seesion,获取后删除
func (this *baseController) GetSessionOnce(key string) interface{} {
	session := this.GetSession(key)
	this.DelSession(key)
	return session
}

// 用户需要登录
func (this *baseController) UserNeedLogin() {
	if !this.IsLogin {
		this.RedirectAndStop(helper.UrlFor("User.Login"), HTTP_STATUS_CODE_302)
	}
}

// View层必须输入的字段.
func (this *baseController) NeedInputField(refer string, fields ...string) {
	for _, v := range fields {
		if len(this.GetString(v)) == 0 {
			switch v {
			case VIEW_INPUT_FIELD_USERNAME:
				{
					this.SetSession(SESSION_KEY_ERROR_MESSAGE, this.L("ErrorMessage.enter_username"))
				}
			case VIEW_INPUT_FIELD_PASSWORD:
				{
					this.SetSession(SESSION_KEY_ERROR_MESSAGE, this.L("ErrorMessage.enter_password"))
				}
			case VIEW_INPUT_FIELD_EMAIL:
				{
					this.SetSession(SESSION_KEY_ERROR_MESSAGE, this.L("ErrorMessage.enter_email"))
				}
			}

			this.RedirectAndStop(refer, HTTP_STATUS_CODE_302)
		}
	}
}

func (this *baseController) CheckInputValid(fields ...string) bool {
	for _, v := range fields {
		switch v {
		case VIEW_INPUT_FIELD_USERNAME:
			{
				return true
			}
		}
	}

	return false
}

// 检查View层输入字段的有效性(Username)
func (this *baseController) CheckInputValidByUsername(username string) bool {
	// 用户名,只允许输入数字,字母,下划线
	reg := regexp.MustCompile("^[0-9a-zA-Z][0-9a-zA-Z-]{1,}[0-9a-zA-Z]$")
	return reg.MatchString(username)
}

// 检查View层输入字段的有效性(Email)
func (this *baseController) CheckInputValidByEmail(email string) bool {
	// 邮箱
	reg := regexp.MustCompile("[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[\\w](?:[\\w-]*[\\w])?")
	return reg.MatchString(email)
}

// 获取语言文件里面相对应的value.
func (this *baseController) L(key string) string {
	return this.Tr(key)
}

// 重定向并且停止
func (this *baseController) RedirectAndStop(redirectUrl string, code int) {
	this.Redirect(redirectUrl, code)
	this.StopRun()
}

// 登陆
func (this *baseController) SetLogin(username, password string) (result *ResultInfo) {
	var err error
	user := &models.User{}
	result = &ResultInfo{Code: 0, Message: "OK", Data: nil}

	if len(username) == 0 || len(password) == 0 {
		result.Code = 200
		result.Message = "username or password is nil"
		result.Data = nil
		return result
	}

	user, err = models.GetUserByUsername(username)
	if err != nil {
		// TODO: 数据库操作出错
		result.Code = 100
		result.Message = err.Error()
		result.Data = nil
		return result
	}

	if user == nil {
		result.Code = 201
		result.Message = "user not exist"
		result.Data = nil
		return result
	}

	if user.Password == helper.MD5(password) {
		this.SetSession(SESSION_KEY_UID, user.Uid)
		this.SetSession(SESSION_KEY_USERNAME, user.Username)

		this.Uid = user.Uid
		this.Username = user.Username
		this.User = user

		result.Code = 0
		result.Message = "OK"
		result.Data = user

		return result
	} else {
		result.Code = 202
		result.Message = "password error"
		result.Data = nil
		return result
	}

	return result
}

// 退出登录
func (this *baseController) SetLogout() {
	this.DelSession(SESSION_KEY_UID)
	this.DelSession(SESSION_KEY_USERNAME)
}

// 检查是否登陆
func (this *baseController) CheckLogin() bool {
	var uid int = 0

	uidTemp := this.GetSession(SESSION_KEY_UID)
	if uidTemp != nil {
		uid = uidTemp.(int)
	}

	if uid <= 0 {
		return false
	}

	user, err := models.GetUserByUid(uid)
	if err != nil {
		// TODO:
	}

	this.User = user

	return true
}

// 防止CSRF攻击
func (this *baseController) SetXSRF() {
	this.Data["xsrfdata"] = template.HTML(this.XsrfFormHtml())
}

// 获取真实的客户端IP
func (this *baseController) GetRealIP() string {
	ip := this.Ctx.Input.Header("X-Real-IP")
	if len(ip) == 0 {
		ip = this.Ctx.Input.IP()
	}

	return ip
}
