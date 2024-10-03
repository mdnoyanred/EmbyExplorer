// ---------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz. No rights reserved.
// Evaluation of Emby DTO & mapping fields to display structures
// Resolve dependencies for TV shows
// ---------------------------------------------------------------------------------------------------------------------

package api

import (
	"Emby_Explorer/models"
	"sort"
	"strconv"
)

const (
	maxActors    = 5
	maxDirectors = 2
	maxStudios   = 1
)

const placeHolder = "-"

func GetFields(collectiontype string) string {
	var m = ""
	switch collectiontype {
	case CollectionMovies:
		m = models.MovieTableDescription.APIFields
	case CollectionTVShows:
		m = models.TVShowTableDescription.APIFields
	case CollectionHomeVideos:
		m = models.HomeVideoTableDescription.APIFields
	default:
	}
	return m
}

func GetMovieDisplayData(dto []BaseItemDto) []models.MovieData {
	result := make([]models.MovieData, 0)
	var movie models.MovieData
	for _, d := range dto {
		movie.Name = d.Name
		movie.OriginalTitle = d.OriginalTitle
		movie.ProductionYear = strconv.Itoa(int(d.ProductionYear))
		movie.Studios = evalStudios(d.Studios)
		movie.Actors, movie.Directors = evalPeople(d.People)
		movie.Genres = evalGenres(d.Genres)
		movie.Container = d.Container
		movie.Resolution = evalResolution(d.Width, d.Height)
		movie.Codecs = evalCodecs(d.MediaSources)
		movie.Runtime = evalRuntime(d.RunTimeTicks)
		movie.Path = d.Path
		result = append(result, movie)
	}
	return result
}

func GetTVShowDisplayData(dto []BaseItemDto) []models.TVShowData {
	result := make([]models.TVShowData, 0)
	series := make([]models.TVShowData, 0)
	seasons := make([]models.TVShowData, 0)
	episodes := make([]models.TVShowData, 0)
	var item models.TVShowData
	for _, d := range dto {
		item = models.TVShowData{}
		switch d.Type_ {
		case seriesType:
			item.Name = d.Name
			item.Actors, _ = evalPeople(d.People)
			item.Genres = evalGenres(d.Genres)
			item.Studios = evalStudios(d.Studios)
			item.Path = d.Path
			item.SeriesID = d.Id
			item.Type_ = d.Type_
			series = append(series, item)
		case seasonType:
			item.Season = d.Name
			item.SeriesID = d.SeriesId
			item.SeasonID = d.Id
			item.SortIndex = d.IndexNumber
			item.Path = d.Path
			item.Type_ = d.Type_
			seasons = append(seasons, item)
		case episodeType:
			item.Episode = d.Name
			item.EpisodeID = d.Id
			item.Runtime = evalRuntime(d.RunTimeTicks)
			item.Container = d.Container
			item.Codecs = evalCodecs(d.MediaSources)
			item.Resolution = evalResolution(d.Width, d.Height)
			item.ProductionYear = strconv.Itoa(int(d.ProductionYear))
			item.Actors, _ = evalPeople(d.People)
			item.SortIndex = d.IndexNumber
			item.Path = d.Path
			item.SeriesID = d.SeriesId
			item.SeasonID = d.SeasonId
			item.Type_ = d.Type_
			episodes = append(episodes, item)
		default:
		}
	}
	// Sort series by Name
	sort.Slice(series, func(i, j int) bool {
		return series[i].Name < series[j].Name
	})
	// Sort seasons by series
	sort.Slice(seasons, func(i, j int) bool {
		return seasons[i].SeriesID < seasons[j].SeriesID
	})
	// Sort episodes by series
	sort.Slice(episodes, func(i, j int) bool {
		return episodes[i].SeriesID < episodes[j].SeriesID
	})
	for _, s := range series {
		result = append(result, s)
		seasonstmp := make([]models.TVShowData, 0)
		// Find seasons for series
		for _, season := range seasons {
			if season.SeriesID == s.SeriesID {
				seasonstmp = append(seasonstmp, season)
			}
		}
		// Sort seasons by IndexNumber
		sort.Slice(seasonstmp, func(i, j int) bool {
			return seasonstmp[i].SortIndex < seasonstmp[j].SortIndex
		})
		for _, n := range seasonstmp {
			// Find episodes for series and season
			episodesstmp := make([]models.TVShowData, 0)
			for _, episode := range episodes {
				if episode.SeriesID == n.SeriesID && episode.SeasonID == n.SeasonID {
					episodesstmp = append(episodesstmp, episode)
				}
			}
			// Sort episodes by IndexNumber
			sort.Slice(episodesstmp, func(i, j int) bool {
				return episodesstmp[i].SortIndex < episodesstmp[j].SortIndex
			})
			for _, e := range episodesstmp {
				e.Name = s.Name
				e.Season = n.Season
				e.Genres = s.Genres
				e.Studios = s.Studios
				if e.Actors == "" {
					e.Actors = s.Actors
				}
				result = append(result, e)
			}
		}
	}
	return result
}

