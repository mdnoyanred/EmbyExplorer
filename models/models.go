// ---------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// Data models for Emby Movies, TV Shows and Home Videos, according to Unison's table model
// Using Unison library (c) Richard A. Wilkes
// https://github.com/richardwilkes/unison
// ---------------------------------------------------------------------------------------------------------------------

package models

import (
	"github.com/richardwilkes/toolbox/errs"
	"github.com/richardwilkes/toolbox/fatal"
	"github.com/richardwilkes/toolbox/tid"
	"github.com/richardwilkes/unison"
)

type TableDescription struct {
	NoOfColumns int
	Captions    []string
	APIFields   string
}

// ---------------------------------------------------------------------------------------------------------------------
// Movies model
// ---------------------------------------------------------------------------------------------------------------------

var _ unison.TableRowData[*MovieRow] = &MovieRow{}
var MovieTable *unison.Table[*MovieRow]
var MovieTableDescription = TableDescription{
	NoOfColumns: 12, //displayed columns only
	Captions: []string{"Title", "Original Title", "Year", "Time", "Actors", "Director", "Studio", "Genre", "Ext.",
		"Codec", "Resolution", "Path"},
	APIFields: "Name,OriginalTitle,MediaSources,Path,Genres,ProductionYear,People,Studios,Width,Height,Container," +
		"Overview,RunTimeTicks,Type_", //no spaces here!
}

type MovieData struct {
	Name           string
	OriginalTitle  string
	ProductionYear string
	Runtime        string
	Actors         string
	Directors      string
	Studios        string
	Genres         string
	Container      string
	Codecs         string
	Resolution     string
	Path           string
	Overview       string
	MovieId        string
}

type MovieRow struct {
	table        *unison.Table[*MovieRow]
	parent       *MovieRow
	children     []*MovieRow
	container    bool
	open         bool
	doubleHeight bool
	id           tid.TID
	M            MovieData
}

func (d *MovieRow) CloneForTarget(target unison.Paneler, newParent *MovieRow) *MovieRow {
	table, ok := target.(*unison.Table[*MovieRow])
	if !ok {
		fatal.IfErr(errs.New("invalid target"))
	}
	clone := *d
	clone.table = table
	clone.parent = newParent
	clone.id = tid.MustNewTID('a')
	return &clone
}

func (d *MovieRow) ID() tid.TID {
	return d.id
}

func (d *MovieRow) Parent() *MovieRow {
	return d.parent
}

func (d *MovieRow) SetParent(parent *MovieRow) {
	d.parent = parent
}

func (d *MovieRow) CanHaveChildren() bool {
	return d.container
}

func (d *MovieRow) Children() []*MovieRow {
	return d.children
}

func (d *MovieRow) SetChildren(children []*MovieRow) {
	d.children = children
}

func (d *MovieRow) CellDataForSort(col int) string {
	switch col {
	case 0:
		return d.M.Name
	case 1:
		return d.M.OriginalTitle
	case 2:
		return d.M.ProductionYear
	case 3:
		return d.M.Runtime
	case 4:
		return d.M.Actors
	case 5:
		return d.M.Directors
	case 6:
		return d.M.Studios
	case 7:
		return d.M.Genres
	case 8:
		return d.M.Container
	case 9:
		return d.M.Codecs
	case 10:
		return d.M.Resolution
	case 11:
		return d.M.Path
	default:
		return ""
	}
}

func (d *MovieRow) ColumnCell(_, col int, foreground, _ unison.Ink, _, _, _ bool) unison.Paneler {
	var text string
	switch col {
	case 0:
		text = d.M.Name
	case 1:
		text = d.M.OriginalTitle
	case 2:
		text = d.M.ProductionYear
	case 3:
		text = d.M.Runtime
	case 4:
		text = d.M.Actors
	case 5:
		text = d.M.Directors
	case 6:
		text = d.M.Studios
	case 7:
		text = d.M.Genres
	case 8:
		text = d.M.Container
	case 9:
		text = d.M.Codecs
	case 10:
		text = d.M.Resolution
	case 11:
		text = d.M.Path
	default:
		text = ""
	}
	wrapper := unison.NewPanel()
	wrapper.SetLayout(&unison.FlexLayout{Columns: 1})
	addText(wrapper, text, foreground, unison.LabelFont)
	return wrapper
}

func (d *MovieRow) IsOpen() bool {
	return d.open
}

func (d *MovieRow) SetOpen(open bool) {
	d.open = open
}

