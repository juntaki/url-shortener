package urlshortener

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mjibson/goon"
	"google.golang.org/appengine/datastore"
)

var letters = []rune("23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz")

func init() {
	rand.Seed(time.Now().UnixNano())

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", adminHandler)
	r.Get("/{ShortURLID}", handler)

	http.Handle("/", r)
}

type ShortURL struct {
	ID  string `datastore:"-" goon:"id"`
	URL string
}

func handler(w http.ResponseWriter, r *http.Request) {
	g := goon.NewGoon(r)

	id := chi.URLParam(r, "ShortURLID")
	su := &ShortURL{ID: id}
	if err := g.Get(su); err != nil {
		w.Write([]byte("404"))
		w.WriteHeader(http.StatusNotFound)
		return
	}
	http.Redirect(w, r, su.URL, http.StatusFound)
}

func randomID(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		w.Write([]byte("Add query ?url="))
		w.WriteHeader(http.StatusOK)
		return
	}

	g := goon.NewGoon(r)
	id := randomID(3)
	for {
		su := &ShortURL{ID: id}
		if err := g.Get(su); err != nil {
			if err == datastore.ErrNoSuchEntity {
				break
			}
		}
		id = randomID(3)
	}
	url := r.URL.Query().Get("url")
	su := &ShortURL{ID: id, URL: url}
	_, err := g.Put(su)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte("http://s.juntaki.com/" + id))
	w.WriteHeader(http.StatusOK)
	return
}
