package points

import (
	//"awesomeProject/fops/points"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	//"strings"

	//"reflect"
	"sort"
	"strconv"
	//"database/sql"
	//"container/heap"
)

type CsvLine struct {
	number             string `json:"number"`
	route_id           string `json:"route_id"`
	route_short_name   string `json:"route_short_name"`
	route_long_name    string `json:"route_long_name"`
	transport_type     string `json:"transport_type"`
	direction          string `json:"direction"`
	stop_id            string `json:"stop_id"`
	next_stop          string `json:"next_stop"`
	stop_distance      string `json:"stop_distance"`
	stop_name          string `json:"stop_name"`
	coordinates        string `json:"coordinates"`
	lon                string `json:"lon"`
	lng                string `json:"lng"`
	timeBetweenStop    string `json:"timeBetweenStop"`
	hasStopThisInteval string `json:"hasStopThisInteval"`
}

// прямой
//var beginlat, beginlon, endlat, endlon float64 = 59.9318, 30.33624, 59.92575, 30.296333
// с пересадкой
//var beginlat, beginlon, endlat, endlon float64 = 59.919147, 30.338238, 59.95316833811,30.218256
// var beginlat, beginlon, endlat, endlon float64 = 59.918147, 30.338238, 59.933811,30.218256
//var beginlat, beginlon, endlat, endlon float64 = 59.920746, 30.387825, 59.936516, 30.302124

// !var begin string = "Church of the Savior on Spilled Blood"
// !var end string = "Grand Choral Synagogue"

var weekDay, fileHour, fileMinute string = neccessaryFile()


//var beginlat, beginlon, endlat, endlon float64 = beginCoords.Lat, beginCoords.Lon, endCoords.Lat, endCoords.Lon

var price int = 0



func bus1(beginlat, beginlon, endlat, endlon float64) (int,[]string, []int, error){

	//var beginCoords, endCoords coords = findLocation(begin, end)
	//var beginlat, beginlon, endlat, endlon float64 = beginCoords.Lat, beginCoords.Lon, endCoords.Lat, endCoords.Lon
	var stopBeginlat, stopBeginlon, stopEndlat, stopEndlon , _, _ = nearest_stops(beginlat, beginlon, endlat, endlon, 0.8)

	min, index, nearBusBegin, nearBusEnd,lenSuitableBuses, _  ,_, _:= directBus(beginlat, beginlon, endlat, endlon, 0.8, 100)

	var price float64

	var quantityTransport int
	quantityTransport = 1
	price = 50

	if min < 1 || min > 500000000 || lenSuitableBuses < 2 {


		quantityTransport = 2
		index1Best, index2Best,commonCoords, potCommonIndex, error := transfer(nearBusBegin, nearBusEnd, stopBeginlat, stopBeginlon, stopEndlat, stopEndlon ,min, beginlat, beginlon,endlat, endlon)

		if error != nil{
			return 0,commonCoords, potCommonIndex, error
		}

		if min != 0 {
			var _, _, stop_distancePath1, route_short_namePath1, transport_typePath1, stop_namePathBegin1, stop_namePathEnd1, coordinatesPath1 = printAllAboutBus(index1Best)
			var _, _, stop_distancePath2, route_short_namePath2, transport_typePath2, stop_namePathBegin2, stop_namePathEnd2, coordinatesPath2 = printAllAboutBus(index2Best)
			fmt.Println("Транспорта №1", "\n", "Расстояние", sumSliceString(stop_distancePath1), "\n", "Название транспорта", route_short_namePath1, "\n", "Вид транспорта", transport_typePath1, "\n", "Остановка начало", stop_namePathBegin1, "\n", "Остановка конец", stop_namePathEnd1, "\n", "Координаты", coordinatesPath1)
			fmt.Println("\n", "Транспорта №2", "\n", "Расстояние", sumSliceString(stop_distancePath2), "\n", "Название транспорта", route_short_namePath2, "\n", "Вид транспорта", transport_typePath2, "\n", "Остановка начало", stop_namePathBegin2, "\n", "Остановка конец", stop_namePathEnd2, "\n", "Координаты", coordinatesPath2)
			fmt.Println("Стоимость проезда", price)
			fmt.Println("Общее время у пути", commonTime(min), "минут")
		}
	} else {
		var _, _, stop_distancePath, route_short_namePath, transport_typePath, stop_namePathBegin, stop_namePathEnd, coordinatesPath = printAllAboutBus(index)
		fmt.Println("Расстояние", sumSliceString(stop_distancePath),"\n","Название транспорта" ,route_short_namePath,"\n","Вид транспорта", transport_typePath, "\n","Остановка начало", stop_namePathBegin,"\n","Остановка конец" ,stop_namePathEnd,"\n","Координаты" ,coordinatesPath)
		fmt.Println("Стоимость проезда", price)
		fmt.Println("Общее время у пути", commonTime(min), "минут")
	}

	return quantityTransport, []string{"1"}, []int{1}, nil
}