func NewMovieRow(id tid.TID, data MovieData) *MovieRow {
	row := &MovieRow{
		table:     MovieTable,
		id:        id,
		container: false,
		open:      false,
		parent:    nil,
		children:  nil,
		M: MovieData{data.Name, data.OriginalTitle, data.ProductionYear,
			data.Runtime, data.Actors, data.Directors, data.Studios,
			data.Genres, data.Container, data.Codecs, data.Resolution,
			data.Path, data.Overview, data.MovieId},
	}
	return row
}

// ---------------------------------------------------------------------------------------------------------------------
// TV shows model
// ---------------------------------------------------------------------------------------------------------------------

var _ unison.TableRowData[*TVShowRow] = &TVShowRow{}
var TVShowTable *unison.Table[*TVShowRow]
var TVShowTableDescription = TableDescription{
	NoOfColumns: 12, //displayed columns only
	Captions: []string{"Series", "Episode", "Season", "Year", "Time", "Actors", "Studio", "Genre", "Ext.", "Codec",
		"Resolution", "Path"},
	APIFields: "Name,MediaSources,Path,Genres,ProductionYear,People,Studios,Width,Height,Container,RunTimeTicks," +
		"Overview,SeriesId,SeasonId,Id,ParentId,IndexNumber,Type_", //no spaces here!
}

type TVShowData struct {
	Name           string
	Episode        string
	Season         string
	ProductionYear string
	Runtime        string
	Actors         string
	Studios        string
	Genres         string
	Container      string
	Codecs         string
	Resolution     string
	Path           string
	Overview       string
	SeriesId       string
	SeasonId       string
	EpisodeId      string
	Type_          string
	SortIndex      int32
}

type TVShowRow struct {
	table        *unison.Table[*TVShowRow]
	parent       *TVShowRow
	children     []*TVShowRow
	container    bool
	open         bool
	doubleHeight bool
	id           tid.TID
	M            TVShowData
}

func (d *TVShowRow) CloneForTarget(target unison.Paneler, newParent *TVShowRow) *TVShowRow {
	table, ok := target.(*unison.Table[*TVShowRow])
	if !ok {
		fatal.IfErr(errs.New("invalid target"))
	}
	clone := *d
	clone.table = table
	clone.parent = newParent
	clone.id = tid.MustNewTID('a')
	return &clone
}

func (d *TVShowRow) ID() tid.TID {
	return d.id
}

func (d *TVShowRow) Parent() *TVShowRow {
	return d.parent
}

func (d *TVShowRow) SetParent(parent *TVShowRow) {
	d.parent = parent
}

func (d *TVShowRow) CanHaveChildren() bool {
	return d.container
}

func (d *TVShowRow) Children() []*TVShowRow {
	return d.children
}

func (d *TVShowRow) SetChildren(children []*TVShowRow) {
	d.children = children
}

func (d *TVShowRow) CellDataForSort(_ int) string {
	return "" // Disable sorting for TV shows (would break dependencies between series and episodes)
}

func (d *TVShowRow) ColumnCell(_, col int, foreground, _ unison.Ink, _, _, _ bool) unison.Paneler {
	var text string
	switch col {
	case 0:
		text = d.M.Name
	case 1:
		text = d.M.Episode
	case 2:
		text = d.M.Season
	case 3:
		text = d.M.ProductionYear
	case 4:
		text = d.M.Runtime
	case 5:
		text = d.M.Actors
	case 6:
		text = d.M.Studios
	case 7:
		text = d.M.Genres
	case 8:
		text = d.M.Container
	case 9:
		text = d.M.Codecs
	case 10:
		text = d.M.Resolution
	case 11:
		text = d.M.Path
	default:
		text = ""
	}
	wrapper := unison.NewPanel()
	wrapper.SetLayout(&unison.FlexLayout{Columns: 1})
	addText(wrapper, text, foreground, unison.LabelFont)
	return wrapper
}

func (d *TVShowRow) IsOpen() bool {
	return d.open
}

func (d *TVShowRow) SetOpen(open bool) {
	d.open = open
}

func NewTVShowRow(id tid.TID, data TVShowData) *TVShowRow {
	row := &TVShowRow{
		table:     TVShowTable,
		id:        id,
		container: false,
		open:      false,
		parent:    nil,
		children:  nil,
		M: TVShowData{data.Name, data.Episode, data.Season, data.ProductionYear,
			data.Runtime, data.Actors, data.Studios, data.Genres, data.Container,
			data.Codecs, data.Resolution, data.Path, data.Overview,
			data.SeriesId, data.SeasonId, data.EpisodeId, data.Type_,
			data.SortIndex},
	}
	return row
}

