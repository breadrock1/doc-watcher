package samba

import (
	"doc-notifier/internal/config"
	"github.com/hirochachacha/go-smb2"
	"log"
	"net"
)

type SambaExporter struct {
	client *smb2.Session
}

func New(sambaConfig config.SambaConfig) *SambaExporter {
	conn, err := net.Dial("tcp", sambaConfig.Address)
	if err != nil {
		log.Println("failed to connect to samba: ", err.Error())
		return nil
	}
	defer conn.Close()

	smbDialer := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     sambaConfig.Username,
			Password: sambaConfig.Password,
		},
	}

	smbSession, err := smbDialer.Dial(conn)
	if err != nil {
		log.Println("failed to connect to samba: ", err.Error())
		return nil
	}

	return &SambaExporter{
		client: smbSession,
	}
}

func (se *SambaExporter) GetListDirs() ([]string, error) {
	names, err := se.client.ListSharenames()
	if err != nil {
		return nil, err
	}

	return names, nil
}

func (se *SambaExporter) Close() error {
	return se.client.Logoff()
}
