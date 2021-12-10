package points

import (
	//"awesomeProject/fops/points"
	"fmt"
	"strconv"

	//"strconv"

	//"strconv"
)

var alon, alat, route_id, stop_distance, direction, route_short_name, transport_type, stop_id, next_stop, stop_name, coordinates, number, timeBetweenBuses, hasStopThisInteval = readFile()


type StopBusStop struct{

	beginTrStop string
	bus string
	endTrStop string

}

type StopBusStopIndex struct{

	beginTrStop int
	bus string
	endTrStop int

}


type AllStopsAndTransfers struct{
	stopBegin int

	beginTrStop int

	endTrStop int

	stopEnd int
}

func bus2(Lat0,Lon0, Lat1, Lon1, coef float64) (float64, [][]int){

	min, index, _, _, _ , error ,  _, _:= directBus(Lat0,Lon0, Lat1, Lon1, 1, 1)
	var index1 [][]int
	index1 = append(index1, index)

	timeBest, indexBest := min, index1





	if error != nil {

		error, min, index := bus11(Lat0, Lon0, Lat1, Lon1, 3, 1)
		timeBest, indexBest = min, index

		if error != nil {

			_, _, nearBusBegin, nearBusEnd, _ ,_, timeToStopByIndexBegin, timeToStopByIndexEnd:= directBus(Lat0,Lon0, Lat1, Lon1, 1.5, 1)


			_, _, _, _ , nearStopBegin, nearStopEnd := nearest_stops(Lat0,Lon0, Lat1, Lon1, 3)


			var commonBus []string
			var stopTransferBegin, stopTransferEnd []string

			var a = 1.0
			for part2 := 0.1; part2 < 1; part2 = part2 + 0.1 {
				for part1 := 0.1; part1 < a; part1 = part1 + 0.1 {
					// смотрим все автобусы, которые ходят в местах наибольших скоплений автобусов, до которых ходят автобусы из начальной точки
					largestStopBegin, busInLargestStopsBegin := tryWithSmallBegin(nearBusBegin, 0, 10, Lat0,Lon0, Lat1, Lon1, part1)
					// смотрим все автобусы, которые ходят в местах наибольших скоплений автобусов, от которых ходят автобусы до конечной точки
					largestStopEnd, busInLargestStopsEnd := tryWithSmallEnd(nearBusEnd, 0, 10, Lat1, Lon1, Lat0,Lon0, part2)

					// либо просто менять либо изменять также начальные точки
					/*
						for k:= range busInLargestStopsEnd{
							k = k
							stopTransferEnd = append(stopTransferEnd, []string{})
						}

					*/

					for k1, listBus1 := range busInLargestStopsBegin {
						for k2, listBus2 := range busInLargestStopsEnd {
							for _, bus1 := range listBus1 {
								for _, bus2 := range listBus2 {

									if bus1 == bus2 {
										commonBus = append(commonBus, bus1)


										stopTransferBegin= append(stopTransferBegin, largestStopBegin[k1])

										stopTransferEnd = append(stopTransferEnd, largestStopEnd[k2])

										//fmt.Println("eeeyyyyyy", bus1, bus2)
									}
								}
							}

						}
						a = a - 0.1
					}
				}
			}





			var setVariantsBusAndTransfers []StopBusStop

			//fmt.Println("busEnd",busEnd)
			for i, l:= range commonBus{
				var current StopBusStop
				current.beginTrStop = stopTransferBegin[i]
				current.endTrStop = stopTransferEnd[i]
				current.bus = l

				setVariantsBusAndTransfers = append(setVariantsBusAndTransfers, current)
			}







			var indexBuses [][]int
			for j, h := range setVariantsBusAndTransfers {
				indexBuses = append(indexBuses, []int{})
				for i, k := range route_id {
					if k == h.bus {
						indexBuses[j] = append(indexBuses[j], i)
					}
				}
			}



			// для
			var stopsIndexes []StopBusStopIndex

			var beg, end int
			for h, h1 := range setVariantsBusAndTransfers{
				stopsIndexes = append(stopsIndexes, StopBusStopIndex{})

				for _, e := range indexBuses[h]{

					if stop_id[e] == h1.beginTrStop{
						beg = e
					}
					if stop_id[e] == h1.endTrStop{
						end = e
					}

					if beg < end{
						stopsIndexes[h].bus = h1.bus
						stopsIndexes[h].beginTrStop = beg
						stopsIndexes[h].endTrStop = end
					} else{
						stopsIndexes[h].bus = "0"
						stopsIndexes[h].beginTrStop = 0
						stopsIndexes[h].endTrStop = 0

					}
				}


			}


			//fmt.Println(stopsIndexes)


			//fmt.Println(uniqueSetVariantsBusAndTransfers)

			uniqueSetVariantsBusAndTransfersFilters := uniqueNumberOfVarIndex(stopsIndexes)
			//fmt.Println("uniqueSetVariantsBusAndTransfersFilters", uniqueSetVariantsBusAndTransfersFilters)



			// !!!!!!!!сделать предпроцессинг рассчиатть для уникальных значений directbus

			//now0 := makeTimestamp()

			var uniqueSetVariantsBusAndTransfersStopsID []StopBusStop

			for _,l := range uniqueSetVariantsBusAndTransfersFilters{
				var w StopBusStop
				w.bus = l.bus
				w.beginTrStop = stop_id[l.beginTrStop]
				w.endTrStop = stop_id[l.endTrStop]
				uniqueSetVariantsBusAndTransfersStopsID = append(uniqueSetVariantsBusAndTransfersStopsID, w)
			}


			timeUniqueAllBeginTransfersStop, indexBestUniqueAllBeginTransfersStop0, _, _:= preprocessingDataForFindDirectBusesToTransfers(uniqueSetVariantsBusAndTransfersStopsID, Lat0, Lon0 ,1,nearStopBegin, nearStopEnd, coef,  timeToStopByIndexBegin, timeToStopByIndexEnd)
			timeUniqueAllEndTransfersStop, indexBestUniqueAllEndTransfersStop0, _, _ := preprocessingDataForFindDirectBusesToTransfers(uniqueSetVariantsBusAndTransfersStopsID, Lat1, Lon1,2, nearStopBegin, nearStopEnd, coef,  timeToStopByIndexBegin, timeToStopByIndexEnd)





			var time11, time12, time13 float64
			var index11, index12, index13 []int
			var indexBestOfTheBest [][] int

			commonTime := 10000000000.0

			for _, k := range uniqueSetVariantsBusAndTransfersFilters {

				//if k.beginTrStop != 0 && k.endTrStop != 0 {


				//fmt.Println(k)

				for q2, w2 := range indexBestUniqueAllEndTransfersStop0 {
					for q1, w1 := range indexBestUniqueAllBeginTransfersStop0 {



						if stop_id[k.beginTrStop] == stop_id[w1[1]] && stop_id[k.endTrStop] == stop_id[w2[0]] && hasStopThisInteval[w1[0]] == "1.0" && hasStopThisInteval[w2[0]] == "1.0" && hasStopThisInteval[k.beginTrStop] == "1.0"{
							//fmt.Println("TRUE")
							index11 = w1
							time11 = timeUniqueAllBeginTransfersStop[q1] / 60
							index12 = w2
							time12 = timeUniqueAllEndTransfersStop[q2] / 60

						}
					}
				}


				time13 = sumSliceString(timeBetweenBuses[k.beginTrStop:k.endTrStop])/60

				index13 = []int{k.beginTrStop, k.endTrStop}

				indexBest := [][]int{index11, index12, index13}




				//fmt.Println(stop_id[indexBest[0][1]]==stop_id[indexBest[2][0]], stop_id[indexBest[2][1]]== stop_id[indexBest[1][0]])

				cTime:= time11+time12+time13
				if commonTime > cTime && time11 != 0 && time12 != 0 && time13 != 0 {
					commonTime = time11+time12+time13
					indexBestOfTheBest = indexBest
				}

				//}
			}

			transferCoordsLat0, _ := strconv.ParseFloat(alat[indexBestOfTheBest[0][0]], 8)
			transferCoordsLon0, _ := strconv.ParseFloat(alon[indexBestOfTheBest[0][0]], 8)
			transferCoordsLat1, _ := strconv.ParseFloat(alat[indexBestOfTheBest[1][1]], 8)
			transferCoordsLon1, _ := strconv.ParseFloat(alon[indexBestOfTheBest[1][1]], 8)

			_, _, timeToBegin := byFoot(Lat0, Lon0, transferCoordsLon0, transferCoordsLat0)
			_, _, timeFromEnd := byFoot(transferCoordsLon1, transferCoordsLat1, Lat1, Lon1)


			timeBest = commonTime + timeToBegin + timeFromEnd
			indexBest = indexBestOfTheBest


			return timeBest, indexBest

		}
	} else {
		fmt.Println(min, index)
		return timeBest, indexBest
	}

	return timeBest, indexBest
}



