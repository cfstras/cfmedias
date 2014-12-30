package db_test

import (
	"github.com/cfstras/cfmedias/db"
	"testing"
)

type item struct {
	file, title, artist           string
	expectedTitle, expectedArtist string
}

func TestTitleMagic(t *testing.T) {
	m := []item{
		{"Olafur Arnalds & Nils Frahm - B1.flac", "B1", "Olafur Arnalds & Nils Frahm",
			"B1", "Olafur Arnalds & Nils Frahm"},
		{"Olafur Arnalds & Nils Frahm - B1.flac", "B1", "",
			"B1", "Olafur Arnalds & Nils Frahm"},
		{"Olafur Arnalds & Nils Frahm - B1.flac", "", "Olafur Arnalds & Nils Frahm",
			"B1", "Olafur Arnalds & Nils Frahm"},
		{"Olafur Arnalds & Nils Frahm - B1.flac", "", "",
			"B1", "Olafur Arnalds & Nils Frahm"},

		{"Ólafur Arnalds & Nils Frahm - B1.flac", "B1", "Ólafur Arnalds & Nils Frahm",
			"B1", "Ólafur Arnalds & Nils Frahm"},
		{"Ólafur Arnalds & Nils Frahm - B1.flac", "B1", "",
			"B1", "Ólafur Arnalds & Nils Frahm"},
		{"Ólafur Arnalds & Nils Frahm - B1.flac", "", "Ólafur Arnalds & Nils Frahm",
			"B1", "Ólafur Arnalds & Nils Frahm"},
		{"Ólafur Arnalds & Nils Frahm - B1.flac", "", "",
			"B1", "Ólafur Arnalds & Nils Frahm"},

		{"01 Olafur Arnalds & Nils Frahm - B1.flac", "B1", "Olafur Arnalds & Nils Frahm",
			"B1", "Olafur Arnalds & Nils Frahm"},
		{"01 Olafur Arnalds & Nils Frahm - B1.flac", "B1", "",
			"B1", "Olafur Arnalds & Nils Frahm"},
		{"01 Olafur Arnalds & Nils Frahm - B1.flac", "", "Olafur Arnalds & Nils Frahm",
			"B1", "Olafur Arnalds & Nils Frahm"},
		{"01 Olafur Arnalds & Nils Frahm - B1.flac", "", "",
			"B1", "Olafur Arnalds & Nils Frahm"},

		{"01. Olafur Arnalds & Nils Frahm - B1.flac", "B1", "Olafur Arnalds & Nils Frahm",
			"B1", "Olafur Arnalds & Nils Frahm"},
		{"01. Olafur Arnalds & Nils Frahm - B1.flac", "B1", "",
			"B1", "Olafur Arnalds & Nils Frahm"},
		{"01. Olafur Arnalds & Nils Frahm - B1.flac", "", "Olafur Arnalds & Nils Frahm",
			"B1", "Olafur Arnalds & Nils Frahm"},
		{"01. Olafur Arnalds & Nils Frahm - B1.flac", "", "",
			"B1", "Olafur Arnalds & Nils Frahm"},

		// never change existing attrs
		{"Britney Spears feat NERD - Boys.mp3", "", "Britney Spears",
			"Britney Spears feat NERD - Boys", "Britney Spears"},

		{"02 - Chopin - Prelude Op. 28.7 in A", "", "Chopin",
			"Prelude Op. 28.7 in A", "Chopin"},
	}
	for _, e := range m {
		title, artist := db.TitleMagic(e.file, e.title, e.artist)
		if e.expectedArtist != artist || e.expectedTitle != title {
			t.Log("NOPE:", "title:", e.title, "artist:", e.artist,
				"\n  file:", e.file, "title:", title, "artist:", artist)
			t.Fail()
		}
	}
}
