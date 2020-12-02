package main

import (
    "fmt"
    "net/http"

    "github.com/gorilla/mux"
)


func main() {
    r := mux.NewRouter()

    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "<h1>Welcome this is the Home Page Hope you enjoy our Encrypt and Decrypt Endpoints</h1>")
    })

    r.HandleFunc("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        title := vars["name"]

        fmt.Fprintf(w, "<h1>Hello, %s!\n</h1>", title)
    })

    http.ListenAndServe(":80", r)
}
