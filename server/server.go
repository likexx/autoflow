package main

import (
    "fmt"
    "net/http"
    "log"
    "io/ioutil"
    "encoding/json"
    "github.com/gorilla/mux"
    autoflow "server/autoflow"
)

type CreateSessionResponse struct {
    Error int
    Result string
}

type QueryNextResponse struct {
    Error int
    Action string
    Parameters autoflow.ActionParameter
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "hello") // send data to client side
}

func autoflowCreateSession(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    flowName := vars["flowname"]
    flow := autoflow.GetInstance()
    sessionId := flow.CreateSession(flowName)
    result := CreateSessionResponse{0, sessionId}

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func autoflowNext(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    sessionId := vars["session_id"]

    var clientParams map[string]interface{}

    bodyBytes, err := ioutil.ReadAll(r.Body)
    if err != nil {
        fmt.Fprintf(w, "missing ")
        return
    }

    if len(bodyBytes) > 0 {
        json.Unmarshal(bodyBytes, &clientParams)
    }

    flow := autoflow.GetInstance()
    action, parameters := flow.QueryNextStep(sessionId, clientParams)

    result := QueryNextResponse{0, action, parameters}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func autoflowError(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    sessionId := vars["session_id"]
    flow := autoflow.GetInstance()
    action, parameters := flow.OnError(sessionId)

    result := QueryNextResponse{0, action, parameters}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func autoflowStop(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    sessionId := vars["session_id"]
    flow := autoflow.GetInstance()
    flow.Stop(sessionId)

    result := CreateSessionResponse{0, ""}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/healthcheck", healthcheck).Methods("GET")

    a := r.PathPrefix("/autoflow").Subrouter()
    a.HandleFunc("/create/{flowname}", autoflowCreateSession).Methods("POST")
    a.HandleFunc("/next/{session_id}", autoflowNext).Methods("POST")
    a.HandleFunc("/error/{session_id}", autoflowError).Methods("POST")
    a.HandleFunc("/stop/{session_id}", autoflowStop).Methods("POST")

    http.Handle("/", r)

    serverAddr := "0.0.0.0:8080"
    err := http.ListenAndServe(serverAddr, nil) // set listen port
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }

}
