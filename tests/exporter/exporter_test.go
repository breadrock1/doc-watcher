package exporter

import (
	"doc-notifier/internal/config"
	"doc-notifier/internal/exporter/samba"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExporter(t *testing.T) {

	smbConfig := config.SambaConfig{
		Address:  "localhost:445",
		Username: "Bread White",
		Password: "O1adush3k",
	}

	t.Run("Connect to local samba", func(t *testing.T) {
		smbLocalConn := samba.New(smbConfig)
		dirs, err := smbLocalConn.GetListDirs()
		assert.NoError(t, err, "failed to connect to smb server")

		fmt.Println(dirs)
	})
}
