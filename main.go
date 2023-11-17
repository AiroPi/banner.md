package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"regexp"
)

type Datas struct {
	Name        string
	Description string
	GithubDatas *GithubRepo
}

type GithubRepo struct {
	Stars        int `json:"stargazers_count"`
	Forks        int `json:"forks_count"`
	Issues       int `json:"open_issues"`
	PullRequests int `json:"open_pull_requests"`
	Owner        struct {
		Login string `json:"login"`
	} `json:"owner"`
}

var (
	githubApiRegex = regexp.MustCompile(`(.+)/(.+)`)
	t              = getTemplate()
)

func getTemplate() *template.Template {
	t := template.Must(template.New("banner.tmpl").ParseFiles("banner.tmpl"))
	return t
}

func bannerHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("title")
	description := r.URL.Query().Get("desc")
	githubRepo := r.URL.Query().Get("repo")

	d := Datas{
		Name:        name,
		Description: description,
	}

	if r.URL.Path != "/banner" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	if githubRepo != "" {
		url := githubApiRegex.ReplaceAllString(githubRepo, "https://api.github.com/repos/$1/$2")
		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, "The github repository provided is not valid. Make sure it is public. Url must be like: https://githubc.com/username/repo", http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			http.Error(w, "The github repository provided is not valid. Make sure it is public. Url must be like: https://githubc.com/username/repo", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Something went wrong.", http.StatusInternalServerError)
			return
		}

		githubRepo := GithubRepo{}
		err = json.Unmarshal(body, &githubRepo)
		if err != nil {
			http.Error(w, "Something went wrong.", http.StatusInternalServerError)
			return
		}

		d.GithubDatas = &githubRepo
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	err := t.Execute(w, d)
	if err != nil {
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

}

func main() {
	http.HandleFunc("/banner", bannerHandler)

	fmt.Printf("Starting server at port 8080\n")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
