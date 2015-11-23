package controllers

import (
	"../helper"
)

var (
	api_field = map[string]string{
		"username": "no_username",
		"password": "no_password",
		"email":    "no_email"}
)

type API struct {
	baseController
}

func (this *API) Login() {
	this.APINeedInputField(VIEW_INPUT_FIELD_USERNAME, VIEW_INPUT_FIELD_PASSWORD)

	tmpKey := CACHE_KEY_LOGIN_NUM + ":" + this.GetString(VIEW_INPUT_FIELD_USERNAME)
	num := helper.GetCacheInt(tmpKey)
	if num == 0 {
		// 缓存过期时间为24小时
		helper.SetCacheInt(tmpKey, 0, 60*60*24)
	}

	// TODO:登录次数比较多的话,是不是应该输入验证码呢?

	result := this.SetLogin(this.GetString(VIEW_INPUT_FIELD_USERNAME), this.GetString(VIEW_INPUT_FIELD_PASSWORD))
	if result.Code != 0 {
		helper.IncCacheInt(tmpKey)
	}

	this.SetResultData(result)
}

func (this *API) APINeedInputField(fields ...string) {
	result := &ResultInfo{Code: 0, Message: "OK", Data: nil}
	for _, field := range fields {

		if len(this.GetString(field)) == 0 {
			result.Code = 900
			result.Message = api_field[field]
			result.Data = nil
			this.SetResultData(result)
			this.StopRun()
		}
	}
}
