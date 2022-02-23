package apiv1

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

type Datasource struct {
	Id                int         `json:"id"`
	Uid               string      `json:"uid"`
	OrgId             int         `json:"orgId"`
	Name              string      `json:"name"`
	Type              string      `json:"type"`
	TypeLogoUrl       string      `json:"typeLogoUrl"`
	Proxy             string      `json:"proxy"`
	Url               string      `json:"url"`
	Password          string      `json:"password"`
	User              string      `json:"user"`
	Database          string      `json:"database"`
	BasicAuth         bool        `json:"basicAuth"`
	BasicAuthUser     string      `json:"basicAuthUser"`
	BasicAuthPassword string      `json:"basicAuthPassword"`
	WithCredentials   bool        `json:"withCredentials"`
	IsDefault         bool        `json:"isDefault"`
	ReadOnly          bool        `json:"readOnly"`
	Version           int         `json:"version"`
	JsonData          interface{} `json:"jsonData"`
	SecureJsonFields  interface{} `json:"secureJsonFields"`
}

func CreateDatasource(datasource *Datasource) (*Datasource, error) {
	if datasource == nil {
		return nil, errors.New("Nil pointer")
	}

	slug := "/api/datasources"
	url := grafanaClientSettings.url + slug

	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(datasource)

	req, err := http.NewRequest(http.MethodPost, url, payloadBuffer)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(grafanaClientSettings.login, grafanaClientSettings.password)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 200 {
		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}

		if data["message"] == "Datasource added" {
			datasource.Id = int(data["id"].(float64))
		}

		return datasource, nil
	}

	return datasource, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

func UpdateDatasource(datasource Datasource) (bool, error) {
	slug := "/api/datasources/" + strconv.Itoa(datasource.Id)
	url := grafanaClientSettings.url + slug

	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(datasource)

	req, err := http.NewRequest(http.MethodPut, url, payloadBuffer)
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(grafanaClientSettings.login, grafanaClientSettings.password)

	res, err := client.Do(req)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	if res.StatusCode == 200 {
		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return false, err
		}

		if data["message"] == "Datasource updated" {
			return true, nil
		}
	}

	return false, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

func DeleteDatasource(datasource *Datasource) (bool, error) {
	if datasource == nil {
		return false, errors.New("Nil pointer")
	}

	slug := "/api/datasources/" + strconv.Itoa(datasource.Id)
	url := grafanaClientSettings.url + slug

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(grafanaClientSettings.login, grafanaClientSettings.password)

	res, err := client.Do(req)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	if res.StatusCode == 200 {
		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return false, err
		}

		if data["message"] == "Data source deleted" {
			return true, nil
		}
	}

	return false, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

func GetDatasource(datasource *Datasource) (*Datasource, error) {
	slug := ""

	if datasource == nil {
		return nil, errors.New("Nil pointer")
	} else if datasource.Id > 0 {
		slug = "/api/datasources/" + strconv.Itoa(datasource.Id)
	} else if datasource.Uid != "" {
		slug = "/api/datasources/uid/" + datasource.Uid
	} else if datasource.Name != "" {
		slug = "/api/datasources/name/" + datasource.Name
	} else {
		return nil, errors.New("No Id, Uid, Name has been set for datasource")
	}

	url := grafanaClientSettings.url + slug

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(grafanaClientSettings.login, grafanaClientSettings.password)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 200 {
		err = json.Unmarshal(body, datasource)
		if err != nil {
			return nil, err
		}
		if datasource.Id == 0 {
			return nil, errors.New("Empty result")
		}

		return datasource, nil
	}

	return nil, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

func GetDatasources() (*[]Datasource, error) {
	slug := "/api/datasources"
	url := grafanaClientSettings.url + slug

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(grafanaClientSettings.login, grafanaClientSettings.password)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 200 {
		var datasources []Datasource
		err = json.Unmarshal(body, &datasources)
		if err != nil {
			return nil, err
		}

		return &datasources, nil
	}

	return nil, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}
