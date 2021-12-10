package points

import (
	"errors"
	"fmt"
	"strconv"
)



func bus11(beginlat, beginlon, endlat, endlon, dist , coef float64)(error, float64, [][]int) {
	var _, _, _, _, _, _, _, _, _, _, _, _, timeBetweenBuses, hasStopThisInteval = readFile()

	var _, _, _, _ , nearBeginIndexes, nearEndIndexes  = nearest_stops(beginlat, beginlon, endlat, endlon, dist)
	min, _, nearBusBegin, nearBusEnd,lenSuitableBuses, _ ,timeToStopByIndexBegin, timeToStopByIndexEnd:= directBus(beginlat, beginlon, endlat, endlon, dist, coef)

	min = 100000000000000.0
	var bestIndex [][]int




	if min < 1 || min > 500000000 || lenSuitableBuses < 1 {

		indexesBusBegin0 := indexesOfParticulaBuses(nearBusBegin)
		indexesBusEnd0 := indexesOfParticulaBuses(nearBusEnd)


		indexesBusBegin := cutListOsStopsForBus(indexesBusBegin0, nearBeginIndexes, "begin")
		indexesBusEnd := cutListOsStopsForBus(indexesBusEnd0, nearEndIndexes, "end")


		var crossStopsIndex [][]int

		//fmt.Println(indexesBusBegin)
		//fmt.Println(indexesBusEnd)


		fmt.Println(crossStopsIndex)
		for begin, beginStopsAll := range indexesBusBegin {
			for _, beginStopsIndex := range beginStopsAll {
				for end, endStopsAll := range indexesBusEnd {
					for _, endStopsIndex := range endStopsAll {


						if stop_id[beginStopsIndex] == stop_id[endStopsIndex] && hasStopThisInteval[beginStopsIndex] == "1.0" && hasStopThisInteval[begin] == "1.0"{

							crossStopsIndex = append(crossStopsIndex,[]int{beginStopsIndex, endStopsIndex,begin, end})

						}


					}
				}
			}
		}



		min := 100000000000000.0
		var bestIndex [][]int



		for _, transfer := range crossStopsIndex{
			//fmt.Println(indexesBusBegin[transfer[1]])


			transferCoordsLatCross0, _ := strconv.ParseFloat(alat[transfer[0]], 8)
			transferCoordsLonCross0, _ := strconv.ParseFloat(alon[transfer[0]], 8)
			transferCoordsLatCross1, _ := strconv.ParseFloat(alat[transfer[1]], 8)
			transferCoordsLonCross1, _ := strconv.ParseFloat(alon[transfer[1]], 8)



			timeBetweenBusesInTransfer:= Distance(transferCoordsLatCross0, transferCoordsLonCross0, transferCoordsLatCross1, transferCoordsLonCross1) / 5 * 60 * coef
			timeToBegin := timeToStopByIndexBegin[indexesBusBegin[transfer[2]][0]] * coef
			timeFromEnd := timeToStopByIndexEnd[indexesBusEnd[transfer[3]][len(indexesBusEnd[transfer[3]])-1]] * coef

			if min > timeBetweenBusesInTransfer +timeToBegin +timeFromEnd +  sumSliceString(timeBetweenBuses[indexesBusBegin[transfer[2]][0] : transfer[0]]) / 60  +   sumSliceString(timeBetweenBuses[transfer[1] : indexesBusEnd[transfer[3]][len(indexesBusEnd[transfer[3]])-1]]) / 60 && stop_id[transfer[0]] == stop_id[transfer[1]]  {
				min =  timeBetweenBusesInTransfer + timeToBegin +timeFromEnd + sumSliceString(timeBetweenBuses[indexesBusBegin[transfer[2]][0] : transfer[0]]) / 60 +   sumSliceString(timeBetweenBuses[transfer[1] : indexesBusEnd[transfer[3]][len(indexesBusEnd[transfer[3]])-1]]) / 60
				bestIndex = [][]int{{indexesBusBegin[transfer[2]][0],transfer[0]},{transfer[1],indexesBusEnd[transfer[3]][len(indexesBusEnd[transfer[3]])-1]}}

			}


		}
		// для получения индексов: первое число в crossStopsIndex: это место первой пересадки, второе число - место второй пересадки, третье число: порядковый номер автобуса в последовательности
		// indexesBusBegin с первой остановки до остановки первой пересадки, 4 число - порядковый номер автобуса в переменной indexesBusEnd с остановки второй пересадки до остановки конечной


		fmt.Println(bestIndex)

		if len(crossStopsIndex) < 1 || min > 100000 {
			return errors.New("try find path with two transfer"), 0, bestIndex
		} else{

			transferCoordsLat0, _ := strconv.ParseFloat(alat[bestIndex[0][0]], 8)
			transferCoordsLon0, _ := strconv.ParseFloat(alon[bestIndex[0][0]], 8)
			transferCoordsLat1, _ := strconv.ParseFloat(alat[bestIndex[1][1]], 8)
			transferCoordsLon1, _ := strconv.ParseFloat(alon[bestIndex[1][1]], 8)




			transferCoordsLatCross0, _ := strconv.ParseFloat(alat[bestIndex[0][1]], 8)
			transferCoordsLonCross0, _ := strconv.ParseFloat(alon[bestIndex[0][1]], 8)
			transferCoordsLatCross1, _ := strconv.ParseFloat(alat[bestIndex[1][0]], 8)
			transferCoordsLonCross1, _ := strconv.ParseFloat(alon[bestIndex[1][0]], 8)

			_,_, timeBetweenBusesInTransfer:= byFoot(transferCoordsLatCross0, transferCoordsLonCross0, transferCoordsLatCross1, transferCoordsLonCross1)
			_, _, timeToBegin := byFoot(beginlat, beginlon, transferCoordsLon0, transferCoordsLat0)
			_, _, timeFromEnd := byFoot(transferCoordsLon1, transferCoordsLat1, endlat, endlon)


			min1 := timeToBegin + timeFromEnd +  timeBetweenBusesInTransfer + sumSliceString(timeBetweenBuses[bestIndex[0][0] : bestIndex[0][1]]) / 60 +   sumSliceString(timeBetweenBuses[bestIndex[1][0] : bestIndex[1][1]]) / 60

			//fmt.Println(coef, min1, bestIndex, "пеший путь")

			return nil, min1, bestIndex
		}

	}



	return nil, min, bestIndex

}



