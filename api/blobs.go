package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jcelliott/lumber"

	"github.com/mu-box/hoarder/backends"
)

// get returns the data corresponding to specified key
func get(rw http.ResponseWriter, req *http.Request) {

	//
	r, err := backends.Read(req.URL.Query().Get(":blob"))
	if err != nil {
		rw.WriteHeader(404)
		rw.Write([]byte(fmt.Sprintf("%s\n", err)))
		return
	}
	defer r.Close() // close the file

	// pipe the file rather than consume the rams
	_, err = io.Copy(rw, r)
	if err != nil {
		lumber.Error("Failed to pipe file to client - %s", err.Error())
		rw.WriteHeader(500)
		rw.Write([]byte(fmt.Sprintf("%s\n", err)))
		return
	}
}

// getHead returns info pertaining to data corresponding to specified key
func getHead(rw http.ResponseWriter, req *http.Request) {

	// get data information
	fi, err := backends.Stat(req.URL.Query().Get(":blob"))
	if err != nil {
		rw.WriteHeader(404)
		rw.Write([]byte(fmt.Sprintf("%s\n", err)))
		return
	}

	// set useful headers
	rw.Header().Set("Content-Length", fmt.Sprintf("%d", fi.Size))
	rw.Header().Set("Last-Modified", fi.ModTime.Format(time.RFC1123))
	rw.Header().Set("Date", time.Now().UTC().Format(time.RFC1123))

	//
	rw.Write(nil)
}

// create writes data corresponding to specified key and returns a success message
func create(rw http.ResponseWriter, req *http.Request) {

	key := req.URL.Query().Get(":blob")

	//
	if err := backends.Write(key, req.Body); err != nil {
		rw.Write([]byte(fmt.Sprintf("%s\n", err)))
		return
	}

	//
	rw.Write([]byte(fmt.Sprintf("'%s' created!\n", key)))
}

// update writes data corresponding to specified key and returns a success message
func update(rw http.ResponseWriter, req *http.Request) {
	create(rw, req)
}

// delete removes key and corresponding data
func delete(rw http.ResponseWriter, req *http.Request) {

	key := req.URL.Query().Get(":blob")

	//
	if err := backends.Remove(key); err != nil {
		rw.Write([]byte(fmt.Sprintf("%s\n", err)))
		return
	}

	//
	rw.Write([]byte(fmt.Sprintf("'%s' destroyed!\n", key)))
}

// list returns a list of all keys with relevant information
func list(rw http.ResponseWriter, req *http.Request) {

	//
	fis, err := backends.List()
	if err != nil {
		lumber.Error("Failed to list backends - %s", err.Error())
		rw.WriteHeader(500)
		rw.Write([]byte(fmt.Sprintf("%s\n", err)))
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	//
	jfis, err := json.Marshal(fis)
	if err != nil {
		lumber.Error("Failed to marshal results - %s", err.Error())
		rw.WriteHeader(500)
		rw.Write([]byte(fmt.Sprintf("%s\n", err)))
		return
	}

	//
	rw.Write(append(jfis, byte('\n')))
}
