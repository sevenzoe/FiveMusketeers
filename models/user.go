package models

import (
	//"errors"
	"time"
)

type User struct {
	Uid         int       `xorm:"int(16) pk autoincr"`
	Username    string    `xorm:"varchar(64) unique notnull"`
	Password    string    `xorm:"varchar(32) notnull"`
	Email       string    `xorm:"varchar(32) notnull"`
	Sex         bool      `xorm:"-"`
	University  string    `xorm:"text notnull"`
	CreateTime  time.Time `xorm:"datetime notnull updated"`
	Intro       string    `xorm:"text notnull"`
	PhoneNumber string    `xorm:"varchar(12) unique notnull"`
	Occupation  string    `xorm:"text notnull"`
}

type AuthUser struct {
	Uid        int       `xorm:"int(16) pk autoincr"`
	Email      string    `xorm:"varchar(32) notnull"`
	UpdateTime time.Time `xorm:"datetime updated notnull"`
}

func GetUser(sql string) ([]User, error) {
	list := make([]User, 0)

	err := db.Table(TABLE_USER).Sql(sql).Find(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func GetAllUser() ([]User, error) {
	list := make([]User, 0)

	err := db.Table(TABLE_USER).Find(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func GetUserByUid(uid int) (*User, error) {
	user := new(User)

	has, err := db.Table(TABLE_USER).Id(uid).Get(user)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, nil
	}

	return user, nil
}

func GetUserByEmail(email string) (*User, error) {
	user := new(User)

	has, err := db.Table(TABLE_USER).Where(TABLE_USER_FIELD_EMAIL + " = " + "'" + email + "'").Get(user)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, nil
	}

	return user, nil
}

func GetUserByUsername(username string) (*User, error) {
	user := new(User)

	has, err := db.Table(TABLE_USER).Where(TABLE_USER_FIELD_USERNAME + " = " + "'" + username + "'").Get(user)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, nil
	}

	return user, nil
}

func QueryUserByEmail(email string) (bool, error) {
	user := new(User)

	has, err := db.Table(TABLE_USER).Where(TABLE_USER_FIELD_EMAIL + " = " + "'" + email + "'").Get(user)
	if err != nil {
		return false, err
	}

	return has, nil
}

func QueryUserByUsername(username string) (bool, error) {
	user := new(User)

	has, err := db.Table(TABLE_USER).Where(TABLE_USER_FIELD_USERNAME + " = " + "'" + username + "'").Get(user)
	if err != nil {
		return false, err
	}

	return has, nil
}

func AddUser(user *User) error {
	_, err := db.Table(TABLE_USER).Insert(user)
	if err != nil {
		return err
	}

	return nil
}

// updata
func EditUser(user *User, updataField ...string) error {
	db.Table(TABLE_USER).Id((*user).Uid).UseBool(updataField...).Update(user)
	return nil
}

func GetAuthUserByEmail(email string) (*AuthUser, error) {
	authUser := new(AuthUser)

	has, err := db.Table(TABLE_AUTH_USER).Where(TABLE_USER_FIELD_EMAIL + " = " + "'" + email + "'").Get(authUser)
	if err != nil {
		return nil, err
	}

	if !has {
		return nil, nil
	}

	return authUser, nil
}

func AddAuthUser(user *AuthUser) error {
	_, err := db.Table(TABLE_AUTH_USER).Insert(user)
	if err != nil {
		return err
	}

	return nil
}
