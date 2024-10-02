// (w) 2024 by Jan Buchholz. No rights reserved.
// UI, helpers and table panels
// Using Unison library (c) Richard A. Wilkes
// https://github.com/richardwilkes/unison

package ui

import (
	"Emby_Explorer/api"
	"Emby_Explorer/assets"
	"Emby_Explorer/models"
	"github.com/richardwilkes/toolbox/tid"
	"github.com/richardwilkes/unison"
	"github.com/richardwilkes/unison/enums/align"
	"github.com/richardwilkes/unison/enums/behavior"
	"github.com/richardwilkes/unison/enums/paintstyle"
)

const (
	toolbuttonX              = 20
	toolbuttonY              = 20
	toolbarFontSize  float32 = 9
	viewsPopupWidth          = 150
	viewsPopupHeight         = 20
)

var viewsPopupMenu *unison.PopupMenu[string]
var prefsBtn *unison.Button
var authBtn *unison.Button
var fetchBtn *unison.Button

var mainContent *unison.Panel

func newSVGButton(svg *unison.SVG) *unison.Button {
	btn := unison.NewButton()
	btn.HideBase = true
	btn.Drawable = &unison.DrawableSVG{
		SVG:  svg,
		Size: unison.NewSize(toolbuttonX, toolbuttonY),
	}
	btn.Font = unison.LabelFont.Face().Font(toolbarFontSize)
	return btn
}

func createButton(title string, svgcontent string) (*unison.Button, error) {
	svg, err := unison.NewSVGFromContentString(svgcontent)
	if err != nil {
		return nil, err
	}
	btn := newSVGButton(svg)
	btn.SetTitle(title)
	btn.SetLayoutData(align.Middle)
	return btn, nil
}

func createSpacer(width float32, panel *unison.Panel) {
	spacer := &unison.Panel{}
	spacer.Self = spacer
	spacer.SetSizer(func(_ unison.Size) (minSize, prefSize, maxSize unison.Size) {
		minSize.Width = width
		prefSize.Width = width
		maxSize.Width = width
		return
	})
	panel.AddChild(spacer)
}

func setFunctions(prefs bool, auth bool, fetch bool) {
	prefsBtn.SetEnabled(prefs)
	authBtn.SetEnabled(auth)
	fetchBtn.SetEnabled(fetch)
}

func createToolbarPanel() *unison.Panel {
	var err error
	panel := unison.NewPanel()
	panel.SetLayout(&unison.FlowLayout{
		HSpacing: 1,
		VSpacing: unison.StdVSpacing,
	})
	prefsBtn, err = createButton(assets.CapPreferences, assets.IconPreferences)
	if err == nil {
		prefsBtn.SetEnabled(true)
		prefsBtn.SetFocusable(false)
		panel.AddChild(prefsBtn)
		prefsBtn.ClickCallback = func() { PreferencesDialog() }
	}
	authBtn, err = createButton(assets.CapAuthenticate, assets.IconLogin)
	if err == nil {
		authBtn.SetEnabled(true)
		authBtn.SetFocusable(false)
		panel.AddChild(authBtn)
		authBtn.ClickCallback = func() { embyAuthenticateUser() }
	}
	createSpacer(25, panel)
	lblItems := unison.NewLabel()
	lblItems.Font = unison.LabelFont.Face().Font(toolbarFontSize)
	lblItems.SetTitle(assets.CapViews)
	lblItems.SetLayoutData(align.Middle)
	panel.AddChild(lblItems)
	createSpacer(5, panel)
	viewsPopupMenu = unison.NewPopupMenu[string]()
	viewsPopupMenu.SetLayoutData(align.Middle)
	viewsPopupMenu.Font = unison.LabelFont.Face().Font(toolbarFontSize)
	viewsPopupSize := unison.NewSize(viewsPopupWidth, viewsPopupHeight)
	viewsPopupMenu.SetSizer(func(_ unison.Size) (minSize, prefSize, maxSize unison.Size) {
		minSize = viewsPopupSize
		prefSize = viewsPopupSize
		maxSize = viewsPopupSize
		return
	})
	viewsPopupMenu.SetFocusable(false)
	panel.AddChild(viewsPopupMenu)
	createSpacer(5, panel)
	fetchBtn, err = createButton(assets.CapFetch, assets.IconFetch)
	if err == nil {
		fetchBtn.SetEnabled(true)
		fetchBtn.SetFocusable(false)
		panel.AddChild(fetchBtn)
		fetchBtn.ClickCallback = func() { embyFetchItemsForUser() }
	}
	return panel
}