func cutListOsStopsForBus(indexesBus [][]int, nearIndexes []int, fromWhatSite string) [][]int {
	var _, _, _, _, direction, _, _, _, _, _, _, _, _,_ = readFile()

	var cutIndex [][] int

	if fromWhatSite == "begin" {
		for _, listIndexBuses := range indexesBus {

			for j, listIndexBus := range listIndexBuses { // индексы всех проходящих вначале автобусов
				for _, stopsIndex := range nearIndexes { // индексы всех остановок призлежащих
					if listIndexBus == stopsIndex && direction[listIndexBus] == direction[stopsIndex] {
						cutIndex = append(cutIndex, listIndexBuses[j:])
					}
				}
			}
		}
	}else{
		for _, listIndexBuses := range indexesBus {

			for j, listIndexBus := range listIndexBuses { // индексы всех проходящих вконце автобусов
				for _, stopsIndex := range nearIndexes { // индексы всех остановок призлежащих
					if listIndexBus == stopsIndex && direction[listIndexBus] == direction[stopsIndex] {
						cutIndex = append(cutIndex, listIndexBuses[:j])
					}
				}
			}
		}
	}


	return cutIndex
}

func indexesOfParticulaBuses(buses []string) [][]int {
	var _, _, route_id, _, _, _, _, _, _, _, _, _, _, _ = readFile()
	var indexesBuses [][] int
	for i, index := range buses{
		indexesBuses = append(indexesBuses, []int{})
		for j, bus := range route_id {
			if index == bus {
				indexesBuses[i] = append(indexesBuses[i], j)
			}
		}
	}

	return indexesBuses
}
