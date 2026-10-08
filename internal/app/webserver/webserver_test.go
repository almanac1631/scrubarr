package webserver

import (
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_removeStaleSocket(t *testing.T) {
	sock := filepath.Join(t.TempDir(), "s.sock")

	// stale: listener killed without unlink
	l, err := net.Listen("unix", sock)
	require.NoError(t, err)
	l.(*net.UnixListener).SetUnlinkOnClose(false)
	require.NoError(t, l.Close())
	require.FileExists(t, sock)
	require.NoError(t, removeStaleSocket(sock))
	_, err = os.Stat(sock)
	require.ErrorIs(t, err, os.ErrNotExist)

	// live: must refuse
	l, err = net.Listen("unix", sock)
	require.NoError(t, err)
	defer l.Close()
	require.Error(t, removeStaleSocket(sock))
	require.FileExists(t, sock)
}
