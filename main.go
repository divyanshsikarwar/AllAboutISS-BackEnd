package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	
)

type Location struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type Response struct {
	Message      string   
	IssPosition  Location  `json:"iss_position"`
}
//Longi and Lati


func GetLocation(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	
	resp, _ := http.Get("http://api.open-notify.org/iss-now.json")
	
	
	var responseFromAPI Response
	json.NewDecoder(resp.Body).Decode(&responseFromAPI)
	
	str,_ := json.Marshal(&responseFromAPI)
	
	res.Write(str)
	
}	






type AddressDetails struct {
	Town string `json:"town"`
	County string `json:"county"`
	Region string `json:"region"`
	State string `json:"state"`
	Country string `json:"country"`
	Country_code string `json:"country_code"`
}

type PhysicsAddress struct {
	Address AddressDetails `json:"address"`
}
//Address on Earth

func ReverseGeo(res http.ResponseWriter, req *http.Request)  {
	fmt.Println("RevGeo")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	
	var longi_lati Location
	err := json.NewDecoder(req.Body).Decode(&longi_lati)
	if err != nil {
		// Handle the error here, for example by logging the error message
		fmt.Println("Error decoding JSON:", err)
		
	}
	fmt.Println(longi_lati)
	resp,_ := http.Get("https://us1.locationiq.com/v1/reverse.php?key=pk.2e2c69c356595ca83c401c67ea119b60&lat=" + longi_lati.Longitude + "+&lon="+ longi_lati.Latitude + "+&format=json")
	fmt.Println(resp.Status)
	
	if (resp.Status != "200 OK"){
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte("500 - Something bad happened!"))
		return
	}

	var response PhysicsAddress
	json.NewDecoder(resp.Body).Decode(&response)
	
	str,_ := json.Marshal(&response)
	res.WriteHeader(http.StatusOK)
	res.Write(str)
	

}





type Astronaut struct {
	Name string `json:"name"`
	Craft string `json:"craft"`
}

type people struct {
	Number int `json:"number"`
	Astronauts []Astronaut `json:"people"`
}

//Names of Astraunuts in Space

func GetAstronauts(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	resp, _ := http.Get("http://api.open-notify.org/astros.json")
	var response people

	json.NewDecoder(resp.Body).Decode(&response)
	var FinalAstro []Astronaut

	for _,Astronaut := range response.Astronauts {
		if(Astronaut.Craft == "ISS"){
			FinalAstro = append(FinalAstro, Astronaut)
		}
	}

	response.Astronauts = FinalAstro
	response.Number = len(FinalAstro)

	ans,_ := json.Marshal(&response)
	res.WriteHeader(http.StatusOK)
	res.Write(ans)

}


func main(){
	router := mux.NewRouter()
	
	router.HandleFunc("/locate", GetLocation).Methods("GET")
	router.HandleFunc("/address", ReverseGeo).Methods("POST")
	router.HandleFunc("/fetchCrew", GetAstronauts).Methods("GET")


	http.ListenAndServe(":8000",router)

}