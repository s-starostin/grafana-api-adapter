package router

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

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

	/*
	   - USERS -
	   Retieving all users:
	   GET
	   .../users

	   Retieving | Deleting single user:
	   GET | DELETE
	   .../users/{id*} (.../users/1 || .../users/id=1 || .../users/login=admin || .../users/email=admin@localhost)
	   .../users/?id=1
	   .../users/?login=test
	   .../users/?email=test@test.test

	   Creating user:
	   POST
	   .../users/ (data: {})
	*/
	f.Group("/users", func() {
		var user grafana.User
		f.Combo("/", func(c flamego.Context) {
			user = grafana.User{}
			user.Id = c.QueryInt64("id")
			user.Login = c.QueryTrim("login")
			user.Email = c.QueryTrim("email")

			if user.Id > 0 || user.Login != "" || user.Email != "" {
				_, err := grafana.GetUser(&user)
				if err != nil && err.Error() == "Empty result" {
					user = grafana.User{}
					c.ResponseWriter().WriteHeader(http.StatusNotFound)
				} else if err != nil {
					log.Print("Got error: " + err.Error())
					c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
				}
			}
		}).Get(func(c flamego.Context) string {
			if user.Id > 0 {
				jsonResponse, err := json.Marshal(user)
				if err != nil {
					log.Print("Got error: " + err.Error())
					c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
					return "null"
				}
				c.ResponseWriter().Header().Add("Content-Type", "application/json")
				return string(jsonResponse)
			}

			c.ResponseWriter().WriteHeader(http.StatusNoContent)
			return "null"
		}).Delete(func(c flamego.Context) string {
			if user.Id == 0 {
				c.ResponseWriter().WriteHeader(http.StatusNotFound)
				return "false"
			}

			status, err := grafana.DeleteUser(&user)
			if err != nil {
				log.Print("Got error: " + err.Error())
				c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
			}
			fmt.Printf("Results: %v\n", status)
			c.ResponseWriter().Header().Add("Content-Type", "application/json")
			return strconv.FormatBool(status)
		}).Post(func(c flamego.Context) string {
			requestBody, err := c.Request().Body().Bytes()
			if err != nil {
				log.Print("Got error: " + err.Error())
			}

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
		f.Combo("/{id}", func(c flamego.Context) {
			user = grafana.User{}
			id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
			if id > 0 {
				user.Id = id
			} else {
				r := regexp.MustCompile(`(?m)^(id|login|email)=([\p{L}\d_!-\.@]+)$`)
				for _, s := range strings.Split(c.Param("id"), ",") {
					parsed := r.FindStringSubmatch(s)
					if parsed == nil || len(parsed) < 3 {
						log.Print("Got error: " + "Unsupported query")
						break
					}
					switch parsed[1] {
					case "id":
						id, err := strconv.ParseInt(parsed[2], 10, 64)
						if err == nil {
							user.Id = id
						} else {
							log.Print("Got error: " + "Unable to parse user id")
						}
					case "login":
						user.Login = parsed[2]
					case "email":
						reMail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
						if reMail.MatchString(parsed[2]) {
							user.Email = parsed[2]
						} else {
							log.Print("Got error: " + "Unable to parse email")
						}
					}
				}
			}

			if user.Id > 0 || user.Login != "" || user.Email != "" {
				_, err := grafana.GetUser(&user)
				if err != nil && err.Error() == "Empty result" {
					user = grafana.User{}
					c.ResponseWriter().WriteHeader(http.StatusNotFound)
				} else if err != nil {
					log.Print("Got error: " + err.Error())
					c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
				}
			}
		}).Get(func(c flamego.Context) string {
			if user.Id == 0 {
				c.ResponseWriter().WriteHeader(http.StatusNotFound)
				return ""
			}
			jsonResponse, err := json.Marshal(user)
			if err != nil {
				log.Print("Got error: " + err.Error())
				c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
				return ""
			}
			c.ResponseWriter().Header().Add("Content-Type", "application/json")
			return string(jsonResponse)
		}).Delete(func(c flamego.Context) string {
			if user.Id == 0 {
				c.ResponseWriter().WriteHeader(http.StatusNotFound)
				return "false"
			}

			status, err := grafana.DeleteUser(&user)
			if err != nil {
				log.Print("Got error: " + err.Error())
				c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
			}
			fmt.Printf("Results: %v\n", status)
			c.ResponseWriter().Header().Add("Content-Type", "application/json")
			return strconv.FormatBool(status)
		})
	})
	f.Get("/users/search/{slug}", func(c flamego.Context) string {
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
	f.Patch("/users/organizations/", func(c flamego.Context) string {
		requestBody, err := c.Request().Body().Bytes()
		if err != nil {
			log.Print("Got error: " + err.Error())
		}

		var userOrganizationsRequest struct {
			User          grafana.User
			Organizations []grafana.UserOrganization
		}

		err = json.Unmarshal(requestBody, &userOrganizationsRequest)
		if err != nil {
			log.Print("Got error: " + err.Error())
		}

		_, err = grafana.GetUser(&userOrganizationsRequest.User)
		if err != nil {
			log.Print("Got error: " + err.Error())
			return "false"
		}

		_, err = grafana.SetUserOrganizations(&userOrganizationsRequest.User, &userOrganizationsRequest.Organizations)
		if err != nil {
			log.Print("Got error: " + err.Error())
		}

		result, err := json.Marshal(userOrganizationsRequest.User)
		if err != nil {
			log.Print("Got error: " + err.Error())
		}

		return string(result)
	})

	/*
	   - ORGANIZATIONS -
	   Retieving all organizations:
	   GET
	   .../organizations

	   Retieving | Deleting single organization:
	   GET | DELETE
	   .../organizations/{id*} (.../organizations/1 || .../organizations/test)
	   .../organizations/?id=1
	   .../organizations/?name=test

	   Creating organization:
	   POST
	   .../organizations/ (data: {})
	*/
	f.Group("/organizations", func() {
		var organization grafana.Organization
		f.Combo("/", func(c flamego.Context) {
			organization = grafana.Organization{}
			organization.Id = c.QueryInt64("id")
			organization.Name = c.QueryTrim("name")

			if organization.Id > 0 || organization.Name != "" {
				_, err := grafana.GetOrganization(&organization)
				if err != nil && err.Error() == "Empty result" {
					organization = grafana.Organization{}
					c.ResponseWriter().WriteHeader(http.StatusNotFound)
				} else if err != nil {
					log.Print("Got error: " + err.Error())
					c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
				}
			}

		}).Get(func(c flamego.Context) string {
			if organization.Id > 0 {
				jsonResponse, err := json.Marshal(organization)
				if err != nil {
					log.Print("Got error: " + err.Error())
					c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
					return "null"
				}
				c.ResponseWriter().Header().Add("Content-Type", "application/json")
				return string(jsonResponse)
			}

			var organizations []grafana.Organization

			organizations, err := grafana.GetOrganizations()
			if err != nil {
				log.Print("Got error: " + err.Error())
				c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
			}
			jsonResponse, err := json.Marshal(organizations)
			if err != nil {
				log.Print("Got error: " + err.Error())
				c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
				return "null"
			}
			c.ResponseWriter().Header().Add("Content-Type", "application/json")
			return string(jsonResponse)
		}).Delete(func(c flamego.Context) string {
			if organization.Id == 0 {
				c.ResponseWriter().WriteHeader(http.StatusNotFound)
				return "false"
			}

			status, err := grafana.DeleteOrganization(&organization)
			if err != nil {
				log.Print("Got error: " + err.Error())
				c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
			}
			fmt.Printf("Results: %v\n", status)
			c.ResponseWriter().Header().Add("Content-Type", "application/json")
			return strconv.FormatBool(status)
		}).Post(func(c flamego.Context) string {
			requestBody, err := c.Request().Body().Bytes()
			if err != nil {
				log.Print("Got error: " + err.Error())
			}

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

		getOrganization := func(uid string) (grafana.Organization, error) {
			organization = grafana.Organization{}
			id, _ := strconv.ParseInt(uid, 10, 64)
			if id > 0 {
				organization.Id = id
				_, err := grafana.GetOrganization(&organization)
				if err != nil {
					return grafana.Organization{}, err
				}
			}

			if len(uid) > 0 && organization.Name == "" {
				reName := regexp.MustCompile(`^([\p{L}\d\s_!-\.@|\]\[\(\)]+)*$`)
				if reName.MatchString(uid) {
					organization.Name = uid
				} else {
					return grafana.Organization{}, errors.New("Unable to parse organization name")
				}

				_, err := grafana.GetOrganization(&organization)
				if err != nil {
					return grafana.Organization{}, err
				}
			}

			return organization, nil
		}

		f.Combo("/{id}", func(c flamego.Context) {
			var err error
			organization, err = getOrganization(c.Param("id"))
			if err != nil && err.Error() == "Empty result" {
				organization = grafana.Organization{}
				c.ResponseWriter().WriteHeader(http.StatusNotFound)
			} else if err != nil {
				log.Print("Got error: " + err.Error())
				c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
			}
		}).Get(func(c flamego.Context) string {
			if organization.Id == 0 {
				c.ResponseWriter().WriteHeader(http.StatusNotFound)
				return ""
			}
			jsonResponse, err := json.Marshal(organization)
			if err != nil {
				log.Print("Got error: " + err.Error())
				c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
				return ""
			}
			c.ResponseWriter().Header().Add("Content-Type", "application/json")
			return string(jsonResponse)
		}).Delete(func(c flamego.Context) string {
			if organization.Id == 0 {
				c.ResponseWriter().WriteHeader(http.StatusNotFound)
				return "false"
			}

			status, err := grafana.DeleteOrganization(&organization)
			if err != nil {
				log.Print("Got error: " + err.Error())
				c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
			}
			fmt.Printf("Results: %v\n", status)
			c.ResponseWriter().Header().Add("Content-Type", "application/json")
			return strconv.FormatBool(status)
		})

		var orgServiceUser grafana.User
		/*
		   - DASHBOARDS FOR ORGANIZATION -

		   Retieving:
		   GET
		   .../organizations/{orgId}/dashboards/ (.../organizations/11/dashboards/ || .../organizations/test/dashboards/)

		   Retieving | Deleting single dashboard:
		   GET | DELETE
		   .../organizations/{orgId}/dashboards/{uid} (.../organizations/11/dashboards/GPXicXZRk || .../organizations/test/dashboards/organization%20title || .../organizations/test/dashboards/23)

		   Creating dashboard:
		   POST
		   .../organizations/{orgId}/dashboards/ (data: {})
		*/

		f.Combo("/{orgId}/dashboards/", func(c flamego.Context) {
			var err error
			organization, err = getOrganization(c.Param("orgId"))
			if err != nil && err.Error() == "Empty result" {
				organization = grafana.Organization{}
				c.ResponseWriter().WriteHeader(http.StatusNotFound)
			} else if err != nil {
				log.Print("Got error: " + err.Error())
				c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
			} else {
				orgServiceUser.Login = "svc" + strconv.FormatInt(organization.Id, 10) + "." + fmt.Sprintf("%x", md5.Sum([]byte(organization.Name)))
				orgServiceUser.Password = util.RandString(12)
				orgServiceUser.OrgId = organization.Id

				_, err = grafana.CreateUser(&orgServiceUser)
				if err != nil && strings.Contains(err.Error(), "already exists") {
					_, err = grafana.GetUser(&orgServiceUser)
					if err != nil {
						log.Print("Got error: " + err.Error())
					}
					_, err = grafana.UpdateUserPassword(&orgServiceUser)
					if err != nil {
						log.Print("Got error: " + err.Error())
					}
				} else if err != nil {
					log.Print("Got error: " + err.Error())
				}

				if orgServiceUser.Id > 0 {
					serviceUserOrgs := []grafana.UserOrganization{}
					serviceUserOrgs = append(serviceUserOrgs, grafana.UserOrganization{
						Id:   organization.Id,
						Name: organization.Name,
						Role: "Admin",
					})

					_, err = grafana.SetUserOrganizations(&orgServiceUser, &serviceUserOrgs)
					if err != nil {
						log.Print("Got error: " + err.Error())
					}
				}
			}
		}).Get(func(c flamego.Context) string {
			if organization.Id == 0 {
				c.ResponseWriter().WriteHeader(http.StatusNotFound)
				return ""
			}

			dashboards, err := grafana.GetDashboardsForUser(&orgServiceUser)

			jsonResponse, err := json.Marshal(dashboards)
			if err != nil {
				log.Print("Got error: " + err.Error())
				c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
				return ""
			}
			c.ResponseWriter().Header().Add("Content-Type", "application/json")
			return string(jsonResponse)
		}).Post(func(c flamego.Context) string {
			requestBody, err := c.Request().Body().Bytes()
			if err != nil {
				log.Print("Got error: " + err.Error())
			}
			var dashboard grafana.Dashboard
			_, err = grafana.GetDashboardForUser(&orgServiceUser, &dashboard)
			if err != nil && err.Error() == "Empty result" {
				dashboard = grafana.Dashboard{}
			}
			err = json.Unmarshal(requestBody, &dashboard)
			if err != nil {
				log.Print("Got error: " + err.Error())
			}
			dashboard.Overwrite = true
			if len(dashboard.Message) == 0 {
				dashboard.Message = "Grafana adapter update " + time.Now().Format("02-01-2006 15:04:05")
			}
			_, err = grafana.UpdateDashboardForUser(&orgServiceUser, &dashboard)
			if err != nil && strings.Contains(err.Error(), "Dashboard not found") {
				c.ResponseWriter().WriteHeader(http.StatusNotFound)
				return "false"
			} else if err != nil && strings.Contains(err.Error(), "Dashboard name cannot be the same") {
				c.ResponseWriter().WriteHeader(http.StatusBadRequest)
				return "false"
			} else if err != nil {
				log.Print("Got error: " + err.Error())
			}

			result, err := json.Marshal(dashboard)
			if err != nil {
				log.Print("Got error: " + err.Error())
			}

			return string(result)
		})
		var dashboard grafana.Dashboard
		f.Combo("/{orgId}/dashboards/{uid}", func(c flamego.Context) {
			var err error
			organization, err = getOrganization(c.Param("orgId"))
			if err != nil && err.Error() == "Empty result" {
				organization = grafana.Organization{}
				c.ResponseWriter().WriteHeader(http.StatusNotFound)
			} else if err != nil {
				log.Print("Got error: " + err.Error())
				c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
			} else {
				orgServiceUser.Login = "svc" + strconv.FormatInt(organization.Id, 10) + "." + fmt.Sprintf("%x", md5.Sum([]byte(organization.Name)))
				orgServiceUser.Password = util.RandString(12)
				orgServiceUser.OrgId = organization.Id

				_, err = grafana.CreateUser(&orgServiceUser)
				if err != nil && strings.Contains(err.Error(), "already exists") {
					_, err = grafana.GetUser(&orgServiceUser)
					if err != nil {
						log.Print("Got error: " + err.Error())
					}
					_, err = grafana.UpdateUserPassword(&orgServiceUser)
					if err != nil {
						log.Print("Got error: " + err.Error())
					}
				} else if err != nil {
					log.Print("Got error: " + err.Error())
				}

				if orgServiceUser.Id > 0 {
					serviceUserOrgs := []grafana.UserOrganization{}
					serviceUserOrgs = append(serviceUserOrgs, grafana.UserOrganization{
						Id:   organization.Id,
						Name: organization.Name,
						Role: "Admin",
					})

					_, err = grafana.SetUserOrganizations(&orgServiceUser, &serviceUserOrgs)
					if err != nil {
						log.Print("Got error: " + err.Error())
					}

					dashboard = grafana.Dashboard{}
					uid := c.Param("uid")
					id, _ := strconv.ParseInt(uid, 10, 64)
					if id > 0 {
						dashboard.Dashboard.Id = id
						_, err := grafana.GetDashboardForUser(&orgServiceUser, &dashboard)
						if err != nil {
							log.Print("Got error: " + err.Error())
						}
					}
					if len(uid) > 0 && dashboard.Dashboard.Title == "" {
						reName := regexp.MustCompile(`^([\p{L}\d\s_!-\.@|\]\[\(\)]+)*$`)
						if reName.MatchString(uid) {
							dashboard.Dashboard.Title = uid
							_, err := grafana.GetDashboardForUser(&orgServiceUser, &dashboard)
							if err != nil {
								log.Print("Got error: " + err.Error())
							}
						}
					}
				}
			}
		}).Get(func(c flamego.Context) string {
			if dashboard.Dashboard.Uid == "" {
				c.ResponseWriter().WriteHeader(http.StatusNotFound)
				return ""
			}
			jsonResponse, err := json.Marshal(&dashboard)
			if err != nil {
				log.Print("Got error: " + err.Error())
				c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
				return ""
			}
			c.ResponseWriter().Header().Add("Content-Type", "application/json")
			return string(jsonResponse)
		}).Delete(func(c flamego.Context) string {
			if dashboard.Dashboard.Uid == "" {
				c.ResponseWriter().WriteHeader(http.StatusNotFound)
				return "false"
			}

			status, err := grafana.DeleteDashboardForUser(&orgServiceUser, &dashboard)
			if err != nil {
				log.Print("Got error: " + err.Error())
				c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
			}

			c.ResponseWriter().Header().Add("Content-Type", "application/json")
			return strconv.FormatBool(status)
		})

		/*
		   - FOLDERS FOR ORGANIZATION -
		   Retieving all folders:
		   GET
		   .../organizations/{orgId}/folders (.../organizations/11/dashboards/)

		   Retieving | Updating single fodler:
		   GET | DELETE
		   .../organizations/{orgId}/folders/{id} (.../organizations/11/folders/nErXDvCkzz | .../organizations/11/folders/11)

		   Creating folder:
		   POST
		   .../organizations/{orgId}/folders/ (data: {})
		*/

		var folder grafana.Folder
		f.Combo("/{orgId}/folders/", func(c flamego.Context) {
			var err error
			organization, err = getOrganization(c.Param("orgId"))
			if err != nil && err.Error() == "Empty result" {
				organization = grafana.Organization{}
				c.ResponseWriter().WriteHeader(http.StatusNotFound)
			} else if err != nil {
				log.Print("Got error: " + err.Error())
				c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
			} else {
				orgServiceUser.Login = "svc" + strconv.FormatInt(organization.Id, 10) + "." + fmt.Sprintf("%x", md5.Sum([]byte(organization.Name)))
				orgServiceUser.Password = util.RandString(12)
				orgServiceUser.OrgId = organization.Id

				_, err = grafana.CreateUser(&orgServiceUser)
				if err != nil && strings.Contains(err.Error(), "already exists") {
					_, err = grafana.GetUser(&orgServiceUser)
					if err != nil {
						log.Print("Got error: " + err.Error())
					}
					_, err = grafana.UpdateUserPassword(&orgServiceUser)
					if err != nil {
						log.Print("Got error: " + err.Error())
					}
				} else if err != nil {
					log.Print("Got error: " + err.Error())
				}

				if orgServiceUser.Id > 0 {
					serviceUserOrgs := []grafana.UserOrganization{}
					serviceUserOrgs = append(serviceUserOrgs, grafana.UserOrganization{
						Id:   organization.Id,
						Name: organization.Name,
						Role: "Admin",
					})

					_, err = grafana.SetUserOrganizations(&orgServiceUser, &serviceUserOrgs)
					if err != nil {
						log.Print("Got error: " + err.Error())
					}
				}
			}
		}).Get(func(c flamego.Context) string {
			if organization.Id == 0 {
				c.ResponseWriter().WriteHeader(http.StatusNotFound)
				return ""
			}
			folders, err := grafana.GetFoldersForUser(&orgServiceUser)

			jsonResponse, err := json.Marshal(folders)
			if err != nil {
				log.Print("Got error: " + err.Error())
				c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
				return ""
			}
			c.ResponseWriter().Header().Add("Content-Type", "application/json")
			return string(jsonResponse)

		}).Post(func(c flamego.Context) string {
			requestBody, err := c.Request().Body().Bytes()
			if err != nil {
				log.Print("Got error: " + err.Error())
			}

			var folder grafana.Folder

			err = json.Unmarshal(requestBody, &folder)
			if err != nil {
				log.Print("Got error: " + err.Error())
			}

			_, err = grafana.CreateFolderForUser(&orgServiceUser, &folder)
			if err != nil && strings.Contains(err.Error(), "Folder name cannot be the same") {
				c.ResponseWriter().WriteHeader(http.StatusBadRequest)
				return "false"
			} else if err != nil && strings.Contains(err.Error(), "already exists") {
				c.ResponseWriter().WriteHeader(http.StatusConflict)
				return "false"
			} else if err != nil {
				log.Print("Got error: " + err.Error())
			}

			result, err := json.Marshal(folder)
			if err != nil {
				log.Print("Got error: " + err.Error())
			}

			return string(result)
		})
		f.Combo("/{orgId}/folders/{id}", func(c flamego.Context) {
			var err error
			organization, err = getOrganization(c.Param("orgId"))
			if err != nil && err.Error() == "Empty result" {
				organization = grafana.Organization{}
				c.ResponseWriter().WriteHeader(http.StatusNotFound)
			} else if err != nil {
				log.Print("Got error: " + err.Error())
				c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
			} else {
				orgServiceUser.Login = "svc" + strconv.FormatInt(organization.Id, 10) + "." + fmt.Sprintf("%x", md5.Sum([]byte(organization.Name)))
				orgServiceUser.Password = util.RandString(12)
				orgServiceUser.OrgId = organization.Id

				_, err = grafana.CreateUser(&orgServiceUser)
				if err != nil && strings.Contains(err.Error(), "already exists") {
					_, err = grafana.GetUser(&orgServiceUser)
					if err != nil {
						log.Print("Got error: " + err.Error())
					}
					_, err = grafana.UpdateUserPassword(&orgServiceUser)
					if err != nil {
						log.Print("Got error: " + err.Error())
					}
				} else if err != nil {
					log.Print("Got error: " + err.Error())
				}

				if orgServiceUser.Id > 0 {
					serviceUserOrgs := []grafana.UserOrganization{}
					serviceUserOrgs = append(serviceUserOrgs, grafana.UserOrganization{
						Id:   organization.Id,
						Name: organization.Name,
						Role: "Admin",
					})

					_, err = grafana.SetUserOrganizations(&orgServiceUser, &serviceUserOrgs)
					if err != nil {
						log.Print("Got error: " + err.Error())
					}

					folder = grafana.Folder{}
					id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
					uid := c.Param("id")
					if id > 0 {
						folder.Id = id
						_, err := grafana.GetFolderByIdForUser(&orgServiceUser, &folder)
						if err != nil && err.Error() == "Empty result" {
							folder = grafana.Folder{}
							c.ResponseWriter().WriteHeader(http.StatusNotFound)
						} else if err != nil {
							log.Print("Got error: " + err.Error())
							c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
						}
					}

					if len(uid) > 0 && folder.Uid == "" {
						folder.Uid = uid
						_, err := grafana.GetFolderForUser(&orgServiceUser, &folder)
						if err != nil && err.Error() == "Empty result" {
							folder = grafana.Folder{}
							c.ResponseWriter().WriteHeader(http.StatusNotFound)
						} else if err != nil {
							log.Print("Got error: " + err.Error())
							c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
						}
					}
				}
			}
		}).Get(func(c flamego.Context) string {
			if folder.Uid == "" {
				c.ResponseWriter().WriteHeader(http.StatusNotFound)
				return ""
			}
			jsonResponse, err := json.Marshal(folder)
			if err != nil {
				log.Print("Got error: " + err.Error())
				c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
				return ""
			}
			c.ResponseWriter().Header().Add("Content-Type", "application/json")
			return string(jsonResponse)
		}).Delete(func(c flamego.Context) string {
			if folder.Uid == "" {
				c.ResponseWriter().WriteHeader(http.StatusNotFound)
				return "false"
			}

			status, err := grafana.DeleteFolderForUser(&orgServiceUser, &folder)
			if err != nil {
				log.Print("Got error: " + err.Error())
				c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
			}

			c.ResponseWriter().Header().Add("Content-Type", "application/json")
			return strconv.FormatBool(status)
		})

		/*
		   - DATASOURCES FOR ORGANIZATION -
		   Retieving all datasources:
		   GET
		   .../organizations/{orgId}/datasources (.../organizations/11/datasources/)

		   Retieving | Updating single datasource:
		   GET | DELETE
		   .../organizations/{orgId}/datasources/{id} (.../organizations/11/datasources/nErXDvCkzz | .../organizations/11/datasources/11)

		   Creating datasource:
		   POST
		   .../organizations/{orgId}/datasources/ (data: {})
		*/

		var datasource grafana.Datasource
		f.Combo("/{orgId}/datasources/", func(c flamego.Context) {
			var err error
			organization, err = getOrganization(c.Param("orgId"))
			if err != nil && err.Error() == "Empty result" {
				organization = grafana.Organization{}
				c.ResponseWriter().WriteHeader(http.StatusNotFound)
			} else if err != nil {
				log.Print("Got error: " + err.Error())
				c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
			} else {
				orgServiceUser.Login = "svc" + strconv.FormatInt(organization.Id, 10) + "." + fmt.Sprintf("%x", md5.Sum([]byte(organization.Name)))
				orgServiceUser.Password = util.RandString(12)
				orgServiceUser.OrgId = organization.Id

				_, err = grafana.CreateUser(&orgServiceUser)
				if err != nil && strings.Contains(err.Error(), "already exists") {
					_, err = grafana.GetUser(&orgServiceUser)
					if err != nil {
						log.Print("Got error: " + err.Error())
					}
					_, err = grafana.UpdateUserPassword(&orgServiceUser)
					if err != nil {
						log.Print("Got error: " + err.Error())
					}
				} else if err != nil {
					log.Print("Got error: " + err.Error())
				}

				if orgServiceUser.Id > 0 {
					serviceUserOrgs := []grafana.UserOrganization{}
					serviceUserOrgs = append(serviceUserOrgs, grafana.UserOrganization{
						Id:   organization.Id,
						Name: organization.Name,
						Role: "Admin",
					})

					_, err = grafana.SetUserOrganizations(&orgServiceUser, &serviceUserOrgs)
					if err != nil {
						log.Print("Got error: " + err.Error())
					}
				}
			}
		}).Get(func(c flamego.Context) string {
			if organization.Id == 0 {
				c.ResponseWriter().WriteHeader(http.StatusNotFound)
				return ""
			}
			datasources, err := grafana.GetDatasourcesForUser(&orgServiceUser)

			jsonResponse, err := json.Marshal(datasources)
			if err != nil {
				log.Print("Got error: " + err.Error())
				c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
				return ""
			}
			c.ResponseWriter().Header().Add("Content-Type", "application/json")
			return string(jsonResponse)

		}).Post(func(c flamego.Context) string {
			requestBody, err := c.Request().Body().Bytes()
			if err != nil {
				log.Print("Got error: " + err.Error())
			}

			var datasource grafana.Datasource

			err = json.Unmarshal(requestBody, &datasource)
			if err != nil {
				log.Print("Got error: " + err.Error())
			}

			_, err = grafana.CreateDatasourceForUser(&orgServiceUser, &datasource)
			if err != nil && strings.Contains(err.Error(), "Required") {
				c.ResponseWriter().WriteHeader(http.StatusUnprocessableEntity)
				return "false"
			} else if err != nil && strings.Contains(err.Error(), "already exists") {
				c.ResponseWriter().WriteHeader(http.StatusConflict)
				return "false"
			} else if err != nil {
				log.Print("Got error: " + err.Error())
			}

			result, err := json.Marshal(datasource)
			if err != nil {
				log.Print("Got error: " + err.Error())
			}

			return string(result)
		})
		f.Combo("/{orgId}/datasources/{id}", func(c flamego.Context) {
			var err error
			organization, err = getOrganization(c.Param("orgId"))
			if err != nil && err.Error() == "Empty result" {
				organization = grafana.Organization{}
				c.ResponseWriter().WriteHeader(http.StatusNotFound)
			} else if err != nil {
				log.Print("Got error: " + err.Error())
				c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
			} else {
				orgServiceUser.Login = "svc" + strconv.FormatInt(organization.Id, 10) + "." + fmt.Sprintf("%x", md5.Sum([]byte(organization.Name)))
				orgServiceUser.Password = util.RandString(12)
				orgServiceUser.OrgId = organization.Id

				_, err = grafana.CreateUser(&orgServiceUser)
				if err != nil && strings.Contains(err.Error(), "already exists") {
					_, err = grafana.GetUser(&orgServiceUser)
					if err != nil {
						log.Print("Got error: " + err.Error())
					}
					_, err = grafana.UpdateUserPassword(&orgServiceUser)
					if err != nil {
						log.Print("Got error: " + err.Error())
					}
				} else if err != nil {
					log.Print("Got error: " + err.Error())
				}

				if orgServiceUser.Id > 0 {
					serviceUserOrgs := []grafana.UserOrganization{}
					serviceUserOrgs = append(serviceUserOrgs, grafana.UserOrganization{
						Id:   organization.Id,
						Name: organization.Name,
						Role: "Admin",
					})

					_, err = grafana.SetUserOrganizations(&orgServiceUser, &serviceUserOrgs)
					if err != nil {
						log.Print("Got error: " + err.Error())
					}

					datasource = grafana.Datasource{}
					id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
					name := c.Param("id")
					if id > 0 {
						datasource.Id = id
						_, err = grafana.GetDatasourceForUser(&orgServiceUser, &datasource)
					}

					if len(name) > 0 && datasource.Uid == "" {
						datasource.Name = name
						_, err = grafana.GetDatasourceForUser(&orgServiceUser, &datasource)
					}

					if len(name) > 0 && datasource.Uid == "" {
						datasource.Uid = name
						_, err = grafana.GetDatasourceForUser(&orgServiceUser, &datasource)
					}

					if err != nil && err.Error() == "Empty result" {
						datasource = grafana.Datasource{}
						c.ResponseWriter().WriteHeader(http.StatusNotFound)
					} else if err != nil {
						log.Print("Got error: " + err.Error())
						c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
					}
				}
			}
		}).Get(func(c flamego.Context) string {
			if datasource.Uid == "" {
				c.ResponseWriter().WriteHeader(http.StatusNotFound)
				return ""
			}
			jsonResponse, err := json.Marshal(datasource)
			if err != nil {
				log.Print("Got error: " + err.Error())
				c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
				return ""
			}
			c.ResponseWriter().Header().Add("Content-Type", "application/json")
			return string(jsonResponse)
		}).Delete(func(c flamego.Context) string {
			if datasource.Uid == "" {
				c.ResponseWriter().WriteHeader(http.StatusNotFound)
				return "false"
			}

			status, err := grafana.DeleteDatasourceForUser(&orgServiceUser, &datasource)
			if err != nil {
				log.Print("Got error: " + err.Error())
				c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
			}

			c.ResponseWriter().Header().Add("Content-Type", "application/json")
			return strconv.FormatBool(status)
		})
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
