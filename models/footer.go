package models

type Footer struct {
	Name string `xorm:"varchar(32) notnull"`
	Link string `xorm:"text notnull"`
	Icon string `xorm:"text notnull"`
	Lang string `xorm:"varchar(10) notnull"`
}

func GetAllFooter() ([]Footer, error) {
	list := make([]Footer, 0)

	err := db.Table(TABLE_FOOTER).Find(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}
