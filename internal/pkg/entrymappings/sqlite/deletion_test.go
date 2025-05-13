package sqlite

import (
	"fmt"
	"github.com/almanac1631/scrubarr/internal/pkg/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_getIdFromApiResp1(t *testing.T) {
	type args struct {
		retrieverInfo common.RetrieverInfo
		apiResp       string
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"parses id from an ArrApp API response",
			args{
				common.RetrieverInfo{
					Category: "arr_app",
				},
				`{"ID":1989302}`,
			},
			int64(1989302),
			assert.NoError,
		},
		{
			"parses id from an TorrentClient API response",
			args{
				common.RetrieverInfo{
					Category: "torrent_client",
				},
				`{"ID":"somehash"}`,
			},
			"somehash",
			assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getIdFromApiResp(tt.args.retrieverInfo, tt.args.apiResp)
			if !tt.wantErr(t, err, fmt.Sprintf("getIdFromApiResp(%v, %v)", tt.args.retrieverInfo, tt.args.apiResp)) {
				return
			}
			assert.Equalf(t, tt.want, got, "getIdFromApiResp(%v, %v)", tt.args.retrieverInfo, tt.args.apiResp)
		})
	}
}
