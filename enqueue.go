package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func handleEnqueuePUT(w http.ResponseWriter, r *http.Request) {

	// TODO add chance to fail

	// check headers
	if !checkHeaders(r) {
		w.Write([]byte("missing expected headers as requested by labscript\n"))
		return
	}

	// read body data
	bodyData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	mid := AppendToQueue(bodyData)
	// append request to queue and format URI with msg id
	uri := fmt.Sprintf(
		"<msg_uri>http://localhost:%d/queue/msg/%d</msg_uri>",
		port, mid)

	// handle headers
	w.Header().Add("Content-Type", "application/xml; charset=utf-8")
	w.Header().Add("Status", "200")
	w.Header().Add("Server", "Trinket")
	w.Header().Add("Connection", "close")
	w.Header().Add("Cache-Control", "private, max-age=0, must-revalidate")

	// write response URI
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(uri))

	log.Printf("PUT request %d in queue\n", mid)
}

func handleEnqueueGET(w http.ResponseWriter, r *http.Request) {
	msgID := mux.Vars(r)["msg_id"]
	username := r.FormValue("username")
	password := r.FormValue("password")

	// TODO handle invalid msg
	// TODO add chance to fail
	// TODO handle msg not processed
	// TODO handle invalid auth

	log.Println(msgID, username, password)
}