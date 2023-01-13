package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
)

const port = ":8080"

type people []struct {
	Id           int      `json:"id"`
	Name         string   `json:"name"`
	Image        string   `json:image`
	Members      []string `json:members`
	Creationdate int      `json:creationDate`
	Firstalbum   string   `json:firstAlbum`
	Localisation []string
	Date         []string
	Relation     []string
	Nameformated string
}

type mapi []struct {
	Id   int                 `json:"id"`
	Loca map[string][]string `json:"datesLocations"`
}
type ind struct {
	Ind mapi `json:"index"`
}

func getdata() people {
	url := "https://groupietrackers.herokuapp.com/api/artists"
	var Rp people
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	Body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	urlr := "https://groupietrackers.herokuapp.com/api/relation"
	re, err := http.Get(urlr)
	if err != nil {
		log.Fatal(err)
	}
	defer re.Body.Close()

	Bod, err := ioutil.ReadAll(re.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(Body, &Rp)
	var ad ind
	json.Unmarshal(Bod, &ad)
	for i := 0; i < len(ad.Ind); i++ {
		a := strings.Split(Rp[i].Name, " ")
		b := ""
		for _, v := range a {
			b = b + strings.ToLower(v)
		}
		Rp[i].Nameformated = b
		for date := range ad.Ind[i].Loca {
			for j := 0; j < len(ad.Ind[i].Loca[date]); j++ {
				Rp[i].Date = append(Rp[i].Date, ad.Ind[i].Loca[date][j])
				Rp[i].Localisation = append(Rp[i].Localisation, date)
				Rp[i].Relation = append(Rp[i].Relation, date+" : "+ad.Ind[i].Loca[date][j])
			}
		}
	}
	return Rp
}

func artist(i int, Data people) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/"+Data[i].Nameformated {
			http.Error(w, "Error 404 page not found", http.StatusNotFound)
			return
		}
		t := template.Must(template.ParseFiles("artist.html"))
		Ap := Data[i]
		t.Execute(w, Ap)
	}
}

func hundel(Data people) {
	for i := 0; i < len(Data); i++ {
		http.HandleFunc("/"+Data[i].Nameformated, artist(i, Data))
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Error 404 page not found", http.StatusNotFound)
		return
	}
	Rp := getdata()
	t := template.Must(template.ParseFiles("index.html"))
	t.Execute(w, Rp)
}

func main() {
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))
	http.HandleFunc("/", index)
	hundel(getdata())
	fmt.Println("(http://localhost:8080): Server started on port", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
