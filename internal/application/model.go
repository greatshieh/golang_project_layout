package application

import "golang_project_layout/pkg/global"

type BaseDemoUser struct {
	global.GVA_MODEL
	Name  string `json:"name" gorm:"comment:用户姓名"`
	Phone string `json:"phone" gorm:"comment:联系方式"`
	Addr  string `json:"addr" gorm:"comment:地址"`
}

func (BaseDemoUser) TableName() string {
	return "base_user"
}

type BaseDemoBook struct {
	global.GVA_MODEL
	Title  string `json:"title" gorm:"comment:书名"`
	Author string `json:"author" gorm:"comment:作者"`
	ISBN   string `json:"ISBN" gorm:"comment:ISBN"`
}

func (BaseDemoBook) TableName() string {
	return "base_book"
}
