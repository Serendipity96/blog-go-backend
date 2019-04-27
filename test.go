package main

import (
	"database/sql"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"net/http"
)


func pong(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("pong"))
}


func main() {

	// 连接数据库
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/bloggo")
	if err != nil {
		fmt.Println("failed connect to mysql", err.Error())
	}

	// 查询tags
	// 为什么要这么写
	rows, err := db.Query("SELECT * FROM tags")
	if err != nil {
		fmt.Println("failed to query", err.Error())
	}
	defer rows.Close()

	tagList := make(map[int]string)
	for rows.Next() {
		var tagId int
		var tagName string
		if err := rows.Scan(&tagId, &tagName); err != nil {
			fmt.Println("failed for", err.Error())
		}
		tagList[tagId] = tagName
	}

	fmt.Println(tagList)

	defer db.Close()

	http.HandleFunc("/", pong)
	http.ListenAndServe(":8081", nil)
}