func preprocessingDataForFindDirectBusesToTransfers(uniqueSetVariantsBusAndTransfers []StopBusStop, Lat, Lon float64, path int, nearStopBegin []int, nearStopEnd []int, coef float64, timeToStopByIndexBegin, timeToStopByIndexEnd map[int]float64) ([]float64, [][]int, []int, []int){




	//здесь Begin также для End можно удалить это из названия
	var allTransfersStop []string


	for _, j := range uniqueSetVariantsBusAndTransfers{
		if path == 1 {
			allTransfersStop = append(allTransfersStop, j.beginTrStop)
		}else{
			allTransfersStop = append(allTransfersStop, j.endTrStop)
		}
	}

	uniqueAllTransfersStop := uniqueString(allTransfersStop)

	var timeUniqueAllTransfersStop []float64
	var indexBestUniqueAllTransfersStop [][]int

	var indexUniqueAllTransfersStop []int


	// все возможные индексы остановок для пересадки
	for _, l := range uniqueAllTransfersStop{
		for q, q1 := range stop_id {
			if q1 == l {
				indexUniqueAllTransfersStop = append(indexUniqueAllTransfersStop, q)
			}
		}
	}
	indexUniqueAllTransfersStop = uniqueInt(indexUniqueAllTransfersStop)

	var indexesIniIndexUniqueAllBeginTransfersStopForIndex []int

	//var indexesPotBegin [][]int

	// то есть у насс получаются все воззможные варианты соединть начальную остановку с останвокй трансфера
	if path == 1 {
		for _, k := range indexUniqueAllTransfersStop {
			for _, k1 := range nearStopBegin {
				if route_id[k] == route_id[k1] && k1 < k && direction[k] == direction[k1] {

					timeByFoot := timeToStopByIndexBegin[k]






					indexBestUniqueAllTransfersStop = append(indexBestUniqueAllTransfersStop, []int{k1, k})

					timeUniqueAllTransfersStop = append(timeUniqueAllTransfersStop, sumSliceString(timeBetweenBuses[k1:k]) + timeByFoot * coef)
					indexesIniIndexUniqueAllBeginTransfersStopForIndex = append(indexesIniIndexUniqueAllBeginTransfersStopForIndex, k)
				}
			}
		}
	} else {
		for _, k2 := range indexUniqueAllTransfersStop {
			for _, k3 := range nearStopEnd {
				if route_id[k2] == route_id[k3] && k2 < k3 && direction[k2] == direction[k3]  {


					timeByFoot := timeToStopByIndexEnd[k2]

					indexBestUniqueAllTransfersStop = append(indexBestUniqueAllTransfersStop, []int{k2, k3})

					timeUniqueAllTransfersStop = append(timeUniqueAllTransfersStop, sumSliceString(timeBetweenBuses[k2:k3]) + timeByFoot * coef)
					indexesIniIndexUniqueAllBeginTransfersStopForIndex = append(indexesIniIndexUniqueAllBeginTransfersStopForIndex, k2)
				}
			}
		}

	}







	return timeUniqueAllTransfersStop, indexBestUniqueAllTransfersStop, indexesIniIndexUniqueAllBeginTransfersStopForIndex, indexUniqueAllTransfersStop

}

