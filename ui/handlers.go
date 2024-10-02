// (w) 2024 by Jan Buchholz. No rights reserved.
// Event handlers

package ui

import (
	"Emby_Explorer/api"
	"Emby_Explorer/assets"
)

var userViews []api.UserView

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
		movies := api.GetMovieDisplayData(dto)
		NewMovieTable(mainContent, movies)
	case api.CollectionTVShows:
		tvshows := api.GetTVShowDisplayData(dto)
		NewTVShowTable(mainContent, tvshows)
	case api.CollectionHomeVideos:
		homevideos := api.GetHomeVideoDisplayData(dto)
		NewHomeVideoTable(mainContent, homevideos)
	default:
	}
}
