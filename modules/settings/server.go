package settings

var Server = struct {
	Port     int
	Host     string
	Login    string
	Password string
}{
	Port:     80,
	Host:     "localhost",
	Login:    "",
	Password: "",
}

func getServerConfigParams() {
	sec := Cfg.Section("server")
	Server.Port = sec.Key("PORT").MustInt(80)
	Server.Host = sec.Key("HOST").MustString("localhost")
	Server.Login = sec.Key("LOGIN").MustString("")
	Server.Password = sec.Key("PASSWORD").MustString("")
}
