package main

import (
	"fmt"
	"net/http"
)

// serve hits on HTTP Fileserver in plain text
func (cfg *apiConfig) handlerMetrics(respw http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		respw.WriteHeader(http.StatusMethodNotAllowed)
		respw.Write([]byte("not get\n"))
		return
	}

	respw.Header().Add("Content-Type", "text/html")
	respw.WriteHeader(http.StatusOK)
	respw.Write([]byte(fmt.Sprintf(
		`<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>`, cfg.fileserverHits.Load())))
}
