package apiv1

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

type DashboardModel struct {
	Id            int64         `json:"id,omitempty"`
	Uid           string        `json:"uid,omitempty"`
	Panels        []interface{} `json:"panels,omitempty"`
	Title         string        `json:"title"`
	Tags          []string      `json:"tags,omitempty"`
	Timezone      string        `json:"timezone,omitempty"`
	SchemaVersion int           `json:"schemaVersion,omitempty"`
	Version       int           `json:"version,omitempty"`
	Refresh       string        `json:"refresh,omitempty"`
}

type DashboardMeta struct {
	FolderId  int64  `json:"folderId,omitempty"`
	FolderUid string `json:"folderUid,omitempty"`
	Url       string `json:"url,omitempty"`
	IsStarred bool   `json:"isStarred,omitempty"`
	IsHome    bool   `json:"isHome,omitempty"`
	CanSave   bool   `json:"canSave,omitempty"`
	CanEdit   bool   `json:"canEdit,omitempty"`
	CanStar   bool   `json:"canStar,omitempty"`
}

type Dashboard struct {
	Dashboard DashboardModel `json:"dashboard,omitempty"`
	Meta      DashboardMeta  `json:"meta,omitempty"`
	FolderId  int64          `json:"folderId,omitempty"`
	FolderUid string         `json:"folderUid,omitempty"`
	Message   string         `json:"message,omitempty"`
	Overwrite bool           `json:"overwrite,omitempty"`
}

func UpdateDashboardForUser(user *User, dashboard *Dashboard) (*Dashboard, error) {
	if user.Login == "" {
		return nil, errors.New("User login must be set")
	}
	if user.Password == "" {
		return nil, errors.New("User password must be set")
	}
	if dashboard == nil {
		return nil, errors.New("Nil pointer")
	}

	slug := "/api/dashboards/db"
	url := grafanaClientSettings.url + slug

	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(dashboard)

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

		if data["message"] == "Dashboard added" {
			dashboard.Dashboard.Id = int64(data["id"].(float64))
		}

		return dashboard, nil
	}

	return dashboard, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

func DeleteDashboardForUser(user *User, dashboard *Dashboard) (bool, error) {
	if user.Login == "" {
		return false, errors.New("User login must be set")
	}
	if user.Password == "" {
		return false, errors.New("User password must be set")
	}
	if dashboard == nil {
		return false, errors.New("Nil pointer")
	}

	slug := "/api/dashboards/uid/" + dashboard.Dashboard.Uid
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

		return true, nil
	}

	return false, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

func GetDashboardForUser(user *User, dashboard *Dashboard) (*Dashboard, error) {
	if user.Login == "" {
		return nil, errors.New("User login must be set")
	}
	if user.Password == "" {
		return nil, errors.New("User password must be set")
	}
	slug := ""
	if dashboard.Dashboard.Id > 0 {
		slug = "/api/search/?dashboardIds=" + strconv.FormatInt(dashboard.Dashboard.Id, 10)
	} else if len(dashboard.Dashboard.Title) > 0 {
		slug = "/api/search/?query=" + dashboard.Dashboard.Title
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
		var dashboardsResponse []DashboardModel
		err = json.Unmarshal(body, &dashboardsResponse)
		if err != nil {
			return nil, err
		}
		if len(dashboardsResponse) == 1 {
			dashboard.Dashboard = dashboardsResponse[0]
			return dashboard, nil
		} else {
			return nil, errors.New("Dashboard not found")
		}
	}

	return nil, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

func GetDashboardForUserByUid(user *User, dashboard *Dashboard) (*Dashboard, error) {
	if user.Login == "" {
		return nil, errors.New("User login must be set")
	}
	if user.Password == "" {
		return nil, errors.New("User password must be set")
	}
	slug := "/api/dashboards/uid/" + dashboard.Dashboard.Uid
	url := grafanaClientSettings.url + slug

	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(dashboard)

	req, err := http.NewRequest(http.MethodGet, url, payloadBuffer)
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
		err = json.Unmarshal(body, dashboard)
		if err != nil {
			return nil, err
		}

		return dashboard, nil
	} else if res.StatusCode == 404 {
		return dashboard, errors.New("Empty result")
	}

	return nil, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

func GetDashboardsForUser(user *User) (*[]Dashboard, error) {
	if user.Login == "" {
		return nil, errors.New("User login must be set")
	}
	if user.Password == "" {
		return nil, errors.New("User password must be set")
	}
	slug := "/api/search/?type=dash-db"
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
		var dashboardsResponse []DashboardModel
		var dashboards []Dashboard
		err = json.Unmarshal(body, &dashboardsResponse)

		if err != nil {
			return nil, err
		}

		for _, entity := range dashboardsResponse {
			dashboards = append(dashboards, Dashboard{Dashboard: entity})
		}

		return &dashboards, nil
	}

	return nil, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}
