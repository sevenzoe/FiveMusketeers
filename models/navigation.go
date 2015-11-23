package models

type NavigationBar struct {
	Id         int64  `xorm:"int(11) pk autoincr"`
	Name       string `xorm:"varchar(32) notnull"`
	Link       string `xorm:"text notnull"`
	OrderValue int64  `xorm:"int(11) notnull"`
	Icon       string `xorm:"text notnull"`
	Lang       string `xorm:"varchar(10) notnull"`
}

func GetAllNavigationBar() ([]NavigationBar, error) {
	list := make([]NavigationBar, 0)

	err := db.Table(TABLE_NAVIGATION_BAR).Find(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}