func transfer(nearBusBegin, nearBusEnd []string, stopBeginlat, stopBeginlon, stopEndlat, stopEndlon []float64,min float64, beginlat, beginlon,endlat, endlon float64) ([]int, []int,[]string, []int, error){

	var commonCoords, lenCommonCoords, potCommonIndex, potCommonBusBegin, potCommonBusEnd = findPotentialBuses(nearBusBegin, nearBusEnd)
	for l1, l2 := range lenCommonCoords {
		if l2 == 0 {
			potCommonBusBegin[l1] = "0"
			potCommonBusEnd[l1] = "0"
		}
	}
	var potCommonBusBegin1 = remove(potCommonBusBegin, "0")
	var potCommonBusEnd1 = remove(potCommonBusEnd, "0")
	var lenCommonCoords1 = removeInt(lenCommonCoords, 0)
	var transfer = suitTransfer(lenCommonCoords1, potCommonIndex, potCommonBusBegin1, potCommonBusEnd1, stopBeginlat, stopBeginlon, stopEndlat, stopEndlon)

	if len(transfer) == 0 {
		return []int{1}, []int{1},commonCoords, potCommonIndex, errors.New("try find path with two transfer")
	}



	var alon, alat, _, _, _, _, _, _, _, _, _, _, _, _ = readFile()
	var index1Best, index2Best []int

	for _, j := range transfer {
		translon, _ := strconv.ParseFloat(alon[j], 8)
		translat, _ := strconv.ParseFloat(alat[j], 8)

		min1,index1,  _, _ , _, _, _, _:= directBus(beginlat, beginlon, translon, translat, 0.8, 1)
		min2,index2,  _, _, _, _,  _, _:= directBus(translon, translat, endlat, endlon, 0.8, 1)


		var dist float64 = 500000000000000000
		min = min1 + min2
		if min1 > 0 && min2 > 0 && min < dist {
			dist = min
			index1Best = index1
			index2Best = index2
		}
	}

	return index1Best, index2Best,commonCoords,potCommonIndex, nil

}


func directBus(beginlat, beginlon, endlat, endlon, dist, coef float64) (float64,[]int,[]string,[]string, int, error, map[int]float64, map[int]float64){
	var stopBeginlat, stopBeginlon, stopEndlat, stopEndlon, nearBeginIndexes, nearEndIndexes  = nearest_stops(beginlat, beginlon, endlat, endlon, dist)

	var nearBusBegin, nearBusEnd []string = nearestBusinStops(stopBeginlat, stopBeginlon, stopEndlat, stopEndlon)
	var suitableBuses = suitableBus(nearBusBegin, nearBusEnd)
	var timeToStopByIndexBegin = preCalculate(nearBeginIndexes, beginlat, beginlon)

	var timeToStopByIndexEnd = preCalculate(nearEndIndexes, endlat, endlon)

	var min, index = bestNumber(stopBeginlat, stopBeginlon, stopEndlat, stopEndlon, nearBusBegin, nearBusEnd, suitableBuses, beginlat, beginlon, endlat, endlon, timeToStopByIndexBegin, timeToStopByIndexEnd, coef)
	lenSuitableBuses := len(suitableBuses)

	if len(index) > 0 {
		return min, index, nearBusBegin, nearBusEnd, lenSuitableBuses, nil, timeToStopByIndexBegin, timeToStopByIndexEnd
	}else{
		return min, index, nearBusBegin, nearBusEnd, lenSuitableBuses, errors.New("try find path with one transfer"), timeToStopByIndexBegin, timeToStopByIndexEnd
	}
	//return 1,[]int{1},[]string{"1"},[]string{"1"}, nil
}

