package matching_manager

/*
func (m *Linker) GetMatchedMedia(page int, sortInfo domain.SortInfo) ([]domain.MatchedMedia, bool, error) {
	if m.linkedMediaCache == nil {
		if err := m.refreshCache(); err != nil {
			return nil, false, err
		}
	}
	hasNext := false
	matchedMedias := make([]domain.MatchedMedia, len(m.linkedMediaCache))
	copy(matchedMedias, m.linkedMediaCache)

	torrentStatusScores := map[string]int{}
	for _, entry := range m.linkedMediaCache {
		totalScore := 0
		for _, part := range entry.Parts {
			totalScore += part.TorrentInformation.GetScore()
		}
		torrentStatusScores[entry.Url] = totalScore / len(entry.Parts)
	}

	slices.SortFunc(matchedMedias, func(a, b domain.MatchedMedia) int {
		var result int
		switch sortInfo.Key {
		case domain.SortKeyName:
			result = strings.Compare(strings.ToLower(a.Title), strings.ToLower(b.Title))
			break
		case domain.SortKeySize:
			result = cmp.Compare(a.Size, b.Size)
			break
		case domain.SortKeyAdded:
			result = cmp.Compare(a.Added.Unix(), b.Added.Unix())
			break
		case domain.SortKeyTorrentStatus:
			result = cmp.Compare(torrentStatusScores[a.Url], torrentStatusScores[b.Url])
			break
		default:
			slog.Error("Received unknown sort key.", "sortKey", sortInfo.Key)
			result = 0 // mark as incomparable
		}
		if sortInfo.Order == domain.SortOrderDesc {
			result = -result
		}
		return result
	})
	if pageSize*page < len(matchedMedias) {
		hasNext = true
		matchedMedias = matchedMedias[pageSize*(page-1) : pageSize*page]
	} else {
		matchedMedias = matchedMedias[pageSize*(page-1):]
	}
	return matchedMedias, hasNext, nil
}

func (m *Linker) GetMatchedMediaBySeriesId(seriesId int64) (media domain.MatchedMedia, err error) {
	return m.getSingleMatchedMediaEntry(domain.MediaTypeSeries, seriesId)
}

func (m *Linker) getSingleMatchedMediaEntry(mediaType domain.MediaType, id int64) (media domain.MatchedMedia, err error) {
	matchedMediaList, err := m.getFilteredMatchedMediaFunc(func(media domain.MatchedMedia) bool {
		return media.Type == mediaType && media.Id == id
	})
	if err != nil {
		return domain.MatchedMedia{}, err
	}
	if len(matchedMediaList) == 0 {
		return domain.MatchedMedia{}, domain.ErrMediaNotFound
	} else if len(matchedMediaList) > 1 {
		return domain.MatchedMedia{}, fmt.Errorf("multiple matched media found with type %s and id %d", mediaType, id)
	}
	return matchedMediaList[0], nil
}

func (m *Linker) getFilteredMatchedMediaFunc(filterFunc func(media domain.MatchedMedia) bool) (media []domain.MatchedMedia, err error) {
	if m.linkedMediaCache == nil {
		if err := m.refreshCache(); err != nil {
			return nil, err
		}
	}
	filteredMediaList := make([]domain.MatchedMedia, 0)
	for _, mediaEntry := range m.linkedMediaCache {
		if filterFunc(mediaEntry) {
			filteredMediaList = append(filteredMediaList, mediaEntry)
		}
	}
	return filteredMediaList, nil
}
*/
