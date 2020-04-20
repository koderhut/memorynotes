package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/koderhut/memorynotes/urlgen"

	"github.com/gorilla/mux"
	"github.com/koderhut/memorynotes/models"
)

// NotesHandler controller
type NotesHandler struct {
	urlGenerator *urlgen.Url
}

type payload struct {
	Content      string `json:"content"`
	NotifyOnRead string `json:"notify-read"`
	Recipient    string `json:"notify-recipient"`
}

type linkReply struct {
	Status bool   `json:"complete"`
	Link   string `json:"link"`
	Id 	   string `json:"note-id"`
}

type contentReply struct {
	Status  bool   `json:"complete"`
	Content string `json:"content"`
}

// Retrieve controller to get the secret note
func (nc NotesHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	status := false
	content := "An error occurred!"

	params := mux.Vars(r)
	note, err := models.Fetch(params["note"])

	if nil == err {
		status = true
		content = note.Content
	}

	reply := contentReply{Status: status, Content: content}

	json.NewEncoder(w).Encode(&reply)
}

// Store controller to save into memory the secret note
func (nc NotesHandler) Store(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var pl payload
	_ = json.NewDecoder(r.Body).Decode(&pl)
	note, err := models.Store(pl.Content)

	status := false
	link := ""

	if nil == err {
		status = true
		link = nc.urlGenerator.Generate("/notes/" + note.ID.String())
	}

	reply := linkReply{Status: status, Link: link, Id: note.ID.String()}

	json.NewEncoder(w).Encode(&reply)
}

// NewNotesHandler initialize a new controller
func NewNotesHandler(u *urlgen.Url) *NotesHandler {
	return &NotesHandler{
		urlGenerator: u,
	}
}
