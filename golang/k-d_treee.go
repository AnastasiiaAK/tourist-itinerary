package points

import (
	"gonum.org/v1/gonum/spatial/kdtree"
	"math"
	"strconv"
)

func findWithTree(currentLat, currentLon, dist float64) []int {
	// Construct a k-d tree of train station locations
	// to identify accessible public transport for the
	// elderly.
	t := kdtree.New(stations, false)

	// Residence.
	q := place{lat: currentLat, lon: currentLon}

	var keep kdtree.Keeper

	// Find all stations within 0.75 of the residence.
	keep = kdtree.NewDistKeeper(dist * dist) // Distances are squared.
	t.NearestSet(keep, q)

	//fmt.Println(`Stations within 750 m of 51.501476N 0.140634W.`)

	var listOfPlaces []int
	for _, c := range keep.(*kdtree.DistKeeper).Heap {
		p := c.Comparable.(place)
		listOfPlaces = append(listOfPlaces, p.name)
		//fmt.Printf("%s: %0.3f km\n", p.name, math.Sqrt(p.Distance(q)))

	}
	//fmt.Println()
	return listOfPlaces
}

// stations is a list of railways stations satisfying the
// kdtree.Interface.

func createTree() kdtree.Interface {
	var alon, alng,_,_, _, _, _,_, _, _, _,number,_, _ = readFile()
	var stop place
	var stops kdtree.Interface
	var sto []place
	for i, j := range number{
		numb, _ := strconv.Atoi(j)
		numb0 := numb - 1
		alon0, _ := strconv.ParseFloat(alon[i], 8)
		alat0, _ := strconv.ParseFloat(alng[i], 8)
		stop = place{name: numb0, lat: alon0, lon: alat0}

		sto = append(sto, stop)

	}
	stops = places(sto)

	return stops
}


var stations = createTree()


// place is a kdtree.Comparable implementations.
type place struct {
	name  int
	lat, lon float64
}

// Compare satisfies the axis comparisons method of the kdtree.Comparable interface.
// The dimensions are:
//  0 = lat
//  1 = lon
func (p place) Compare(c kdtree.Comparable, d kdtree.Dim) float64 {
	q := c.(place)
	switch d {
	case 0:
		return p.lat - q.lat
	case 1:
		return p.lon - q.lon
	default:
		panic("illegal dimension")
	}
}

// Dims returns the number of dimensions to be considered.
func (p place) Dims() int { return 2 }

// Distance returns the distance between the receiver and c.
func (p place) Distance(c kdtree.Comparable) float64 {
	q := c.(place)
	d := haversine(p.lat, p.lon, q.lat, q.lon)
	return d * d
}

// haversine returns the distance between two geographic coordinates.
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const r = 6371 // km
	sdLat := math.Sin(radians(lat2-lat1) / 2)
	sdLon := math.Sin(radians(lon2-lon1) / 2)
	a := sdLat*sdLat + math.Cos(radians(lat1))*math.Cos(radians(lat2))*sdLon*sdLon
	d := 2 * r * math.Asin(math.Sqrt(a))
	return d // km
}

func radians(d float64) float64 {
	return d * math.Pi / 180
}

// places is a collection of the place type that satisfies kdtree.Interface.
type places []place

func (p places) Index(i int) kdtree.Comparable         { return p[i] }
func (p places) Len() int                              { return len(p) }
func (p places) Pivot(d kdtree.Dim) int                { return plane{places: p, Dim: d}.Pivot() }
func (p places) Slice(start, end int) kdtree.Interface { return p[start:end] }

// plane is required to help places.
type plane struct {
	kdtree.Dim
	places
}

func (p plane) Less(i, j int) bool {
	switch p.Dim {
	case 0:
		return p.places[i].lat < p.places[j].lat
	case 1:
		return p.places[i].lon < p.places[j].lon
	default:
		panic("illegal dimension")
	}
}
func (p plane) Pivot() int { return kdtree.Partition(p, kdtree.MedianOfMedians(p)) }
func (p plane) Slice(start, end int) kdtree.SortSlicer {
	p.places = p.places[start:end]
	return p
}
func (p plane) Swap(i, j int) {
	p.places[i], p.places[j] = p.places[j], p.places[i]
}
