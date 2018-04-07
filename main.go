package urlshortener

import (
	"context"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"encoding/base64"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mjibson/goon"
	qrcode "github.com/skip2/go-qrcode"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

var letters = []rune("23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz")
var baseURL string
var siteKey string
var indexTemplate *template.Template
var statsTemplate *template.Template

func init() {
	baseURL = os.Getenv("BASE_URL")
	siteKey = os.Getenv("RECAPTCHA_SITE_KEY")
	statsTemplate = template.Must(template.ParseFiles("stats.html"))
	indexTemplate = template.Must(template.ParseFiles("index.html"))
	rand.Seed(time.Now().UnixNano())

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", indexHandler)
	r.Post("/", adminHandler)
	r.Get("/{ShortURLID}", handler)

	http.Handle("/", r)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := indexTemplate.ExecuteTemplate(w, "index.html", siteKey)
	if err != nil {
		log.Errorf(appengine.NewContext(r), err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

type ShortURL struct {
	ID  string `datastore:"-" goon:"id"`
	URL string
}

func (s *ShortURL) Base64QRCode() template.URL {
	var png []byte
	png, err := qrcode.Encode(s.URL, qrcode.Low, 256)
	if err != nil {
		return "data:image/png;base64,"
	}

	result := "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)

	return template.URL(result)
}

func (s *ShortURL) ShortURL() template.URL {
	return template.URL(path.Join(baseURL, s.ID))
}

func (s *ShortURL) statsHandler(w http.ResponseWriter, r *http.Request) {
	err := statsTemplate.ExecuteTemplate(w, "stats.html", s)
	if err != nil {
		log.Errorf(appengine.NewContext(r), err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	g := goon.NewGoon(r)

	id := chi.URLParam(r, "ShortURLID")
	su := &ShortURL{ID: id}
	if err := g.Get(su); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	noRedirect := r.URL.Query().Get("noredirect")
	if noRedirect == "true" {
		su.statsHandler(w, r)
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

type RecaptchaResponse struct {
	Success     bool      `json:"success"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

func validate(ctx context.Context, r *http.Request) bool {
	response := r.FormValue("g-recaptcha-response")
	client := urlfetch.Client(ctx)
	result, err := client.PostForm("https://www.google.com/recaptcha/api/siteverify",
		url.Values{
			"secret":   {os.Getenv("RECAPTCHA_SECRET")},
			"remoteip": {r.RemoteAddr},
			"response": {response},
		})
	if err != nil {
		return false
	}
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return false
	}

	rr := RecaptchaResponse{}
	err = json.Unmarshal(body, &rr)
	if err != nil {
		return false
	}
	return rr.Success
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	validated := validate(appengine.NewContext(r), r)
	if !validated {
		w.Write([]byte("Recaptcha validation failed"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	u := r.FormValue("url")
	if _, err := url.ParseRequestURI(u); err != nil {
		w.Write([]byte("Bad URL"))
		w.WriteHeader(http.StatusBadRequest)
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
	su := &ShortURL{ID: id, URL: u}
	_, err := g.Put(su)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, string(su.ShortURL())+"?noredirect=true", http.StatusFound)
}
