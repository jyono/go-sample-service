package main

import (
	"log"

	"github.com/jyono/go-sample-server/internal/server"
)

/* Example requests
curl -X POST localhost:8080/api/v1/log -d '{"record": {"value": "log1"}}'
curl -X GET localhost:8080/api/v1/log -d '{"offset": 0}'
*/
func main() {
	server := server.NewHttpServer(":8080")
	log.Fatal(server.ListenAndServe())
}