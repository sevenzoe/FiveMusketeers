package controllers

import (
	"../models"
	"fmt"
)

type Default struct {
	baseController
}

func (this *Default) Index() {

	about, err := models.GetAbout()
	if err != nil {
		fmt.Println(err)
	}

	user, err := models.GetAllUser()
	if err != nil {
		fmt.Println(err)
	}

	u, _ := models.GetUserByEmail("304501847@qq.com")

	fmt.Println("email : ", u)

	this.Data["About"] = about
	this.Data["UserList"] = user

	this.Data["IsLogin"] = this.IsLogin

	this.TplNames = "default/index.html"
}

func (this *Default) Search() {
}
