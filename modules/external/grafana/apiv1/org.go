package apiv1

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

type Organization struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name"`
}

type OrganizationUser struct {
	Id         int    `json:"userId"`
	Email      string `json:"email"`
	Name       string `json:"name"`
	Login      string `json:"login"`
	OrgId      int    `json:"orgId"`
	Role       string `json:"role"`
	LastSeenAt string `json:"lastSeenAt"`
}

func GetOrganizations() (*[]Organization, error) {
	slug := "/api/orgs"
	url := grafanaClientSettings.url + slug

	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

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
		var organizations []Organization
		err = json.Unmarshal(body, &organizations)
		if err != nil {
			return nil, err
		}

		return &organizations, nil
	}

	return nil, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

func GetOrganization(organization *Organization) (*Organization, error) {
	slug := ""

	if organization == nil {
		return nil, errors.New("Nil pointer")
	} else if organization.Id > 0 {
		slug = "/api/orgs/" + strconv.Itoa(organization.Id)
	} else if organization.Name != "" {
		slug = "/api/orgs/name/" + organization.Name
	} else {
		return nil, errors.New("No Id, Name has been set for organization")
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
		err = json.Unmarshal(body, organization)
		if err != nil {
			return nil, err
		}
		if organization.Id == 0 {
			return nil, errors.New("Empty result")
		}

		return organization, nil
	}

	return nil, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

func SwitchCurrentOrganizationForUser(user *User, orgId int) (bool, error) {
	if user == nil {
		return false, errors.New("Nil pointer")
	}

	slug := "/api/users/" + strconv.Itoa(user.Id) + "/using/" + strconv.Itoa(orgId)
	url := grafanaClientSettings.url + slug

	req, err := http.NewRequest(http.MethodPost, url, http.NoBody)
	if err != nil {
		return false, err
	}

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

		if data["message"] == "Active organization changed" {
			return true, nil
		}
	}

	return false, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

//Note: The api will work in the following two ways

//Need to set GF_USERS_ALLOW_ORG_CREATE=true
//Set the config value users.allow_org_create to true in ini file

func CreateOrganization(organization *Organization) (*Organization, error) {
	if organization == nil {
		return nil, errors.New("Nil pointer")
	}

	slug := "/api/orgs"
	url := grafanaClientSettings.url + slug

	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(organization)

	req, err := http.NewRequest(http.MethodPost, url, payloadBuffer)
	if err != nil {
		return organization, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(grafanaClientSettings.login, grafanaClientSettings.password)

	res, err := client.Do(req)
	if err != nil {
		return organization, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return organization, err
	}

	if res.StatusCode == 200 {
		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return organization, err
		}

		if data["message"] == "Organization created" {
			organization.Id = int(data["orgId"].(float64))
		}

		return organization, nil
	}

	return organization, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

func UpdateOrganization(organization *Organization) (bool, error) {
	if organization == nil {
		return false, errors.New("Nil pointer")
	}

	slug := "/api/orgs/" + strconv.Itoa(organization.Id)
	url := grafanaClientSettings.url + slug

	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(organization)

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

		if data["message"] == "Organization updated" {
			return true, nil
		}
	}

	return false, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

func DeleteOrganization(organization *Organization) (bool, error) {
	if organization == nil {
		return false, errors.New("Nil pointer")
	}

	slug := "/api/orgs/" + strconv.Itoa(organization.Id)
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

		if data["message"] == "Organization deleted" {
			return true, nil
		}
	}

	return false, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

func GetUsersInOrganization(organization *Organization) (*[]OrganizationUser, error) {
	if organization == nil {
		return nil, errors.New("Nil pointer")
	}

	slug := "/api/orgs/" + strconv.Itoa(organization.Id) + "/users"
	url := grafanaClientSettings.url + slug

	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(organization)

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
		var users []OrganizationUser
		err = json.Unmarshal(body, &users)
		if err != nil {
			return nil, err
		}

		return &users, nil
	}

	return nil, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}
