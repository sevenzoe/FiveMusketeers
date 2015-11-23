package models

import (
	"errors"
)

type About struct {
	Id   int    `xorm:"int(11) notnull"`
	Name string `xorm:"varchar(100) notnull"`
	//ChairmanSpeak string `xorm:"text notnull"`
	//CompanyInfo     string `xorm:"text notnull"`
	//Honor           string `xorm:"text notnull"`
	//ImportantThings string `xorm:"text notnull"`
}

func GetAbout() (*About, error) {
	about := new(About)

	has, err := db.Table(TABLE_ABOUT).Get(about)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, errors.New("not exist")
	}

	return about, nil
}
