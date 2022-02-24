package apiv1

import (
	"net/http"
	"time"
)

var grafanaClientSettings = struct {
	url      string
	timeout  time.Duration
	login    string
	password string
}{
	url:      "http://localhost:3000",
	timeout:  30,
	login:    "admin",
	password: "admin",
}

var client http.Client

func NewClient(url, login, password string) {
	grafanaClientSettings.url = url
	grafanaClientSettings.login = login
	grafanaClientSettings.password = password
	buildClient()
}

func buildClient() {
	client = http.Client{Timeout: grafanaClientSettings.timeout * time.Second}
}
