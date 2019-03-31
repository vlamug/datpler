package input

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func NewAPI() func(listenAddr, path string, lines chan<- string) {
	return func(listenAddr, path string, lines chan<- string) {
		http.HandleFunc(path, func(rw http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				rw.WriteHeader(http.StatusBadRequest)
				log.Fatalf("")
				return
			}

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Printf("could not read body: %s\n", err)
				rw.WriteHeader(http.StatusBadRequest)
				return
			}

			req := &Request{}
			err = json.Unmarshal(body, req)
			if err != nil {
				log.Printf("could unmarshal request: %s\n", err)
				rw.WriteHeader(http.StatusBadRequest)
				return
			}

			lines <- req.Data
		})
		log.Printf("Listening syslog input on %s%s", listenAddr, path)
		http.ListenAndServe(listenAddr, nil)
	}
}

type Request struct {
	Type string `json:"type"`
	Data string `json:"data"`
}
