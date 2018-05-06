package main

import (
    "html/template"
	"net/http"
	"github.com/gorilla/mux"
    "fmt"
)


// dongik's main function
func indexHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("URL.Path = %q\n", r.URL.Path)
    fmt.Printf("URL = %q\n", r.RemoteAddr)
    t := template.New("main")
    t, _ = template.ParseFiles("order.html")  // Parse template file.
    t.Execute(w,"Hello world")  // merge.
}
func orderHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
}


func main(){
	r := mux.NewRouter()
    r.HandleFunc("/",indexHandler);
	r.HandleFunc("/orders", orderHandler).Methods("POST")
	http.ListenAndServe(":8080", r); 
}
