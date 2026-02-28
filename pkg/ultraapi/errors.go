package ultraapi

import "fmt"

type ErrUnexpectedApiResp struct {
	RespCode int
	Resp     []byte
}

func (err ErrUnexpectedApiResp) Error() string {
	return fmt.Sprintf("invalid api response (status code: %d): %q", err.RespCode, string(err.Resp))
}
