package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/go-rod/rod/lib/utils"
)

var (
	// MirrorChromium to fetch the latest chromium version
	MirrorChromium = "https://npm.taobao.org/mirrors/chromium-browser-snapshots/Linux_x64/"
	// MirrorChromiumRegExp to match the MirrorChromium html source
	MirrorChromiumRegExp = regexp.MustCompile(`\Q"/mirrors/chromium-browser-snapshots/Linux_x64/\E(\d+)`)
)

var slash = filepath.FromSlash

func main() {
	res, err := http.Get(MirrorChromium)
	utils.E(err)

	matchs := MirrorChromiumRegExp.FindAllStringSubmatch(utils.MustReadString(res.Body), -1)
	if len(matchs) <= 0 {
		utils.E(fmt.Errorf("cannot match version of the latest chromium from %s", MirrorChromium))
	}

	revision := matchs[len(matchs)-1][1]

	if revision == "" {
		utils.E(fmt.Errorf("empty version of the latest chromium %s", revision))
	}

	build := utils.S(`// generated by running "go generate" on project root

package launcher

// DefaultRevision for chrome
// curl -s -S https://www.googleapis.com/download/storage/v1/b/chromium-browser-snapshots/o/Mac%2FLAST_CHANGE\?alt\=media
const DefaultRevision = {{.revision}}
`,
		"revision", revision,
	)

	utils.E(utils.OutputFile(slash("lib/launcher/revision.go"), build))

}
