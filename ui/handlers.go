// ---------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz. No rights reserved.
// Event handlers
// ---------------------------------------------------------------------------------------------------------------------

package ui

import (
	"Emby_Explorer/api"
	"Emby_Explorer/assets"
	"Emby_Explorer/models"
	"fmt"
)

var userViews []api.UserView
var movieData []models.MovieData
var tvshowData []models.TVShowData
var homevideoData []models.HomeVideoData

func embyAuthenticateUser() {
	userViews = nil
	err := api.AuthenticateUserInt()
	if err != nil {
		DialogToDisplaySystemError(assets.ErrAuthFailed, err)
		return
	} else {
		userViews, err = api.UserGetViewsInt()
		if err != nil {
			DialogToDisplaySystemError(assets.ErrFetchViewsFailed, err)
			return
		}
		viewsPopupMenu.RemoveAllItems()
		for i, v := range userViews {
			viewsPopupMenu.AddItem(v.Name)
			if i == 0 {
				viewsPopupMenu.SelectIndex(i)
				setFunctions(false, false, true, true)
			}
		}
	}
}

func embyFetchItemsForUser() {
	index := viewsPopupMenu.SelectedIndex()
	view := userViews[index]
	dto, err := api.UserGetItenmsInt(view.Id, view.CollectionType)
	if err != nil {
		DialogToDisplaySystemError(assets.ErrFetchItemsFailed, err)
		return
	}
	mainContent.RemoveAllChildren()
	switch view.CollectionType {
	case api.CollectionMovies:
		movieData = api.GetMovieDisplayData(dto)
		newMovieTable(mainContent, movieData)
		if len(movieData) > 0 {
			models.MovieTable.SelectByIndex(0)
		}
	case api.CollectionTVShows:
		tvshowData = api.GetTVShowDisplayData(dto)
		newTVShowTable(mainContent, tvshowData)
		if len(tvshowData) > 0 {
			models.TVShowTable.SelectByIndex(0)
		}
	case api.CollectionHomeVideos:
		homevideoData = api.GetHomeVideoDisplayData(dto)
		newHomeVideoTable(mainContent, homevideoData)
		if len(homevideoData) > 0 {
			models.HomeVideoTable.SelectByIndex(0)
		}
	default:
	}
}

func embyFetchDetails() {
	if len(userViews) > 0 {
		index := viewsPopupMenu.SelectedIndex()
		view := userViews[index]
		switch view.CollectionType {
		case api.CollectionMovies:
			movie := models.MovieTable.SelectedRows(true)
			for _, mo := range movie {
				fmt.Println(mo.M.Overview)
				break
			}
		case api.CollectionTVShows:
			tvshow := models.TVShowTable.SelectedRows(true)
			for _, tv := range tvshow {
				if tv.M.Type_ == api.EpisodeType {
					fmt.Println(tv.M.Overview)
					break
				}
			}
		case api.CollectionHomeVideos:

		default:
		}
	}
}
