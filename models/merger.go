package models

import (
	//"encoding/json"
	//"fmt"
)

type Point struct {
    X    int64
    Y    float64
   
}

func mergerArray(allData map[int64][]Point,startTime int64,endTime int64,isService bool)[]Point{
	index := (endTime-startTime)/(5*60*1000)	
	mergerData := make([]Point,index)
	
	var j int64 = 0
	for j = 0; j < index; j++ {	
		x := startTime+(5*60*1000)*int64(j)
		mergerData[j].X = x	
	}

	var i int64 = 0
	var k int64 = 0
	for k=0;k<int64(len(allData));k++{
		value := allData[k]
		for i = 0;i<index;i++{
			mergerData[i].Y += value[i].Y
			
		}
	}
//	for _,value := range allData{
//		for i = 0;i<index;i++{
//			mergerData[i].Y += value[i].Y
//		}
//	}

	interval := getInterval(endTime-startTime)
	point := interval/(5*60)
	//(5*60)/8=37.5
	var coefficient float64 = 37.5
	if !isService {
		coefficient = float64(interval)/8
	}
	if point > 1 {
		result := make([]Point,index/point)
		if isService {	
			for i = 0;i<index/point;i++ {
				var maxValue float64 = 0
				result[i].X = mergerData[i*point].X
			
				for j=0;j<point;j++ {
					if maxValue < mergerData[i*point+j].Y{
						maxValue = mergerData[i*point+j].Y
					}
				}
				result[i].Y = maxValue/coefficient
			}
			return result
		}else{
			for i = 0;i<index/point;i++ {
				var value float64 = 0
				result[i].X = mergerData[i*point].X
			
				for j=0;j<point;j++ {
					value += mergerData[i*point+j].Y
				}
				result[i].Y = value/coefficient
			}
			return result
		}
			
	}else{
		for i = 0;i<index;i++{
			mergerData[i].Y = (mergerData[i].Y)/coefficient
			
		}
		return mergerData
		
	}
	
}


func mergerMap(allData map[int64]map[int64]float64,startTime int64,endTime int64) []Point{
	//数据点，比如1个月8000多个点
	index := (endTime-startTime)/(5*60*1000)

	resultData := make([]Point, index)
	var j int64 = 0
	for j = 0; j < index; j++ {
		x := startTime+(5*60*1000)*int64(j)
		resultData[j].X = x
	
	}
	var i int64  = 0
	for i=0;i<int64(len(allData));i++{
		value := allData[i]
		var k int64 = 0
		for k=0;k<index;k++{
			x := resultData[k].X
			resultData[k].Y += value[x]
		}
	}
	
	return resultData

}

func getInterval(time int64)int64{
	var fl float64 = (24 * 60 * 60 * 1000 * 1.0) 
	fl = (float64(time) / fl)
	if fl<1{
		fl = 1
	}
	return 5 * 60 * int64(fl)
}

//func mergerLocalArray(allData []Point,startTime int64,endTime int64)[]Point{
//	index := (endTime-startTime)/(5*60*1000)

//	interval := getInterval(endTime-startTime)
//	point := interval/(5*60*1000)
//	fmt.Println(point)
//	fmt.Println(point)
//	var result = make([]Point,index/point)
	
//	if point > 1 {
//		var i int64 = 0
//		for i = 0;i<index/point;i++ {
//			var maxValue int64 = 0
//			result[i].X = allData[i*point].X
//			var j int64 = 0
//			for j=0;j<point;j++ {
//				if maxValue < allData[i*point+j].Y{
//					maxValue = allData[i*point+j].Y
//				}
//			}
//			result[i].Y = maxValue
//		}
//	}
//	return result
	
//}