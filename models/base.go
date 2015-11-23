package models

import (
	//"fmt"
	"../helper"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"strconv"
	//"time"
)

const (
	TABLE_USER           = "user"
	TABLE_ABOUT          = "about"
	TABLE_AUTH_USER      = "authuser"
	TABLE_NAVIGATION_BAR = "navigation_bar"
	TABLE_FOOTER         = "footer"

	TABLE_USER_FIELD_USERNAME = "username"
	TABLE_USER_FIELD_PASSWORD = "password"
	TABLE_USER_FIELD_EMAIL    = "email"
)

var (
	RedisOn              bool
	RedisAddress         string
	SignupMaxOnDay       int
	FindPasswordMaxOnDay int
	SmtpFrom             string
	SmtpServer           string
	SmtpAccount          string
	SmtpPassword         string
	SmtpProtocol         string
	SmtpDomain           string

	NavBarMap map[string][]NavigationBar
	FooterMap map[string][]Footer
)

var db *xorm.Engine

func init() {
	initDatabase()
	initAppConfig()
	initNavigation()
}

func initDatabase() error {
	var err error

	username := beego.AppConfig.String("database::username")
	password := beego.AppConfig.String("database::password")
	ip := beego.AppConfig.String("database::ip")
	port := beego.AppConfig.String("database::port")
	name := beego.AppConfig.String("database::name")

	dataSourceName := username + ":" + password + "@tcp(" + ip + ":" + port + ")/" + name + "?charset=utf8"
	db, err = xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		db.ShowErr = true
		return err
	}

	maxOpenConns := beego.AppConfig.String(("database::maxOpenConns"))
	maxIdleConns := beego.AppConfig.String(("database::maxIdleConns"))

	maxOC, _ := strconv.Atoi(maxOpenConns)
	macIC, _ := strconv.Atoi(maxIdleConns)

	db.SetMaxOpenConns(maxOC)
	db.SetMaxIdleConns(macIC)

	db.Sync2(new(About), new(User), new(NavigationBar))
	db.ShowSQL = true

	return nil
}

func initAppConfig() error {
	var err error

	RedisOn, err = beego.AppConfig.Bool("redisOn")
	RedisAddress = beego.AppConfig.String("redisAddress")
	SignupMaxOnDay, err = beego.AppConfig.Int("signup::signupMaxOnDay")
	FindPasswordMaxOnDay, err = beego.AppConfig.Int("findpassword::findpasswordMaxOnday")

	SmtpFrom = beego.AppConfig.String("smtp::from")
	SmtpServer = beego.AppConfig.String("smtp::server")
	SmtpAccount = beego.AppConfig.String("smtp::account")
	SmtpPassword = beego.AppConfig.String("smtp::password")
	SmtpProtocol = beego.AppConfig.String("smtp::protocol")
	SmtpDomain = beego.AppConfig.String("smtp::domain")

	if err != nil {
		return err
	}

	return nil
}

func initNavigation() error {
	NavBarMap = make(map[string][]NavigationBar)

	NavList, err := GetAllNavigationBar()
	if err != nil {
		// TODO:
		return err
	}

	for i, _ := range NavList {
		if len(NavList[i].Lang) > 0 {
			lang := NavList[i].Lang
			NavBarMap[lang] = append(NavBarMap[lang], NavList[i])
		}
	}

	return nil
}

func GetChangePasswordToken(username, email, now string) string {
	return helper.MD5(username + email + now)
}

func GetChangePasswordUrl(token string) string {
	return SmtpProtocol + "://" + SmtpDomain + "/user/findpasswordauthemailpost?token=" + token
}
