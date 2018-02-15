package main

import (
	"encoding/csv"
	"encoding/json"
	"github.com/jasn/goors"
	"html/template"
	"math"
	"net/http"
	"os"
	"strconv"
)

type Foodtruck struct {
	Lat, Lng float64
	Name     string
}

type MyHandler struct {
	foodtrucks           []Foodtruck
	rangeSearchStructure *goors.RangeSearchAdvanced
}

func NewMyHandler(foodtrucks []Foodtruck) *MyHandler {
	var handler MyHandler
	handler.foodtrucks = foodtrucks
	handler.preprocessFoodtrucks()
	return &handler
}

func (self *MyHandler) preprocessFoodtrucks() {
	points := make([]goors.Point, len(self.foodtrucks))
	for i, truck := range self.foodtrucks {
		points[i] = goors.MakePoint(truck.Lat, truck.Lng)
	}
	self.rangeSearchStructure = goors.NewRangeSearchAdvanced(points)
	self.rangeSearchStructure.Build()
}

func (self *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	getValueMap := r.URL.Query()
	x0, errx0 := strconv.ParseFloat(getValueMap.Get("x0"), 64)
	x1, errx1 := strconv.ParseFloat(getValueMap.Get("x1"), 64)
	y0, erry0 := strconv.ParseFloat(getValueMap.Get("y0"), 64)
	y1, erry1 := strconv.ParseFloat(getValueMap.Get("y1"), 64)

	if errx0 != nil || errx1 != nil || erry0 != nil || erry1 != nil {
		invalidCoordinate(w)
		return
	}

	bottomLeft := goors.MakePoint(math.Min(x0, x1), math.Min(y0, y1))
	topRight := goors.MakePoint(math.Max(x0, x1), math.Max(y0, y1))
	res := self.rangeSearchStructure.Query(bottomLeft, topRight)
	foodtrucksRes := make([]Foodtruck, len(res))
	for i := 0; i < len(res); i++ {
		foodtrucksRes[i] = self.foodtrucks[res[i]]
	}

	b, err := json.Marshal(foodtrucksRes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write(b)
	}
}

func invalidCoordinate(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}

func readFoodtrucks() []Foodtruck {
	filename := "data/cleaned_trucks.csv"
	f, _ := os.Open(filename)
	csvReader := csv.NewReader(f)
	records, _ := csvReader.ReadAll()
	var res []Foodtruck
	for _, v := range records[1:] {
		lat, _ := strconv.ParseFloat(v[0], 64)
		lng, _ := strconv.ParseFloat(v[1], 64)
		res = append(res, Foodtruck{lat, lng, v[2]})
	}
	return res
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("html/index.html")
	t.Execute(w, nil)
}

func main() {
	foodtrucks := readFoodtrucks()
	http.Handle("/", NewMyHandler(foodtrucks))
	http.HandleFunc("/index", serveIndex)
	http.ListenAndServe(":8080", nil)
}
