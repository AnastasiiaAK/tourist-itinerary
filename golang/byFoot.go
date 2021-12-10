package points

import (
	hp "container/heap"
	"encoding/json"
	"io/ioutil"
	"math"
	"strconv"
	"time"
	"fmt"
)

type JsonGeometry struct{
	Geometry []JsonCoords `json:"geometry"`
}

type JsonElements struct {
	Elements []JsonGeometry `json:"elements"`
}

type JsonCoords struct{
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type buildGeometry struct {
	point0 JsonCoords
	point1 JsonCoords
}

/*
type maxminCoords struct {
	district string
	coords []float64
}

*/



//var begin string = "Church of the Savior on Spilled Blood" //! 2 пересадки и 1 пересадка и прямой
//var begin string = "Grand Choral Synagogue"
//var begin string = "Protection of the Holy Virgin Temple" //! 1 пересадка

//var end string = "Temple Theodore Icon of the Mother of God" // прямой
//var end string = "Ruinny Bridge" // 2 пересадки
//var end string = "GeniumPark"
//var beginCoords, endCoords coords = findLocation(begin, end)




func main(){



	//var begin string = "Church of the Savior on Spilled Blood" //! 2 пересадки и 1 пересадка и прямой
	//var end string = "Grand Choral Synagogue"
	//var end string = "Protection of the Holy Virgin Temple" //! 1 пересадка 2 пересадки с 1 км

	//var end string = "Temple Theodore Icon of the Mother of God" // прямой

	//var end string = "Ruinny Bridge" // 2 пересадки
	//var end string = "GeniumPark"
	//var end string = "Kirov Central Park of Culture and Recreation"
	//var end string = "Trinity Cathedral" // 2 пересадки с 1 км

	//var begin string = "Tekhnologichesky Institut"
	//var end string = "Vasilyevsky Island"
	//var beginCoords, endCoords coords = findLocation(begin, end)


	var beginCoords, endCoords coords
	beginCoords.Lat, beginCoords.Lon = 59.916389,30.318611
	endCoords.Lat, endCoords.Lon = 59.9386111111,30.2561111111


	beginCoords.Lat, beginCoords.Lon = 59.9386111111,30.2561111111
	endCoords.Lat, endCoords.Lon = 59.936024308,30.3259956416

	beginCoords.Lat, beginCoords.Lon = 59.936024308,30.3259956416
	endCoords.Lat, endCoords.Lon = 59.93461,30.33253

	beginCoords.Lat, beginCoords.Lon = 59.93461,30.33253
	endCoords.Lat, endCoords.Lon = 59.940119,30.328904

	beginCoords.Lat, beginCoords.Lon =59.940119,30.328904
	endCoords.Lat, endCoords.Lon = 59.95006,30.315913


	beginCoords.Lat, beginCoords.Lon =59.95006,30.315913
	endCoords.Lat, endCoords.Lon = 59.948699,30.327187

	beginCoords.Lat, beginCoords.Lon =59.948699,30.327187
	endCoords.Lat, endCoords.Lon = 59.941715,30.299208





	//fmt.Println(beginCoords, endCoords)
	//var beginCoords, endCoords coords

	//beginCoords.Lat, beginCoords.Lon = 60.0672000874539,30.3365993
	//beginCoords.Lat, beginCoords.Lon = 59.9828565251344,30.4092999
	//beginCoords.Lat = 59.9828565251344
	//beginCoords.Lon = 30.4092999
	//endCoords.Lat, endCoords.Lon = 59.8466200283376,30.4831524
	//endCoords.Lat, endCoords.Lon = 59.984301827335,30.2436981

	//beginCoords.Lat, beginCoords.Lon = 60.027851,30.222640000000002 // "ПР. АВИАКОНСТРУКТОРОВ, 47" 2 пересадки
	//endCoords.Lat, endCoords.Lon = 59.9767876172422,30.4015617//"ЗАМШИНА УЛ."
	//endCoords.Lat, endCoords.Lon = 59.9217134316208,30.4777584 // "ТОВАРИЩЕСКИЙ ПР." 2 пересадки
	//endCoords.Lon = 30.2436981

	//beginCoords.Lat, beginCoords.Lon = 59.851786588881,30.2269344 //МОРСКОЙ ТЕХНИЧЕСКИЙ УНИВЕРСИТЕТ 1 пересадка прямой
	//endCoords.Lat, endCoords.Lon = 59.78005896907,30.1043625//ПОС. ТОРИКИ
	//endCoords.Lat, endCoords.Lon = 59.8834751064445,30.2649879 //ул. Возрождения прямой

	//endCoords.Lat, endCoords.Lon = 59.9533626866534,30.3244419//"ТРОИЦКАЯ ПЛ." 1 пересадка


	//beginCoords.Lat, beginCoords.Lon = 60.058298, 30.331365 //frpm south to north
	//endCoords.Lat, endCoords.Lon = 59.834088, 30.349234 //frpm south to north

	//beginCoords.Lat, beginCoords.Lon = 59.999986, 30.206868 //frpm south to north
	//endCoords.Lat, endCoords.Lon = 59.990715, 30.441014 //frpm south to north

	//beginCoords.Lat, beginCoords.Lon = 59.999986, 30.206868 //frpm севео восток
	//endCoords.Lat, endCoords.Lon = 59.990715, 30.441014 //frpm северо запад

	//beginCoords.Lat, beginCoords.Lon = 59.932627, 30.265919
	//endCoords.Lat, endCoords.Lon = 59.924713, 30.288578


	fmt.Println(beginCoords, endCoords)


	now0 := makeTimestamp()

			time := mainPath(beginCoords.Lat,beginCoords.Lon, endCoords.Lat, endCoords.Lon)

	now1 := makeTimestamp()

	fmt.Println("!Время выполнения программы",now1 - now0)
	fmt.Println("timeByBus",time)


	var distTaxi = mainTaxi(beginCoords, endCoords)
	var timeTaxi = distTaxi/30 * 60
	fmt.Println("timeByTaxi", timeTaxi)
}



func mainPath(Lat0,Lon0, Lat1, Lon1 float64) (float64, [][]int) {




	var path, dist, timeByFoot = byFoot(Lat0,Lon0, Lat1, Lon1)

	fmt.Println("timeByFoot",timeByFoot)
	timeBest := timeByFoot

	if timeByFoot > 20 || timeByFoot == 0{

		now0 := makeTimestamp()

		timeBest, indexBest := bus2(Lat0, Lon0, Lat1, Lon1, 1)
		fmt.Println(timeBest, indexBest)
		now1 := makeTimestamp()

		fmt.Println("Время выполнения программы",now1 - now0)

		return timeBest, indexBest



	} else{
		fmt.Println("Coordinates by foot", path, "\n", "Distance", Round(dist, 10), "km", "\n", "Time", timeByFoot, "minuts")
	}

	return timeBest, [][]int{}


}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func neccessaryFile()(string, string, string){

	startTime := time.Now()
	hour := startTime.Format("15")
	minute := startTime.Format("04")
	weekDay := startTime.Format("Monday")
	minute = "02"
	hour = "17"
	weekDay = "Thursday"
	fileHour := hour
	fileMinute:= string(minute[0])


	if string(minute[0]) == "0"{
		fileMinute = string(minute[0])
	} else {
		fileMinute = string(minute[0]) + "0"
	}




	return weekDay, fileHour, fileMinute
}


func changeFileForTime(plusTime float64)(string, string, string){
	var weekDay, fileHour, fileMinute string = neccessaryFile()
	fileMinute0, _ := strconv.ParseInt(fileMinute, 10, 64)
	fileHour0, _ := strconv.ParseInt(fileHour, 10, 64)
	plusTime0 := int64(plusTime)

	fileMinute1 := fileMinute0 + plusTime0
	//fmt.Println("fileMinute1 / 60 >= 1", fileMinute1 / 60 >= 1)
	if fileMinute1 / 60 >= 1 {
		//fmt.Println("fileHour01", fileHour0)
		fileHour0 = fileHour0 + fileMinute1 / 60
		//fmt.Println("fileHour0", fileHour0)
		//fmt.Println("fileMinute1 / 60", fileMinute1 / 60)
		fileMinute0 = fileMinute1 - fileMinute1 / 60 * 60
		//fmt.Println("fileMinute1", fileMinute1)
		//fmt.Println("fileMinute1 / 60 * 60", fileMinute1 / 60 * 60)
		//fmt.Println("fileMinute0", fileMinute0)

	} else {
		fileMinute0 = fileMinute0 + plusTime0
		fileHour0 = fileHour0
	}


	fileHour = strconv.FormatInt(fileHour0, 10)
	//fmt.Println("fileHour",fileHour)
	fileMinute2 := strconv.FormatInt(fileMinute0, 10)

	if string(fileMinute2[0]) == "0"{
		fileMinute = string(fileMinute2[0])
	} else {
		fileMinute = string(fileMinute2[0]) + "0"
	}
	//fmt.Println("fileMinute", fileMinute)

	return weekDay, fileHour, fileMinute
}

func byFoot(Lat0,Lon0, Lat1, Lon1 float64) ([]JsonCoords, float64, float64) {

	var district =  []string {"centralnii_district.json", "petrogradskii_district.json",
		"kalininskii_district.json", "primorskii_district.json",
		"admiralteiskii_district.json", "nevskii_district.json", "kirovskii_district.json",
		"frunzenskii_district.json", "vyborgskii_district.json",
		"vasileostrovskii_district.json","moskovskii_district.json"}





	var coordsAll0 = unionCoords(district)
	var beginCoordinates JsonCoords
	beginCoordinates.Lat, beginCoordinates.Lon = Lat0, Lon0
	var endCoordinates JsonCoords
	endCoordinates.Lat, endCoordinates.Lon = Lat1, Lon1

	var coordsAll1 = insertLocation(beginCoordinates, coordsAll0)
	var coordsAll = insertLocation(endCoordinates, coordsAll1)

	var edgesForBuild, _ = forBuild(coordsAll)

	//fmt.Println("edgesForBuild",edgesForBuild)

	graph := newGraph()


	for _, k := range edgesForBuild{
		//fmt.Println(k.point0)
		graph.addEdge(k.point0, k.point1, Distance(k.point0.Lat, k.point0.Lon, k.point1.Lat, k.point1.Lon))
	}

	var dist , path = graph.getPath(beginCoordinates, endCoordinates)

	var timeByFoot = Round(dist / 5 * 60, 4)
	// fmt.Println("Coordinates by foot", path, "\n", "Distance", Round(dist, 10),"km","\n", "Time", timeByFoot, "minuts")

	return path, dist, timeByFoot

}



func forBuild(coordsAll []JsonGeometry) ([]buildGeometry, []JsonCoords) {

	var coordsForBuild []buildGeometry
	var nodeForBuild []JsonCoords
	for _, coordinates := range coordsAll{
		if len(coordinates.Geometry) > 0{
			for i , _ := range coordinates.Geometry[:len(coordinates.Geometry) - 1] {

				var kopl buildGeometry
				kopl.point0, kopl.point1 = coordinates.Geometry[i], coordinates.Geometry[i+1]
				coordsForBuild = append(coordsForBuild, kopl)

				nodeForBuild = append(nodeForBuild, coordinates.Geometry[i])
			}
		}
	}

	return coordsForBuild, nodeForBuild

}



func insertLocation(beCoordinates JsonCoords, coordsAll []JsonGeometry) []JsonGeometry {


	//fmt.Println(beCoordinates.Lat, beCoordinates.Lon)
	var near float64 = 1000000000000
	var point int = 0
	var findRoad int = 0
	for k, road := range coordsAll{
		for i := range road.Geometry{
			var a = Distance(beCoordinates.Lat, beCoordinates.Lon, road.Geometry[i].Lat,road.Geometry[i].Lon )
			if a < near{
				near = a
				point = i
				findRoad = k
			}
		}
	}


	var coordsZero JsonCoords
	coordsZero.Lat, coordsZero.Lon = 0, 0
	//fmt.Println(coordsAll[findRoad].Geometry)

	//fmt.Println(coordsAll[findRoad].Geometry)
	coordsAll[findRoad].Geometry = append(coordsAll[findRoad].Geometry, coordsZero)
	//fmt.Println(coordsAll[findRoad].Geometry)

	copy(coordsAll[findRoad].Geometry[point + 2:], coordsAll[findRoad].Geometry[point + 1:])
	coordsAll[findRoad].Geometry[point + 1] =  beCoordinates

	//fmt.Println(coordsAll[findRoad].Geometry)


	return coordsAll

}



func unionCoords(district []string) []JsonGeometry {
	var coordsAll []JsonGeometry
	for _, distr := range district {
		//fmt.Println(distr)
		fileName := fmt.Sprintf("points/%s" , distr)
		file, _ := ioutil.ReadFile(fileName)
		var nastya JsonElements
		json.Unmarshal(file, &nastya)

		coordsAll = append(coordsAll, nastya.Elements...)

	}
	return coordsAll
}

type JsonSights struct {
	Lat                       float64 `json:"lat"`
	Lng                       float64 `json:"lng"`
	Title                     string `json:"title"`
}

type coords struct {
	Title string
	Lat float64
	Lon float64
}

func findLocation(begin, end string) (coords, coords) {

	type Title struct {
		Title string `json:"title"`
	}
	type Lng struct {
		Lng float64 `json:"lng"`
	}
	type Lat struct {
		Lat float64 `json:"lat"`
	}


	file, _ := ioutil.ReadFile("Saint+Petersburg-finalized.json")


	var nastya []JsonSights
	json.Unmarshal(file, &nastya)


	var coordinates []coords

	for i:=0; i < len(nastya); i++{
		coordinates = append(coordinates, coords{Title: nastya[i].Title,Lat: nastya[i].Lat, Lon:nastya[i].Lng})
	}

	var beginCoords coords
	var endCoords coords

	for _, l := range coordinates{
		if l.Title == begin{
			beginCoords = l
		}
		if l.Title == end{
			endCoords = l
		}

		if beginCoords.Lat > 0 && endCoords.Lat > 0{
			break
		}
	}
	return beginCoords, endCoords
}

func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h)) / 1000}