func GetHomeVideoDisplayData(dto []BaseItemDto) []models.HomeVideoData {
	result := make([]models.HomeVideoData, 0)
	var video models.HomeVideoData
	for _, d := range dto {
		video.Name = d.Name
		video.Container = d.Container
		video.Resolution = evalResolution(d.Width, d.Height)
		video.Codecs = evalCodecs(d.MediaSources)
		video.Runtime = evalRuntime(d.RunTimeTicks)
		video.Path = d.Path
		result = append(result, video)
	}
	return result
}

func evalStudios(studios []NameLongIdPair) string {
	var s = ""
	for i, studio := range studios {
		i++
		if i > maxStudios {
			break
		}
		s = commaString(s, studio.Name)
	}
	return s
}

func evalPeople(people []BaseItemPerson) (string, string) {
	var actors = ""
	var directors = ""
	var countactors = 0
	var countdirectors = 0
	for _, p := range people {
		if *p.Type_ == ACTOR_PersonType {
			countactors++
			if countactors <= maxActors {
				actors = commaString(actors, p.Name)
			}
		}
		if *p.Type_ == DIRECTOR_PersonType {
			countdirectors++
			if countdirectors <= maxDirectors {
				directors = commaString(directors, p.Name)
			}
		}
		if countactors > maxActors && countdirectors > maxDirectors {
			break
		}
	}
	return actors, directors
}

func evalGenres(genres []string) string {
	var s = ""
	for _, genre := range genres {
		s = commaString(s, genre)
	}
	return s
}

func evalRuntime(ticks int64) string {
	var s = ""
	if ticks > 0 {
		r := ticks / 10000000
		hours := r / 3600
		minutes := (r % 3600) / 60
		if hours > 0 {
			s = strconv.Itoa(int(hours)) + "h"
		}
		if minutes > 0 {
			s = s + strconv.Itoa(int(minutes)) + "m"
		}
	}
	return s
}

func evalCodecs(media []MediaSourceInfo) string {
	var codecs = ""
	var codecVideo = ""
	var codecAudio = ""
	for _, m := range media {
		for _, s := range m.MediaStreams {
			if *s.Type_ == VIDEO_MediaStreamType {
				codecVideo = s.Codec
			}
			if *s.Type_ == AUDIO_MediaStreamType {
				codecAudio = s.Codec
			}
		}
		if codecVideo == "" {
			codecVideo = placeHolder
		}
		if codecAudio == "" {
			codecAudio = placeHolder
		}
		codecs = codecVideo + ", " + codecAudio
		break
	}
	return codecs
}

func evalResolution(w int32, h int32) string {
	var r = ""
	if w > 0 && h > 0 {
		r = strconv.Itoa(int(w)) + "x" + strconv.Itoa(int(h))
	}
	return r
}

func commaString(source string, append string) string {
	s := source
	if s != "" {
		s = s + ", " + append
	} else {
		s = append
	}
	return s
}