func createTablePanel() *unison.Panel {
	mainContent = unison.NewPanel()
	mainContent.SetLayout(&unison.FlexLayout{
		Columns:  1,
		HSpacing: 1,
		VSpacing: 1,
	})
	mainContent.SetLayoutData(&unison.FlexLayoutData{
		HAlign: align.Fill,
		VAlign: align.Fill,
		HGrab:  true,
		VGrab:  true,
	})
	logo := createEmbyLogo()
	if logo != nil {
		mainContent.AddChild(logo)
	}
	return mainContent.AsPanel()
}

func createEmbyLogo() *unison.Panel {
	var svg *unison.SVG
	var err error
	if svg, err = getEmbyLogo(); err != nil {
		return nil
	}
	panel := unison.NewPanel()
	panel.SetLayoutData(&unison.FlexLayoutData{
		MinSize: unison.NewSize(50, 50),
		HSpan:   1,
		VSpan:   1,
		HAlign:  align.Fill,
		VAlign:  align.Fill,
		HGrab:   true,
		VGrab:   true,
	})
	panel.DrawCallback = func(gc *unison.Canvas, dirty unison.Rect) {
		gc.DrawRect(dirty, unison.ThemeSurface.Light.Paint(gc, dirty, paintstyle.Fill))
		svg.DrawInRectPreservingAspectRatio(gc, panel.ContentRect(false), nil, nil)
	}
	return panel
}

func getEmbyLogo() (*unison.SVG, error) {
	logo, err := unison.NewSVGFromContentString(assets.EmbyLogo)
	if err != nil {
		return nil, err
	}
	return logo, nil
}

func NewMovieTable(content *unison.Panel, movieData []api.MovieData) {
	models.MovieTable = unison.NewTable[*models.MovieRow](&unison.SimpleTableModel[*models.MovieRow]{})
	models.MovieTable.Columns = make([]unison.ColumnInfo, movieNumberOfColumns)
	for i := range models.MovieTable.Columns {
		models.MovieTable.Columns[i].ID = i
		models.MovieTable.Columns[i].Minimum = 20
		models.MovieTable.Columns[i].Maximum = 10000
	}
	rows := make([]*models.MovieRow, 0)
	for _, m := range movieData {
		r := models.NewMovieRow(tid.MustNewTID('a'), m)
		rows = append(rows, r)
	}
	models.MovieTable.SetRootRows(rows)
	models.MovieTable.SizeColumnsToFit(true)
	header := unison.NewTableHeader[*models.MovieRow](models.MovieTable,
		unison.NewTableColumnHeader[*models.MovieRow](movieCaptions[0], ""),
		unison.NewTableColumnHeader[*models.MovieRow](movieCaptions[1], ""),
		unison.NewTableColumnHeader[*models.MovieRow](movieCaptions[2], ""),
		unison.NewTableColumnHeader[*models.MovieRow](movieCaptions[3], ""),
		unison.NewTableColumnHeader[*models.MovieRow](movieCaptions[4], ""),
		unison.NewTableColumnHeader[*models.MovieRow](movieCaptions[5], ""),
		unison.NewTableColumnHeader[*models.MovieRow](movieCaptions[6], ""),
		unison.NewTableColumnHeader[*models.MovieRow](movieCaptions[7], ""),
		unison.NewTableColumnHeader[*models.MovieRow](movieCaptions[8], ""),
		unison.NewTableColumnHeader[*models.MovieRow](movieCaptions[9], ""),
		unison.NewTableColumnHeader[*models.MovieRow](movieCaptions[10], ""),
		unison.NewTableColumnHeader[*models.MovieRow](movieCaptions[11], ""),
	)
	header.SetLayoutData(&unison.FlexLayoutData{
		HAlign: align.Fill,
		VAlign: align.Fill,
		HGrab:  true,
	})
	scrollArea := unison.NewScrollPanel()
	scrollArea.SetContent(models.MovieTable, behavior.Fill, behavior.Fill)
	scrollArea.SetLayoutData(&unison.FlexLayoutData{
		HAlign: align.Fill,
		VAlign: align.Fill,
		HGrab:  true,
		VGrab:  true,
	})
	scrollArea.SetColumnHeader(header)
	content.AddChild(scrollArea)
}

