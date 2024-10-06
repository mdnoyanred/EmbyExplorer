// ---------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// Emby REST API
// ---------------------------------------------------------------------------------------------------------------------

package api

import (
	"bytes"
	"encoding/json"
	"errors"
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
	substItemId   = "$itemid"
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
	GETImages            = "/Items/" + substItemId + "/Images"
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
	client          = "EmbyExplorer"
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
	paraFormat    = "format="
	paraMaxWidth  = "MaxWidth="
	paraMaxHeight = "MaxHeight="
	apiKey        = "api_key="
)

// Supported Emby collection types
const (
	CollectionMovies     = "movies"
	CollectionTVShows    = "tvshows"
	CollectionHomeVideos = "homevideos"
)

var AllowedCollectionTypes = []string{CollectionMovies, CollectionTVShows, CollectionHomeVideos}

// Emby item types
const (
	VideoType   = "Video"
	SeriesType  = "Series"
	SeasonType  = "Season"
	EpisodeType = "Episode"
	MovieType   = "Movie"
	FolderType  = "Folder"
)

const statusCodeOK = 200

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

type ImageFormat string

const (
	ImageFormatBmp ImageFormat = "bmp"
	ImageFormatGif ImageFormat = "gif"
	ImageFormatJpp ImageFormat = "jpp"
	ImageFormatPng ImageFormat = "png"
)

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

func CreateRestUrlForUser(endpoint string, userid string) string {
	url := CreateRestUrl(endpoint)
	url = strings.Replace(url, substUserId, userid, 1)
	return url
}

func CreateRestUrlForPrimaryImage(endpoint string, itemid string) string {
	url := CreateRestUrl(endpoint)
	url = strings.Replace(url, substItemId, itemid, 1)
	url = url + "/" + string(PRIMARY_ImageType) + "/0"
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
	clnt := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jbody))
	if err != nil {
		return err
	}
	req.Header.Add(contentType, contentTypeJSON)
	req.Header.Add(authHeader, header)
	response, err := clnt.Do(req)
	if err != nil {
		return err
	}
	if response.StatusCode != statusCodeOK {
		return errors.New(response.Status)
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

func UserGetViews(userid string, accesstoken string) ([]UserView, error) {
	var userViews = make([]UserView, 0)
	result := QueryResultBaseItemDto{}
	url := CreateRestUrlForUser(GETViews, userid)
	url = url + "?" + apiKey + accesstoken
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != statusCodeOK {
		return nil, errors.New(response.Status)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return userViews, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	for _, item := range result.Items {
		for _, collectionType := range AllowedCollectionTypes {
			if item.CollectionType == collectionType {
				var v = UserView{
					Name:           item.Name,
					CollectionType: item.CollectionType,
					Id:             item.Id,
				}
				userViews = append(userViews, v)
				break
			}
		}
	}
	return userViews, nil
}

func UserGetItems(userid string, collectionid string, collectiontype string, accesstoken string) ([]BaseItemDto, error) {
	var tmp QueryResultBaseItemDto
	var result = make([]BaseItemDto, 0)
	url := CreateRestUrlForUser(GETItems, userid)
	url = url + "?" + apiKey + accesstoken
	url = url + "&" + paraRecursive + "true"
	url = url + "&" + paraParentId + collectionid
	url = url + "&" + paraFields + GetFields(collectiontype) //fields to fetch
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != statusCodeOK {
		return nil, errors.New(response.Status)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &tmp)
	if err != nil {
		return nil, err
	}
	for _, item := range tmp.Items {
		switch collectiontype {
		case CollectionMovies:
			if item.Type_ == MovieType {
				result = append(result, item)
			}
		case CollectionTVShows:
			if item.Type_ == SeriesType || item.Type_ == SeasonType || item.Type_ == EpisodeType {
				result = append(result, item)
			}
		case CollectionHomeVideos:
			if item.Type_ == VideoType || item.Type_ == FolderType {
				result = append(result, item)
			}
		default:
		}
	}
	return result, nil
}

func GetPrimaryImageForItem(itemid string, format ImageFormat, maxwidth string, maxheight string, accesstoken string) ([]byte, error) {
	url := CreateRestUrlForPrimaryImage(GETImages, itemid)
	url = url + "?" + apiKey + accesstoken
	if format == ImageFormatBmp || format == ImageFormatGif || format == ImageFormatJpp || format == ImageFormatPng {
		url = url + "&" + paraFormat + string(format)
	}
	if maxwidth != "" {
		url = url + "&" + paraMaxWidth + maxwidth
	}
	if maxheight != "" {
		url = url + "&" + paraMaxHeight + maxheight
	}
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != statusCodeOK {
		return nil, errors.New(response.Status)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err == nil {
		return body, nil
	}
	return nil, err
}

func AuthenticateUserInt() error {
	return AuthenticateUserByCredentials(embyPreferences.EmbyUser, embyPreferences.EmbyPassword)
}

func UserGetViewsInt() ([]UserView, error) {
	return UserGetViews(EmbySession.User.Id, EmbySession.AccessToken)
}

func UserGetItenmsInt(collectionid string, collectiontype string) ([]BaseItemDto, error) {
	return UserGetItems(EmbySession.User.Id, collectionid, collectiontype, EmbySession.AccessToken)
}

func GetPrimaryImageForItemInt(itemid string, format ImageFormat, maxwidth string, maxheight string) ([]byte, error) {
	return GetPrimaryImageForItem(itemid, format, maxwidth, maxheight, EmbySession.AccessToken)
}

func createPair(key string, value string) string {
	const qu = `"`
	return key + "=" + qu + value + qu
}

func createHeader(userid string) string {
	var h string
	host, _ := os.Hostname()
	h = authType + " " + createPair(authKeyUserId, userid) + ", " + createPair(authKeyClient, client) + ", " +
		createPair(authKeyDevice, runtime.GOOS) + ", " + createPair(authKeyDeviceId, host) + ", " +
		createPair(authKeyVersion, "1.0.0.0")
	return h
}
