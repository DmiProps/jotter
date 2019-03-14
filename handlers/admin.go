package handlers

import (
	"net/http"
	"io"
)

func MainPage(w http.ResponseWriter, r *http.Request) {

	io.WriteString(w, "We running!")

	/*flusher, _ := w.(http.Flusher)

	rdata := requestData{}
	wdata := responseData{}
	json.NewDecoder(r.Body).Decode(&rdata)
	r.Body.Close()

	if rdata.Method == "" {
		// Request from browser.
		io.WriteString(w, "We working!")
	} else if rdata.Method == "stop" {
		// Request method "stop".
		wdata.Response = "OK, we closed."
		json.NewEncoder(w).Encode(wdata)
		flusher.Flush()
		os.Exit(0)
	} else {
		// Unknown method.
		wdata.Response = "Unknown method: " + rdata.Method + ". We running."
		json.NewEncoder(w).Encode(wdata)
	}*/

}