func NewTVShowTable(content *unison.Panel, tvshowData []api.TVShowData) {
	models.TVShowTable = unison.NewTable[*models.TVShowRow](&unison.SimpleTableModel[*models.TVShowRow]{})
	models.TVShowTable.Columns = make([]unison.ColumnInfo, tvshowNumberOfColumns)
	for i := range models.TVShowTable.Columns {
		models.TVShowTable.Columns[i].ID = i
		models.TVShowTable.Columns[i].Minimum = 20
		models.TVShowTable.Columns[i].Maximum = 10000
	}
	rows := make([]*models.TVShowRow, 0)
	for _, m := range tvshowData {
		r := models.NewTVShowRow(tid.MustNewTID('a'), m)
		rows = append(rows, r)
	}
	models.TVShowTable.SetRootRows(rows)
	models.TVShowTable.SizeColumnsToFit(true)
	header := unison.NewTableHeader[*models.TVShowRow](models.TVShowTable,
		unison.NewTableColumnHeader[*models.TVShowRow](tvshowCaptions[0], ""),
		unison.NewTableColumnHeader[*models.TVShowRow](tvshowCaptions[1], ""),
		unison.NewTableColumnHeader[*models.TVShowRow](tvshowCaptions[2], ""),
		unison.NewTableColumnHeader[*models.TVShowRow](tvshowCaptions[3], ""),
		unison.NewTableColumnHeader[*models.TVShowRow](tvshowCaptions[4], ""),
		unison.NewTableColumnHeader[*models.TVShowRow](tvshowCaptions[5], ""),
		unison.NewTableColumnHeader[*models.TVShowRow](tvshowCaptions[6], ""),
		unison.NewTableColumnHeader[*models.TVShowRow](tvshowCaptions[7], ""),
		unison.NewTableColumnHeader[*models.TVShowRow](tvshowCaptions[8], ""),
		unison.NewTableColumnHeader[*models.TVShowRow](tvshowCaptions[9], ""),
		unison.NewTableColumnHeader[*models.TVShowRow](tvshowCaptions[10], ""),
		unison.NewTableColumnHeader[*models.TVShowRow](tvshowCaptions[11], ""),
	)
	header.SetLayoutData(&unison.FlexLayoutData{
		HAlign: align.Fill,
		VAlign: align.Fill,
		HGrab:  true,
	})
	scrollArea := unison.NewScrollPanel()
	scrollArea.SetContent(models.TVShowTable, behavior.Fill, behavior.Fill)
	scrollArea.SetLayoutData(&unison.FlexLayoutData{
		HAlign: align.Fill,
		VAlign: align.Fill,
		HGrab:  true,
		VGrab:  true,
	})
	scrollArea.SetColumnHeader(header)
	content.AddChild(scrollArea)
}

func NewHomeVideoTable(content *unison.Panel, homevideoData []api.HomeVideoData) {
	models.HomeVideoTable = unison.NewTable[*models.HomeVideoRow](&unison.SimpleTableModel[*models.HomeVideoRow]{})
	models.HomeVideoTable.Columns = make([]unison.ColumnInfo, homevideoNumberOfColumns)
	for i := range models.HomeVideoTable.Columns {
		models.HomeVideoTable.Columns[i].ID = i
		models.HomeVideoTable.Columns[i].Minimum = 20
		models.HomeVideoTable.Columns[i].Maximum = 10000
	}
	rows := make([]*models.HomeVideoRow, 0)
	for _, m := range homevideoData {
		r := models.NewHomeVideoRow(tid.MustNewTID('a'), m)
		rows = append(rows, r)
	}
	models.HomeVideoTable.SetRootRows(rows)
	models.HomeVideoTable.SizeColumnsToFit(true)
	header := unison.NewTableHeader[*models.HomeVideoRow](models.HomeVideoTable,
		unison.NewTableColumnHeader[*models.HomeVideoRow](homevideoCaptions[0], ""),
		unison.NewTableColumnHeader[*models.HomeVideoRow](homevideoCaptions[1], ""),
		unison.NewTableColumnHeader[*models.HomeVideoRow](homevideoCaptions[2], ""),
		unison.NewTableColumnHeader[*models.HomeVideoRow](homevideoCaptions[3], ""),
		unison.NewTableColumnHeader[*models.HomeVideoRow](homevideoCaptions[4], ""),
		unison.NewTableColumnHeader[*models.HomeVideoRow](homevideoCaptions[5], ""),
	)
	header.SetLayoutData(&unison.FlexLayoutData{
		HAlign: align.Fill,
		VAlign: align.Fill,
		HGrab:  true,
	})
	scrollArea := unison.NewScrollPanel()
	scrollArea.SetContent(models.HomeVideoTable, behavior.Fill, behavior.Fill)
	scrollArea.SetLayoutData(&unison.FlexLayoutData{
		HAlign: align.Fill,
		VAlign: align.Fill,
		HGrab:  true,
		VGrab:  true,
	})
	scrollArea.SetColumnHeader(header)
	content.AddChild(scrollArea)
}
