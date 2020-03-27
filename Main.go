package main

import (
	"bufio"
	"github.com/tdewolff/canvas"
	"os"
	"strings"
)

type Section struct {
	lines []Line
	tags  []Tag
}

type Line struct {
	lyrics string
	chords []Chord
}

type Chord struct {
	name        string
	charOffset  int
	pixelOffset float64
}

type Tag struct {
	name  string
	value string
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

var fontFamily *canvas.FontFamily
var canva *canvas.Canvas
var context *canvas.Context

func main() {
	file, err := os.Open("overcome-A.cho")
	handle(err)
	scanner := bufio.NewScanner(file)

	var sections []Section

	var section = Section{}

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "CCLI") {
			break
		} else if strings.HasPrefix(scanner.Text(), "{") {
			section.tags = append(section.tags, parseTag(scanner.Text()))
		} else if len(scanner.Text()) > 0 {
			section.lines = append(section.lines, parseLine(scanner.Text()))
		} else {
			sections = append(sections, section)
			section = Section{}
		}

	}

	initCanvas(3840, 1770)
	renderSections(sections, context)
}

func parseTag(byteLine string) (tag Tag) {
	raw := strings.Split(byteLine, ": ")
	return Tag{raw[0], raw[1]}
}

func parseLine(byteLine string) (line Line) {
	var lyricRaw []byte
	for i := 0; i < len(byteLine); i++ {
		if byteLine[i] == '[' {
			var chordName []byte
			for j := i + 1; j < len(byteLine); j++ {
				if byteLine[j] != ']' {
					chordName = append(chordName, byteLine[j])
				} else {
					i = j
					break
				}
			}
			chord := Chord{string(chordName), i, 0.0}
			line.chords = append(line.chords, chord)
		} else {
			if byteLine[i] != '\r' {
				lyricRaw = append(lyricRaw, byteLine[i])
			}
		}
	}

	line.lyrics = string(lyricRaw)

	return line
}

//
// Grapphical Stuff and things
//

func initCanvas(width float64, height float64) {
	fontFamily = canvas.NewFontFamily("Ubuntu")
	fontFamily.Use(canvas.CommonLigatures)
	if err := fontFamily.LoadFontFile("C:\\Windows\\Fonts\\Ubuntu-M.ttf", canvas.FontRegular); err != nil {
		panic(err)
	}

	canva = canvas.New(width, height)
	context = canvas.NewContext(canva)
}

func renderSections(sections [][]Line, c *canvas.Context) {
	for i := range sections {
		if sections[i][0].tags[0].name == "comment" {
			renderSection(sections[i], c)
		}
	}
}

func renderSection(section []Line, c *canvas.Context) {
	//setUp canvas
	c.SetFillColor(canvas.White)
	//fontSize, hMax, wMax := calcFontSize(section)
	//
	//calcPixelOffset(&section, fontSize)
	//
	//face := fontFamily.Face(fontSize, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	//line :=
}

func getTextBoxBounds(fontSize float64, str string) canvas.Rect {
	face := fontFamily.Face(fontSize, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	box := canvas.NewTextBox(face, str, canva.H, canva.W, canvas.Left, canvas.Top, 0.0, 0.0)

	return box.Bounds()
}

func calcFontSize(section []Line) (pnt, hMax, wMax float64) {

	fontSize := 12.0
	fontHeight := 0.0
	fontWidth := 0.0

	longestLine := ""
	for i := range section {
		if len(section[i].lyrics) > len(longestLine) {
			longestLine = section[i].lyrics
		}
	}

	if !strings.ContainsAny(longestLine, "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM") {
		return fontSize, 0.0, 0.0
	}

	for fontWidth < canva.W && (fontHeight*2.0*float64(len(section)-1) < canva.H) {
		size := getTextBoxBounds(fontSize, longestLine)

		fontHeight = size.H
		fontWidth = size.W

		fontSize += 3
		//fmt.Printf("Testing font %f \n", fontSize)
	}
	return fontSize, fontHeight, fontWidth
}

func calcPixelOffset(section *[]Line, fontSize float64) {

}