func diff(a []int, b []int) []int {
	// Turn b into a map
	var m map[int]bool
	m = make(map[int]bool, len(b))
	for _, s := range b {
		m[s] = false
	}
	// Append values from the longest slice that don't exist in the map
	var diff []int
	for _, s := range a {
		if _, ok := m[s]; !ok {
			diff = append(diff, s)
			continue
		}
		m[s] = true
	}
	// Sort the resulting slice
	return diff
}



func timeForTwoTransferFromBegin(nearBus []string, IndexTransferBegin , indexNearestStop []int) []float64{
	//var _, _, route_id, _, _, _, _, stop_id, _, _, _, _, timeBetweenStops , _= readFile()

	var time1 []float64
	//var a int = len(nearBusBegin)
	var currentIndexBuses [100] [] int
	// начальная остановка

	for i, near_begin_bus := range nearBus{
		time1 = append(time1, 0)
		for index, _ := range stop_id {
			if near_begin_bus == route_id[index]{
				currentIndexBuses[i] = append(currentIndexBuses[i],index)

			}
		}
	}



	for i, _ := range currentIndexBuses{
		for _, l1 := range IndexTransferBegin{
			for _, k2 := range indexNearestStop{
				if intInSlice(l1,currentIndexBuses[i]) && intInSlice(k2,currentIndexBuses[i]){
					if k2 < l1 {
						time1[i] = sumSliceString(timeBetweenBuses[k2:l1])
					} else {
						time1[i] = sumSliceString(timeBetweenBuses[l1:k2])
					}
				}
			}

		}
	}

	return time1
}


