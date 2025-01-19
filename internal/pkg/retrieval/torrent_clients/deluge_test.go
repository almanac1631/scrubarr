package torrent_clients

import (
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval"
	delugeclient "github.com/gdm85/go-libdeluge"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestDelugeEntryRetriever_parseDelugeTorrentStatus(t *testing.T) {
	type fields struct {
		allowedFileEndings []string
	}
	type args struct {
		torrentStatus *delugeclient.TorrentStatus
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []retrieval.Entry
	}{
		{
			"can parse standard torrent status",
			fields{[]string{".mkv"}},
			args{&delugeclient.TorrentStatus{
				ActiveTime:       0,
				CompletedTime:    1725444490,
				Ratio:            1.74,
				SavePath:         "/home/someguy/downloads/deluge",
				DownloadLocation: "/home/someguy/downloads/deluge",
				Name:             "Family.Guy.S22E14.1080p.WEB.H264",
				SeedingTime:      751680, // 8.7 days
				TrackerHost:      "sometrakkr.co.uk",
				Files: []delugeclient.File{
					{
						Size: 1800000000,
						Path: "Somesubfolder/series/Episode1.mkv",
					},
				},
			}},
			[]retrieval.Entry{
				{
					Name: "Episode1.mkv",
					AdditionalData: TorrentClientEntry{
						TorrentClientName: "deluge",
						TorrentName:       "Family.Guy.S22E14.1080p.WEB.H264",
						DownloadFilePath:  "/home/someguy/downloads/deluge/Somesubfolder/series/Episode1.mkv",
						DownloadedAt:      time.Date(2024, time.September, 4, 10, 8, 10, 0, time.UTC),
						Ratio:             1.74,
						FileSizeBytes:     1800000000,
						TrackerHost:       "sometrakkr.co.uk",
					},
				},
			},
		},
		{
			"can filter multiple files in a torrent status",
			fields{[]string{".mkv"}},
			args{&delugeclient.TorrentStatus{
				ActiveTime:       0,
				CompletedTime:    1725444490,
				Ratio:            1.74,
				SavePath:         "/home/someguy/downloads/deluge",
				DownloadLocation: "/home/someguy/downloads/deluge",
				Name:             "Family.Guy.S22E14.1080p.WEB.H264",
				SeedingTime:      751680, // 8.7 days
				TrackerHost:      "sometrakkr.co.uk",
				Files: []delugeclient.File{
					{
						Size: 1800000000,
						Path: "Somesubfolder/series/Episode1.mp4",
					},
					{
						Size: 19021903,
						Path: "Somesubfolder/series/Episode2.mkv",
					},
				},
			}},
			[]retrieval.Entry{
				{
					Name: "Episode2.mkv",
					AdditionalData: TorrentClientEntry{
						TorrentClientName: "deluge",
						TorrentName:       "Family.Guy.S22E14.1080p.WEB.H264",
						DownloadFilePath:  "/home/someguy/downloads/deluge/Somesubfolder/series/Episode2.mkv",
						DownloadedAt:      time.Date(2024, time.September, 4, 10, 8, 10, 0, time.UTC),
						Ratio:             1.74,
						FileSizeBytes:     19021903,
						TrackerHost:       "sometrakkr.co.uk",
					},
				},
			},
		},
		{
			"can parse zero files in a torrent status",
			fields{nil},
			args{&delugeclient.TorrentStatus{
				ActiveTime:       0,
				CompletedTime:    1725444490,
				Ratio:            1.74,
				SavePath:         "/home/someguy/downloads/deluge",
				DownloadLocation: "/home/someguy/downloads/deluge",
				Name:             "Family.Guy.S22E14.1080p.WEB.H264",
				SeedingTime:      751680, // 8.7 days
				TrackerHost:      "sometrakkr.co.uk",
				Files:            []delugeclient.File{},
			}},
			[]retrieval.Entry{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DelugeEntryRetriever{
				nil, tt.fields.allowedFileEndings,
			}
			if got := d.parseDelugeTorrentStatus(tt.args.torrentStatus); !reflect.DeepEqual(got, tt.want) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
