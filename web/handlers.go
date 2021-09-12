package web

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/saarmornel/reading-list/misc"
	"github.com/saarmornel/reading-list/repo"
	"net/http"
)

func GetUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Not protected!\n")
}

func GetMyBookmarks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.Context().Value("Username").(string)

	b, err := repo.GetBookmarks(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bookmarksJson, err := misc.GetJsonFromJsonObjs(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(bookmarksJson))
}

func GetBookmarks(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	username := p.ByName("username")

	b, err := repo.GetBookmarks(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bookmarksJson, err := misc.GetJsonFromJsonObjs(b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(bookmarksJson))
}

func deleteBookmark(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	url := ps.ByName("url")
	username := r.Context().Value("Username")
	status, err := repo.DeleteBookmark(username.(string), url)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	statusJson, err := misc.GetJsonFromJsonObjs(status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(statusJson))
}

func CreateBookmark(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var b repo.Bookmark
	username := r.Context().Value("Username").(string)
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newBookmark, err := repo.CreateBookmark(username, &b)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	bookmarkJson, err := misc.GetJsonFromJsonObjs(newBookmark)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(bookmarkJson))
}

func GetMyUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Protected!\n")
}
