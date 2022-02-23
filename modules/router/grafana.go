package router

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/flamego/auth"
	"github.com/flamego/flamego"

	grafana "grafana-adapter/modules/external/grafana/apiv1"
	"grafana-adapter/modules/settings"
	"grafana-adapter/modules/util"
)

func init() {
	settings.GetFromDefaultConf()
}

func Start() {
	f := flamego.Classic()
	f.Map(log.New(os.Stdout, "[grafana-adapter] ", 0))

	//users
	f.Get("/user/", func(c flamego.Context) string {
		requestBody, err := c.Request().Body().Bytes()
		if err != nil {
			log.Print("Got error: " + err.Error())
		}

		var user grafana.User
		err = json.Unmarshal(requestBody, &user)
		if err != nil {
			log.Print("Got error: " + err.Error())
		}

		if user.Id == 0 {
			return "null"
		}
		res, err := grafana.GetUser(&user)
		jsonResponse, err := json.Marshal(res)
		if err != nil {
			log.Print("Got error: " + err.Error())
		}
		c.ResponseWriter().Header().Add("Content-Type", "application/json")
		return string(jsonResponse)
	})
	f.Delete("/user/", func(c flamego.Context) string {
		requestBody, err := c.Request().Body().Bytes()
		if err != nil {
			log.Print("Got error: " + err.Error())
		}

		var user grafana.User
		err = json.Unmarshal(requestBody, &user)
		if err != nil {
			log.Print("Got error: " + err.Error())
		}
		_, err = grafana.GetUser(&user)
		if err != nil {
			log.Print("Got error: " + err.Error())
			return "false"
		}

		status, err := grafana.DeleteUser(&user)
		if err != nil {
			log.Print("Got error: " + err.Error())
		}
		fmt.Printf("Results: %v\n", status)
		c.ResponseWriter().Header().Add("Content-Type", "application/json")
		return strconv.FormatBool(status)
	})
	f.Post("/user/", func(c flamego.Context) string {
		requestBody, err := c.Request().Body().Bytes()
		if err != nil {
			log.Print("Got error: " + err.Error())
		}

		var user grafana.User
		err = json.Unmarshal(requestBody, &user)
		if err != nil {
			log.Print("Got error: " + err.Error())
		}

		if user.Password == "" {
			user.Password = util.RandString(12)
		}

		_, err = grafana.CreateUser(&user)
		if err != nil {
			log.Print("Got error: " + err.Error())
		}

		if user.Login == "" {
			user.Login = user.Email
		}

		result, err := json.Marshal(user)
		if err != nil {
			log.Print("Got error: " + err.Error())
		}

		return string(result) //strconv.FormatBool(user.Id > 0);
	})
	f.Get("/user/search/{slug}", func(c flamego.Context) string {
		res, err := grafana.SearchUsers(c.Param("query"))
		if err != nil {
			log.Print("Got error: " + err.Error())
		}
		jsonResponse, err := json.Marshal(res)
		if err != nil {
			log.Print("Got error: " + err.Error())
		}
		c.ResponseWriter().Header().Add("Content-Type", "application/json")
		return string(jsonResponse)
	})
	f.Patch("/user/organizations/", func(c flamego.Context) string {
		requestBody, err := c.Request().Body().Bytes()
		if err != nil {
			log.Print("Got error: " + err.Error())
		}

		var user grafana.User
		err = json.Unmarshal(requestBody, &user)
		if err != nil {
			log.Print("Got error: " + err.Error())
		}
		_, err = grafana.GetUser(&user)
		if err != nil {
			log.Print("Got error: " + err.Error())
			return "false"
		}

		_, err = grafana.SetUserOrganizations(&user)
		if err != nil {
			log.Print("Got error: " + err.Error())
		}

		result, err := json.Marshal(user)
		if err != nil {
			log.Print("Got error: " + err.Error())
		}

		return string(result)
	})

	//organizations
	f.Post("/organization/", func(c flamego.Context) string {
		requestBody, err := c.Request().Body().Bytes()
		if err != nil {
			log.Print("Got error: " + err.Error())
		}

		var organization grafana.Organization
		err = json.Unmarshal(requestBody, &organization)
		if err != nil {
			log.Print("Got error: " + err.Error())
		}

		_, err = grafana.CreateOrganization(&organization)
		if err != nil {
			log.Print("Got error: " + err.Error())
		}

		result, err := json.Marshal(organization)
		if err != nil {
			log.Print("Got error: " + err.Error())
		}

		return string(result)
	})
	f.Delete("/organization/", func(c flamego.Context) string {
		requestBody, err := c.Request().Body().Bytes()
		if err != nil {
			log.Print("Got error: " + err.Error())
		}

		var organization grafana.Organization
		err = json.Unmarshal(requestBody, &organization)
		if err != nil {
			log.Print("Got error: " + err.Error())
		}

		_, err = grafana.GetOrganization(&organization)
		if err != nil {
			log.Print("Got error: " + err.Error())
		}

		result, err := grafana.DeleteOrganization(&organization)
		if err != nil {
			log.Print("Got error: " + err.Error())
		}

		return strconv.FormatBool(result)
	})

	//index
	f.Get("/", func(c flamego.Context) string {
		c.ResponseWriter().Header().Add("Content-Type", "text/html")
		return `<html><head><title>Grafana Adapter</title></head><body style="display: -webkit-flex; display: flex; -webkit-align-items: center; align-items: center; -webkit-justify-content: center; justify-content: center;"><svg xmlns="http://www.w3.org/2000/svg" viewBox="22.4 66 51 51" width="255" height="255"><g transform="translate(-11.520834,168.12362)" style="fill:none;stroke:#181b1f;stroke-opacity:1"><path style="fill:none;fill-opacity:1;stroke:#181b1f;stroke-width:1.16779;stroke-miterlimit:4;stroke-dasharray:none;stroke-opacity:1" d="m 76.721537,-79.583115 1.20611,-0.206027 2.400742,0.413747 3.709327,2.134183 v 1.270057 l -3.709327,2.134183 -2.400742,0.413747 -1.20611,-0.206028 -0.846379,3.158659 1.147605,0.424625 1.872252,1.558704 2.145246,3.702886 -0.635039,1.09993 -4.279473,-0.0063 -2.285954,-0.84212 -0.941542,-0.781486 -2.312334,2.312334 0.781486,0.941541 0.84212,2.285956 0.0063,4.279473 -1.09993,0.635039 -3.702886,-2.145246 -1.558704,-1.872252 -0.424625,-1.147605 -3.158659,0.846379 0.206028,1.206112 -0.413747,2.40074 -2.134184,3.709327 h -1.270056 l -2.134184,-3.709327 -0.413747,-2.40074 0.206028,-1.206112 -3.158659,-0.846379 -0.424625,1.147605 -1.558704,1.872252 -3.702886,2.145246 -1.09993,-0.635039 0.0063,-4.279473 0.84212,-2.285956 0.781486,-0.941541 -2.312334,-2.312334 -0.941542,0.781486 -2.285954,0.84212 -4.279473,0.0063 -0.635039,-1.09993 2.145246,-3.702886 1.872252,-1.558704 1.147605,-0.424625 -0.846379,-3.158659 -1.20611,0.206028 -2.400742,-0.413747 -3.709327,-2.134183 v -1.270057 l 3.709327,-2.134183 2.400742,-0.413747 1.20611,0.206027 0.846379,-3.158658 -1.147605,-0.424625 -1.872252,-1.558704 -2.145246,-3.702886 0.635039,-1.09993 4.279473,0.0063 2.285954,0.84212 0.941542,0.781487 2.312334,-2.312335 -0.781486,-0.941543 -0.84212,-2.285952 -0.0063,-4.279474 1.09993,-0.635038 3.702886,2.145244 1.558704,1.872252 0.424625,1.147606 3.158659,-0.846379 -0.206028,-1.20611 0.413747,-2.400742 2.134184,-3.709328 h 1.270056 l 2.134184,3.709328 0.413747,2.400742 -0.206028,1.20611 3.158659,0.846379 0.424625,-1.147606 1.558704,-1.872252 3.702886,-2.145244 1.09993,0.635038 -0.0063,4.279474 -0.84212,2.285952 -0.781486,0.941543 2.312334,2.312335 0.941542,-0.781487 2.285954,-0.84212 4.279473,-0.0063 0.635039,1.09993 -2.145246,3.702886 -1.872252,1.558704 -1.147605,0.424625 z"/><path style="fill:none;fill-rule:evenodd;stroke:#181b1f;stroke-width:0.3;stroke-linecap:butt;stroke-linejoin:miter;stroke-miterlimit:4;stroke-dasharray:1.2, 0.6, 0.3, 0.6;stroke-dashoffset:0;stroke-opacity:1" d="m 65.296341,-55.946924 c 0,0 4.006185,-1.028639 6.875738,-3.345164 3.112133,-2.512353 6.03531,-5.231873 7.839866,-11.174642 1.814931,-5.976939 0.245709,-11.529764 0.245709,-11.529764"/> </g> <g transform="translate(90.604454,163.0733)" style="fill:none;stroke:#181b1f;stroke-opacity:1"> <path style="fill:#181b1f;fill-opacity:1;stroke:#181b1f;stroke-width:0.244842;stroke-opacity:1" d="m -25.403747,-74.532795 c 0.995674,0.259703 4.450542,-0.520083 4.133957,0.66916 -1.476606,5.546838 -3.391312,3.457763 -4.859548,6.346846 0.81556,-5.143325 0.374266,-12.942475 -9.500961,-19.296605 -3.31201,-1.292251 -7.804906,-1.710191 -11.444118,-0.551766 -3.815007,1.214384 -7.693639,3.811627 -9.455049,7.406963 -2.926055,5.972575 -4.115684,15.823906 -0.09631,19.952237 4.019377,4.12833 16.783249,6.057468 16.783249,6.057468 l 0.208539,1.190747 -0.418793,2.370163 -2.160214,3.662081 h -1.285549 l -2.160213,-3.662081 -0.418793,-2.370163 0.208539,-1.190747 -3.197184,-0.835598 -0.429803,1.132988 -1.577717,1.848404 -3.74805,2.117921 -1.113346,-0.62695 0.0065,-4.224965 0.852391,-2.256836 0.791018,-0.929549 -2.34054,-2.282883 -0.953025,0.771532 -2.313835,0.831394 -4.331668,0.0063 -0.642783,-1.085921 2.171409,-3.65572 1.895087,-1.538851 1.161603,-0.419215 -0.856702,-3.118427 -1.220821,0.203402 -2.430022,-0.408476 -3.754568,-2.107 v -1.25388 l 3.754568,-2.107 2.430022,-0.408476 1.220821,0.203402 0.856702,-3.118427 -1.161603,-0.419215 -1.895087,-1.538851 -2.171409,-3.65572 0.642783,-1.085921 4.331668,0.0063 2.313835,0.831394 0.953025,0.771532 2.34054,-2.282883 -0.791018,-0.929549 -0.852391,-2.256836 -0.0065,-4.224965 1.113346,-0.62695 3.74805,2.117921 1.577717,1.848404 0.429803,1.132988 3.197184,-0.835598 -0.208539,-1.190747 0.418793,-2.370163 2.160213,-3.662081 h 1.285549 l 2.160214,3.662081 0.418793,2.370163 -0.208539,1.190747 3.197183,0.835598 0.429803,-1.132988 1.577715,-1.848404 3.748046,-2.117921 1.113347,0.62695 -0.0065,4.224965 -0.852392,2.256836 -0.791017,0.929549 2.340539,2.282883 0.953025,-0.771532 2.313833,-0.831394 c 0,0 3.864978,-0.464225 4.363438,0.158551 0,0 -1.774762,4.071397 -2.703826,4.512106 l -1.913259,2.02275 z"/> </g> <path d="m 60.945554,79.828945 c -3.56334,-4.16028 -9.55703,-5.98062 -14.42394,-5.62389 -4.8669,0.35673 -9.07393,3.2743 -11.98739,6.83327 -2.87675,3.51414 -4.60608,8.40632 -4.27654,12.93579 0.2986,4.10429 2.47934,8.276235 5.53084,11.037155 7.32934,6.1694 14.21778,7.67727 18.49458,3.15022 5.39942,-4.52307 6.97817,-13.242135 3.08851,-18.571145 -4.42531,-6.06288 -13.67462,-4.9983 -16.01631,1.81537 -0.64012,1.95159 -0.25617,4.50649 0.87661,6.09897 1.13278,1.59248 3.29474,2.519745 4.55026,2.549135 0.2488,0.002 0.49282,-0.0249 0.72414,-0.05599 l 0.20056,-0.0362 c 0.25564,-0.0525 0.42014,-0.3025 0.36741,-0.55804 -0.0467,-0.22039 -0.24509,-0.3752 -0.47021,-0.36647 l -0.1672,0.003 c -0.19143,-0.004 -0.38386,-0.01 -0.5783,-0.0456 -1.63939,-0.30275 -3.61141,-2.60788 -3.6627,-4.74429 -0.0513,-2.1364 2.19753,-4.81163 4.47931,-5.24314 2.28177,-0.43152 5.25781,0.75646 6.49275,3.17014 1.23494,2.41368 0.48466,6.67475 -1.5365,8.91072 -3.01055,3.33051 -8.494872,3.74582 -11.84404,1.58931 -6.216826,-4.00298 -6.930182,-11.704962 -3.65667,-17.37485 2.992971,-5.183978 10.919238,-8.548925 17.64642,-6.48077 3.26282,1.05432 5.7649,2.51848 7.80348,5.58608 0.62978,0.94767 1.0755,1.93042 1.4304,3.09939 0.27695,0.91224 0.82549,-0.59449 1.01429,-2.23306 0.3006,-2.60887 -1.18244,-1.94087 -4.07976,-5.44511 z" style="fill:#181b1f;fill-opacity:1;stroke-width:0.491871"/></svg></body></html>`
	})

	grafana.NewClient("http://"+settings.GrafanaBackend.Host+":"+strconv.Itoa(settings.GrafanaBackend.Port),
		settings.GrafanaBackend.Login, settings.GrafanaBackend.Password)

	if settings.Server.Login != "" && settings.Server.Password != "" {
		f.Use(auth.Basic(settings.Server.Login, settings.Server.Password))
	}
	f.Run(settings.Server.Host, settings.Server.Port)
}
