package cws

import (
	"fmt"
	"net/http"
)

type Server struct {
	Port          string
	GitHubHandler http.HandlerFunc
}

func (s Server) Start() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "Hello from Concourse Webhook Server. See http://github.com/springernature/concourse-webhook-server")
	})
	http.HandleFunc("/github", s.GitHubHandler)
	return http.ListenAndServe(":"+s.Port, nil)
}