type nearStopsTime struct {
	stop_id string
	indexesOfStops []int
	timeToStopByFoot float64
}


func preCalculate(nearBeginIndexes []int, beginlat, beginlon float64) map[int] float64 {
	//var _, _, _, _ , nearBeginIndexes, _  = nearest_stops(beginlat, beginlon, endlat, endlon, dist)

	var stopTimeBegin map[int] float64

	//var stopBeginIndexes nearStopsTime

	stopTimeBegin = make(map[int]float64)


	var stopIndexBegin map[string] []int
	stopIndexBegin = make(map[string] []int)

	for _, l := range nearBeginIndexes{
		for k, n := range stop_id{
			if n == stop_id[l] && k == l {
				stopIndexBegin[n] = append(stopIndexBegin[n], l)
			}

		}
	}

	for _, l := range stopIndexBegin{

		lonStop1, _ := strconv.ParseFloat(alat[l[0]], 8)
		latStop1, _ := strconv.ParseFloat(alon[l[0]], 8)


		timeToStop := Distance(beginlat, beginlon, latStop1, lonStop1) / 5 * 60


		for _, index := range l {
			stopTimeBegin[index] = timeToStop
		}

	}

	return stopTimeBegin
}



func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}


func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			k := append(s[:i], s[i+1:]...)
			return remove(k, r)
		}
	}

	return s
}

func removeInt(s []int, r int) []int {
	for i, v := range s {
		if v == r {
			k := append(s[:i], s[i+1:]...)
			return removeInt(k, r)
		}
	}
	return s
}


func suitTransfer(lenCommonCoords, potCommonIndex[] int, potCommonBusBegin, potCommonBusEnd []string, stopBeginlat, stopBeginlon, stopEndlat, stopEndlon []float64 ) [] int {



	var beginOk []int
	var endOk []int
	var d string

	var alon, alat, route_id, _, _, _, _, _, _, _, _, _, _, _ = readFile()


	for l, valBegin := range potCommonBusBegin {

		var a = sumSlice(lenCommonCoords[:l])
		var b = sumSlice(lenCommonCoords[:l+1])

		var nastyaBegin = potCommonIndex[a:b]
		for _, m := range nastyaBegin {
			for u, _ := range alat {
				for i, k := range stopBeginlat {
					lon, _ := strconv.ParseFloat(alon[u], 8)
					lat, _ := strconv.ParseFloat(alat[u], 8)

					if route_id[u] == valBegin && Distance(lat, lon, k, stopBeginlon[i]) < 0.1{
						//lat == k && lon == stopBeginlon[i] {

						/*
							for _, j := range route_id{
								if j == route_id[u] && direction[u] == "Прямое"{
									d = "direct"
								}
							}

						*/
						if d == "direct"{
							//if m > u || m < u {
							if m < u{
								beginOk = append(beginOk, m)
							}
						} else {
							if m < u {
								beginOk = append(beginOk, m)
							}
						}
					}
				}

			}
		}
	}




	for l, valEnd := range potCommonBusEnd {

		var a = sumSlice(lenCommonCoords[:l])
		var b = sumSlice(lenCommonCoords[:l+1])
		var nastyaEnd = potCommonIndex[a:b]
		for _, m := range nastyaEnd {
			for u, _ := range alat {
				for i, k := range stopEndlat {
					lon, _ := strconv.ParseFloat(alon[u], 8)
					lat, _ := strconv.ParseFloat(alat[u], 8)
					if route_id[u] == valEnd && Distance(lat, lon, k, stopEndlon[i]) < 0.1 {
						//lat == k && lon == stopEndlon[i] {
						/*
							for _, j := range route_id{
								if j == route_id[u] && direction[u] == "Прямое"{
									d = "direct"
								}
							}

						*/
						if d == "direct"{
							//if m > u || m < u {
							if m > u {
								endOk = append(endOk, m)
							}
						} else {
							if m > u {
								endOk = append(endOk, m)
							}
						}
					}
				}

			}
		}
	}



	beginOk = uniqueInt(beginOk)
	endOk = uniqueInt(endOk)

	var ok []int
	for _, j := range beginOk{
		for _,k := range endOk{
			if j == k{
				ok = append(ok, j)
			}
		}
	}

	ok = uniqueInt(ok)



	var OkOk []int
	var routeok [] string
	for _, s := range ok{
		if stringInSlice(route_id[s], routeok){
		} else{
			routeok = append(routeok, route_id[s])
			OkOk = append(OkOk, s)
		}
	}





	return OkOk

}

