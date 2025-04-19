package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/samokw/workout/internal/api"
	"github.com/samokw/workout/internal/store"
	"github.com/samokw/workout/migrations"
)

// Apllication holds all the applications shared resources
// We can just pass a single struct into the http handlers
type Application struct {
	Logger *log.Logger
	// This is responsible for handling all workouts
	WorkoutHandler *api.WorkoutHandler
	// This is how we will deal with database connections
	DB *sql.DB
}

func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}

	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// This line creates a brand new logger that writes to standard output and prefixes every message with current date and time
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Our stores will go here
	workoutStore := store.NewPostgresWorkoutStore(pgDB)
	// Our handlers will go here
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		DB:             pgDB,
	}

	return app, nil
}

/*
w - http.ResponseWriter when we need to communicate back to the caller if we need to
r - *http.Request this is what the caller is sending us we might make changes to this in our middleware so we need a pointer
*/
func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "status is available\n")
}
