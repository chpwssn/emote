package main

import (
	"encoding/json"
	"net/http"

	"github.com/chpwssn/emote/emotestore"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	store := emotestore.Emotestore{
		Rootpath: "data",
	}
	store.Init()
	fs := http.FileServer(http.Dir("static/"))
	r.Handle("/static/", http.StripPrefix("/static/", fs))

	r.HandleFunc("/local/list", func(w http.ResponseWriter, r *http.Request) {
		local := store.AllEmotes()
		json.NewEncoder(w).Encode(local)
	})

	r.HandleFunc("/local/meta/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		result, err := store.GetEmoteRecord(name)
		if err == nil {
			w.WriteHeader(404)
			return
		}
		json.NewEncoder(w).Encode(result)
	})

	r.HandleFunc("/local/{name}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		result, err := store.GetEmoteFileContents(name)
		if err == nil {
			w.WriteHeader(404)
			return
		}
		w.Write(result)
	})

	r.HandleFunc("/local", func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("file")
		result, err := store.StoreNewEmote(r.FormValue("name"), r.FormValue("credit"), file, *header)
		if err != nil {
			w.WriteHeader(403)
			w.Write([]byte(err.Error()))
		}
		json.NewEncoder(w).Encode(result)
	}).Methods("PUT")

	http.ListenAndServe(":80", r)
}