func findPotentialBuses(nearBusBegin, nearBusEnd []string)([]string,[]int,[]int, []string, []string){
	var _, _, route_id, _, _, _, _, _, _, _, coords, _, _, _= readFile()
	var commonBusBegin, commonBusEnd []string
	//var buses []string
	var commonCoords []string
	var lenCommonCoords, commonIndex []int
	// var begCoords, endCoords []string



	for _, busBegin := range nearBusBegin {
		for _, busEnd := range nearBusEnd {
			var a []int
			var b []int
			var begCoords, endCoords []string
			var commonCoords1 [] string
			for i, route := range route_id {
				if busBegin == route {
					a = append(a, i)
					begCoords = append(begCoords, coords[i])

				}
				if busEnd == route {
					b = append(b, i)
					endCoords = append(endCoords, coords[i])

				}
			}

			for n, val1 := range begCoords{
				for _, val2 := range endCoords{
					var lat1, lon1, lat2, lon2 float64
					_, _ = fmt.Sscanf(val1, "%b,%b", &lat1, &lon1)
					_, _ = fmt.Sscanf(val2, "%b,%b", &lat2, &lon2)


					if Distance(lat1, lon1,lat2, lon2) < 0.1{
						commonCoords = append(commonCoords, val1)
						commonIndex = append(commonIndex, a[n])
						commonCoords1 = append(commonCoords1, val1)
					}
				}
			}


			lenCommonCoords = append(lenCommonCoords, len(commonCoords1))
			commonBusBegin = append(commonBusBegin,busBegin)
			commonBusEnd = append(commonBusEnd,busEnd)

		}
	}


	return commonCoords,lenCommonCoords,commonIndex, commonBusBegin, commonBusEnd
}




func printAllAboutBus(index []int) ([]string,[]string, []string, string, string, string, string, []string  ) {


	var alon, alat,_,stop_distance, _, route_short_name, transport_type,_, _, stop_name, coordinates,_, _, _ = readFile()

	var begin int = index[1]
	var end int = index[0]
	var lonPath [] string = alon[begin:(end + 1)]
	var latPath [] string = alat[begin:(end+1)]

	var stop_distancePath [] string = stop_distance[begin:(end+1)]

	var route_short_namePath [] string = route_short_name[begin:(end+1)]
	var transport_typePath [] string = transport_type[begin:(end+1)]

	var stop_namePath [] string = stop_name[begin:(end+1)]
	var coordinatesPath [] string = coordinates[begin:(end+1)]

	return lonPath,latPath, stop_distancePath, route_short_namePath[0], transport_typePath[0], stop_namePath[0], stop_namePath[len(stop_namePath) - 1], coordinatesPath



}




