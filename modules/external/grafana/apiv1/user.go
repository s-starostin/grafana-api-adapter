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

type User struct {
	Id             int64     `json:"id,omitempty"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	Login          string    `json:"login"`
	Password       string    `json:"password,omitempty"`
	Theme          string    `json:"theme"`
	OrgId          int64     `json:"orgId"`
	IsGrafanaAdmin bool      `json:"isGrafanaAdmin"`
	IsDisabled     bool      `json:"isDisabled"`
	IsExternal     bool      `json:"isExternal"`
	AuthLabels     []string  `json:"authLabels"`
	UpdatedAt      time.Time `json:"updatedAt"`
	CreatedAt      time.Time `json:"createdAt"`
	AvatarUrl      string    `json:"avatarUrl"`
}

type UserOrganization struct {
	Id   int64  `json:"orgId,omitempty"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func SearchUsers(query string) (*[]User, error) {
	slug := "/api/users/search"
	url := grafanaClientSettings.url + slug

	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(grafanaClientSettings.login, grafanaClientSettings.password)

	q := req.URL.Query()

	q.Add("query", query)
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
		var data = struct {
			Users []User `json:"users"`
		}{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}

		if len(data.Users) == 0 {
			return nil, errors.New("Empty result")
		}
		return &data.Users, nil
	}

	return nil, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

func GetUser(user *User) (*User, error) {
	slug := ""

	if user == nil {
		return nil, errors.New("Nil pointer")
	} else if user.Id > 0 {
		slug = "/api/users/" + strconv.FormatInt(user.Id, 10)
	} else if user.Email != "" || user.Login != "" {
		slug = "/api/users/lookup"
	} else {
		return nil, errors.New("No Id, Login, Email has been set for user")
	}

	url := grafanaClientSettings.url + slug

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return user, err
	}

	if user.Id == 0 && (user.Email != "" || user.Login != "") {
		q := req.URL.Query()

		q.Add("loginOrEmail",
			func() string {
				if user.Email != "" {
					return user.Email
				} else {
					return user.Login
				}
			}())
		req.URL.RawQuery = q.Encode()
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
		err = json.Unmarshal(body, user)
		if err != nil {
			return nil, err
		}

		return user, nil
	} else if res.StatusCode == 404 {
		return user, errors.New("Empty result")
	}

	return nil, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

func UpdateUser(user *User) (*User, error) {
	slug := "/api/users/" + strconv.FormatInt(user.Id, 10)
	url := grafanaClientSettings.url + slug

	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(user)

	req, err := http.NewRequest(http.MethodPut, url, payloadBuffer)
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

		if data["message"] == "User updated" {
			return user, nil
		}
	}

	return nil, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

func UpdateUserPassword(user *User) (*User, error) {
	if user.Id == 0 {
		return nil, errors.New("No user id provided")
	}
	slug := "/api/admin/users/" + strconv.FormatInt(user.Id, 10) + "/password"
	url := grafanaClientSettings.url + slug
	if user.Password == "" {
		return nil, errors.New("No user password provided")
	}
	jsonReq := `{"password":"` + user.Password + `"}`

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer([]byte(jsonReq)))
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

		if data["message"] == "User password updated" {
			return user, nil
		}
	}

	return nil, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

func SetUserGrafanaAdmin(user *User, isAdmin bool) (bool, error) {
	if user == nil {
		return false, errors.New("Nil pointer")
	}

	slug := "/api/admin/users/" + strconv.FormatInt(user.Id, 10) + "/permissions"
	url := grafanaClientSettings.url + slug

	jsonReq := `{"isGrafanaAdmin": ` + strconv.FormatBool(isAdmin) + "}"

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer([]byte(jsonReq)))
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

		if data["message"] == "User permissions updated" {
			return true, nil
		}
	}

	return false, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

func CreateUser(user *User) (*User, error) {
	if user == nil {
		return nil, errors.New("Nil pointer")
	}

	slug := "/api/admin/users"
	url := grafanaClientSettings.url + slug

	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(user)

	req, err := http.NewRequest(http.MethodPost, url, payloadBuffer)
	if err != nil {
		return user, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(grafanaClientSettings.login, grafanaClientSettings.password)

	res, err := client.Do(req)
	if err != nil {
		return user, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return user, err
	}

	if res.StatusCode == 200 {
		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return nil, err
		}

		if data["message"] == "User created" {
			user.Id = int64(data["id"].(float64))
		}

		return user, nil
	}

	return user, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))

}

func DeleteUser(user *User) (bool, error) {
	if user == nil {
		return false, errors.New("Nil pointer")
	}

	slug := "/api/admin/users/" + strconv.FormatInt(user.Id, 10)
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

		if data["message"] == "User deleted" {
			return true, nil
		}
	}

	return false, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))

}

