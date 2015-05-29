package models

import (
	"github.com/astaxie/beego"
	//"fmt"
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

func GetData(sql string,startTime int64,endTime int64,errors *bool)[]float64{
	
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
	
	var ch chan []float64 = make(chan []float64)
	for i=0;i<int64(len(dbAddress));i++{
		value := dbAddress[i]
		go Search(sql,value,ch,errors,startTime,endTime)
	}
	allData := map[int64][]float64{}
	
	var j int64 = 0
	for j = 0;j<i;j++{
		allData[j] = <- ch
	}

	return mergerMap(allData,startTime,endTime)
}

type Call struct {
	Band func(string,int64,int64,bool,bool) (string) 
	
}

type Result struct {
    Data    []Point
    Maxy    float64
   
}

func GetAllData(sql string,startTime int64,endTime int64,isService bool,errors *bool)string{	
	pcAddress := map[int64]string{}
	var i int64 = 0
	for  {
		pc := beego.AppConfig.String("pc_"+strconv.FormatInt(i,10))
		if  "" != pc {
			//fmt.Println(pc)
			pcAddress[i] = pc
			i++		
		}else{
			break
		}
		
	}
	var chpc chan []float64 = make(chan []float64)
	var k int64 = 0

	for k=0;k<int64(len(pcAddress));k++{
		value := pcAddress[k]
		go  callPc(sql,startTime,endTime,value,chpc,isService,errors)
	}

	allData := map[int64][]float64{}

	arrayData := GetData(sql,startTime,endTime,errors)

	var j int64 = 0
	for j = 0;j<i;j++{
		allData[j] = <- chpc
	}
	allData[i] = arrayData
	data,maxY := mergerArray(allData,startTime,endTime,isService)
	
	result := Result{data,maxY}
	bytes, _ := json.Marshal(result)  
	jsonStr := string(bytes)
	return jsonStr
}

func callPc(sql string,startTime int64,endTime int64,pcAddress string,chpc chan []float64,isService bool,errors *bool){
	index := (endTime-startTime)/(5*60*1000)
	arrayData := make([]float64,index)
	defer func(){ 
        if err:=recover();err!=nil{
			*errors = true	
            log.Println(pcAddress," callPCerr:",time.Now(), err)
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