func timeForTwoTransferBetween(nearBus []string, nearestStopTransfer1, nearestStopTransfer2 []string) []float64{
	//var _, _, route_id, _, _, _, _, stop_id, _, _, _, _, timeBetweenStops, _ = readFile()

	var time1 []float64
	//var a int = len(nearBusBegin)
	var currentIndexBuses [100] [] int
	// начальная остановка

	for i, near_begin_bus := range nearBus{
		time1 = append(time1, 0)
		for index, _ := range stop_id {
			if near_begin_bus == route_id[index]{
				currentIndexBuses[i] = append(currentIndexBuses[i],index)

			}
		}
	}

	var indexStops1 []int
	for j, k1 := range stop_id {
		for _, o2 := range nearestStopTransfer1{
			if k1 == o2 {
				indexStops1 = append(indexStops1, j)
			}


		}
	}


	var indexStops2 []int
	for j, k1 := range stop_id {
		for _, o2 := range nearestStopTransfer2{
			if k1 == o2 {
				indexStops2 = append(indexStops2, j)
			}


		}
	}




	for i, _ := range currentIndexBuses{
		for _, l1 := range indexStops1{
			for _, k2 := range indexStops2{
				if intInSlice(l1,currentIndexBuses[i]) && intInSlice(k2,currentIndexBuses[i]){
					if k2 < l1 {
						time1[i] = sumSliceString(timeBetweenBuses[k2:l1])
					} else{
						time1[i] = sumSliceString(timeBetweenBuses[l1:k2])
					}
				}
			}

		}
	}

	return time1
}

func tryWithSmallBegin(nearBusBegin []string, minR,maxR, lat1, lon1, lat2, lon2, part float64) ([]string,[][]string) {



	largestStopBegin := stopsInParticularRadius(nearBusBegin, minR, maxR, lat1, lon1, lat2, lon2, part)

	// выбрать 3 остановки с наибольшим количеством и попробовать от этих точек до конечных построить прямой маршрут
	// если нет (или сразу сделать с другого конца), то тоже самое сделать с другого конца и попробвать найти общие автобусы которые ходят на этих
	// если не нашла, то увеличить радиус поиска с 3 до 5 и сделать тоже самое


	return largestStopBegin, busInStopsRadius(largestStopBegin,nearBusBegin)
}

func tryWithSmallEnd(nearBusEnd []string, minR,maxR, lat1, lon1, lat2, lon2, part float64) ([]string,[][]string) {


	//fmt.Println(nearBusEnd)

	largestStopEnd := stopsInParticularRadius(nearBusEnd, minR,maxR, lat1, lon1, lat2, lon2, part)
	//fmt.Println(largestStopBegin)

	//fmt.Println(largestStopEnd)
	//fmt.Println(largestStopBegin)



	// выбрать 3 остановки с наибольшим количеством и попробовать от этих точек до конечных построить прямой маршрут
	// если нет (или сразу сделать с другого конца), то тоже самое сделать с другого конца и попробвать найти общие автобусы которые ходят на этих
	// если не нашла, то увеличить радиус поиска с 3 до 5 и сделать тоже самое



	return 	largestStopEnd, busInStopsRadius(largestStopEnd,nearBusEnd)

}



