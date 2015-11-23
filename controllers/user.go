package controllers

import (
	"../helper"
	"../models"
	"fmt"
	"time"
)

type User struct {
	baseController
}

func (this *User) Index() {
	this.UserNeedLogin()

	fmt.Println("index")
	this.Data["User"] = this.User
	this.SetXSRF()
}

func (this *User) Login() {

	if this.IsLogin {
		// TODO: 已经登陆
		this.RedirectAndStop(helper.UrlFor("User.Index"), HTTP_STATUS_CODE_302)
	}

	this.Data["ErrorMessage"] = this.GetSessionOnce(SESSION_KEY_ERROR_MESSAGE)
	this.SetXSRF()
}

func (this *User) LoginPost() {
	this.NeedInputField(helper.UrlFor("User.Login"), VIEW_INPUT_FIELD_USERNAME, VIEW_INPUT_FIELD_PASSWORD)

	result := this.SetLogin(this.GetString(VIEW_INPUT_FIELD_USERNAME), this.GetString(VIEW_INPUT_FIELD_PASSWORD))
	if result.Code != 0 {
		switch result.Code {
		case 200:
			{
				this.SetSession(SESSION_KEY_ERROR_MESSAGE, this.L("ErrorMessage.username_or_passwrod_nil"))
			}
		case 201:
			{
				this.SetSession(SESSION_KEY_ERROR_MESSAGE, this.L("ErrorMessage.user_not_exist"))
			}
		case 202:
			{
				this.SetSession(SESSION_KEY_ERROR_MESSAGE, this.L("ErrorMessage.passwrod_error"))
			}
		default:
			{
				this.SetSession(SESSION_KEY_ERROR_MESSAGE, this.L("ErrorMessage.unknow"))
			}
		}
		this.RedirectAndStop(helper.UrlFor("User.Login"), HTTP_STATUS_CODE_302)
	} else {
		this.RedirectAndStop(helper.UrlFor("User.Index"), HTTP_STATUS_CODE_302)
	}
}

func (this *User) Logout() {
	this.SetLogout()
	this.RedirectAndStop("/", HTTP_STATUS_CODE_302)
}

func (this *User) Signup() {

	if this.IsLogin {
		// TODO: 已经注册
		this.RedirectAndStop(helper.UrlFor("User.Index"), HTTP_STATUS_CODE_302)
	}

	if !this.EnableSignup {
		this.SetSession(SESSION_KEY_ERROR_MESSAGE, this.L("ErrorMessage.signup_limit_on_day"))
	}

	this.Data["ErrorMessage"] = this.GetSessionOnce(SESSION_KEY_ERROR_MESSAGE)
	this.SetXSRF()
}

