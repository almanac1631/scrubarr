package inventory

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/almanac1631/scrubarr/internal/app/webserver"
	"github.com/almanac1631/scrubarr/pkg/domain"
)

type mediaId struct {
	MediaType domain.MediaType
	Id        int64
	FileId    int64
	Season    int
}

func (m mediaId) getMatchingLinkedMediaIndexes(linkedMediaFiles []LinkedMediaFile) []int {
	indexes := make([]int, 0)
	for i, linkedMediaFile := range linkedMediaFiles {
		if m.FileId != 0 {
			if linkedMediaFile.Id == m.FileId {
				indexes = append(indexes, i)
			}
			continue
		} else if m.Season != 0 {
			if linkedMediaFile.Season == m.Season {
				indexes = append(indexes, i)
			}
			continue
		}
		indexes = append(indexes, i)
	}
	return indexes
}

func (m mediaId) String() string {
	idStr := fmt.Sprintf("%s-%d", m.MediaType, m.Id)
	if m.FileId != 0 {
		return fmt.Sprintf("%s-%d", idStr, m.FileId)
	} else if m.Season != 0 {
		return fmt.Sprintf("%s-s-%d", idStr, m.Season)
	}
	return idStr
}

func parseMediaId(rawId string) (mediaId, error) {
	idSplit := strings.Split(rawId, "-")
	mediaType := idSplit[0]
	if mediaType != "movie" && mediaType != "series" {
		return mediaId{}, webserver.ErrMalformedMediaId
	}
	id, err := strconv.ParseInt(idSplit[1], 10, 64)
	if err != nil {
		return mediaId{}, webserver.ErrMalformedMediaId
	}
	var fileId int64
	var season int
	if len(idSplit) == 3 {
		// movie-10-8
		fileId, err = strconv.ParseInt(idSplit[2], 10, 64)
		if err != nil {
			return mediaId{}, webserver.ErrMalformedMediaId
		}
	} else if len(idSplit) == 4 {
		// series-1337-s-2
		if idSplit[2] != "s" {
			return mediaId{}, webserver.ErrMalformedMediaId
		}
		season, err = strconv.Atoi(idSplit[3])
		if err != nil {
			return mediaId{}, webserver.ErrMalformedMediaId
		}
	}
	return mediaId{
		MediaType: domain.MediaType(mediaType),
		Id:        id,
		FileId:    fileId,
		Season:    season,
	}, nil
}
