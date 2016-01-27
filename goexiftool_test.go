package goexiftool

import (
	"errors"
	"fmt"
	"os"
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
		"paf",
		errors.New("does not exist"),
		0,
		"",
		"",
		"",
		time.Time{},
		false,
	},
	{
		"test_files/2016-01-02-13h19m03s_A7RII.85mm1.8Z_03943.jpg",
		nil,
		239,
		"Zeiss Batis 85mm F1.8",
		"ILCE-7RM2",
		"Off, Did not fire",
		time.Date(2016, time.January, 02, 13, 19, 03, 0, time.UTC),
		true,
	},
	{
		"test_files/2015-12-20-11h05m19s_A7RII.55mm1.8ZA_02566.jpg",
		nil,
		239,
		"Sony FE 55mm F1.8 ZA",
		"ILCE-7RM2",
		"Off, Did not fire",
		time.Date(2015, time.December, 20, 11, 05, 19, 0, time.UTC),
		true,
	},
	{
		"test_files/2015-12-21-13h31m44s_A7RII.35mm1.4_02725.jpg",
		nil,
		231,
		"Sigma 35mm f/1.4 DG HSM",
		"ILCE-7RM2",
		"Off, Did not fire",
		time.Date(2015, time.December, 21, 13, 31, 44, 0, time.UTC),
		false,
	},
}

func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewMediaFile(testImages[1].filename)
	}
}

func Test(t *testing.T) {
	for _, i := range testImages {
		m, err := NewMediaFile(i.filename)
		fmt.Println("Testing " + i.filename)
		if i.analyzeError != err {
			if os.IsNotExist(err) && i.analyzeError.Error() != "does not exist" {
				t.Errorf("Analyze(%s) returned %s, expected %s!", i.filename, err.Error(), i.analyzeError.Error())
				return
			}
		}
		if m != nil {
			// Analyze
			if numInfos := len(m.Info); numInfos != i.numInfos {
				t.Errorf("Analyze(%s) returned %d infos, expected %d!", i.filename, numInfos, i.numInfos)
			}
			// Get with Flash exiftool tag
			flash, err := m.Get("Flash")
			if flash != i.flash {
				t.Errorf("Get(%s)(Flash) returned %s, expected %s!", i.filename, flash, i.flash)
			}
			// GetLens
			lens, err := m.GetLens()
			if err != nil {
				t.Errorf("GetLens(%s) returned error %s!", i.filename, err.Error())
			}
			if lens != i.lens {
				t.Errorf("GetLens(%s) returned %s, expected %s!", i.filename, lens, i.lens)
			}
			// GetCamera
			camera, err := m.GetCamera()
			if err != nil {
				t.Errorf("GetCamera(%s) returned error %s!", i.filename, err.Error())
			}
			if camera != i.camera {
				t.Errorf("GetCamera(%s) returned %s, expected %s!", i.filename, camera, i.camera)
			}
			// GetDate
			date, err := m.GetDate()
			if err != nil {
				t.Errorf("GetDate(%s) returned error %s!", i.filename, err.Error())
			}
			if date != i.date {
				t.Errorf("GetDate(%s) returned %s, expected %s!", i.filename, date.String(), i.date.String())
			}
			// IsGeoTagged
			geoTagged := m.IsGeoTagged()
			if geoTagged != i.geoTagged {
				t.Errorf("IsGeoTagged(%s) returned %v, expected %v!", i.filename, geoTagged, i.geoTagged)
			}
		}
	}
}
