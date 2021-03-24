package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const client_id = ""
const client_secret = ""

type StravaOAuthAccessResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	Athlete      struct {
		ID int `json:"id"`
	} `json:"athlete"`
}

func main() {
	http.HandleFunc("/stravaauthenticate", stravaAuthenticate)
	http.HandleFunc("/stravaredirect", stravaRedirect)
	http.ListenAndServe(":8080", nil)
}

func stravaAuthenticate(w http.ResponseWriter, r *http.Request) {
	redirectURL := fmt.Sprintf("http://www.strava.com/oauth/authorize?client_id=%s&response_type=code&redirect_uri=http://localhost:8080/stravaredirect&approval_prompt=force&scope=read", client_id)
	http.Redirect(w, r, redirectURL, http.StatusPermanentRedirect)
}

func stravaRedirect(w http.ResponseWriter, r *http.Request) {
	httpClient := http.Client{}
	err := r.ParseForm()
	if err != nil {
		log.Println("Can't parse form")
		w.WriteHeader(http.StatusBadRequest)
	}
	//get the code from
	code := r.FormValue("code")

	reqURL := fmt.Sprintf("https://www.strava.com/oauth/token?client_id=%s&client_secret=%s&code=%s&grant_type=authorization_code", client_id, client_secret, code)
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		log.Printf("Problem to create HTTP request %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	//We want json..
	req.Header.Set("accept", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Not possible to send HTTP request %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer res.Body.Close()

	// Parse the request
	var s StravaOAuthAccessResponse
	if err := json.NewDecoder(res.Body).Decode(&s); err != nil {
		log.Printf("JSON parse problem: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	log.Printf("Got AccessToken %v, RefreshToken %v for Athlete %v", s.AccessToken, s.RefreshToken, s.Athlete.ID)
}
