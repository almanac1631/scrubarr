package inventory

import (
	"testing"
	"time"

	"github.com/almanac1631/scrubarr/internal/app/webserver"
	"github.com/almanac1631/scrubarr/pkg/domain"
	"github.com/almanac1631/scrubarr/pkg/util"
	"github.com/stretchr/testify/require"
)

var testTorrentInformationMissing = webserver.TorrentInformation{
	LinkStatus: webserver.TorrentLinkMissing,
	Tracker:    domain.Tracker{},
	Ratio:      -1.0,
	Age:        time.Duration(-1),
}

func Test_generateRawFileBasedMediaRow(t *testing.T) {
	torrentEntry1 := &domain.TorrentEntry{
		Ratio: 2.6,
		Added: util.MustParseDate("2022-08-12 00:00:00"),
	}
	torrentEntry2 := &domain.TorrentEntry{
		Ratio: 0.9,
		Added: util.MustParseDate("2023-08-11 00:00:00"),
	}
	now = func() time.Time {
		return util.MustParseDate("2023-08-12 00:00:00")
	}
	type args struct {
		media enrichedLinkedMedia
	}
	tests := []struct {
		name string
		args args
		want webserver.MediaRow
	}{
		{
			name: "media entry with single complete file",
			args: args{
				media: enrichedLinkedMedia{
					linkedMedia: LinkedMedia{
						MediaMetadata: domain.MediaMetadata{
							Id:    1337,
							Type:  domain.MediaTypeMovie,
							Title: "Some movie title",
							Url:   "http://example.com/movie.mp4",
							Added: util.MustParseDate("2021-08-12 00:00:00"),
						},
						Files: []LinkedMediaFile{
							{
								MediaFile: domain.MediaFile{
									Id:               10,
									OriginalFilePath: "some/path/anywhere/Some-movie-title-1080p.mp4",
									Size:             int64(8232),
								},
								TorrentEntry: &domain.TorrentEntry{
									Ratio: 1.4,
									Added: util.MustParseDate("2022-08-12 00:00:00"),
								},
							},
						},
					},
					size:  3827948,
					added: util.MustParseDate("2021-08-12 00:00:00"),
				},
			},
			want: webserver.MediaRow{
				Id:    "movie-1337",
				Type:  domain.MediaTypeMovie,
				Title: "Some movie title",
				Url:   "http://example.com/movie.mp4",
				Size:  3827948,
				Added: util.MustParseDate("2021-08-12 00:00:00"),
				TorrentInformation: webserver.TorrentInformation{
					LinkStatus: webserver.TorrentLinkPresent,
					Ratio:      1.4,
					Age:        time.Hour * 24 * 365,
				},
				ChildMediaRows: []webserver.MediaRow{{
					Id:    "movie-1337-10",
					Title: "Some-movie-title-1080p.mp4",
					Size:  int64(8232),
					Added: util.MustParseDate("2022-08-12 00:00:00"),
					TorrentInformation: webserver.TorrentInformation{
						LinkStatus: webserver.TorrentLinkPresent,
						Ratio:      1.4,
						Age:        time.Hour * 24 * 365,
					},
					ChildMediaRows: []webserver.MediaRow{},
				}},
			},
		},
		{
			name: "media entry with multiple complete files",
			args: args{
				enrichedLinkedMedia{
					size:  9000,
					added: util.MustParseDate("2020-08-12 00:00:00"),
					linkedMedia: LinkedMedia{
						MediaMetadata: domain.MediaMetadata{
							Id:    1337,
							Type:  domain.MediaTypeSeries,
							Title: "Breaking Bad",
							Url:   "https://some-series.com/series/breaking-bad",
						},
						Files: []LinkedMediaFile{
							{
								MediaFile: domain.MediaFile{
									Id:               1337_1,
									OriginalFilePath: "Breaking Bad S1 E1.mp4",
									Season:           1,
									Size:             1000,
								},
								TorrentEntry: torrentEntry1,
							},
							{
								MediaFile: domain.MediaFile{
									Id:               1337_2,
									OriginalFilePath: "Breaking Bad S1 E2.mp4",
									Season:           1,
									Size:             1000,
								},
								TorrentEntry: torrentEntry1,
							},
							{
								MediaFile: domain.MediaFile{
									Id:               1337_3,
									OriginalFilePath: "Breaking Bad S2 E1.mp4",
									Season:           2,
									Size:             1500,
								},
								TorrentEntry: torrentEntry2,
							},
						},
					},
				},
			},
			want: webserver.MediaRow{
				Id:    "series-1337",
				Type:  domain.MediaTypeSeries,
				Title: "Breaking Bad",
				Url:   "https://some-series.com/series/breaking-bad",
				Size:  9000,
				Added: util.MustParseDate("2020-08-12 00:00:00"),
				TorrentInformation: webserver.TorrentInformation{
					LinkStatus: webserver.TorrentLinkPresent,
					Ratio:      -1.0,
					Age:        time.Duration(-1),
				},
				ChildMediaRows: []webserver.MediaRow{
					{
						Id:    "series-1337-13371",
						Title: "Breaking Bad S1 E1.mp4",
						Size:  1000,
						Added: util.MustParseDate("2022-08-12 00:00:00"),
						TorrentInformation: webserver.TorrentInformation{
							LinkStatus: webserver.TorrentLinkPresent,
							Ratio:      torrentEntry1.Ratio,
							Age:        now().Sub(util.MustParseDate("2022-08-12 00:00:00")),
						},
						ChildMediaRows: []webserver.MediaRow{},
					},
					{
						Id:    "series-1337-13372",
						Title: "Breaking Bad S1 E2.mp4",
						Size:  1000,
						Added: util.MustParseDate("2022-08-12 00:00:00"),
						TorrentInformation: webserver.TorrentInformation{
							LinkStatus: webserver.TorrentLinkPresent,
							Ratio:      torrentEntry1.Ratio,
							Age:        now().Sub(util.MustParseDate("2022-08-12 00:00:00")),
						},
						ChildMediaRows: []webserver.MediaRow{},
					},
					{
						Id:    "series-1337-13373",
						Title: "Breaking Bad S2 E1.mp4",
						Size:  1500,
						Added: util.MustParseDate("2023-08-11 00:00:00"),
						TorrentInformation: webserver.TorrentInformation{
							LinkStatus: webserver.TorrentLinkPresent,
							Ratio:      torrentEntry2.Ratio,
							Age:        now().Sub(util.MustParseDate("2023-08-11 00:00:00")),
						},
						ChildMediaRows: []webserver.MediaRow{},
					},
				},
			},
		},
		{
			name: "media entry with single incomplete file",
			args: args{
				media: enrichedLinkedMedia{
					linkedMedia: LinkedMedia{
						MediaMetadata: domain.MediaMetadata{
							Id:    1337,
							Type:  domain.MediaTypeMovie,
							Title: "Some movie title",
							Url:   "http://example.com/movie.mp4",
							Added: util.MustParseDate("2021-08-12 00:00:00"),
						},
						Files: []LinkedMediaFile{
							{
								MediaFile: domain.MediaFile{
									Id:               10,
									OriginalFilePath: "some/path/anywhere/Some-movie-title-1080p.mp4",
									Size:             int64(8232),
								},
								TorrentEntry: nil,
							},
						},
					},
					size:  3827948,
					added: util.MustParseDate("2021-08-12 00:00:00"),
				},
			},
			want: webserver.MediaRow{
				Id:                 "movie-1337",
				Type:               domain.MediaTypeMovie,
				Title:              "Some movie title",
				Url:                "http://example.com/movie.mp4",
				Size:               3827948,
				Added:              util.MustParseDate("2021-08-12 00:00:00"),
				TorrentInformation: testTorrentInformationMissing,
				ChildMediaRows: []webserver.MediaRow{{
					Id:                 "movie-1337-10",
					Title:              "Some-movie-title-1080p.mp4",
					Size:               int64(8232),
					TorrentInformation: testTorrentInformationMissing,
					ChildMediaRows:     []webserver.MediaRow{},
				}},
			},
		},
		{
			name: "media entry with multiple incomplete files",
			args: args{
				enrichedLinkedMedia{
					size:  9000,
					added: util.MustParseDate("2020-08-12 00:00:00"),
					linkedMedia: LinkedMedia{
						MediaMetadata: domain.MediaMetadata{
							Id:    1337,
							Type:  domain.MediaTypeSeries,
							Title: "Breaking Bad",
							Url:   "https://some-series.com/series/breaking-bad",
						},
						Files: []LinkedMediaFile{
							{
								MediaFile: domain.MediaFile{
									Id:               1337_1,
									OriginalFilePath: "Breaking Bad S1 E1.mp4",
									Season:           1,
									Size:             1000,
								},
								TorrentEntry: torrentEntry1,
							},
							{
								MediaFile: domain.MediaFile{
									Id:               1337_2,
									OriginalFilePath: "Breaking Bad S1 E2.mp4",
									Season:           1,
									Size:             1000,
								},
								TorrentEntry: nil,
							},
							{
								MediaFile: domain.MediaFile{
									Id:               1337_3,
									OriginalFilePath: "Breaking Bad S2 E1.mp4",
									Season:           2,
									Size:             1500,
								},
								TorrentEntry: torrentEntry2,
							},
						},
					},
				},
			},
			want: webserver.MediaRow{
				Id:    "series-1337",
				Type:  domain.MediaTypeSeries,
				Title: "Breaking Bad",
				Url:   "https://some-series.com/series/breaking-bad",
				Size:  9000,
				Added: util.MustParseDate("2020-08-12 00:00:00"),
				TorrentInformation: webserver.TorrentInformation{
					LinkStatus: webserver.TorrentLinkIncomplete,
					Ratio:      -1.0,
					Age:        time.Duration(-1),
				},
				ChildMediaRows: []webserver.MediaRow{
					{
						Id:    "series-1337-13371",
						Title: "Breaking Bad S1 E1.mp4",
						Size:  1000,
						Added: util.MustParseDate("2022-08-12 00:00:00"),
						TorrentInformation: webserver.TorrentInformation{
							LinkStatus: webserver.TorrentLinkPresent,
							Ratio:      torrentEntry1.Ratio,
							Age:        now().Sub(util.MustParseDate("2022-08-12 00:00:00")),
						},
						ChildMediaRows: []webserver.MediaRow{},
					},
					{
						Id:                 "series-1337-13372",
						Title:              "Breaking Bad S1 E2.mp4",
						Size:               1000,
						Added:              time.Time{},
						TorrentInformation: testTorrentInformationMissing,
						ChildMediaRows:     []webserver.MediaRow{},
					},
					{
						Id:    "series-1337-13373",
						Title: "Breaking Bad S2 E1.mp4",
						Size:  1500,
						Added: util.MustParseDate("2023-08-11 00:00:00"),
						TorrentInformation: webserver.TorrentInformation{
							LinkStatus: webserver.TorrentLinkPresent,
							Ratio:      torrentEntry2.Ratio,
							Age:        now().Sub(util.MustParseDate("2023-08-11 00:00:00")),
						},
						ChildMediaRows: []webserver.MediaRow{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateRawMediaRowFromLinkedMedia(tt.args.media)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_applyEvaluationReport(t *testing.T) {
	tracker := &domain.Tracker{
		Name:     "mock-tracker",
		MinRatio: 0.0,
		MinAge:   time.Duration(0),
	}
	tracker2 := &domain.Tracker{
		Name:     "mock-tracker-2",
		MinRatio: 1.0,
		MinAge:   time.Hour * 24 * 30,
	}
	torrentEntry1 := &domain.TorrentEntry{
		Id:    "some-torrent-entry-1",
		Added: util.MustParseDate("2020-08-12 00:00:00"),
	}
	torrentEntry2 := &domain.TorrentEntry{
		Id:    "some-torrent-entry-2",
		Added: util.MustParseDate("2020-08-13 00:00:00"),
	}
	torrentInfoPresent1 := webserver.TorrentInformation{
		LinkStatus: webserver.TorrentLinkPresent,
		Ratio:      1.0,
		Age:        time.Hour * 24 * 30,
	}
	torrentInfoPresentTracker1 := torrentInfoPresent1
	torrentInfoPresentTracker1.Tracker = *tracker
	torrentInfoPresent2 := webserver.TorrentInformation{
		LinkStatus: webserver.TorrentLinkPresent,
		Tracker:    *tracker,
		Ratio:      2.5,
		Age:        time.Hour * 24 * 15,
	}
	torrentInfoPresentTracker2 := torrentInfoPresent2
	torrentInfoPresentTracker2.Tracker = *tracker
	type args struct {
		media enrichedLinkedMedia
		row   webserver.MediaRow
	}
	tests := []struct {
		name string
		args args
		want webserver.MediaRow
	}{
		{
			name: "single file movie",
			args: args{
				media: enrichedLinkedMedia{
					linkedMedia: LinkedMedia{
						MediaMetadata: domain.MediaMetadata{
							Id: 10,
						},
						Files: []LinkedMediaFile{
							{
								MediaFile: domain.MediaFile{
									Id: 1337,
								},
							},
						},
					},
					evaluationReport: EvaluationReport{
						Result:  EvaluationReportPart{Decision: domain.DecisionSafeToDelete},
						Seasons: map[int]EvaluationReportPart{},
						Files: map[int64]EvaluationReportPart{
							1337: {Decision: domain.DecisionSafeToDelete},
						},
					},
				},
				row: webserver.MediaRow{
					Id: "movie-10",
					TorrentInformation: webserver.TorrentInformation{
						LinkStatus: webserver.TorrentLinkPresent,
					},
					ChildMediaRows: []webserver.MediaRow{
						{
							Id: "movie-10-1337",
							TorrentInformation: webserver.TorrentInformation{
								LinkStatus: webserver.TorrentLinkPresent,
							},
						},
					},
				},
			},
			want: webserver.MediaRow{
				Id:       "movie-10",
				Decision: domain.DecisionSafeToDelete,
				TorrentInformation: webserver.TorrentInformation{
					LinkStatus: webserver.TorrentLinkPresent,
				},
				ChildMediaRows: []webserver.MediaRow{{
					Id: "movie-10-1337",
					TorrentInformation: webserver.TorrentInformation{
						LinkStatus: webserver.TorrentLinkPresent,
					},
					Decision: domain.DecisionSafeToDelete,
				}},
			},
		},
		{
			name: "multi file series",
			args: args{
				media: enrichedLinkedMedia{
					linkedMedia: LinkedMedia{
						MediaMetadata: domain.MediaMetadata{
							Id: 10,
						},
						Files: []LinkedMediaFile{
							{
								MediaFile: domain.MediaFile{
									Id:     1337_1,
									Season: 1,
								},
							},
							{
								MediaFile: domain.MediaFile{
									Id:     1337_2,
									Season: 1,
								},
							},
							{
								MediaFile: domain.MediaFile{
									Id:     1337_3,
									Season: 2,
								},
							},
							{
								MediaFile: domain.MediaFile{
									Id: 1337_4,
								},
							},
						},
					},
					evaluationReport: EvaluationReport{
						Result: EvaluationReportPart{Decision: domain.DecisionSafeToDelete},
						Seasons: map[int]EvaluationReportPart{
							1: {Decision: domain.DecisionSafeToDelete},
							2: {Decision: domain.DecisionSafeToDelete},
						},
						Files: map[int64]EvaluationReportPart{
							1337_1: {Decision: domain.DecisionSafeToDelete},
							1337_2: {Decision: domain.DecisionSafeToDelete},
							1337_3: {Decision: domain.DecisionSafeToDelete},
							1337_4: {Decision: domain.DecisionSafeToDelete},
						},
					},
				},
				row: webserver.MediaRow{
					Id:                 "series-10",
					TorrentInformation: testTorrentInformationMissing,
					ChildMediaRows: []webserver.MediaRow{
						{
							Id:                 "series-10-13371",
							TorrentInformation: testTorrentInformationMissing,
						},
						{
							Id:                 "series-10-13372",
							TorrentInformation: testTorrentInformationMissing,
						},
						{
							Id:                 "series-10-13373",
							TorrentInformation: testTorrentInformationMissing,
						},
						{
							Id:                 "series-10-13374",
							TorrentInformation: testTorrentInformationMissing,
						},
					},
				},
			},
			want: webserver.MediaRow{
				Id:                 "series-10",
				Decision:           domain.DecisionSafeToDelete,
				TorrentInformation: testTorrentInformationMissing,
				ChildMediaRows: []webserver.MediaRow{
					{
						Id:                 "series-10-s-1",
						Title:              "Season 1",
						TorrentInformation: testTorrentInformationMissing,
						Decision:           domain.DecisionSafeToDelete,
						ChildMediaRows: []webserver.MediaRow{
							{
								Id:                 "series-10-13371",
								TorrentInformation: testTorrentInformationMissing,
								Decision:           domain.DecisionSafeToDelete,
							},
							{
								Id:                 "series-10-13372",
								TorrentInformation: testTorrentInformationMissing,
								Decision:           domain.DecisionSafeToDelete,
							},
						},
					},
					{
						Id:                 "series-10-s-2",
						Title:              "Season 2",
						TorrentInformation: testTorrentInformationMissing,
						Decision:           domain.DecisionSafeToDelete,
						ChildMediaRows: []webserver.MediaRow{
							{
								Id:                 "series-10-13373",
								TorrentInformation: testTorrentInformationMissing,
								Decision:           domain.DecisionSafeToDelete,
							},
						},
					},
					{
						Id:                 "series-10-13374",
						TorrentInformation: testTorrentInformationMissing,
						Decision:           domain.DecisionSafeToDelete,
					},
				},
			},
		},
		{
			name: "calc season size attribute",
			args: args{
				media: enrichedLinkedMedia{
					linkedMedia: LinkedMedia{
						MediaMetadata: domain.MediaMetadata{
							Id: 10,
						},
						Files: []LinkedMediaFile{
							{
								MediaFile: domain.MediaFile{
									Id:     1337_1,
									Season: 1,
									Size:   1000,
								},
							},
							{
								MediaFile: domain.MediaFile{
									Id:     1337_2,
									Season: 1,
									Size:   1500,
								},
							},
						},
					},
					evaluationReport: EvaluationReport{
						Result: EvaluationReportPart{
							Decision: domain.DecisionSafeToDelete,
						},
						Seasons: map[int]EvaluationReportPart{
							1: {
								Decision: domain.DecisionSafeToDelete,
							},
						},
						Files: map[int64]EvaluationReportPart{
							1337_1: {
								Decision: domain.DecisionSafeToDelete,
							},
							1337_2: {
								Decision: domain.DecisionSafeToDelete,
							},
						},
					},
				},
				row: webserver.MediaRow{
					Id:                 "series-10",
					TorrentInformation: testTorrentInformationMissing,
					ChildMediaRows: []webserver.MediaRow{
						{
							Id:                 "series-10-13371",
							TorrentInformation: testTorrentInformationMissing,
							Size:               1000,
						},
						{
							Id:                 "series-10-13372",
							TorrentInformation: testTorrentInformationMissing,
							Size:               1500,
						},
					},
				},
			},
			want: webserver.MediaRow{
				Id:                 "series-10",
				Decision:           domain.DecisionSafeToDelete,
				TorrentInformation: testTorrentInformationMissing,
				ChildMediaRows: []webserver.MediaRow{
					{
						Id:                 "series-10-s-1",
						Title:              "Season 1",
						TorrentInformation: testTorrentInformationMissing,
						Decision:           domain.DecisionSafeToDelete,
						Size:               2500,
						ChildMediaRows: []webserver.MediaRow{
							{
								Id:                 "series-10-13371",
								TorrentInformation: testTorrentInformationMissing,
								Size:               1000,
								Decision:           domain.DecisionSafeToDelete,
							},
							{
								Id:                 "series-10-13372",
								TorrentInformation: testTorrentInformationMissing,
								Size:               1500,
								Decision:           domain.DecisionSafeToDelete,
							},
						},
					},
				},
			},
		},
		{
			name: "calc season torrent attributes - same torrent",
			args: args{
				media: enrichedLinkedMedia{
					linkedMedia: LinkedMedia{
						MediaMetadata: domain.MediaMetadata{
							Id: 10,
						},
						Files: []LinkedMediaFile{
							{
								MediaFile: domain.MediaFile{
									Id:     1337_1,
									Season: 1,
								},
								TorrentEntry: torrentEntry1,
							},
							{
								MediaFile: domain.MediaFile{
									Id:     1337_2,
									Season: 1,
								},
								TorrentEntry: torrentEntry1,
							},
						},
					},
					evaluationReport: EvaluationReport{
						Result: EvaluationReportPart{Decision: domain.DecisionSafeToDelete},
						Seasons: map[int]EvaluationReportPart{
							1: {Decision: domain.DecisionSafeToDelete},
						},
						Files: map[int64]EvaluationReportPart{
							1337_1: {
								Decision: domain.DecisionSafeToDelete,
								Tracker:  tracker,
							},
							1337_2: {
								Decision: domain.DecisionSafeToDelete,
								Tracker:  tracker,
							},
						},
					},
				},
				row: webserver.MediaRow{
					Id:                 "series-10",
					TorrentInformation: torrentInfoPresent1,
					ChildMediaRows: []webserver.MediaRow{
						{
							Id:                 "series-10-13371",
							TorrentInformation: torrentInfoPresent1,
							Added:              util.MustParseDate("2020-08-12 00:00:00"),
						},
						{
							Id:                 "series-10-13372",
							TorrentInformation: torrentInfoPresent1,
							Added:              util.MustParseDate("2020-08-12 00:00:00"),
						},
					},
				},
			},
			want: webserver.MediaRow{
				Id:                 "series-10",
				Decision:           domain.DecisionSafeToDelete,
				TorrentInformation: torrentInfoPresentTracker1,
				ChildMediaRows: []webserver.MediaRow{
					{
						Id:                 "series-10-s-1",
						Title:              "Season 1",
						TorrentInformation: torrentInfoPresentTracker1,
						Added:              util.MustParseDate("2020-08-12 00:00:00"),
						Decision:           domain.DecisionSafeToDelete,
						ChildMediaRows: []webserver.MediaRow{
							{
								Id:                 "series-10-13371",
								TorrentInformation: torrentInfoPresentTracker1,
								Added:              util.MustParseDate("2020-08-12 00:00:00"),
								Decision:           domain.DecisionSafeToDelete,
							},
							{
								Id:                 "series-10-13372",
								TorrentInformation: torrentInfoPresentTracker1,
								Added:              util.MustParseDate("2020-08-12 00:00:00"),
								Decision:           domain.DecisionSafeToDelete,
							},
						},
					},
				},
			},
		},
		{
			name: "calc season torrent attributes - different torrents",
			args: args{
				media: enrichedLinkedMedia{
					linkedMedia: LinkedMedia{
						MediaMetadata: domain.MediaMetadata{
							Id: 10,
						},
						Files: []LinkedMediaFile{
							{
								MediaFile: domain.MediaFile{
									Id:     1337_1,
									Season: 1,
								},
								TorrentEntry: torrentEntry1,
							},
							{
								MediaFile: domain.MediaFile{
									Id:     1337_2,
									Season: 1,
								},
								TorrentEntry: torrentEntry2,
							},
						},
					},
					evaluationReport: EvaluationReport{
						Result: EvaluationReportPart{Decision: domain.DecisionPending},
						Seasons: map[int]EvaluationReportPart{
							1: {Decision: domain.DecisionPending},
						},
						Files: map[int64]EvaluationReportPart{
							1337_1: {
								Decision: domain.DecisionSafeToDelete,
								Tracker:  tracker,
							},
							1337_2: {
								Decision: domain.DecisionPending,
								Tracker:  tracker,
							},
						},
					},
				},
				row: webserver.MediaRow{
					Id: "series-10",
					TorrentInformation: webserver.TorrentInformation{
						LinkStatus: webserver.TorrentLinkPresent,
						Ratio:      -1.0,
						Age:        time.Duration(-1),
					},
					ChildMediaRows: []webserver.MediaRow{
						{
							Id:                 "series-10-13371",
							TorrentInformation: torrentInfoPresent1,
							Added:              util.MustParseDate("2020-08-12 00:00:00"),
						},
						{
							Id:                 "series-10-13372",
							TorrentInformation: torrentInfoPresent2,
							Added:              util.MustParseDate("2020-08-13 00:00:00"),
						},
					},
				},
			},
			want: webserver.MediaRow{
				Id:       "series-10",
				Decision: domain.DecisionPending,
				TorrentInformation: webserver.TorrentInformation{
					LinkStatus: webserver.TorrentLinkPresent,
					Tracker:    *tracker,
					Ratio:      -1.0,
					Age:        time.Duration(-1),
				},
				ChildMediaRows: []webserver.MediaRow{
					{
						Id:    "series-10-s-1",
						Title: "Season 1",
						TorrentInformation: webserver.TorrentInformation{
							LinkStatus: webserver.TorrentLinkPresent,
							Tracker:    *tracker,
							Ratio:      -1.0,
							Age:        time.Duration(-1),
						},
						Decision: domain.DecisionPending,
						ChildMediaRows: []webserver.MediaRow{
							{
								Id:                 "series-10-13371",
								TorrentInformation: torrentInfoPresentTracker1,
								Added:              util.MustParseDate("2020-08-12 00:00:00"),
								Decision:           domain.DecisionSafeToDelete,
							},
							{
								Id:                 "series-10-13372",
								TorrentInformation: torrentInfoPresentTracker2,
								Added:              util.MustParseDate("2020-08-13 00:00:00"),
								Decision:           domain.DecisionPending,
							},
						},
					},
				},
			},
		},
		{
			name: "calc season torrent attributes - different trackers",
			args: args{
				media: enrichedLinkedMedia{
					linkedMedia: LinkedMedia{
						MediaMetadata: domain.MediaMetadata{
							Id: 10,
						},
						Files: []LinkedMediaFile{
							{
								MediaFile: domain.MediaFile{
									Id:     1337_1,
									Season: 1,
								},
								TorrentEntry: torrentEntry1,
							},
							{
								MediaFile: domain.MediaFile{
									Id:     1337_2,
									Season: 1,
								},
								TorrentEntry: torrentEntry2,
							},
						},
					},
					evaluationReport: EvaluationReport{
						Result: EvaluationReportPart{Decision: domain.DecisionPending},
						Seasons: map[int]EvaluationReportPart{
							1: {Decision: domain.DecisionPending},
						},
						Files: map[int64]EvaluationReportPart{
							1337_1: {
								Decision: domain.DecisionSafeToDelete,
								Tracker:  tracker,
							},
							1337_2: {
								Decision: domain.DecisionPending,
								Tracker:  tracker2,
							},
						},
					},
				},
				row: webserver.MediaRow{
					Id: "series-10",
					TorrentInformation: webserver.TorrentInformation{
						LinkStatus: webserver.TorrentLinkPresent,
						Ratio:      -1.0,
						Age:        time.Duration(-1),
					},
					ChildMediaRows: []webserver.MediaRow{
						{
							Id:                 "series-10-13371",
							TorrentInformation: torrentInfoPresent1,
							Added:              util.MustParseDate("2020-08-12 00:00:00"),
						},
						{
							Id:                 "series-10-13372",
							TorrentInformation: torrentInfoPresent2,
							Added:              util.MustParseDate("2020-08-13 00:00:00"),
						},
					},
				},
			},
			want: webserver.MediaRow{
				Id:       "series-10",
				Decision: domain.DecisionPending,
				TorrentInformation: webserver.TorrentInformation{
					LinkStatus: webserver.TorrentLinkPresent,
					Ratio:      -1.0,
					Age:        time.Duration(-1),
				},
				ChildMediaRows: []webserver.MediaRow{
					{
						Id:    "series-10-s-1",
						Title: "Season 1",
						TorrentInformation: webserver.TorrentInformation{
							LinkStatus: webserver.TorrentLinkPresent,
							Ratio:      -1.0,
							Age:        time.Duration(-1),
						},
						Decision: domain.DecisionPending,
						ChildMediaRows: []webserver.MediaRow{
							{
								Id:                 "series-10-13371",
								TorrentInformation: torrentInfoPresentTracker1,
								Added:              util.MustParseDate("2020-08-12 00:00:00"),
								Decision:           domain.DecisionSafeToDelete,
							},
							{
								Id: "series-10-13372",
								TorrentInformation: webserver.TorrentInformation{
									LinkStatus: webserver.TorrentLinkPresent,
									Tracker:    *tracker2,
									Ratio:      2.5,
									Age:        time.Hour * 24 * 15,
								},
								Added:    util.MustParseDate("2020-08-13 00:00:00"),
								Decision: domain.DecisionPending,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := applyEvaluationReport(tt.args.media, tt.args.row)
			require.Equal(t, tt.want, got)
		})
	}
}