func busInStopsRadius(largestStop []string,nearBus []string) [][]string{
	//var _, _, route_id,_, _, _, _,stop_id, _, _, _, _, _ , _= readFile()
	var ind [][]string
	for l, k := range largestStop{
		ind = append(ind, []string{})
		for i, h :=range stop_id {
			if h==k{
				ind[l] = append(ind[l], route_id[i])
			}


		}
	}

	return ind
}



func stopsInParticularRadius(nearBus []string , minRadius, maxRadius float64,lat1,lon1,lat2,lon2, part float64) []string{


	var listOfIndex []int
	var listofStops []string
	// попытаться найти ближайшие точки остановок с большим количеством автобусов
	// подсчитать все остановки от текущей до текущей + 5, какое кол-во автобусов их пересекает.
	// вывести 3 с наибольшим кол-вом приходящих автобусов
	// найти все остановки в радиусе 3 км от начальной точки, на которые ходят текушие автобусы, и подсчитать сколько автобусов других там останавливаются

	lat := part * lat1 + (1-part) * lat2
	lon := part * lon1 + (1-part) * lon2
	listOfIndexStops := findWithTree(lat, lon, maxRadius)




	for _, routeNumber1 := range nearBus {
		for indRoute2, routeNumber2 := range route_id {
			if routeNumber1 == routeNumber2{

				if intInSlice(indRoute2, listOfIndexStops){
					//if Distance(lat, lon,transferCoordsLat,transferCoordsLon) < maxRadius && Distance(lat, lon,transferCoordsLat,transferCoordsLon) > minRadius {
					listOfIndex = append(listOfIndex, indRoute2) // список всех остановок вблизи пути автобуса вбизи определенной точки
					listofStops = append(listofStops, stop_id[indRoute2])
				}
			}
		}

	}

	var mainListOfStopsIndex []int
	var mainListOfStops []string
	listofStops = uniqueString(listofStops)

	for indexStop, stop := range stop_id {
		for _,l := range listofStops {
			//fmt.Println(stop)
			//fmt.Println(l)
			if stop == l {
				mainListOfStopsIndex = append(mainListOfStopsIndex, indexStop)
				mainListOfStops = append(mainListOfStops, stop)
			}
		}
	}



	stops, stopsCount := count_words(mainListOfStops)

	// определяем 10 остановок где в этом радиусе ходт самое большое колво автобусов

	var largestStop []string
	i := 0
	for k := 0; k < 500; k++ {
		i = MaxIndex(stopsCount)

		if len(stops) > k && i != 0 {
			largestStop = append(largestStop, stops[i])
			stopsCount = removeInt(stopsCount, i)
			stops = remove(stops, stops[i])
		}
	}


	//  если это кол во меньше 1 то увеличиваем радиус поиска
	if len(largestStop) < 1{
		return stopsInParticularRadius(nearBus , minRadius, maxRadius + 1, lat1, lon1, lat2, lon2, part)
	}

	// возвращаем остановки с самым большим колвом
	return largestStop
}


func MaxIndex(stopsCount []int) int {
	m := 0
	ind := 0
	for i, e := range stopsCount {
		if e > m {
			m = e
			ind = i

		}
	}

	return ind
}



func count_words (words []string) ([] string ,[]int) {
	uniqueWords := uniqueString(words)
	var countWords []int
	for i, word1 := range uniqueWords{
		countWords = append(countWords, 0)
		for _, word2 := range words{
			if word1 == word2 {
				countWords[i]++
			}
		}
	}


	/*
		word_counts := make(map[string]int)
		for _, word :=range words{
			word_counts[word]++
		}

	*/

	return uniqueWords, countWords
}


func uniqueNumberOfVar(stringSlice []StopBusStop) []StopBusStop {
	keys := make(map[StopBusStop]bool)
	list := []StopBusStop{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}


func uniqueNumberOfVarIndex(stringSlice []StopBusStopIndex) []StopBusStopIndex {
	keys := make(map[StopBusStopIndex]bool)
	list := []StopBusStopIndex{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
