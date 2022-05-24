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
	Id                int64       `json:"id,omitempty"`
	Uid               string      `json:"uid,omitempty"`
	OrgId             int64       `json:"orgId,omitempty"`
	Name              string      `json:"name"`
	Type              string      `json:"type"`
	TypeLogoUrl       string      `json:"typeLogoUrl,omitempty"`
	Proxy             string      `json:"proxy,omitempty"`
	Access            string      `json:"access,omitempty"`
	Url               string      `json:"url"`
	Password          string      `json:"password,omitempty"`
	User              string      `json:"user,omitempty"`
	Database          string      `json:"database,omitempty"`
	BasicAuth         bool        `json:"basicAuth,omitempty"`
	BasicAuthUser     string      `json:"basicAuthUser,omitempty"`
	BasicAuthPassword string      `json:"basicAuthPassword,omitempty"`
	WithCredentials   bool        `json:"withCredentials,omitempty"`
	IsDefault         bool        `json:"isDefault,omitempty"`
	ReadOnly          bool        `json:"readOnly,omitempty"`
	Version           int         `json:"version,omitempty"`
	JsonData          interface{} `json:"jsonData,omitempty"`
	SecureJsonFields  interface{} `json:"secureJsonFields,omitempty"`
}

func CreateDatasourceForUser(user *User, datasource *Datasource) (*Datasource, error) {
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
	req.SetBasicAuth(user.Login, user.Password)

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
			datasource.Id = int64(data["id"].(float64))
		}

		return datasource, nil
	}

	return datasource, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

func UpdateDatasourceForUser(user *User, datasource *Datasource) (bool, error) {
	slug := "/api/datasources/" + strconv.FormatInt(datasource.Id, 10)
	url := grafanaClientSettings.url + slug

	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(datasource)

	req, err := http.NewRequest(http.MethodPut, url, payloadBuffer)
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(user.Login, user.Password)

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

func DeleteDatasourceForUser(user *User, datasource *Datasource) (bool, error) {
	slug := ""

	if datasource == nil {
		return false, errors.New("Nil pointer")
	} else if datasource.Id > 0 {
		slug = "/api/datasources/" + strconv.FormatInt(datasource.Id, 10)
	} else if datasource.Uid != "" {
		slug = "/api/datasources/uid/" + datasource.Uid
	} else if datasource.Name != "" {
		slug = "/api/datasources/name/" + datasource.Name
	} else {
		return false, errors.New("No Id, Uid, Name has been set for datasource")
	}

	url := grafanaClientSettings.url + slug

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(user.Login, user.Password)

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

func GetDatasourceForUser(user *User, datasource *Datasource) (*Datasource, error) {
	slug := ""

	if datasource == nil {
		return nil, errors.New("Nil pointer")
	} else if datasource.Id > 0 {
		slug = "/api/datasources/" + strconv.FormatInt(datasource.Id, 10)
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
	req.SetBasicAuth(user.Login, user.Password)

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

func GetDatasourcesForUser(user *User) (*[]Datasource, error) {
	slug := "/api/datasources"
	url := grafanaClientSettings.url + slug

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(user.Login, user.Password)

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
