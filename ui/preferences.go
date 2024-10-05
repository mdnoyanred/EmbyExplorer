// ---------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// Save/load preferences (Window size/position & Emby access data)
// ---------------------------------------------------------------------------------------------------------------------

package ui

import (
	"Emby_Explorer/assets"
	"Emby_Explorer/settings"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

const preferencesFileName = "org.janbuchholz.embyexplorer.json"

func SavePreferences() error {
	s := settings.GetPreferences()
	j, err := json.Marshal(s)
	if err == nil {
		dir, _ := os.UserConfigDir()
		dir = filepath.Join(dir, assets.AppName)
		_, err := os.Stat(dir)
		if err != nil {
			if err := os.Mkdir(dir, os.ModePerm); err != nil {
				return err
			}
		}
		fname := filepath.Join(dir, preferencesFileName)
		err = os.WriteFile(fname, j, 0644)
	}
	if err == nil {
		settings.SetPreferences(s)
	}
	return err
}

func LoadPreferences() error {
	var s settings.Settings
	dir, err := os.UserConfigDir()
	dir = filepath.Join(dir, assets.AppName)
	fname := filepath.Join(dir, preferencesFileName)
	j, err := os.Open(fname)
	if err == nil {
		byteValue, _ := io.ReadAll(j)
		_ = j.Close()
		err = json.Unmarshal(byteValue, &s)
	}
	if err == nil {
		settings.SetPreferences(s)
	}
	return err
}
