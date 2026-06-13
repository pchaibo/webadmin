package main

import (
    "fmt"
    "webadmin/config"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/schema"
)

func main() {
    dsn := config.Get("mysqldb")
    if dsn == "" {
        dsn = "root:root@tcp(127.0.0.1:3306)/webadmin?charset=utf8mb4&parseTime=True&loc=Local"
    }
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        NamingStrategy: schema.NamingStrategy{TablePrefix: "tu_", SingularTable: true},
    })
    if err != nil {
        fmt.Println("open error:", err)
        return
    }

    rows, err := db.Raw("SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE, IFNULL(COLUMN_DEFAULT,''), IFNULL(COLUMN_COMMENT,'') FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'webadmin' AND TABLE_NAME = 'tu_shellgroup' ORDER BY ORDINAL_POSITION").Rows()
    if err != nil {
        fmt.Println("query error:", err)
        rows2, _ := db.Raw("SHOW TABLES").Rows()
        if rows2 != nil {
            fmt.Println("Existing tables:")
            for rows2.Next() {
                var t string
                rows2.Scan(&t)
                fmt.Println(" -", t)
            }
            rows2.Close()
        }
        return
    }
    defer rows.Close()

    fmt.Println("tu_shellgroup columns:")
    for rows.Next() {
        var name, colType, nullable, defaultVal, comment string
        rows.Scan(&name, &colType, &nullable, &defaultVal, &comment)
        fmt.Printf("  %-20s %-20s nullable=%-5s default=%-10s comment=%s\n", name, colType, nullable, defaultVal, comment)
    }
}
