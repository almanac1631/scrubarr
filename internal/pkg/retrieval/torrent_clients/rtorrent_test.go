package torrent_clients

import (
	"github.com/almanac1631/scrubarr/internal/pkg/retrieval"
	"github.com/autobrr/go-rtorrent"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRtorrentEntryRetriever_parseTorrentFileList(t *testing.T) {
	type fields struct {
		allowedFileEndings []string
	}
	type args struct {
		torrent         rtorrent.Torrent
		torrentFileList []rtorrent.File
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[retrieval.EntryName]retrieval.Entry
	}{
		{
			"can parse a single torrent file",
			fields{[]string{".mkv"}},
			args{rtorrent.Torrent{
				Path:     "/some/dir/Some cheeky torrent file bundle",
				Name:     "Some cheeky torrent file bundle",
				Ratio:    0.84,
				Finished: time.Date(2024, time.January, 13, 14, 10, 39, 0, time.UTC),
			}, []rtorrent.File{
				{"Movie1.mkv", 10921},
			}},
			map[retrieval.EntryName]retrieval.Entry{
				"Movie1.mkv": {
					Name: "Movie1.mkv",
					AdditionalData: TorrentClientEntry{
						TorrentClientName: "rtorrent",
						TorrentName:       "Some cheeky torrent file bundle",
						DownloadFilePath:  "/some/dir/Some cheeky torrent file bundle/Movie1.mkv",
						DownloadedAt:      time.Date(2024, time.January, 13, 14, 10, 39, 0, time.UTC),
						Ratio:             0.84,
						FileSizeBytes:     10921,
						TrackerHost:       "<unknown>",
					},
				},
			},
		},
		{
			"can filter a torrent file with invalid file extensions",
			fields{[]string{".mkv"}},
			args{rtorrent.Torrent{
				Path:     "/some/dir/Some cheeky torrent file bundle",
				Name:     "Some cheeky torrent file bundle",
				Ratio:    0.84,
				Finished: time.Date(2024, time.January, 13, 14, 10, 39, 0, time.UTC),
			}, []rtorrent.File{
				{"Movie1.mkv", 10921},
				{"Movie2.mp4", 89428920},
			}},
			map[retrieval.EntryName]retrieval.Entry{
				"Movie1.mkv": {
					Name: "Movie1.mkv",
					AdditionalData: TorrentClientEntry{
						TorrentClientName: "rtorrent",
						TorrentName:       "Some cheeky torrent file bundle",
						DownloadFilePath:  "/some/dir/Some cheeky torrent file bundle/Movie1.mkv",
						DownloadedAt:      time.Date(2024, time.January, 13, 14, 10, 39, 0, time.UTC),
						Ratio:             0.84,
						FileSizeBytes:     10921,
						TrackerHost:       "<unknown>",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RtorrentEntryRetriever{
				nil, tt.fields.allowedFileEndings,
			}
			assert.Equalf(t, tt.want, r.parseTorrentFileList(tt.args.torrent, tt.args.torrentFileList), "parseTorrentFileList(%v, %v)", tt.args.torrent, tt.args.torrentFileList)
		})
	}
}