func bestNumber(stopBeginlat, stopBeginlon, stopEndlat, stopEndlon []float64,nearBusBegin, nearBusEnd []string,suitableBuses []string, Lat0,Lon0, Lat1, Lon1  float64, timeToStopByIndexBegin, timeToStopByIndexEnd map[int]float64, coef float64) (float64, []int) {
	var alon, alat, route_id, _, direction, _, _, _, _, _, _, _, timeBetweenStop, _ = readFile()



	var timeBus float64 = 1e+22
	var indexBest [] int
	for _, buses := range suitableBuses {

		//var minDist int = 10000000
		//var suitNumber []string
		var suitIndexDirect []int
		var suitIndexReverse []int
		//var busPath [] string

		//var rev string
		// так у нас есть индексы прямого движения автобуа и обратного
		for j, val := range route_id {
			if buses == val {
				if direction[j] == "Прямое" {
					suitIndexDirect = append(suitIndexDirect, j)
				} else {
					suitIndexReverse = append(suitIndexReverse, j)
				}

				//if direction[j] == "Обратное"{
				//	rev = "t"
				//}
			}
		}

		var indexBeginDir []int
		var indexEndDir []int
		var indexBeginRev []int
		var indexEndRev []int
		var findBestRev []int
		var findBestDir []int

		for _, index := range suitIndexDirect {
			lon, _ := strconv.ParseFloat(alon[index], 8)
			lat, _ := strconv.ParseFloat(alat[index], 8)
			for k := range stopBeginlon {
				if lon == stopBeginlon[k] && lat == stopBeginlat[k] {
					//var indexBeginDir1 []int
					indexBeginDir = append(indexBeginDir, index)
					//indexBeginDir  = indexBeginDir1

				}
			}

			for m := range stopEndlon {
				if lon == stopEndlon[m] && lat == stopEndlat[m] {
					//var indexEndDir1 []int
					indexEndDir = append(indexEndDir, index)
					//indexEndDir  = indexEndDir1

				}
			}
		}

		// итак у нас есть список из индексов пути автобусов и индексы их начала и конца

		if len(indexBeginDir) > 0 && len(indexEndDir) > 0 {
			if indexBeginDir[0] < indexEndDir[0] {
				var min, _ = findMinAndMax(indexBeginDir)
				var _, max = findMinAndMax(indexEndDir)
				findBestDir = append(findBestDir, min, max)
			} //else {
			//	if rev != "t"{
			//		var min, _ = findMinAndMax(indexEndDir)
			//		var _, max = findMinAndMax(indexBeginDir)
			//		findBestDir = append(findBestDir, max, min)
			//	}
			//}
		}

		for _, index := range suitIndexReverse {
			lon, _ := strconv.ParseFloat(alon[index], 8)
			lat, _ := strconv.ParseFloat(alat[index], 8)
			for k := range stopBeginlon {

				if lon == stopBeginlon[k] && lat == stopBeginlat[k] {
					indexBeginRev = append(indexBeginRev, index)
				}
			}

			for m := range stopEndlon {
				if lon == stopEndlon[m] && lat == stopEndlat[m] {
					indexEndRev = append(indexEndRev, index)
				}
			}

		}

		if len(indexBeginRev) > 0 && len(indexEndRev) > 0 {
			if indexBeginRev[0] < indexEndRev[0] {
				var min, _ = findMinAndMax(indexBeginRev)
				var _, max = findMinAndMax(indexEndRev)
				findBestRev = append(findBestRev, min, max)
			}
		}





		if len(findBestDir) > 0 {


			timeToStop := timeToStopByIndexBegin[findBestDir[0]] * coef
			timeFromStop := timeToStopByIndexEnd[findBestDir[1]] * coef



			//fmt.Println(Lat0, Lon0, latStop1, lonStop1, latStop2, lonStop2, Lat1, Lon1)

			//if sumSliceString(timeBetweenStop[findBestDir[0]:findBestDir[1]]) / 60  < timeBus {
			if sumSliceString(timeBetweenStop[findBestDir[0]:findBestDir[1]])/60 + timeToStop + timeFromStop < timeBus {



				timeBus = sumSliceString(timeBetweenStop[findBestDir[0]:findBestDir[1]])/60 + timeToStop + timeFromStop

				//timeBus = sumSliceString(timeBetweenStop[findBestDir[0]:findBestDir[1]]) / 60
				indexBest = findBestDir

			}
		}

		if len(findBestRev) > 0 {


			//fmt.Println(Lat0, Lon0, latStop1, lonStop1, latStop2, lonStop2, Lat1, Lon1)


			timeToStop := timeToStopByIndexBegin[findBestRev[0]] * coef
			timeFromStop := timeToStopByIndexEnd[findBestRev[1]] * coef

			if sumSliceString(timeBetweenStop[findBestRev[0]:findBestRev[1]])/60+timeToStop+timeFromStop < timeBus {
				//if sumSliceString(timeBetweenStop[findBestRev[0]:findBestRev[1]]) / 60 < timeBus{
				//fmt.Println("timeToStop", timeToStop,timeFromStop, sumSliceString(timeBetweenStop[findBestRev[0]:findBestRev[1]])/ 60)
				timeBus = sumSliceString(timeBetweenStop[findBestRev[0]:findBestRev[1]])/60+ timeToStop+ timeFromStop
				//timeUsual =  sumSliceString(timeBetweenStop[findBestRev[0]:findBestRev[1]])/60 + timeToStop + timeFromStop
				//timeOnlyBus = sumSliceString(timeBetweenStop[findBestRev[0]:findBestRev[1]])/60

				//timeBus = sumSliceString(timeBetweenStop[findBestRev[0]:findBestRev[1]]) / 60
				indexBest = findBestRev
			}
		}



	}



	if len(indexBest) != 0 {




		//timeToStop := timeToStopByIndexBegin[indexBest[0]]
		//timeFromStop := timeToStopByIndexEnd[indexBest[1]]
		//timeOnlyBus := sumSliceString(timeBetweenStop[indexBest[0]:indexBest[1]])/60
		//timeBusTr := timeOnlyBus + timeToStop + timeFromStop
		//fmt.Println(coef,timeToStop  + timeFromStop,timeBusTr, timeOnlyBus, "comfrot", timeBus )
		//timeBus = timeBus
	}



	return timeBus, indexBest

}



