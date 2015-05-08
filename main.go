package main

import (
	_ "merger_point/routers"
	"runtime"
	"github.com/hprose/hprose-go/hprose"
	"os"
	"fmt"
	//"strconv"
	"merger_point/models"
)

func main() {

	server := hprose.NewTcpServer("tcp4://0.0.0.0:4321/")
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)
	server.ThreadCount = runtime.NumCPU() * 2
	server.AddFunction("band", band)
	server.Start()
	b := make([]byte, 1)
	os.Stdin.Read(b)
	server.Stop()
	//beego.Run()

}

func band(sql string,startTime int64,endTime int64,golang bool,isService bool)string {
	var errors bool = false;
	if(golang){		
		result := models.GetStr(sql,startTime,endTime,&errors)
		fmt.Println("main:getstr",errors)
		if errors {
			return "true"
		}else{
			return result
		}
	}else{		
		result := models.GetAllData(sql,startTime,endTime,isService,&errors)
		fmt.Println("main:getall",errors)
		if errors {
			return "true"
		}else{
			return result
		}
	}
	
}

