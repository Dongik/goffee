package main

import (
    "encoding/json"
    "html/template"
    "bufio"
	"net/http"
	"github.com/gorilla/mux"
    "fmt"
    "encoding/csv"
    "io"
    "os"
    "log"
    "strconv"

    //. "github.com/Dongik/goffee/models"
)
type Drink struct {
    Name    string  `json:"name"`
    Price   int     `json:"price"`
}
type Request struct {
    Order string    `json:"order"`
    Number int      `json:"number"`
}

var ipUser map[string]string
var orderNum map[string]int
var drinkPrice map[string]int
var drinks []Drink

func respondWithError(w http.ResponseWriter, code int, msg string) {
    respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

// dongik's main function
func indexHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("URL.Path = %q\n", r.URL.Path)
    fmt.Printf("URL = %q\n", r.RemoteAddr)
    t := template.New("main")
    t, _ = template.ParseFiles("static/order.html")  // Parse template file.
    t.Execute(w,"Hello world")  // merge.
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    var request Request
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    order := request.Order
    if number, ok := orderNum[order]; ok {
    //do something here
        orderNum[order] = number + 1
    }else {
        orderNum[order] = 1
    }
    respondWithJson(w, http.StatusOK, map[string]string{"result":"success"})
}

func drinksHandler(w http.ResponseWriter, r *http.Request) {
    respondWithJson(w, http.StatusOK, drinks)
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
    var result []Request
    for order := range orderNum{
        result = append(result,Request{
            Order:order,
            Number:orderNum[order],
        })
    }
    respondWithJson(w, http.StatusOK, result)
}

func loadDrinks(){
    //load list
    csvFile, err := os.Open("list.csv")
    if err != nil { fmt.Print(err) }

    reader := csv.NewReader(bufio.NewReader(csvFile))

    for {
        line, err := reader.Read()
        if err == io.EOF {
            break
        } else if err != nil {
            log.Fatal(err)
        }
        name := line[0]
        price, _ := strconv.Atoi(line[1])
        if err != nil { fmt.Print(err)}
        drinks = append(drinks, Drink{
            Name: name,
            Price: price,
            })
    }
}

func init(){
    //initialize map
    ipUser = make(map[string]string)
    orderNum = make(map[string]int)
    loadDrinks()
}

func main(){
    //router
	r := mux.NewRouter()
    //drink
    r.HandleFunc("/drinks", drinksHandler).Methods("GET")//get all of drinks

    //order
    r.HandleFunc("/results", resultHandler).Methods("GET")
	r.HandleFunc("/orders", orderHandler).Methods("PUT")//update order
    //r.HandleFunc("/orders", ordersHandler).Methods("GET")//get all orders
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

    err := http.ListenAndServe(":8080", r)
    if err != nil {
        fmt.Printf(err.Error())
    }
}
