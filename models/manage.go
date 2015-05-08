package models

import (
	"github.com/astaxie/beego"
	"fmt"
	"github.com/hprose/hprose-go/hprose"
	"log"
	"strconv"
	"encoding/json"
	"time"
	"strings"
)



func GetStr(sql string,startTime int64,endTime int64,errors *bool)string{
	resultData := GetData(sql,startTime,endTime,errors)

	bytes, _ := json.Marshal(resultData)  
	jsonStr := string(bytes)
	return jsonStr
}

func GetData(sql string,startTime int64,endTime int64,errors *bool)[]Point{
	dbAddress := map[int64]string{}
	var i int64 = 0
	for  {
		db := beego.AppConfig.String("portal_"+strconv.FormatInt(i,10))
		if "" != db {
			dbAddress[i] = db
			i++
		}else{
			break
		}
	}
	var ch chan map[int64]float64 = make(chan map[int64]float64)
	for i=0;i<int64(len(dbAddress));i++{
		value := dbAddress[i]
		go Search(sql,value,ch,errors)
	}
		
	allData := map[int64]map[int64]float64{}
	var j int64 = 0
	for j = 0;j<i;j++{
		allData[j] = <- ch
	}

	return mergerMap(allData,startTime,endTime)
}

type Call struct {
	Band func(string,int64,int64,bool,bool) (string) 
	
}



func GetAllData(sql string,startTime int64,endTime int64,isService bool,errors *bool)string{	
	pcAddress := map[int64]string{}
	var i int64 = 0
	for  {
		pc := beego.AppConfig.String("pc_"+strconv.FormatInt(i,10))
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
		go  callPc(sql,startTime,endTime,value,chpc,isService,errors)
	}

	allData := map[int64][]Point{}

	arrayData := GetData(sql,startTime,endTime,errors)

	var j int64 = 0
	for j = 0;j<i;j++{
		allData[j] = <- chpc
	}
	allData[i] = arrayData
	data := mergerArray(allData,startTime,endTime,isService)
	
	bytes, _ := json.Marshal(data)  
	jsonStr := string(bytes)
	return jsonStr
}

func callPc(sql string,startTime int64,endTime int64,pcAddress string,chpc chan []Point,isService bool,errors *bool){
	index := (endTime-startTime)/(5*60*1000)
	arrayData := make([]Point,index)
	defer func(){ 
        if err:=recover();err!=nil{
			*errors = true	
            log.Println("callPCerr:",time.Now(), err)
			chpc <- arrayData
        }

    }()
	client := hprose.NewClient(pcAddress)
	var call *Call
	client.UseService(&call)
	data := call.Band(sql,startTime,endTime,true,isService)
	if strings.EqualFold(data, "true"){
		*errors = true
		chpc <- arrayData
		panic("call pc error")
	}
	dataByte := []byte(data)

	json.Unmarshal(dataByte, &arrayData)

	chpc <- arrayData
}