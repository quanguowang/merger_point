package models

import (
	"database/sql"
	//"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	//"strconv"
	"time"
)



func  Search(hql string,dbAddress string,ch chan []float64,errors *bool,startTime int64,endTime int64){
	defer func(){ // 必须要先声明defer，否则不能捕获到panic异常
        if err:=recover();err!=nil{	
			*errors = true	 
            log.Println("SearchDBerr:",time.Now(), err)
        }
    }()
	indexs := (endTime-startTime)/(5*60*1000)
	data := make([]float64, indexs)
	//start := time.Now()
	db, err := sql.Open("mysql", dbAddress)
	//end := time.Now()
	//fmt.Println("open conn: ", end.Sub(start))
	log.Println("sql:", hql)
	rows, err := db.Query(hql)
	//queryTime := time.Now()
	//fmt.Println("query time:", queryTime.Sub(end))
	if err != nil {
		ch <- data
		panic(err)
	}

	for rows.Next() {
		var record_time time.Time
		var value float64 
		rows.Scan(&record_time,&value)
		record_long := (record_time.Unix()-(8*60*60))*1000
		index := (record_long-startTime)/(5*60*1000)
		data[index] = value
		
	}
	db.Close()
	//fmt.Println("all query:", time.Now().Sub(queryTime))
	ch <- data
	

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

