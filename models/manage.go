package models

import (
	"github.com/astaxie/beego/config"
	"fmt"
	"github.com/hprose/hprose-go/hprose"
	//"log"
	"strconv"
	"encoding/json"
	"time"
)



func GetStr(sql string,startTime int64,endTime int64) string{
	start := time.Now()
	resultData := GetData(sql,startTime,endTime)
	fmt.Println("getData time: ", time.Now().Sub(start))

	bytes, _ := json.Marshal(resultData)  
	jsonStr := string(bytes)
	return jsonStr
}

func GetData(sql string,startTime int64,endTime int64)[]Point{
	iniconf, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
	    //.Fatal(err)
		fmt.Println(err)
	}
	dbAddress := map[int64]string{}
	var i int64 = 0
	for  {
		db := iniconf.String("portal_"+strconv.FormatInt(i,10))
		if "" != db {
			dbAddress[i] = db
			i++
		}else{
			break
		}
	}
	var ch chan map[int64]int64 = make(chan map[int64]int64)
	for i=0;i<int64(len(dbAddress));i++{
		value := dbAddress[i]
		go Search(sql,value,ch)
	}

	
	allData := map[int64]map[int64]int64{}
	var j int64 = 0
	for j = 0;j<i;j++{
		allData[j] = <- ch
	}
	
	return mergerMap(allData,startTime,endTime)
}

type Call struct {
	Band func(string,int64,int64,bool) (string) 
	
}



func GetAllData(sql string,startTime int64,endTime int64)string{
	iniconf, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
	    //.Fatal(err)
		fmt.Println(err)
	}
	
	pcAddress := map[int64]string{}
	var i int64 = 0
	for  {
		pc := iniconf.String("pc_"+strconv.FormatInt(i,10))
		if  "" != pc {
			fmt.Println(pc)
			pcAddress[i] = pc
			i++
		
		}else{
			break
		}
		
	}
	var chpc chan []Point = make(chan []Point)
	var k int64 = 0
	for k=0;k<int64(len(pcAddress));k++{
		value := pcAddress[k]
		go  callPc(sql,startTime,endTime,value,chpc)
	}

	allData := map[int64][]Point{}

	arrayData := GetData(sql,startTime,endTime)
	var j int64 = 0
	for j = 0;j<i;j++{
		allData[j] = <- chpc
	}
	allData[i] = arrayData

	data := mergerArray(allData,startTime,endTime)
	
	bytes, _ := json.Marshal(data)  
	jsonStr := string(bytes)
	//fmt.Println(jsonStr)
	return jsonStr
}

func callPc(sql string,startTime int64,endTime int64,pcAddress string,chpc chan []Point){
	client := hprose.NewClient(pcAddress)
	var call *Call
	client.UseService(&call)
	data := call.Band(sql,startTime,endTime,true)
	dataByte := []byte(data)
	
	//数据点，比如1个月8000多个点
	index := (endTime-startTime)/(5*60*1000)
	arrayData := make([]Point,index)
	
	json.Unmarshal(dataByte, &arrayData)
	
	//return arrayData
	chpc <- arrayData
}