package apiv1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Folder struct {
	Id        int64     `json:"id,omitempty"`
	Uid       string    `json:"uid,omitempty"`
	Title     string    `json:"title"`
	Url       string    `json:"url,omitempty"`
	HasAcl    bool      `json:"hasAcl,omitempty"`
	CanSave   bool      `json:"canSave,omitempty"`
	CanEdit   bool      `json:"canEdit,omitempty"`
	CanAdmin  bool      `json:"canAdmin,omitempty"`
	CreatedBy string    `json:"createdBy,omitempty"`
	Created   time.Time `json:"created,omitempty"`
	UpdatedBy string    `json:"updatedBy,omitempty"`
	Update    time.Time `json:"update,omitempty"`
	Version   int       `json:"version,omitempty"`
	Overwrite bool      `json:"overwrite,omitempty"`
}

func CreateFolderForUser(user *User, folder *Folder) (*Folder, error) {
	if folder == nil {
		return nil, errors.New("Nil pointer")
	}

	slug := "/api/folders"
	url := grafanaClientSettings.url + slug

	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(folder)

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
		err = json.Unmarshal(body, &folder)
		if err != nil {
			return nil, err
		}

		return folder, nil
	}

	return folder, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

func UpdateFolderForUser(user *User, folder *Folder) (bool, error) {
	slug := "/api/folder/" + folder.Uid
	url := grafanaClientSettings.url + slug

	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(folder)

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

		if data["message"] == "Folder updated" {
			return true, nil
		}
	}

	return false, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

func DeleteFolderForUser(user *User, folder *Folder) (bool, error) {
	if folder == nil {
		return false, errors.New("Nil pointer")
	}

	slug := "/api/folders/" + folder.Uid
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

	fmt.Printf("Results: %v\n", res)

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
		return true, nil
	} else if res.StatusCode == 404 {
		return false, errors.New("Folder not found")
	}

	return false, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))

}

func GetFoldersForUser(user *User, limitOptional ...int) ([]Folder, error) {
	limit := 1000

	folders := make([]Folder, 0)

	if len(limitOptional) > 0 {
		limit = limitOptional[0]
	}

	slug := "/api/folders/"
	url := grafanaClientSettings.url + slug

	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(grafanaClientSettings.login, grafanaClientSettings.password)

	q := req.URL.Query()
	q.Add("limit", strconv.Itoa(limit))
	req.URL.RawQuery = q.Encode()

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
		err = json.Unmarshal(body, &folders)
		if err != nil {
			return nil, err
		}
		return folders, nil
	}

	return nil, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

func GetFolderForUser(user *User, folder *Folder) (*Folder, error) {
	slug := "/api/folders/" + folder.Uid
	url := grafanaClientSettings.url + slug

	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(folder)

	req, err := http.NewRequest(http.MethodGet, url, payloadBuffer)
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
		err = json.Unmarshal(body, folder)
		if err != nil {
			return nil, err
		}

		return folder, nil
	} else if res.StatusCode == 404 {
		return folder, errors.New("Empty result")
	}

	return nil, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

func GetFolderByIdForUser(user *User, folder *Folder) (*Folder, error) {
	slug := "/api/folders/id/" + strconv.FormatInt(folder.Id, 10)
	url := grafanaClientSettings.url + slug

	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(folder)

	req, err := http.NewRequest(http.MethodGet, url, payloadBuffer)
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
		err = json.Unmarshal(body, folder)
		if err != nil {
			return nil, err
		}

		return folder, nil
	} else if res.StatusCode == 404 {
		return folder, errors.New("Empty result")
	}

	return nil, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}
