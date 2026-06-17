package model

import (
	"webadmin/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var Db *gorm.DB

func init() {
	dsn := config.Get("mysqldb")
	if dsn == "" {
		dsn = "root:root@tcp(127.0.0.1:3306)/webadmin?charset=utf8mb4&parseTime=True&loc=Local"
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "tu_",
			SingularTable: true,
		},
	})
	if err != nil {
		return
	}
	Db = db
	// if err := db.AutoMigrate(&Admin{}, &AuthGroup{}, &AuthRule{}, &ShellGroup{}, &Shell{}, &ShellMax{}, &ShellMin{}, &Coin{}); err != nil {
	// 	fmt.Println("auto migrate failed:", err)
	// }
}
