// Package goexiftool is a very simple wrapper around the excellent exiftool (http://www.sno.phy.queensu.ca/~phil/exiftool/).
package goexiftool

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// MediaFile holds all the metadata provided by exiftool in a map
type MediaFile struct {
	Filename string
	Info     map[string]string
}

// String displays all metadata
func (m *MediaFile) String() string {
	txt := m.Filename + ":\n"
	for k, v := range m.Info {
		txt += "\t" + k + " = " + v + "\n"
	}
	return txt
}

// Analyze calls exiftool on the file and parses its output.
func (m *MediaFile) Analyze() (err error) {
	cmdName, err := exec.LookPath("exiftool")
	if err != nil {
		return errors.New("exiftool is not installed")
	}
	cmdArgs := []string{m.Filename}
	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			res := strings.SplitN(scanner.Text(), ":", 2)
			key := strings.TrimSpace(res[0])
			value := strings.TrimSpace(res[1])
			m.Info[key] = value
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		os.Exit(1)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		os.Exit(1)
	}
	return
}

// Get an exiftool value.
// exiftool has its own entry names, sometimes aggregating several exif/xmp/iptc tags.
func (m *MediaFile) Get(tagLabel string) (tagValue string, err error) {
	tagValue, ok := m.Info[tagLabel]
	if !ok {
		err = errors.New("Unknown tag " + tagLabel)
	}
	return
}

// GetLens used for a MediaFile.
func (m *MediaFile) GetLens() (lens string, err error) {
	lens, ok := m.Info["Lens ID"]
	if !ok {
		err = errors.New("Unknown exiftool value : Lens ID")
	}
	return
}

// GetCamera used for a MediaFile.
func (m *MediaFile) GetCamera() (camera string, err error) {
	camera, ok := m.Info["Camera Model Name"]
	if !ok {
		err = errors.New("Unknown exiftool value : Camera Model Name")
	}
	return
}

// GetDate of creation of this MediaFile.
func (m *MediaFile) GetDate() (date time.Time, err error) {
	dateString, ok := m.Info["Date/Time Original"]
	if !ok {
		err = errors.New("Unknown exiftool value : Date/Time Original")
	}
	date, err = time.Parse("2006:01:02 15:04:05", dateString)
	if err != nil {
		err = errors.New("Date has unexpected format: " + dateString)
	}
	return
}

// IsGeoTagged returns if GPS data is found.
func (m *MediaFile) IsGeoTagged() (isGeoTagged bool) {
	_, isGeoTagged = m.Info["GPS Position"]
	return
}

// getExistingPath ensures a path actually exists, and returns an existing absolute path or an error.
func getExistingPath(path string) (existingPath string, err error) {
	// check root exists or pwd+root exists
	if filepath.IsAbs(path) {
		existingPath = path
	} else {
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		existingPath = filepath.Join(pwd, path)
	}
	// check root exists
	_, err = os.Stat(existingPath)
	return
}

// NewMediaFile initializes a MediaFile and parses its metadata with exiftool.
func NewMediaFile(filename string) (mf *MediaFile, err error) {
	filename, err = getExistingPath(filename)
	if os.IsNotExist(err) {
		return nil, err
	}
	mf = &MediaFile{Filename: filename, Info: make(map[string]string)}
	err = mf.Analyze()
	return
}
