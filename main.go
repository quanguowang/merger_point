package main

import (
	_ "merger_point/routers"

	"github.com/hprose/hprose-go/hprose"
	"os"
	//"fmt"
	//"strconv"
	"merger_point/models"
)

func main() {
	server := hprose.NewTcpServer("tcp4://0.0.0.0:4321/")
	server.AddFunction("band", band)
	server.Start()
	b := make([]byte, 1)
	os.Stdin.Read(b)
	server.Stop()
	//beego.Run()

}

func band(sql string,startTime int64,endTime int64,golang bool) string {
	if(golang){		
		return models.GetStr(sql,startTime,endTime)
	}else{		
		return models.GetAllData(sql,startTime,endTime)
	}
	
}




