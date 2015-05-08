package models

import (
	"database/sql"
	//"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	//"strconv"
	"time"
)



func  Search(hql string,dbAddress string,ch chan map[int64]float64,errors *bool){
	defer func(){ // 必须要先声明defer，否则不能捕获到panic异常
        if err:=recover();err!=nil{	
			*errors = true	 
            log.Println("SearchDBerr:",time.Now(), err)
        }
    }()
	data := map[int64]float64{}
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
		oldValue,ok := data[record_long]
		if ok {
			 data[record_long] = oldValue+value
		}else{
			data[record_long] = value
		}
		
	}
	db.Close()
	//fmt.Println("all query:", time.Now().Sub(queryTime))
	ch <- data
	//return data


}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

