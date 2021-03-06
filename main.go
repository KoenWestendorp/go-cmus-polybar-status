package main

// A cmus status module program
// Outputs a single line string with the status of the track currently playing
// on your cmus player, through cmus-remote.
//
// Koen Westendorp, 2020
// koenw.gitlab.io
// GitHub: KoenWestendorp

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"math"
	"strconv"
)

func main() {
    stat := parseStatus(getStatus())

    /*--- CONFIG ---*/
	// The separator placed in between the different elements.
	sep := "  "

    // Set the length of the progress bar.
    barLength := 7

	// The elements are defined here. You can add your own if you like.
	disp := stat.artist + " \u2014 " + stat.title // \u2014 represents an em dash.
	ind := statusIndicator(stat.playing)
	prog := progressIndicator(stat.duration, stat.position, barLength)
	dur := formatDuration(stat.duration)
	pos := formatDuration(stat.position)
	// album := "(" + stat.album + ")" 

	// This array can be rearranged, in order to modify the output.
	output := []string {ind, pos, prog, dur, disp}

	fmt.Print(strings.Join(output, sep))
}

func getStatus() []string {
	cmd := exec.Command("cmus-remote", "-Q")

	status, err := cmd.CombinedOutput()
	if err != nil {
		os.Exit(1)
	}

	output := strings.Split(string(status), "\n")

	return output
}

type status struct {
	playing bool

	title string
	artist string
	album string

	duration int
	position int
}

func parseStatus(s []string) (status) {
	var stat status

	offset := 0

	playing := strings.TrimPrefix(s[0], "status ")
	if playing == "playing" {
		stat.playing = true
	} else {
		if playing == "stopped" {
			offset = -2
		}

		stat.playing = false
	}

	stat.title = strings.TrimPrefix(s[4 + offset], "tag title ")
	stat.artist = strings.TrimPrefix(s[5 + offset], "tag artist ")
	stat.album = strings.TrimPrefix(s[6 + offset], "tag album ")

	var err1, err2 error

	stat.duration, err1 = strconv.Atoi(strings.TrimPrefix(s[2 + offset], "duration "))
	stat.position, err2 = strconv.Atoi(strings.TrimPrefix(s[3 + offset], "position "))

	if err1 != nil || err2 != nil {
		os.Exit(1)
	}

	return stat
}

func parseDuration(seconds int) (int, int){
	if seconds < 0 {
		return 0, 0
	} else {
		minutes := float64(seconds) / 60.0
		min := math.Floor(minutes)
		sec := math.Floor((minutes - min) * 60)

		return int(min), int(sec)
	}
}

func formatDuration(seconds int) string {
	sep := ":"

	min, sec := parseDuration(seconds)

    // Format: '04:20'.
	return fmt.Sprintf("%02d%s%02d", min, sep, sec)
}

func statusIndicator(playing bool) string {
	playchar := ">"
	pausechar := "\""

	if playing {
		return playchar
	} else {
		return pausechar
	}
}

func progressIndicator(dur, pos, len int) string {
	linechar := "-"
	pointerchar := "|"

	progress := float64(pos) / float64(dur)

	pre := strings.Repeat(linechar, int(progress * float64(len)))
	suf := strings.Repeat(linechar, int((1 - progress) * float64(len)))

	return pre + pointerchar + suf
}
