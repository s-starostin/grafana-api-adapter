package settings

import (
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	ini "gopkg.in/ini.v1"

	"grafana-adapter/modules/util"
)

var (
	// AppVer is the version of the current build. It is set in main.go from main.Version.
	AppVer string
	// AppBuiltWith represents a human readable version go runtime build version and build tags. (See main.go formatBuiltWith().)
	AppBuiltWith string
	// AppStartTime store time the app has started
	AppStartTime time.Time
	// AppPath represents the path to the binary
	AppPath string
	// AppWorkPath is the "working directory" of the app.
	AppWorkPath string

	Cfg          *ini.File
	CustomPath   string // Custom directory path
	CustomConf   string
	PIDFile      = "/run/grafana-adapter.pid"
	WritePIDFile bool
)

// loadFromConf initializes configuration context.
// NOTE: do not print any log except error.
func loadFromConf(allowEmpty bool, extraConfig string) {
	Cfg = ini.Empty()

	if WritePIDFile && len(PIDFile) > 0 {
		createPIDFile(PIDFile)
	}

	isFile, err := util.IsFile(CustomConf)
	if err != nil {
		log.Fatal("Unable to check if %s is a file. Error: %v", CustomConf, err)
	}
	if isFile {
		if err := Cfg.Append(CustomConf); err != nil {
			log.Fatal("Failed to load custom conf '%s': %v", CustomConf, err)
		}
	} else if !allowEmpty {
		log.Fatal("Unable to find configuration file: %q.\nEnsure you are running in the correct environment or set the correct configuration file with -c.", CustomConf)
	} // else: no config file, a config file might be created at CustomConf later (might not)

	if extraConfig != "" {
		if err = Cfg.Append([]byte(extraConfig)); err != nil {
			log.Fatal("Unable to append more config: %v", err)
		}
	}
}

func getAppPath() (string, error) {
	var AppPath string
	var err error
	AppPath, err = exec.LookPath(os.Args[0])

	if err != nil {
		return "", err
	}
	AppPath, err = filepath.Abs(AppPath)
	if err != nil {
		return "", err
	}
	// Note: we don't use path.Dir here because it does not handle case
	//	which path starts with two "/" in Windows: "//psf/Home/..."
	return strings.ReplaceAll(AppPath, "\\", "/"), err
}

func getWorkPath(AppPath string) string {
	workPath := AppWorkPath

	if grafanaAdapterWorkPath, ok := os.LookupEnv("GRAFANA_ADAPTER_WORK_DIR"); ok {
		workPath = grafanaAdapterWorkPath
	}
	if len(workPath) == 0 {
		i := strings.LastIndex(AppPath, "/")
		if i == -1 {
			workPath = AppPath
		} else {
			workPath = AppPath[:i]
		}
	}
	return strings.ReplaceAll(workPath, "\\", "/")
}

func createPIDFile(pidPath string) {
	currentPid := os.Getpid()
	if err := os.MkdirAll(filepath.Dir(pidPath), os.ModePerm); err != nil {
		log.Fatal("Failed to create PID folder: %v", err)
	}

	file, err := os.Create(pidPath)
	if err != nil {
		log.Fatal("Failed to create PID file: %v", err)
	}
	defer file.Close()
	if _, err := file.WriteString(strconv.FormatInt(int64(currentPid), 10)); err != nil {
		log.Fatal("Failed to write PID information: %v", err)
	}
}

func SetConf(providedCustom, providedConf, providedWorkPath string) {
	if len(providedWorkPath) != 0 {
		AppWorkPath = filepath.ToSlash(providedWorkPath)
	}
	if grafanaAdapterConf, ok := os.LookupEnv("GRAFANA_ADAPTER_CONF"); ok {
		CustomPath = grafanaAdapterConf
	}
	if len(providedCustom) != 0 {
		CustomPath = providedCustom
	}
	if len(CustomPath) == 0 {
		CustomPath = path.Join(AppWorkPath, "")
	} else if !filepath.IsAbs(CustomPath) {
		CustomPath = path.Join(AppWorkPath, CustomPath)
	}

	if len(providedConf) != 0 {
		CustomConf = providedConf
	}
	if len(CustomConf) == 0 {
		CustomConf = path.Join(CustomPath, "config.ini")
	} else if !filepath.IsAbs(CustomConf) {
		CustomConf = path.Join(CustomPath, CustomConf)
		log.Print("Using 'custom' directory as relative origin for configuration file: '%s'", CustomConf)
	}
}

func getParams() {
	getServerConfigParams()
	getGrafanaBackendConfigParams()
}

func GetFromDefaultConf() {
	SetConf("", "", "")
	loadFromConf(true, "")
	getParams()
}

func init() {
	//log.NewLogger(0, "console", "console", fmt.Sprintf(`{"level": "info", "colorize": %t, "stacktraceLevel": "none"}`, log.CanColorStdout))

	var err error
	if AppPath, err = getAppPath(); err != nil {
		log.Fatal("Failed to get app path: %v", err)
	}
	AppWorkPath = getWorkPath(AppPath)
}
