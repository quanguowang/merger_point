package models

import (
	"database/sql"
	//"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	//"strconv"
	"time"
	//"fmt"
)



func  Search(hql string,dbAddress string){
	data := map[int64]int64{}
	//start := time.Now()
	db, err := sql.Open("mysql", dbAddress)
	//end := time.Now()
	//fmt.Println("open conn: ", end.Sub(start))
	checkErr(err)
	log.Println("sql:", hql)
	rows, err := db.Query(hql)
	//queryTime := time.Now()
	//fmt.Println("query time:", queryTime.Sub(end))
	checkErr(err)
	
	//data := map[int64]int64{}
	for rows.Next() {
		var record_time time.Time
		var value int64 
		rows.Scan(&record_time,&value)
		record_long := (record_time.Unix())*1000
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



