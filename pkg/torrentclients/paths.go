package torrentclients

import "strings"

// joinRemotePath joins a directory reported by a torrent client host using the separator of that host.
func joinRemotePath(dir string, name string) string {
	if dir == "" {
		return ""
	}
	if strings.HasSuffix(dir, "/") || strings.HasSuffix(dir, `\`) {
		return dir + name
	}
	separator := "/"
	if strings.Contains(dir, `\`) && !strings.Contains(dir, "/") {
		separator = `\`
	}
	return dir + separator + name
}
