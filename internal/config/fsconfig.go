package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	ProgramName              = "wizefs"
	ProgramVersion           = "0.0.5"
	FilesystemCurrentVersion = 1
	FilesystemConfigFilename = "wizefs.conf"
)

type FSType int

const (
	FSHack FSType = iota - 1 // -1
	FSNone
	FSLoopback
	FSZip
)

type FilesystemConfig struct {
	// Creator is the WizeFS version string
	Creator string `json:"creator"`
	// Version is the version this filesystem uses
	Version    uint16 `json:"version"`
	Origin     string `json:"origin"`
	OriginPath string `json:"originpath"`
	Type       FSType `json:"type"`

	filename string
}

func NewFilesystemConfig(origin, originPath string, itype FSType) *FilesystemConfig {
	return &FilesystemConfig{
		Creator:    ProgramName + " ver. " + ProgramVersion,
		Version:    FilesystemCurrentVersion,
		Origin:     origin,
		OriginPath: originPath,
		Type:       itype,
		filename:   filepath.Join(origin, FilesystemConfigFilename),
	}
}

func (c *FilesystemConfig) Save() error {
	tmp := c.filename + ".tmp"
	// 0400 permissions: wizefs.conf should be kept secret and never be written to.
	fd, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0400)
	if err != nil {
		return err
	}
	js, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	// For convenience for the user, add a newline at the end.
	js = append(js, '\n')
	_, err = fd.Write(js)
	if err != nil {
		return err
	}
	err = fd.Sync()
	if err != nil {
		return err
	}
	err = fd.Close()
	if err != nil {
		return err
	}
	err = os.Rename(tmp, c.filename)
	return err
}

// TODO: load FilesystemConfig
