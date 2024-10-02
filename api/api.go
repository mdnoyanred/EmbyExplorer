// (w) 2024 by Jan Buchholz. No rights reserved.
// Emby REST API

package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
)

const (
	substHostname = "$hostname"
	substPort     = "$port"
	substUserId   = "$userid"
)

const (
	standardURL = "http://" + substHostname + ":" + substPort + "/emby"
	secureURL   = "https://" + substHostname + ":" + substPort + "/emby"
)

const (
	GETUsersPublic       = "/Users/Public"
	POSTAuthenticateUser = "/Users/AuthenticateByName"
	GETViews             = "/Users/" + substUserId + "/Views"
	GETItems             = "/Users/" + substUserId + "/Items"
)

// Fields for auth. request
const (
	authType        = "Emby"
	authHeader      = "Authorization"
	authKeyUserId   = "UserId"
	authKeyClient   = "Client"
	authKeyDevice   = "Device"
	authKeyDeviceId = "DeviceId"
	authKeyVersion  = "Version"
)

const (
	contentType     = "Content-Type"
	contentTypeJSON = "application/json"
)

// URL parameters
const (
	paraParentId  = "ParentId="
	paraRecursive = "Recursive="
	paraFields    = "Fields="
	apiKey        = "api_key="
)

// Supported Emby collection types
const (
	CollectionMovies     = "movies"
	CollectionTVShows    = "tvshows"
	CollectionHomeVideos = "homevideos"
)

var allowedCollectionTypes = []string{CollectionMovies, CollectionTVShows, CollectionHomeVideos}

// Emby item types
const (
	videoType   = "Video"
	seriesType  = "Series"
	seasonType  = "Season"
	episodeType = "Episode"
	movieType   = "Movie"
)

// Body for auth. REST call
type authBody struct {
	Username string
	Pw       string
}

// Emby connection settings
type emby struct {
	EmbySecure   bool
	EmbyServer   string
	EmbyPort     string
	EmbyUser     string
	EmbyPassword string
}

// UserView Emby views for current user
type UserView struct {
	Name           string
	CollectionType string
	Id             string
}

var BasicUrl string
var EmbySession AuthenticationResult
var embyPreferences emby

func InitApiPreferences(secure bool, server string, port string, user string, password string) {
	embyPreferences.EmbySecure = secure
	embyPreferences.EmbyServer = server
	embyPreferences.EmbyPort = port
	embyPreferences.EmbyUser = user
	embyPreferences.EmbyPassword = password
	CreateBasicUrl(embyPreferences.EmbySecure, embyPreferences.EmbyServer, embyPreferences.EmbyPort)
}

func CreateBasicUrl(secure bool, hostname string, port string) {
	var url string
	if secure {
		url = secureURL
	} else {
		url = standardURL
	}
	BasicUrl = strings.Replace(url, substHostname, hostname, 1)
	BasicUrl = strings.Replace(BasicUrl, substPort, port, 1)
}

func CreateRestUrl(endpoint string) string {
	return BasicUrl + endpoint
}

func CreateRestUrlForUser(endpoint string, id string) string {
	url := CreateRestUrl(endpoint)
	url = strings.Replace(url, substUserId, id, 1)
	return url
}

func FindUserIdByName(username string) (string, error) {
	var users []UserDto
	var response *http.Response
	var err error
	var body []byte
	response, err = http.Get(CreateRestUrl(GETUsersPublic))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err = io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal(body, &users)
	if err != nil {
		return "", err
	}
	for _, user := range users {
		if strings.ToUpper(user.Name) == strings.ToUpper(username) {
			return user.Id, nil
		}
	}
	return "", nil
}

func AuthenticateUserByCredentials(username string, password string) error {
	id, err := FindUserIdByName(username)
	if err != nil {
		return err
	}
	var result AuthenticationResult
	body := authBody{username, password}
	jbody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	url := CreateRestUrl(POSTAuthenticateUser)
	header := createHeader(id)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jbody))
	if err != nil {
		return err
	}
	req.Header.Add(contentType, contentTypeJSON)
	req.Header.Add(authHeader, header)
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return err
	}
	EmbySession = result
	return nil
}

func UserGetViews(id string, accesstoken string) ([]UserView, error) {
	var userViews = make([]UserView, 0)
	result := QueryResultBaseItemDto{}
	url := CreateRestUrlForUser(GETViews, id)
	url = url + "?" + apiKey + accesstoken
	response, err := http.Get(url)
	if err != nil {
		return userViews, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return userViews, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return userViews, err
	}
	for _, item := range result.Items {
		for _, collectionType := range allowedCollectionTypes {
			if item.CollectionType == collectionType {
				var v = UserView{
					Name:           item.Name,
					CollectionType: item.CollectionType,
					Id:             item.Id,
				}
				userViews = append(userViews, v)
			}
		}
	}
	return userViews, nil
}

func UserGetItenms(id string, collectionid string, collectiontype string, accesstoken string) ([]BaseItemDto, error) {
	var tmp QueryResultBaseItemDto
	var result = make([]BaseItemDto, 0)
	url := CreateRestUrlForUser(GETItems, id)
	url = url + "?" + apiKey + accesstoken
	url = url + "&" + paraRecursive + "true"
	url = url + "&" + paraParentId + collectionid
	url = url + "&" + paraFields + GetFields(collectiontype) //fields to fetch
	response, err := http.Get(url)
	if err != nil {
		return result, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &tmp)
	if err != nil {
		return result, err
	}
	for _, item := range tmp.Items {
		switch collectiontype {
		case CollectionMovies:
			if item.Type_ == movieType {
				result = append(result, item)
			}
		case CollectionTVShows:
			if item.Type_ == seriesType || item.Type_ == seasonType || item.Type_ == episodeType {
				result = append(result, item)
			}
		case CollectionHomeVideos:
			if item.Type_ == videoType {
				result = append(result, item)
			}
		default:
		}
	}
	return result, nil
}

func AuthenticateUserInt() error {
	return AuthenticateUserByCredentials(embyPreferences.EmbyUser, embyPreferences.EmbyPassword)
}

func UserGetViewsInt() ([]UserView, error) {
	return UserGetViews(EmbySession.User.Id, EmbySession.AccessToken)
}

func UserGetItenmsInt(collectionid string, collectiontype string) ([]BaseItemDto, error) {
	return UserGetItenms(EmbySession.User.Id, collectionid, collectiontype, EmbySession.AccessToken)
}

func createPair(key string, value string) string {
	const qu = `"`
	return key + "=" + qu + value + qu
}

func createHeader(id string) string {
	var h string
	host, _ := os.Hostname()
	h = authType + " " + createPair(authKeyUserId, id) + ", " + createPair(authKeyClient, "PC") + ", " +
		createPair(authKeyDevice, runtime.GOOS) + ", " + createPair(authKeyDeviceId, host) + ", " +
		createPair(authKeyVersion, "1.0.0.0")
	return h
}
