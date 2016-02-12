package main

import (
	"encoding/gob"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"html/template"
	"log"
	"net/http"
)

type App struct {
	db        *DB
	store     *Store
	templates map[string]*template.Template
}

func Init() *App {
	// allows user to be encoded so that it can be stored in a session
	gob.Register(&User{})

	db := InitDB()
	store := NewStore("todo-proper-secret")
	templates := LoadTemplates()

	app := &App{db, store, templates}

	return app
}

func main() {
	app := Init()

	authMiddleWare := alice.New(app.AuthorizedHandler)

	router := mux.NewRouter()
	router.HandleFunc("/", app.handleGetSubjects)
	router.HandleFunc("/topics/{subjectName}", app.handleGetTopics)
	router.HandleFunc("/threads/{subjectName}/{topicName}", app.handleGetThreads)
	router.HandleFunc("/thread/{subjectName}/{topicName}/{threadID}", app.handleGetThread)

	router.HandleFunc("/login/{username}", app.handleLogin)
	router.HandleFunc("/logout", app.handleLogout)

	router.Handle("/upvote/{threadID}", authMiddleWare.ThenFunc(app.handleAddUpvote)).Methods("POST")
	router.Handle("/upvote/{threadID}", authMiddleWare.ThenFunc(app.handleRemoveUpvote)).Methods("DELETE")

	http.Handle("/", router)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
