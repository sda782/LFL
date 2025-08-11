package main

import (
	"apfern/lfl/scraper"
	"html/template"
	"net/http"
	"strings"
)

type PageData struct {
	Title   string
	Chapter template.HTML
	Prev    *string
	Next    *string
}

func main() {
	http.HandleFunc("/read", func(w http.ResponseWriter, r *http.Request) {
		chapter_url := r.URL.Query().Get("chapter")
		if chapter_url == "" {
			w.Write([]byte("chapter url is required"))
			return
		}
		tmpl, err := template.ParseFiles("read.html")
		if err != nil {
			w.Write([]byte("Error parsing template"))
			return
		}

		chapter := scraper.GetNovel(&chapter_url)

		next_local_url := r.URL.Host + "/read?chapter=" + *chapter.Next
		prev_local_url := r.URL.Host + "/read?chapter=" + *chapter.Prev

		if strings.Contains(next_local_url, "/null") {
			next_local_url = ""
		}
		if strings.Contains(prev_local_url, "/null") {
			prev_local_url = ""
		}

		data := PageData{
			Title:   chapter.Title,
			Chapter: template.HTML(chapter.Text),
			Next:    &next_local_url,
			Prev:    &prev_local_url,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			w.Write([]byte("Error executing template"))
			return
		}
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("use the /read endpoint to read with a chapter query url"))
	})

	http.ListenAndServe(":8080", nil)
}