// ---------------------------------------------------------------------------------------------------------------------
// Home videos model
// ---------------------------------------------------------------------------------------------------------------------

var _ unison.TableRowData[*HomeVideoRow] = &HomeVideoRow{}
var HomeVideoTable *unison.Table[*HomeVideoRow]
var HomeVideoTableDescription = TableDescription{
	NoOfColumns: 7, //displayed columns only
	Captions:    []string{"Title", "Folder", "Time", "Ext.", "Codec", "Resolution", "Path"},
	APIFields:   "Name,MediaSources,Path,Width,Height,Container,RunTimeTicks,ParentId,Type_", //no spaces here!
}

type HomeVideoData struct {
	Name       string
	Folder     string
	Runtime    string
	Container  string
	Codecs     string
	Resolution string
	Path       string
	FolderId   string
	ParentId   string
}

type HomeVideoRow struct {
	table        *unison.Table[*HomeVideoRow]
	parent       *HomeVideoRow
	children     []*HomeVideoRow
	container    bool
	open         bool
	doubleHeight bool
	id           tid.TID
	M            HomeVideoData
}

func (d *HomeVideoRow) CloneForTarget(target unison.Paneler, newParent *HomeVideoRow) *HomeVideoRow {
	table, ok := target.(*unison.Table[*HomeVideoRow])
	if !ok {
		fatal.IfErr(errs.New("invalid target"))
	}
	clone := *d
	clone.table = table
	clone.parent = newParent
	clone.id = tid.MustNewTID('a')
	return &clone
}

func (d *HomeVideoRow) ID() tid.TID {
	return d.id
}

func (d *HomeVideoRow) Parent() *HomeVideoRow {
	return d.parent
}

func (d *HomeVideoRow) SetParent(parent *HomeVideoRow) {
	d.parent = parent
}

func (d *HomeVideoRow) CanHaveChildren() bool {
	return d.container
}

func (d *HomeVideoRow) Children() []*HomeVideoRow {
	return d.children
}

func (d *HomeVideoRow) SetChildren(children []*HomeVideoRow) {
	d.children = children
}

func (d *HomeVideoRow) CellDataForSort(col int) string {
	switch col {
	case 0:
		return d.M.Name
	case 1:
		return d.M.Folder
	case 2:
		return d.M.Runtime
	case 3:
		return d.M.Container
	case 4:
		return d.M.Codecs
	case 5:
		return d.M.Resolution
	case 6:
		return d.M.Path
	default:
		return ""
	}
}

func (d *HomeVideoRow) ColumnCell(_, col int, foreground, _ unison.Ink, _, _, _ bool) unison.Paneler {
	var text string
	switch col {
	case 0:
		text = d.M.Name
	case 1:
		text = d.M.Folder
	case 2:
		text = d.M.Runtime
	case 3:
		text = d.M.Container
	case 4:
		text = d.M.Codecs
	case 5:
		text = d.M.Resolution
	case 6:
		text = d.M.Path
	default:
		text = ""
	}
	wrapper := unison.NewPanel()
	wrapper.SetLayout(&unison.FlexLayout{Columns: 1})
	addText(wrapper, text, foreground, unison.LabelFont)
	return wrapper
}

func (d *HomeVideoRow) IsOpen() bool {
	return d.open
}

func (d *HomeVideoRow) SetOpen(open bool) {
	d.open = open
}

func NewHomeVideoRow(id tid.TID, data HomeVideoData) *HomeVideoRow {
	row := &HomeVideoRow{
		table:     HomeVideoTable,
		id:        id,
		container: false,
		open:      false,
		parent:    nil,
		children:  nil,
		M: HomeVideoData{data.Name, data.Folder, data.Runtime, data.Container, data.Codecs,
			data.Resolution, data.Path, data.FolderId, data.ParentId},
	}
	return row
}

func addText(parent *unison.Panel, text string, ink unison.Ink, font unison.Font) {
	tx := unison.NewText(text, &unison.TextDecoration{Font: font})
	label := unison.NewLabel()
	label.Font = font
	label.LabelTheme.OnBackgroundInk = ink
	label.SetTitle(tx.String())
	parent.AddChild(label)
}
