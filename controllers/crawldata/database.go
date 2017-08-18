package crawldata

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	//"strconv"
	s "strings"
)

func OpenDatabase() (*sql.DB, error) {
	// 连接数据库
	db, err := sql.Open("mysql", "root:chenghai3c@tcp(localhost:3306)/wish_catch?charset=utf8")
	if err != nil {
		return nil, err
	}
	return db, nil
}

/*
   该函数将获取的数据存储到数据库
*/
func InsertData(datas *ImageDatas) {
	imageDatas := *datas
	// 连接数据库
	db, err := OpenDatabase()
	if err != nil {
		fmt.Printf(s.Join([]string{"连接数据库失败", err.Error()}, "-->"))
	}
	defer db.Close()

	for i := 0; i < len(imageDatas); i++ {
		imageData := imageDatas[i]
		// Prepare statement for inserting data
		imgIns, err := db.Prepare("INSERT INTO gratisography (img_url, type_name, title, width, height) VALUES( ?, ?, ?, ?, ? )") // ? = placeholder
		if err != nil {
			fmt.Println(s.Join([]string{"拼装数据格式", err.Error()}, "-->"))
		}
		defer imgIns.Close() // Close the statement when we leave main()

		//预格式是否存在http://www.gratisography.com
		if len(imageData.Src) < 30 {
			_, err := imgIns.Exec(s.Join([]string{"http://www.gratisography.com", imageData.Src}, "/"), imageData.Tp, imageData.Title, imageData.Width, imageData.Height)
			if err != nil {
				fmt.Println(s.Join([]string{"插入数据失败", err.Error()}, "-->"))
			}
		} else {
			_, err := imgIns.Exec(s.Join([]string{"", imageData.Src}, ""), imageData.Tp, imageData.Title, imageData.Width, imageData.Height)
			if err != nil {
				fmt.Println(s.Join([]string{"插入数据失败", err.Error()}, "-->"))
			}
		}
	}
}

func GetAllImages() (imageDatas ImageDatas, err error) {
	// 连接数据库
	db, err := OpenDatabase()
	if err != nil {
		fmt.Printf(s.Join([]string{"连接数据库失败", err.Error()}, "-->"))
		return nil, err
	}
	defer db.Close()

	// Prepare statement for inserting data
	imgOut, err := db.Query("SELECT * FROM gratisography")
	if err != nil {
		fmt.Println(s.Join([]string{"获取数据失败", err.Error()}, "-->"))
		return nil, err
	}
	defer imgOut.Close()

	// 定义扫描select到的数据库字段的变量
	var (
		id          int
		img_url     string
		type_name   string
		title       string
		width       int
		height      int
		create_time string
	)
	for imgOut.Next() {
		// db.Query()中select几个字段就需要Scan多少个字段
		err := imgOut.Scan(&id, &img_url, &type_name, &title, &width, &height, &create_time)
		if err != nil {
			fmt.Println(s.Join([]string{"查询数据失败", err.Error()}, "-->"))
			return nil, err
		} else {
			imageData := ImageData{img_url, type_name, title, width, height}
			imageDatas = append(imageDatas, imageData)
		}
	}

	return imageDatas, nil
}
