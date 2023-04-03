package application

import (
	"fmt"
	"golang_project_layout/pkg/global"
)

func Run() {
	newUser := BaseDemoUser{Name: "Jack", Phone: "12345678901", Addr: "BJ"}
	global.GVA_DB.Create(&newUser)

	user := BaseDemoUser{}
	global.GVA_DB.First(&user)
	fmt.Printf("%+v\n", user)
}
