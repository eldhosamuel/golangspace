package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type server struct{}

func ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	//resp.WriteHeader(http.StatusOK)
	//resp.Write([]byte(`{"Name": "Eldho Samuel"}`))
	// fmt.Fprintf(w, "This is an example server.\n")
	// io.WriteString(w, "This is an example server.\n")

	/*switch req.Method {
	case "GET":
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(`{"message": "get called"}`))
	case "POST":
		resp.WriteHeader(http.StatusCreated)
		resp.Write([]byte(`{"message": "post called"}`))
	case "PUT":
		resp.WriteHeader(http.StatusAccepted)
		resp.Write([]byte(`{"message": "put called"}`))
	case "DELETE":
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(`{"message": "delete called"}`))
	default:
		resp.WriteHeader(http.StatusNotFound)
		resp.Write([]byte(`{"message": "not found"}`))
	}*/

	if req.Method == "POST" {
		resp.WriteHeader(http.StatusOK)
		resp.Write([]byte(`{"Name": "Eldho Samuel"}`))
	} else {
		resp.WriteHeader(http.StatusForbidden)
	}

}

func main() {
	//http.HandleFunc("/getLoanOffers", ServeHTTP)
	//http.HandleFunc("/test", testOffers)

	http.HandleFunc("/consumeapi", postRequestHTTP)
	err := http.ListenAndServeTLS(":8443", "server.crt", "server.key", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

type employee struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func testOffers(w http.ResponseWriter, r *http.Request) {
	headerContentTtype := r.Header.Get("Content-Type")
	if headerContentTtype != "application/json" {
		apiResponse(w, "Error", "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	var e employee
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&e)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			apiResponse(w, "Error", "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			apiResponse(w, "Error", "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}

	//apiResponse(w, "Status", "Request Accepted", http.StatusOK)
	apiResponse(w, "ReturnCode", "000", http.StatusOK)
	return
}

func apiResponse(w http.ResponseWriter, jsonKey string, jsonValue string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)

	resp[jsonKey] = jsonValue

	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func postRequestHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	if req.Method == "POST" {

		fmt.Println("2. Performing Http Post...")
		todo := Todo{1, 2, "lorem ipsum dolor sit amet", true}
		jsonReq, err := json.Marshal(todo)
		resp, err := http.Post("https://jsonplaceholder.typicode.com/todos", "application/json; charset=utf-8", bytes.NewBuffer(jsonReq))
		if err != nil {
			log.Fatalln(err)
		}

		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)

		// Convert response body to string
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)

		// Convert response body to Todo struct
		var todoStruct Todo
		json.Unmarshal(bodyBytes, &todoStruct)
		fmt.Printf("%+v\n", todoStruct)

	} else {
		resp.WriteHeader(http.StatusForbidden)
	}

}