func GetOrganizationsByUser(user *User) (*[]UserOrganization, error) {
	if user == nil {
		return nil, errors.New("Nil pointer")
	}

	slug := "/api/users/" + strconv.FormatInt(user.Id, 10) + "/orgs"
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
		var organizations []UserOrganization
		err = json.Unmarshal(body, &organizations)
		if err != nil {
			return nil, err
		}

		return &organizations, nil
	}

	return nil, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
}

func SetUserOrganizations(user *User, organizations *[]UserOrganization) (bool, error) {
	if user == nil {
		return false, errors.New("Nil pointer")
	}

	if len(*organizations) == 0 {
		return false, errors.New("Empty organization list")
	}

	currentOrganizations, err := GetOrganizationsByUser(user)
	if err != nil {
		return false, err
	}

EXIST:
	for _, userOrganization := range *organizations {
		organization := Organization{
			Id:   userOrganization.Id,
			Name: userOrganization.Name,
		}

		_, err = GetOrganization(&organization)
		if organization.Id == 0 {
			return false, errors.New("Organization " + organization.Name + " doesn't exist")
		}

		if userOrganization.Role == "" {
			userOrganization.Role = "Viewer"
		}

		userAlreadyInOrg := false

		for _, currentOrganization := range *currentOrganizations {
			if organization.Name == currentOrganization.Name {
				userAlreadyInOrg = true
				if userOrganization.Role == currentOrganization.Role {
					continue EXIST
				}
			}
		}

		httpMethod := ""
		slug := ""
		url := ""
		jsonReq := ""

		if userAlreadyInOrg == true {
			slug = "/api/orgs/" + strconv.FormatInt(organization.Id, 10) + "/users/" + strconv.FormatInt(user.Id, 10)
			url = grafanaClientSettings.url + slug

			jsonReq = `{"role": "` + userOrganization.Role + `"}`
			httpMethod = http.MethodPatch
		} else {
			slug = "/api/orgs/" + strconv.FormatInt(organization.Id, 10) + "/users/"
			url = grafanaClientSettings.url + slug

			jsonReq = `{"loginOrEmail": "` + user.Email + `", "role": "` + userOrganization.Role + `"}`
			httpMethod = http.MethodPost
		}

		req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer([]byte(jsonReq)))
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
		if res.StatusCode != 200 {
			return false, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))
		}
	}

	deleteFromOrganizations := []UserOrganization{}

OK:
	for _, userCurrentOrganization := range *currentOrganizations {
		for _, userOrganization := range *organizations {
			if userOrganization.Name == userCurrentOrganization.Name {
				continue OK
			}
		}
		deleteFromOrganizations = append(deleteFromOrganizations, userCurrentOrganization)
	}

	for _, userOrganization := range deleteFromOrganizations {
		organization := Organization{
			Id:   userOrganization.Id,
			Name: userOrganization.Name,
		}
		_, err := DeleteUserFromOrganization(user, &organization)
		if err != nil {
			if err.Error() == "Cannot remove last organization admin" {
				fmt.Printf("Got error: %v\n", err.Error()+" ("+organization.Name+")")
				continue
			}
			return false, err
		}
	}

	currentOrganizations, err = GetOrganizationsByUser(user)
	if err != nil {
		return false, err
	}

	return true, nil
}

func DeleteUserFromOrganization(user *User, organization *Organization) (bool, error) {
	if user == nil || organization == nil {
		return false, errors.New("Nil pointer")
	}

	slug := "/api/orgs/" + strconv.FormatInt(organization.Id, 10) + "/users/" + strconv.FormatInt(user.Id, 10)
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

		if data["message"] == "User removed from organization" {
			return true, nil
		}
	} else if res.StatusCode == 400 {
		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return false, err
		}

		if data["message"] == "Cannot remove last organization admin" {
			err = errors.New("Cannot remove last organization admin")
			return false, err
		}
	}

	return false, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))

}

func CreateUserApiToken(user *User, organization *Organization) (bool, error) {
	if user == nil || organization == nil {
		return false, errors.New("Nil pointer")
	}

	slug := "/api/orgs/" + strconv.FormatInt(organization.Id, 10) + "/users/" + strconv.FormatInt(user.Id, 10)
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

		if data["message"] == "User removed from organization" {
			return true, nil
		}
	} else if res.StatusCode == 400 {
		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			return false, err
		}

		if data["message"] == "Cannot remove last organization admin" {
			err = errors.New("Cannot remove last organization admin")
			return false, err
		}
	}

	return false, errors.New("Got response: " + strconv.Itoa(res.StatusCode) + ", body: " + string(body))

}
