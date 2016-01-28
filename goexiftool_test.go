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
	lens         string
	camera       string
	flash        string
	date         time.Time
	geoTagged    bool
}{
	{
		"paf",
		errors.New("does not exist"),
		"",
		"",
		"",
		time.Time{},
		false,
	},
	{
		"test_files/2016-01-02-13h19m03s_A7RII.85mm1.8Z_03943.jpg",
		nil,
		"Zeiss Batis 85mm F1.8",
		"ILCE-7RM2",
		"Off, Did not fire",
		time.Date(2016, time.January, 02, 13, 19, 03, 0, time.UTC),
		true,
	},
	{
		"test_files/2015-12-20-11h05m19s_A7RII.55mm1.8ZA_02566.jpg",
		nil,
		"Sony FE 55mm F1.8 ZA",
		"ILCE-7RM2",
		"Off, Did not fire",
		time.Date(2015, time.December, 20, 11, 05, 19, 0, time.UTC),
		true,
	},
	{
		"test_files/2015-12-21-13h31m44s_A7RII.35mm1.4_02725.jpg",
		nil,
		"Sigma 35mm f/1.4 DG HSM",
		"ILCE-7RM2",
		"Off, Did not fire",
		time.Date(2015, time.December, 21, 13, 31, 44, 0, time.UTC),
		false,
	},
	{
		"test_files/2009-07-03-16h34m02s_40D.100mm2.8Macro_8702-1.jpg",
		nil,
		"Canon EF 100mm f/2.8 Macro USM",
		"Canon EOS 40D",
		"Off, Did not fire",
		time.Date(2009, time.July, 03, 16, 34, 02, 30000000, time.UTC),
		false,
	},
	{
		"test_files/2009-03-01-10h06m27s_40D.100mm2_7107-1.jpg",
		nil,
		"Canon EF 100mm f/2 USM",
		"Canon EOS 40D",
		"Off, Did not fire",
		time.Date(2009, time.March, 01, 10, 06, 27, 610000000, time.UTC),
		false,
	},
	{
		"test_files/2008-08-05-18h26m36s_40D.70-200mm4L-IS_4296.jpg",
		nil,
		"Canon EF 70-200mm f/4L IS",
		"Canon EOS 40D",
		"Off, Did not fire",
		time.Date(2008, time.August, 05, 18, 26, 36, 290000000, time.UTC),
		false,
	},
	{
		"test_files/2008-06-01-10h25m18s_40D.17-55mm2.8-IS_2300.jpg",
		nil,
		"Canon EF-S 17-55mm f/2.8 IS USM",
		"Canon EOS 40D",
		"Off, Did not fire",
		time.Date(2008, time.June, 01, 10, 25, 18, 780000000, time.UTC),
		false,
	},
	{
		"test_files/2006-02-01-23h59m14s_350D.18-55mm_1042.jpg",
		nil,
		"Unknown 18-55mm",
		"Canon EOS 350D DIGITAL",
		"On, Fired",
		time.Date(2006, time.February, 01, 23, 59, 14, 0, time.UTC),
		true,
	},
	{
		"test_files/2006-01-21-16h08m37s_350D.70-300mm_0860.jpg",
		nil,
		"Unknown 70-300mm",
		"Canon EOS 350D DIGITAL",
		"Off, Did not fire",
		time.Date(2006, time.January, 21, 16, 8, 37, 0, time.UTC),
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
			if len(m.Info) == 0 {
				t.Errorf("Analyze(%s) returned no metadata!", i.filename)
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
