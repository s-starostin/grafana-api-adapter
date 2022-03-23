package settings

var GrafanaBackend = struct {
	Port     int
	Host     string
	Login    string
	Password string
}{
	Port:     3000,
	Host:     "localhost",
	Login:    "admin",
	Password: "admin",
}

func getGrafanaBackendConfigParams() {
	sec := Cfg.Section("grafana")
	GrafanaBackend.Port = sec.Key("PORT").MustInt(3000)
	GrafanaBackend.Host = sec.Key("HOST").MustString("localhost")
	GrafanaBackend.Login = sec.Key("LOGIN").MustString("admin")
	GrafanaBackend.Password = sec.Key("PASSWORD").MustString("admin")
}
