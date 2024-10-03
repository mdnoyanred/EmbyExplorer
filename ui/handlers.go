// ---------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz. No rights reserved.
// Event handlers
// ---------------------------------------------------------------------------------------------------------------------

package ui

import (
	"Emby_Explorer/api"
	"Emby_Explorer/assets"
	"Emby_Explorer/models"
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
				setFunctions(false, false, true)
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
	case api.CollectionTVShows:
		tvshowData = api.GetTVShowDisplayData(dto)
		newTVShowTable(mainContent, tvshowData)
	case api.CollectionHomeVideos:
		homevideoData = api.GetHomeVideoDisplayData(dto)
		newHomeVideoTable(mainContent, homevideoData)
	default:
	}
}