func suitableBus (nearBusBegin, nearBusEnd []string) []string {
	var suitableBus []string

	for i := range nearBusBegin {
		for j := range nearBusEnd {
			if nearBusBegin[i] == nearBusEnd[j] {
				suitableBus = append(suitableBus, nearBusBegin[i])
			}
		}
	}

	return suitableBus
}



func nearestBusinStops(stopBeginlat, stopBeginlon, stopEndlat, stopEndlon []float64 )([]string ,[]string){
	var alon, alat,route_id,_, _, _, _,_, _, _, _,_,_, hasStopThisInteval= readFile()



	var nearBusBegin []string
	var nearBusEnd []string
	//alon = alon[1:len(alon)]
	//alat = alat[1:len(alat)]
	//route_id = route_id[1:len(route_id)]

	for i := range alon {
		lon, _ := strconv.ParseFloat(alon[i], 8)
		lat, _ := strconv.ParseFloat(alat[i], 8)

		for j := range stopBeginlon{
			if lon == stopBeginlon[j] && lat == stopBeginlat[j] && hasStopThisInteval[i] == "1.0"{ // условие равенства единицы задает, наличие автобуса в этом временном интевале на останвоке
				// то что
				//fmt.Println(route_id[i],hasStopThisInteval[i] )
				nearBusBegin = append(nearBusBegin, route_id[i])
			}
		}
	}

	for i := range alon {
		lon, _ := strconv.ParseFloat(alon[i], 8)
		lat, _ := strconv.ParseFloat(alat[i], 8)

		for j := range stopEndlon {

			if lon == stopEndlon[j] && lat == stopEndlat[j] {
				nearBusEnd = append(nearBusEnd, route_id[i])

			}
		}
	}

	var busBegin = uniqueString(nearBusBegin)
	var busEnd = uniqueString(nearBusEnd)



	return busBegin, busEnd
}





func nearest_stops(beginlat,beginlon, endlat, endlon, dist float64) ([]float64,[]float64,[]float64,[]float64, []int, []int) {
	var alon, alng,_,_, _, _, _,_, _, _, _,_, _,_= readFile()



	var nearBeginlat, nearBeginlon, nearEndlat, nearEndlon []float64

	nearBeginIndexes := findWithTree(beginlat,beginlon, dist)
	nearEndIndexes:= findWithTree(endlat,endlon, dist)


	for _, index:= range nearBeginIndexes{

		lon1, _ := strconv.ParseFloat(alon[index], 8)
		lat1, _ := strconv.ParseFloat(alng[index], 8)
		nearBeginlat = append(nearBeginlat, lat1)
		nearBeginlon = append(nearBeginlon, lon1)
	}

	for _, index:= range nearEndIndexes{

		lon1, _ := strconv.ParseFloat(alon[index], 8)
		lat1, _ := strconv.ParseFloat(alng[index], 8)

		nearEndlat = append(nearEndlat, lat1)
		nearEndlon = append(nearEndlon, lon1)
	}



	var nearBeginlat1, nearBeginlon1 [] float64
	for a, b := range nearBeginlat{
		if floatInSlice(b, nearBeginlat1){
		}else{
			nearBeginlat1 = append(nearBeginlat1, b)
			nearBeginlon1 = append(nearBeginlon1,  nearBeginlon[a])
		}
	}


	var nearEndlat1, nearEndlon1 [] float64
	for a, b := range nearEndlat{
		if floatInSlice(b, nearEndlat1) && floatInSlice(nearEndlon[a], nearEndlon1)  {
		} else {
			nearEndlat1 = append(nearEndlat1, b)
			nearEndlon1 = append(nearEndlon1,  nearEndlon[a])
		}
	}



	//nearBeginlat1 := nearBeginlat
	//nearBeginlon1 := nearBeginlon

	//nearEndlat1 := nearEndlat
	//nearEndlon1 := nearEndlon

	return nearBeginlat1, nearBeginlon1, nearEndlat1, nearEndlon1, nearBeginIndexes, nearEndIndexes

}