func (this *User) SignupPost() {
	this.NeedInputField(helper.UrlFor("User.Signup"), VIEW_INPUT_FIELD_USERNAME, VIEW_INPUT_FIELD_PASSWORD, VIEW_INPUT_FIELD_EMAIL)

	ip := this.GetRealIP()
	tmpKey := CACHE_KEY_SIGNUP_NUM + ":" + ip
	num := helper.GetCacheInt(tmpKey)
	if num == 0 {
		// 缓存过期时间为24小时
		helper.SetCacheInt(tmpKey, 0, 60*60*24)
	}

	// 注册次数超过限制
	if num > models.SignupMaxOnDay {
		this.EnableSignup = false
		this.SetSession(SESSION_KEY_ERROR_MESSAGE, this.L("ErrorMessage.signup_limit_on_day"))
		this.RedirectAndStop("/", HTTP_STATUS_CODE_302)
	} else {
		this.EnableSignup = true
		helper.IncCacheInt(tmpKey)
	}

	user := models.User{}
	user.Username = this.GetString(VIEW_INPUT_FIELD_USERNAME)
	user.Password = this.GetString(VIEW_INPUT_FIELD_PASSWORD)
	user.Email = this.GetString(VIEW_INPUT_FIELD_EMAIL)

	if !this.CheckInputValidByUsername(user.Username) {
		this.SetSession(SESSION_KEY_ERROR_MESSAGE, this.L("ErrorMessage.username_valid"))
		this.RedirectAndStop(helper.UrlFor("User.Signup"), HTTP_STATUS_CODE_302)
	}

	if !this.CheckInputValidByEmail(user.Email) {
		this.SetSession(SESSION_KEY_ERROR_MESSAGE, this.L("ErrorMessage.email_valid"))
		this.RedirectAndStop(helper.UrlFor("User.Signup"), HTTP_STATUS_CODE_302)
	}

	has, err := models.QueryUserByUsername(user.Username)
	if err != nil {
		// TODO: 数据库操作出错
	}

	if has {
		this.SetSession(SESSION_KEY_ERROR_MESSAGE, this.L("ErrorMessage.username_has_been_signup"))
		this.RedirectAndStop(helper.UrlFor("User.Signup"), HTTP_STATUS_CODE_302)
	}

	has, err = models.QueryUserByEmail(user.Email)
	if err != nil {
		// TODO: 数据库操作出错
	}

	if has {
		this.SetSession(SESSION_KEY_ERROR_MESSAGE, this.L("ErrorMessage.email_has_been_signup"))
		this.RedirectAndStop(helper.UrlFor("User.Signup"), HTTP_STATUS_CODE_302)
	}

	user.Password = helper.MD5(user.Password)

	err = models.AddUser(&user)
	if err != nil {
		// TODO: 数据库操作出错
		this.RedirectAndStop(helper.UrlFor("Error.Database"), HTTP_STATUS_CODE_302)
	}

	this.RedirectAndStop(helper.UrlFor("User.Index"), HTTP_STATUS_CODE_302)
}

func (this *User) FindPassword() {
	this.Data["ErrorMessage"] = this.GetSessionOnce(SESSION_KEY_ERROR_MESSAGE)
	this.SetXSRF()
}

func (this *User) FindPasswordPost() {
	this.NeedInputField(helper.UrlFor("User.FindPassword"), VIEW_INPUT_FIELD_USERNAME)

	username := this.GetString(VIEW_INPUT_FIELD_USERNAME)

	tmpKey := CACHE_KEY_FIND_PASSWORD_NUM + ":" + username
	num := helper.GetCacheInt(tmpKey)
	if num == 0 {
		// 缓存过期时间为24小时
		helper.SetCacheInt(tmpKey, 0, 60*60*24)
	}

	// 找回密码超过限制
	if num > models.FindPasswordMaxOnDay {
		this.SetSession(SESSION_KEY_ERROR_MESSAGE, this.L("ErrorMessage.findpassword_limit_on_day"))
		this.RedirectAndStop("/", HTTP_STATUS_CODE_302)
	} else {
		helper.IncCacheInt(tmpKey)
	}

	user, err := models.GetUserByUsername(username)
	if err != nil {
		// TODO:数据库操作错误
	}

	if user == nil {
		// 用户不存在
		this.SetSession(SESSION_KEY_ERROR_MESSAGE, this.L("ErrorMessage.user_not_exist"))
		this.RedirectAndStop(helper.UrlFor("User.FindPassword"), HTTP_STATUS_CODE_302)
	}

	email := helper.EncodeEmail(user.Email)
	this.SetSession(SESSION_KEY_AUTH_USERNAME, user.Username)
	this.SetSession(SESSION_KEY_AUTH_EMAIL, user.Email)
	this.SetSession(SESSION_KEY_ENCODE_EMAIL, email)
	this.RedirectAndStop(helper.UrlFor("User.FindPasswordSendEmail"), HTTP_STATUS_CODE_302)
}