func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}






type edge struct {
	node   JsonCoords
	weight float64
}

type graph struct {
	nodes map[JsonCoords][]edge
}

func newGraph() *graph {
	return &graph{nodes: make(map[JsonCoords][]edge)}
}

func (g *graph) addEdge(origin, destiny JsonCoords, weight float64) {
	g.nodes[origin] = append(g.nodes[origin], edge{node: destiny, weight: weight})
	g.nodes[destiny] = append(g.nodes[destiny], edge{node: origin, weight: weight})
}

func (g *graph) getEdges(node JsonCoords) []edge {
	return g.nodes[node]
}

func (g *graph) getPath(origin, destiny JsonCoords) (float64, []JsonCoords) {
	h := newHeap()
	h.push(path{value: 0, nodes: []JsonCoords{origin}})
	visited := make(map[JsonCoords]bool)

	for len(*h.values) > 0 {
		// Find the nearest yet to visit node
		p := h.pop()
		node := p.nodes[len(p.nodes)-1]

		if visited[node] {
			continue
		}

		if node == destiny {
			return p.value, p.nodes
		}

		for _, e := range g.getEdges(node) {
			if !visited[e.node] {
				// We calculate the total spent so far plus the cost and the path of getting here
				h.push(path{value: p.value + e.weight, nodes: append([]JsonCoords{}, append(p.nodes, e.node)...)})
			}
		}

		visited[node] = true
	}

	return 0, nil
}