func nearest_stops_brute(beginlat,beginlon, dist float64) ( []int) {
	var alon, alng,_,_, _, _, _,_, _, _, _,_,_, _ = readFile()


	//alon = alon[1:len(alon)]
	//
	//alng = alng[1:len(alng)]
	//fmt.Println(dist)
	var distBegin float64

	var nearBeginlat []float64
	var nearBeginlng []float64

	var nearBeginIndexes []int

	for i := 0; i < len(alon); i++ {
		var blon = alon[i]
		var blng = alng[i]
		lat1, _ := strconv.ParseFloat(blng, 8)
		lon1, _ := strconv.ParseFloat(blon, 8)
		distBegin = Distance(lon1,lat1,beginlat,beginlon)


		if distBegin <= dist {
			nearBeginlat = append(nearBeginlat, lat1)
			nearBeginlng = append(nearBeginlng, lon1)
			nearBeginIndexes = append(nearBeginIndexes, i)
		}

	}



	//nearBeginlat1 := nearBeginlat
	//nearBeginlon1 := nearBeginlon

	//nearEndlat1 := nearEndlat
	//nearEndlon1 := nearEndlon

	return  nearBeginIndexes

}


func ReadCsv(filename string) ([][]string, error) {

	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}





func unique(floatSlice [] float64) []float64 {
	keys := make(map[float64]bool)
	list := []float64{}
	for _, entry := range floatSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func uniqueInt(intSlice [] int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func uniqueString(stringSlice [] string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}


func findMinAndMax(a []int) (min int, max int) {
	min = a[0]
	max = a[0]
	for _, value := range a {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}
	return min, max
}


func findMinAndMaxFloat(a []float64) (min float64) {
	min = 10000000000
	for _, value := range a {

		if value < min && value > 0 {
			min = value
		}
	}
	return min
}

func readFile()([]string,[]string, []string, []string,[]string,[]string, []string, []string,[]string,[]string, []string, []string, []string, []string ){
	//lines, err := ReadCsv("2020-10-22.08:00-08:10")
	fileName := fmt.Sprintf("points/traffic/%s/%s:%s" , weekDay, fileHour, fileMinute)

	//fmt.Println(fileName)
	lines, err := ReadCsv(fileName)
	if err != nil {
		panic(err)
	}
	var alat []string
	var alon []string
	var route_id []string
	var stop_distance []string
	var direction, route_short_name,transport_type,stop_id,next_stop,stop_name,coordinates, number, timeBetweenStop, hasStopThisInteval []string

	for _, line := range lines {
		data := CsvLine{
			number:           line[0],
			route_id:         line[1],
			route_short_name: line[2],
			route_long_name:  line[3],
			transport_type:   line[4],
			direction:        line[5],
			stop_id:          line[6],
			next_stop:        line[7],
			stop_distance:    line[8],
			stop_name:        line[9],
			coordinates:      line[10],
			lon:              line[11],
			lng:              line[12],
			timeBetweenStop:  line[13],
			hasStopThisInteval: line[14],
		}

		var lon = data.lon
		var lat = data.lng
		var aroute_id = data.route_id
		var astop_distance = data.stop_distance
		var adirection = data.direction
		var aroute_short_name = data.route_short_name
		var atransport_type = data.transport_type
		var astop_id = data.stop_id
		var anext_stop = data.next_stop
		var astop_name = data.stop_name
		var acoordinates = data.coordinates
		var anumber = data.number
		var atimeBetweenStop = data.timeBetweenStop
		var ahasStopThisInteval = data.hasStopThisInteval

		alon = append(alon, lon)
		alat = append(alat, lat)
		route_id = append(route_id, aroute_id)
		stop_distance = append(stop_distance, astop_distance)
		direction = append(direction,adirection)
		route_short_name = append(route_short_name, aroute_short_name)
		transport_type = append(transport_type, atransport_type)
		stop_id = append(stop_id, astop_id)
		next_stop = append(next_stop, anext_stop)
		stop_name = append(stop_name, astop_name)
		coordinates = append(coordinates,acoordinates)
		number = append(number, anumber)
		timeBetweenStop = append(timeBetweenStop, atimeBetweenStop)
		hasStopThisInteval = append(hasStopThisInteval, ahasStopThisInteval)
	}

	alon = alon[1:len(alon)]
	alat = alat[1:len(alat)]
	route_id = route_id[1:len(route_id)]
	stop_distance = stop_distance[1:len(stop_distance)]
	direction = direction[1:len(direction)]
	route_short_name = route_short_name[1:len(route_short_name)]
	transport_type = transport_type[1:len(transport_type)]
	stop_id = stop_id[1:len(stop_id)]
	next_stop = next_stop[1:len(next_stop)]
	stop_name = stop_name[1:len(stop_name)]
	coordinates = coordinates[1:len(coordinates)]
	number = number[1:len(number)]
	timeBetweenStop = timeBetweenStop[1:len(timeBetweenStop)]
	hasStopThisInteval = hasStopThisInteval[1:len(hasStopThisInteval)]

	return alon, alat, route_id,stop_distance, direction, route_short_name, transport_type,stop_id, next_stop, stop_name, coordinates, number, timeBetweenStop, hasStopThisInteval

}

/*
func intersection(a []int, b []int) (inter []int) {
	low, high := a, b
	if len(a) > len(b) {
		low = b
		high = a

	}

	done := false
	for i, l := range low {
		for j, h := range high {
			f1 := i + 1
			f2 := j + 1
			if l == h {
				inter = append(inter, h)
				if f1 < len(low) && f2 < len(high) {
					if low[f1] != high[f2] {
						done = true
					}
				}
				high = high[:j+copy(high[j:], high[j+1:])]
				break
			}
		}
		if done {
			break
		}
	}
	return
}
*/

func inter(arrs ...[]int) []int {
	res := []int{}
	x := arrs[0][0]
	i := 1
	for {
		off := sort.SearchInts(arrs[i], x)
		if off == len(arrs[i]) {
			// we emptied one slice, we're done.
			break
		}
		if arrs[i][off] == x {
			i++
			if i == len(arrs) {
				// x was in all the slices
				res = append(res, x)
				x++ // search for the next possible x.
				i = 0
			}
		} else {
			x = arrs[i][off]
			i = 0 // This can be done a bit more optimally.
		}
	}
	return res
}



func intersection(a []string, b []string) (inter []string) {
	// interacting on the smallest list first can potentailly be faster...but not by much, worse case is the same
	low, high := a, b
	if len(a) > len(b) {
		low = b
		high = a
	}

	done := false
	for i, l := range low {
		for j, h := range high {
			// get future index values
			f1 := i + 1
			f2 := j + 1
			if l == h {
				inter = append(inter, h)
				if f1 < len(low) && f2 < len(high) {
					// if the future values aren't the same then that's the end of the intersection
					if low[f1] != high[f2] {
						done = true
					}
				}
				// we don't want to interate on the entire list everytime, so remove the parts we already looped on will make it faster each pass
				high = high[:j+copy(high[j:], high[j+1:])]
				break
			}
		}
		// nothing in the future so we are done
		if done {
			break
		}
	}
	return
}

/* !
func Round(x, unit float64) float64 {
	return (math.Floor(x * unit) /unit)
}

*/

func floatInSlice(a float64, list []float64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sumSlice(num []int) int {
	var sum int

	for _, j := range num{
		sum += j
	}

	return sum

}

func sumSliceString(num []string) float64 {
	var sum float64

	for _, j := range num{

		jj, _ := strconv.ParseFloat(j, 8)

		sum += jj
	}

	return sum

}


func commonTime(min float64) float64{
	var time float64 = min / 60
	return time
}



