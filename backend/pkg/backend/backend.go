package backend

import (
	"fmt"
	"log"
	"net/http"

	"github.com/amadeusitgroup/miniapp/backend/pkg/db"
	"github.com/gorilla/mux"
)

// App contains the App string
type App struct {
	port  string
	mongo string
}

// NewApplication creates and initializes a backend App
func NewApplication(port, mongo string) *App {
	p := fmt.Sprintf(":%s", port)
	return &App{
		port: p,
	}
}

// Run runs the backend application
func (a *App) Run() {
	r := mux.NewRouter()
	// TODO: clean this with Handle pattern and gorilla context
	r.HandleFunc("/airlines", db.HandleReadAirlines).Methods("GET")
	r.HandleFunc("/routes", db.HandleReadRoutes).Methods("GET")
	r.HandleFunc("/airports", db.HandleReadAirports).Methods("GET")
	if err := http.ListenAndServe(a.port, r); err != nil {
		log.Fatal(err)
	}

}