type path struct {
	value float64
	nodes []JsonCoords
}

type minPath []path

func (h minPath) Len() int  { return len(h) }
func (h minPath) Less(i, j int) bool { return h[i].value < h[j].value }
func (h minPath) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *minPath) Push(x interface{}) {
	*h = append(*h, x.(path))
}

func (h *minPath) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type heap1 struct {
	values *minPath
}

func newHeap() *heap1 {
	return &heap1{values: &minPath{}}
}

func (h *heap1) push(p path) {
	hp.Push(h.values, p)
}

func (h *heap1) pop() path {
	i := hp.Pop(h.values)
	return i.(path)
}


func Round(x, unit float64) float64 {
	return (math.Floor(x * unit) /unit)
}



/*
	var maxlat, minlat, maxlon, minlon float64 = -1000, 1000, -1000, 1000
	for _,district := range allDistrict {
		var coords = unionCoords([]string {district})
		for _, j := range coords{
			if j.Geometry[0].Lat > maxlat{
				maxlat = j.Geometry[0].Lat
			}
			if j.Geometry[0].Lat < minlat{
				minlat = j.Geometry[0].Lat
			}

			if j.Geometry[0].Lon > maxlon{
				maxlon = j.Geometry[0].Lon
			}
			if j.Geometry[0].Lon < minlon{
				minlon = j.Geometry[0].Lon
			}

		}


		fmt.Println(coords)
	}

*/
