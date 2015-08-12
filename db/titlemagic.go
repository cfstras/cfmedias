package db

import (
	"regexp"
	"strings"

	log "github.com/cfstras/cfmedias/logger"
)

func TitleMagic(file, titleOld, artistOld string) (title, artist string) {
	title, artist = titleOld, artistOld

	title, artist = doRegex(title, artist, file)
	if title == "" {
		title = removeExtensions(file)
	}
	if titleOld != title || artistOld != artist {
		log.Log.Println("corrected file:", file,
			"\n   title:", titleOld,
			"artist:", artistOld,
			"\nto title:", title, "artist:", artist)
	}
	return
}

func removeExtensions(str string) string {
	for _, s := range Types {
		if strings.HasSuffix(str, "."+s) {
			return str[:len(str)-len(s)-1]
		}
	}
	return str
}

const (
	wordRe = `(?U:[\p{L}\p{Nd} !"#$%&'()*+,.:;<=>?@]{3,})`
)

func doRegex(oldTitle, oldArtist, file string) (title, artist string) {
	var artistRe, titleRe string
	title, artist = oldTitle, oldArtist
	if artist == "" {
		artistRe = `(?P<artist>` + wordRe + `)`
	} else {
		escapedArtist := regexp.QuoteMeta(artist)
		artistRe = `(?P<artist>(?i:` + escapedArtist + `))`
	}
	if title == "" {
		titleRe = `(?P<title>` + wordRe + `)`
	} else {
		escapedTitle := regexp.QuoteMeta(title)
		titleRe = `(?P<title>(?i:` + escapedTitle + `))`
	}

	re := regexp.MustCompile(`^(?:(?P<tracknum>\d{1,4})(?:. ?| | - ))?` +
		artistRe +
		`(?: ?- ?| )` +
		titleRe + `$`)
	fileN := removeExtensions(file)
	n1 := re.SubexpNames()
	matches := re.FindAllStringSubmatch(fileN, -1)
	if len(matches) == 1 {
		r2 := matches[0]
		md := map[string]string{}
		for i, n := range r2 {
			md[n1[i]] = n
		}
		if oldArtist == "" {
			artist = strings.TrimSpace(md["artist"])
		}
		if oldTitle == "" {
			title = strings.TrimSpace(md["title"])
		}
	}
	return

}
