package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"./database"
	"./retrieval"
)

var tpl *template.Template

type ResultPage struct {
	Id           int64    `json:"Id"`
	Score        float64  `json:"Score"`
	Title        string   `json:"Title"`
	Url          string   `json:"Url"`
	LastModified string   `json:"LastModified"`
	PageSize     string   `json:"PageSize"`
	Keywords     []string `json:"Keywords"`
	Parents      []string `json:"Parents"`
	Children     []string `json:"Children"`
}

type TheResult struct {
	ResultPages []ResultPage `json:PageScore`
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*html"))
}

func index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.html", nil)
}
func processinput(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//create json files of the result
	fmt.Println(time.Now())
	retrieval.RetrievalFunction(r.FormValue("searchInput"))
	fmt.Println(time.Now())

	//extract the json files
	jsonFile, err := os.Open("search_output.json")

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("connected to search_output.json")
	fmt.Print(jsonFile)
	defer jsonFile.Close()
	fmt.Print(jsonFile)
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var results TheResult
	fmt.Println("STEP 2")
	json.Unmarshal(byteValue, &results)
	fmt.Print(results)
	fmt.Println("STEP 4")

	for i := 0; i < 2; i++ {
		fmt.Println(results.ResultPages[i].Title)
	}
	/* 	d := struct {
	   		Time string
	   		//String string
	   	}{
	   		Time: "what time is it",
	   		//String: t,
	   	} */

	tpl.ExecuteTemplate(w, "result.html", nil)
}
func main() {
	fmt.Println("Now Listening on 8000")
	database.OpenAllDb()
	http.HandleFunc("/", index)
	http.HandleFunc("/index", index)
	http.HandleFunc("/result", processinput)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
