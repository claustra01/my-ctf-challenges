package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
)

func sanitizer(msg string) (string, error) {
	if len(msg) > 128 {
		return "", fmt.Errorf("too long message")
	}

	if utf8.ValidString(msg) == false {
		return "", fmt.Errorf("invalid character")
	}

	sanitized := bluemonday.StrictPolicy().Sanitize(msg)

	// &lt;\w+
	sanitized = strings.ReplaceAll(sanitized, "&lt;", "<")
	sanitized = strings.ReplaceAll(sanitized, "&gt;", ">")
	var reHTML = regexp.MustCompile(`<(/)?\w+`)

	sanitized = reHTML.ReplaceAllString(sanitized, "")

	return sanitized, nil
}

func generateID() uuid.UUID {
	return uuid.New()
}

func validateID(id string) error {
	_, err := uuid.Parse(id)
	return err
}

func main() {
	mux := http.NewServeMux()

	// top
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open("/app/index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		stat, err := f.Stat()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.ServeContent(w, r, "index.html", stat.ModTime(), f)
	})

	// create note
	mux.HandleFunc("POST /create", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form data", http.StatusBadRequest)
			return
		}

		sanitized, err := sanitizer(r.FormValue("message"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id := generateID()
		filePath := fmt.Sprintf("/app/notes/%s", id)

		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		if _, err := f.Write([]byte(sanitized)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/notes/%s", id.String()), http.StatusSeeOther)
	})

	// read note
	mux.HandleFunc("GET /notes/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if err := validateID(id); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		filePath := fmt.Sprintf("/app/notes/%s", id)
		data, err := os.ReadFile(filePath)
		if err != nil {
			http.Error(w, "note not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		w.Write(data)
	})

	// edit note
	mux.HandleFunc("PUT /notes/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if err := validateID(id); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form data", http.StatusBadRequest)
			return
		}

		sanitized, err := sanitizer(r.FormValue("message"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		filePath := fmt.Sprintf("/app/notes/%s", id)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.Error(w, "note not found", http.StatusNotFound)
			return
		}

		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		if _, err := f.Write([]byte(sanitized)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/notes/%s", id), http.StatusSeeOther)
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