func (this *User) FindPasswordSendEmail() {
	this.Data["ErrorMessage"] = this.GetSessionOnce(SESSION_KEY_ERROR_MESSAGE)
	this.Data["EncodeEmail"] = this.GetSessionOnce(SESSION_KEY_ENCODE_EMAIL)
	this.SetXSRF()
}

func (this *User) FindPasswordSendEmailPost() {
	email := helper.GetStringFromInterface(this.GetSessionOnce(SESSION_KEY_AUTH_EMAIL))
	username := helper.GetStringFromInterface(this.GetSessionOnce(SESSION_KEY_AUTH_USERNAME))

	if len(email) == 0 || len(username) == 0 {
		this.RedirectAndStop(helper.UrlFor("User.PasswordFind"), HTTP_STATUS_CODE_302)
	}

	// sendurl -> /user/findpasswordauthemailpost?token=*******************
	token := models.GetChangePasswordToken(username, email, time.Now().String())
	sendUrl := models.GetChangePasswordUrl(token)
	msgHead := this.L("SendMessage.find_password_head")
	msgFooter := this.L("SendMessage.find_password_footer")
	message := msgHead + "<br><a href='" + sendUrl + "' target='_blank'>" + sendUrl + "</a><br>" + msgFooter

	err := helper.SendEmail(models.SmtpFrom, models.SmtpAccount, models.SmtpPassword, models.SmtpServer, email, this.L("SnedMessage.find_password_title"), helper.HtmlUnicode(message), "html")
	if err != nil {
		this.SetSession(SESSION_KEY_ERROR_MESSAGE, this.L("ErrorMessage.send_email_error"))
		this.RedirectAndStop(helper.UrlFor("User.FindPassword"), HTTP_STATUS_CODE_302)
	}

	// token 缓存过期时间为5分钟
	// TODO: 不知道redis空字符串是不是返回""
	str := helper.GetCacheString(CACHE_KEY_FIND_PASSWORD_TOKEN)
	if len(str) == 0 {
		helper.SetCacheString(CACHE_KEY_FIND_PASSWORD_TOKEN, token, 60*5)
	}

	this.SetSession(SESSION_KEY_RESET_USERNAME, username)
	this.SetSession(SESSION_KEY_RESET_EMAIL, email)

	/*
		authUser := models.AuthUser{}
		authUser.Email = email
		err = models.AddAuthUser(authUser)
		if err != nil {
			// TODO:
		}*/
}

func (this *User) FindPasswordAuthEmail() {
	this.Data["ErrorMessage"] = this.GetSessionOnce(SESSION_KEY_ERROR_MESSAGE)
	this.Data["SuccessMessage"] = this.GetSessionOnce(SESSION_KEY_SUCCESS_MESSAGE)
	if this.IsAuthFindPassword {
		this.RedirectAndStop(helper.UrlFor("User.ResetPassword"), HTTP_STATUS_CODE_302)
	} else {
		this.RedirectAndStop(helper.UrlFor("User.FindPassword"), HTTP_STATUS_CODE_302)
	}
}

func (this *User) FindPasswordAuthEmailPost() {
	this.IsAuthFindPassword = false

	this.NeedInputField(helper.UrlFor("User.FindPassword"), VIEW_INPUT_FIELD_TOKEN)

	token := this.GetString(VIEW_INPUT_FIELD_TOKEN)

	str := helper.GetCacheString(CACHE_KEY_FIND_PASSWORD_TOKEN)

	if len(str) == 0 {
		this.SetSession(SESSION_KEY_ERROR_MESSAGE, this.L("ErrorMessage.token_auth_expired"))
		this.RedirectAndStop(helper.UrlFor("User.FindPassword"), HTTP_STATUS_CODE_302)
	}

	if str != token {
		this.SetSession(SESSION_KEY_ERROR_MESSAGE, this.L("ErrorMessage.token_auth_failed"))
		this.RedirectAndStop(helper.UrlFor("User.FindPassword"), HTTP_STATUS_CODE_302)
	}

	this.IsAuthFindPassword = true
	this.SetSession(SESSION_KEY_SUCCESS_MESSAGE, this.L("SuccessMessage.reset_password"))
}

