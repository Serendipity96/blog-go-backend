package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"io/ioutil"
	"net/http"
)

type tagsJson struct {
	TagId   int    `json:"tag_id"`
	TagName string `json:"tag_name"`
}
type articleJson struct {
	ArticleId int `json:"article_id"`
	ArticleTitle string `json:"article_title"`
	ArticleTime string `json:"article_time"`
	ArticleAbstract string `json:"article_abstract"`
}

func main() {

	// 连接数据库
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/bloggo")
	if err != nil {
		fmt.Println("failed connect to mysql", err.Error())
		return
	}

	rows, err := db.Query("SELECT * FROM tag")
	if err != nil {
		fmt.Println("failed to query", err.Error())
		return
	}
	defer rows.Close()

	var tJson tagsJson
	var tList []tagsJson
	for rows.Next() {
		var tagId int
		var tagName string
		if err := rows.Scan(&tagId, &tagName); err != nil {
			fmt.Println("failed for", err.Error())
			return
		}
		tJson = tagsJson{tagId, tagName}
		tList = append(tList, tJson)
	}

	tRes, err := json.Marshal(tList)
	if err != nil {
		fmt.Println("convert json failed", err)
		return
	}

	defer db.Close()

	http.HandleFunc("/tags", func(resWriter http.ResponseWriter, request *http.Request) {
		resWriter.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
		resWriter.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
		resWriter.Header().Set("content-type", "application/json")
		resWriter.Write(tRes)


	})

	http.HandleFunc("/articleList", func(resWriter http.ResponseWriter, request *http.Request) {
		resWriter.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
		resWriter.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
		resWriter.Header().Set("content-type", "application/json")
		result, err := ioutil.ReadAll(request.Body)
		if err != nil {
			fmt.Println("request.Body decoded failed", err)
			return
		}
		reqTagId := string(result)

		if reqTagId != ""{
			queryStr := "SELECT id,title,timestamp,abstract FROM article where tag_id = " + reqTagId
			articleList, err := db.Query(queryStr)
			if err != nil {
				fmt.Println("failed to query: ", err.Error())
				return
			}
			defer articleList.Close()

			var aJson articleJson
			var aList []articleJson

			for articleList.Next() {
				var artId int
				var artTitle string
				var artTime []uint8
				var artAbstract string
				if err := articleList.Scan(&artId, &artTitle,&artTime,&artAbstract); err != nil {
					fmt.Println("failed for:::::", err.Error())
					return
				}
				artTime2 :=string(artTime)
				aJson = articleJson{artId, artTitle,artTime2,artAbstract}
				aList = append(aList, aJson)
			}

			aRes, err := json.Marshal(aList)
			if err != nil {
				fmt.Println("convert json failed", err)
				return
			}
			fmt.Println("tt",tList)
			fmt.Println("aa",aList)
			resWriter.Write(aRes)
		}

	})
	http.ListenAndServe(":8081", nil)
}
