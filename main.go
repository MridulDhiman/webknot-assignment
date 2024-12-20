package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Rating struct {
	Average float64 `json:"average"`
	Reviews int `json:"reviews"`
}

type Response struct {
	Price string `json:"price"`
	Name string `json:"name"`
	Rating Rating
	Image string `json:"image"`
	Id int `json:"id"`
}

type ResponseType []Response 

func (a ResponseType) Len() int {
	 return len(a)
	 }

func (a ResponseType) Less(i, j int) bool { 
	price1 := strings.TrimPrefix(a[i].Price, "$")
	price2 := strings.TrimPrefix(a[j].Price, "$")
	priceFloat, err:= strconv.ParseFloat(price1, 64)
	if err != nil {
		log.Fatal(err)
	}
	price2Float, err := strconv.ParseFloat(price2, 64)
	if err != nil {
		log.Fatal(err)
	}
	return priceFloat < price2Float;
}

func (a ResponseType) Swap(i int, j int) {
	 a[i], a[j] = a[j], a[i]
	 }

func sortHandler(w http.ResponseWriter, r* http.Request) {
	resp, err := http.Get("https://api.sampleapis.com/beers/ale")
	if err != nil {
		log.Fatalln(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
	   log.Fatalln(err)
	}
 //Convert the body to type string
	var arr []Response
	if err := json.Unmarshal(body, &arr); err != nil {
		log.Fatal("(error) could not unmarshal json")
	}

	// sort by price 
	sort.Sort(ResponseType(arr))
	out, err := json.Marshal(&arr)
	if err != nil {
		log.Fatal("could not convert arr into json: ", err)
	}
	w.Write(out)
}

func filterHandler(w http.ResponseWriter, r* http.Request) {
	resp, err := http.Get("https://api.sampleapis.com/beers/ale")
	if err != nil {
		log.Fatalln(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
	   log.Fatalln(err)
	}
 //Convert the body to type string
	
	var arr []Response
	if err := json.Unmarshal(body, &arr); err != nil {
		log.Fatal("(error) could not unmarshal json")
	}

	
	// filter by rating 
	rating := r.URL.Query().Get("rating")
	// get avg. rating <= rating
	fmt.Println("rating: ", rating)
	if rating == "" {
		w.Write([]byte("could not find rating"))
	}
	ratingFloat, err := strconv.ParseFloat(rating, 64)
	if err != nil {
		log.Fatal("could not convert rating to int64")
	}
	var ans []Response
	for  _, element:= range arr {
		if element.Rating.Average > ratingFloat {
			ans = append(ans, element)
		}
	}
	out, err := json.Marshal(&ans)
	if err != nil {
		log.Fatal("could not convert arr into json: ", err)
	}
	w.Write(out)

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/sort", sortHandler).Methods("GET");
	router.HandleFunc("/filter", filterHandler).Methods("GET");
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalln("Unexpected Error Occurred: ", err)
	}
}