func (this *User) ResetPassword() {
	if !this.IsAuthFindPassword {
		this.RedirectAndStop(helper.UrlFor("User.FindPassword"), HTTP_STATUS_CODE_302)
	}
}

func (this *User) ResetPasswordPost() {
	this.NeedInputField(helper.UrlFor("User.ResetPassword"), VIEW_INPUT_FIELD_NEW_PASSWORD, VIEW_INPUT_FIELD_CONFIRM_PASSWORD)

	email := helper.GetStringFromInterface(this.GetSessionOnce(SESSION_KEY_AUTH_EMAIL))
	username := helper.GetStringFromInterface(this.GetSessionOnce(SESSION_KEY_AUTH_USERNAME))

	newPassword := this.GetString(VIEW_INPUT_FIELD_NEW_PASSWORD)
	confirmPassword := this.GetString(VIEW_INPUT_FIELD_CONFIRM_PASSWORD)

	if newPassword != confirmPassword {
		this.SetSession(SESSION_KEY_ERROR_MESSAGE, this.L("ErrorMessage.password_not_equal"))
		this.RedirectAndStop(helper.UrlFor("User.ResetPassword"), HTTP_STATUS_CODE_302)
	}

	user, err := models.GetUserByUsername(username)
	if err != nil {
		// TODO
	}

	if user.Email != email {
		// TODO
	}

	err = models.EditUser(user)
	if err != nil {

	}

	this.SetSession(SESSION_KEY_SUCCESS_MESSAGE, this.L("SuccessMessage.reset_password_success"))
}

func (this *User) ChangePassword() {
	this.UserNeedLogin()
	this.SetXSRF()
}

func (this *User) ChangePasswordPost() {
	this.UserNeedLogin()
	this.NeedInputField(helper.UrlFor("User.ChangePassword"), VIEW_INPUT_FIELD_OLD_PASSWORD, VIEW_INPUT_FIELD_NEW_PASSWORD, VIEW_INPUT_FIELD_CONFIRM_PASSWORD)

}

func (this *User) Space() {
	var err error
	var user *models.User

	username := this.Ctx.Input.Param(":username")

	if !this.CheckInputValidByUsername(username) {
		this.Abort("404")
	}

	if len(username) > 0 {
		user, err = models.GetUserByUsername(username)
		if err != nil {
			// TODO:
		}

		if user == nil {
			this.Abort("404")
		}
	}

	this.Data["User"] = user
	this.SetXSRF()
}

func (this *User) Profile() {
	this.UserNeedLogin()

	// baseController有Username,User也有username
	user, err := models.GetUserByUsername(this.User.Username)
	if err != nil {
		// TODO:
	}

	this.Data["UserProfile"] = user
	this.SetXSRF()
}

func (this *User) ProfilePost() {
	this.UserNeedLogin()

	user := models.User{}

	user.Uid = this.Uid
	user.Intro = this.GetString(VIEW_INPUT_FIELD_INTRODUCTION)
	user.University = this.GetString(VIEW_INPUT_FIELD_UNIVERSITY)
	user.Occupation = this.GetString(VIEW_INPUT_FIELD_OCCUPATION)
	user.PhoneNumber = this.GetString(VIEW_INPUT_FIELD_PHONENUMBER)

	err := models.EditUser(&user, VIEW_INPUT_FIELD_INTRODUCTION, VIEW_INPUT_FIELD_UNIVERSITY, VIEW_INPUT_FIELD_OCCUPATION, VIEW_INPUT_FIELD_PHONENUMBER)
	if err != nil {
		// TODO:
	}

	this.RedirectAndStop(helper.UrlFor("User.Profile"), HTTP_STATUS_CODE_302)
	// intro
	// phonenumber
	// occupation
	// university
}
