package main

import (
	"task3/example"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	example.Run(db)

	// dbConnStr := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"

	// // 使用sqlx连接到数据库
	// db2, err := sqlx.Connect("mysql", dbConnStr)
	// if err != nil {
	// 	log.Fatal("无法连接到数据库:", err)
	// }
	// example.RunSqlX(db2)
	// defer db2.Close()

	// fmt.Println("成功连接到数据库")
}
