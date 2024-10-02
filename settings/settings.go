// (w) 2024 by Jan Buchholz. No rights reserved.
// Preferences, type and access functions

package settings

import "github.com/richardwilkes/unison"

type Settings struct {
	WindowRect   unison.Rect
	EmbySecure   bool
	EmbyServer   string
	EmbyPort     string
	EmbyUser     string
	EmbyPassword string
}

var settings Settings

func SetPreferencesDetail(rect unison.Rect, secure bool, server string, port string, user string, password string) {
	settings.WindowRect = rect
	settings.EmbySecure = secure
	settings.EmbyServer = server
	settings.EmbyPort = port
	settings.EmbyUser = user
	settings.EmbyPassword = password
}

func SetPreferences(s Settings) {
	settings = s
}

func GetPreferences() Settings {
	return settings
}

func GetPreferencesDetail() (rect unison.Rect, secure bool, server string, port string, user string, password string) {
	return settings.WindowRect, settings.EmbySecure, settings.EmbyServer, settings.EmbyPort,
		settings.EmbyUser, settings.EmbyPassword
}

func Valid() bool {
	return settings.EmbyServer != "" && settings.EmbyPort != "" && settings.EmbyUser != "" && settings.EmbyPassword != "" &&
		settings.WindowRect.Width > 0 && settings.WindowRect.Height > 0
}
