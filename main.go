package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbUser := os.Getenv("BAAS_DB_USER")
	dbPass := os.Getenv("BAAS_DB_PASS")
	dbName := os.Getenv("BAAS_DB_NAME")
	db, err := connectToDb(dbUser, dbPass, dbName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// redisServer := os.Getenv("BAAS_REDIS_SERVER")
	// redisPass := os.Getenv("BAAS_REDIS_PASS")
	// redis := connectToRedis(redisServer, redisPass)

	publicPtr := flag.String("public_key", "jwt.rsa.pub", "The public key used for encoding JWT tokens")
	privatePtr := flag.String("private_key", "jwt.rsa", "The private key used for decoding JWT tokens")
	portPtr := flag.Int("port", 22222, "The port used to serve the service")
	flag.Parse()

	public_key := []byte(*publicPtr)
	private_key := []byte(*privatePtr)

	router := createRouter(db, public_key, private_key)

	log.Fatal(http.ListenAndServe(fmt.Sprint(":%d", *portPtr), router))
}
