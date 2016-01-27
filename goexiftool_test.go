package goexiftool

import (
	"testing"
	"time"
)

var testImages = []struct {
	filename     string
	analyzeError error
	numInfos     int
	lens         string
	camera       string
	flash        string
	date         time.Time
	geoTagged    bool
}{
	{
		"test/DSC03696.jpg",
		nil,
		230,
		"Zeiss Batis 85mm F1.8",
		"ILCE-7RM2",
		"Off, Did not fire",
		time.Date(2015, time.December, 30, 11, 18, 29, 0, time.UTC),
		false,
	},
}

func TestAnalyze(t *testing.T) {
	for _, i := range testImages {
		m, err := NewMediaFile(i.filename)
		if err != i.analyzeError {
			t.Errorf("Analyze(%s) returned %s, expected %s!", i.filename, err.Error(), i.analyzeError.Error())
		}
		if numInfos := len(m.Info); numInfos != i.numInfos {
			t.Errorf("Analyze(%s) returned %d infos, expected %d!", i.filename, numInfos, i.numInfos)
		}
	}
}

func TestGet(t *testing.T) {
	for _, i := range testImages {
		m, err := NewMediaFile(i.filename)
		if err != nil {
			t.Errorf("Get(%s) returned %s, expected %s!", i.filename, err.Error(), i.analyzeError.Error())
		}
		flash, err := m.Get("Flash")
		if flash != i.flash {
			t.Errorf("Get(%s)(Flash) returned %s, expected %s!", i.filename, flash, i.flash)
		}
	}
}

func TestGetLens(t *testing.T) {
	for _, i := range testImages {
		m, err := NewMediaFile(i.filename)
		if err != nil {
			t.Errorf("GetLens(%s) returned %s, expected %s!", i.filename, err.Error(), i.analyzeError.Error())
		}
		lens, err := m.GetLens()
		if err != nil {
			t.Errorf("GetLens(%s) returned error %s!", i.filename, err.Error())
		}
		if lens != i.lens {
			t.Errorf("GetLens(%s) returned %s, expected %s!", i.filename, lens, i.lens)
		}
	}
}

func TestGetCamera(t *testing.T) {
	for _, i := range testImages {
		m, err := NewMediaFile(i.filename)
		if err != nil {
			t.Errorf("GetCamera(%s) returned %s, expected %s!", i.filename, err.Error(), i.analyzeError.Error())
		}
		camera, err := m.GetCamera()
		if err != nil {
			t.Errorf("GetCamera(%s) returned error %s!", i.filename, err.Error())
		}
		if camera != i.camera {
			t.Errorf("GetCamera(%s) returned %s, expected %s!", i.filename, camera, i.camera)
		}
	}
}

func TestGetDate(t *testing.T) {
	for _, i := range testImages {
		m, err := NewMediaFile(i.filename)
		if err != nil {
			t.Errorf("GetDate(%s) returned %s, expected %s!", i.filename, err.Error(), i.analyzeError.Error())
		}
		date, err := m.GetDate()
		if err != nil {
			t.Errorf("GetDate(%s) returned error %s!", i.filename, err.Error())
		}
		if date != i.date {
			t.Errorf("GetDate(%s) returned %s, expected %s!", i.filename, date.String(), i.date.String())
		}
	}
}

func TestIsGeoTagged(t *testing.T) {
	for _, i := range testImages {
		m, err := NewMediaFile(i.filename)
		if err != nil {
			t.Errorf("IsGeoTagged(%s) returned %s, expected %s!", i.filename, err.Error(), i.analyzeError.Error())
		}
		geoTagged := m.IsGeoTagged()
		if geoTagged != i.geoTagged {
			t.Errorf("IsGeoTagged(%s) returned %v, expected %v!", i.filename, geoTagged, i.geoTagged)
		}
	}
}
