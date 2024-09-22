package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/kumersun/bnovo/repository"
)

var guestRepo *repository.GuestRepository

type Config struct {
	Port             string
	PostgresDB       string
	PostgresUser     string
	PostgresPassword string
	PostgresHost     string
	PostgresURL      string
}

var conf Config
var start = time.Now()

func init() {
	godotenv.Load()

	conf.Port = os.Getenv("PORT")
	conf.PostgresDB = os.Getenv("POSTGRES_DB")
	conf.PostgresUser = os.Getenv("POSTGRES_USER")
	conf.PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	conf.PostgresHost = os.Getenv("POSTGRES_HOST")

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		conf.PostgresUser,
		conf.PostgresPassword,
		conf.PostgresHost,
		conf.PostgresDB,
	)

	conf.PostgresURL = connStr
}

func getDebugMemory() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return strconv.FormatUint(m.Sys/1024, 10)
}

func getDebugTime() string {
	return strconv.FormatInt(time.Since(start).Milliseconds(), 10)
}

func main() {
	ctx := context.Background()

	dbpool, err := pgxpool.New(ctx, conf.PostgresURL)
	if err != nil {
		log.Fatal("Unable to create database connection pool:", err)
	}

	defer dbpool.Close()

	err = dbpool.Ping(ctx)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	guestRepo = repository.NewGuestRepository(dbpool)

	http.HandleFunc("/guest", func(w http.ResponseWriter, r *http.Request) {
		start = time.Now()

		switch r.Method {
		case "GET":
			getGuests(w, r)
		case "POST":
			createGuest(w, r)
		default:
			sendErrorResponse(w, errors.New("Method not allowed"), http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/guest/{id}", func(w http.ResponseWriter, r *http.Request) {
		start = time.Now()

		switch r.Method {
		case "GET":
			getGuest(w, r)
		case "PUT":
			updateGuest(w, r)
		case "DELETE":
			deleteGuest(w, r)
		default:
			sendErrorResponse(w, errors.New("Method not allowed"), http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		start = time.Now()

		switch r.Method {
		case "GET":
			w.Write([]byte("OK"))
		default:
			sendErrorResponse(w, errors.New("Method not allowed"), http.StatusMethodNotAllowed)
		}
	})

	fmt.Printf("Server started on port %v\n", conf.Port)
	log.Fatal(http.ListenAndServe(":"+conf.Port, nil))
}
