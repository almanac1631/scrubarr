package torrent_clients

import (
	"testing"
	"time"

	"github.com/almanac1631/scrubarr/internal/pkg/common"
	"github.com/autobrr/go-rtorrent"
	"github.com/stretchr/testify/assert"
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
		want   map[common.EntryName]common.Entry
	}{
		{
			"can parse multiple torrent files",
			fields{[]string{".mkv"}},
			args{rtorrent.Torrent{
				Hash:     "somehash",
				Path:     "/some/dir/Some cheeky torrent file bundle",
				Name:     "Some cheeky torrent file bundle",
				Ratio:    0.84,
				Finished: time.Date(2024, time.January, 13, 14, 10, 39, 0, time.UTC),
			}, []rtorrent.File{
				{"Movie1.mkv", 10921},
				{"Movie2.mkv", 687492},
			}},
			map[common.EntryName]common.Entry{
				"Movie1.mkv": {
					Name:     "Movie1.mkv",
					FilePath: "/some/dir/Some cheeky torrent file bundle/Movie1.mkv",
					ParentId: "somehash",
					AdditionalData: TorrentClientEntry{
						ID:                "somehash",
						TorrentClientName: "rtorrent",
						TorrentName:       "Some cheeky torrent file bundle",
						DownloadFilePath:  "/some/dir/Some cheeky torrent file bundle/Movie1.mkv",
						DownloadedAt:      time.Date(2024, time.January, 13, 14, 10, 39, 0, time.UTC),
						Ratio:             0.84,
						FileSizeBytes:     10921,
						TrackerHost:       "<unknown>",
					},
				},
				"Movie2.mkv": {
					Name:     "Movie2.mkv",
					FilePath: "/some/dir/Some cheeky torrent file bundle/Movie2.mkv",
					ParentId: "somehash",
					AdditionalData: TorrentClientEntry{
						ID:                "somehash",
						TorrentClientName: "rtorrent",
						TorrentName:       "Some cheeky torrent file bundle",
						DownloadFilePath:  "/some/dir/Some cheeky torrent file bundle/Movie2.mkv",
						DownloadedAt:      time.Date(2024, time.January, 13, 14, 10, 39, 0, time.UTC),
						Ratio:             0.84,
						FileSizeBytes:     687492,
						TrackerHost:       "<unknown>",
					},
				},
			},
		},
		{
			"can filter torrent files with invalid file extensions",
			fields{[]string{".mkv"}},
			args{rtorrent.Torrent{
				Hash:     "someotherhash",
				Path:     "/some/dir/Some cheeky torrent file bundle",
				Name:     "Some cheeky torrent file bundle",
				Ratio:    0.84,
				Finished: time.Date(2024, time.January, 13, 14, 10, 39, 0, time.UTC),
			}, []rtorrent.File{
				{"Movie1.mkv", 10921},
				{"Movie3.mp4", 89428920},
				{"Movie2.mkv", 687492},
			}},
			map[common.EntryName]common.Entry{
				"Movie1.mkv": {
					Name:     "Movie1.mkv",
					FilePath: "/some/dir/Some cheeky torrent file bundle/Movie1.mkv",
					ParentId: "someotherhash",
					AdditionalData: TorrentClientEntry{
						ID:                "someotherhash",
						TorrentClientName: "rtorrent",
						TorrentName:       "Some cheeky torrent file bundle",
						DownloadFilePath:  "/some/dir/Some cheeky torrent file bundle/Movie1.mkv",
						DownloadedAt:      time.Date(2024, time.January, 13, 14, 10, 39, 0, time.UTC),
						Ratio:             0.84,
						FileSizeBytes:     10921,
						TrackerHost:       "<unknown>",
					},
				},
				"Movie2.mkv": {
					Name:     "Movie2.mkv",
					FilePath: "/some/dir/Some cheeky torrent file bundle/Movie2.mkv",
					ParentId: "someotherhash",
					AdditionalData: TorrentClientEntry{
						ID:                "someotherhash",
						TorrentClientName: "rtorrent",
						TorrentName:       "Some cheeky torrent file bundle",
						DownloadFilePath:  "/some/dir/Some cheeky torrent file bundle/Movie2.mkv",
						DownloadedAt:      time.Date(2024, time.January, 13, 14, 10, 39, 0, time.UTC),
						Ratio:             0.84,
						FileSizeBytes:     687492,
						TrackerHost:       "<unknown>",
					},
				},
			},
		},
		{
			"uses the torrent name if only one valid file exists in the torrent",
			fields{[]string{".mkv"}},
			args{rtorrent.Torrent{
				Hash:     "yetanotherhash",
				Path:     "/some/dir/Some cheeky torrent file bundle",
				Name:     "Some cheeky torrent file bundle.mkv",
				Ratio:    0.84,
				Finished: time.Date(2024, time.January, 13, 14, 10, 39, 0, time.UTC),
			}, []rtorrent.File{
				{"Movie1.mkv", 10921},
				{"Movie2.mp4", 89428920},
			}},
			map[common.EntryName]common.Entry{
				"Some cheeky torrent file bundle.mkv": {
					Name:     "Some cheeky torrent file bundle.mkv",
					FilePath: "/some/dir/Some cheeky torrent file bundle/Movie1.mkv",
					ParentId: "yetanotherhash",
					AdditionalData: TorrentClientEntry{
						ID:                "yetanotherhash",
						TorrentClientName: "rtorrent",
						TorrentName:       "Some cheeky torrent file bundle.mkv",
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
			"uses the torrent name if only one valid file exists in the torrent and adds the file extension",
			fields{[]string{".mkv"}},
			args{rtorrent.Torrent{
				Hash:     "yetanotherhash",
				Path:     "/some/dir/Some cheeky torrent file bundle",
				Name:     "Some cheeky torrent file bundle",
				Ratio:    0.84,
				Finished: time.Date(2024, time.January, 13, 14, 10, 39, 0, time.UTC),
			}, []rtorrent.File{
				{"Movie1.mkv", 10921},
				{"Movie2.mp4", 89428920},
			}},
			map[common.EntryName]common.Entry{
				"Some cheeky torrent file bundle.mkv": {
					Name:     "Some cheeky torrent file bundle.mkv",
					FilePath: "/some/dir/Some cheeky torrent file bundle/Movie1.mkv",
					ParentId: "yetanotherhash",
					AdditionalData: TorrentClientEntry{
						ID:                "yetanotherhash",
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
