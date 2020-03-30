package main

import (
	"bufio"
	"github.com/tdewolff/canvas"
	"math"
	"os"
	"strings"
)

type Section struct {
	lines []Line
	tags  []Tag
}

type Line struct {
	lyrics string
	chords []*Chord
}

type Chord struct {
	name        string
	charOffset  int
	pixelOffset float64
}

func (chord *Chord) calcPixelOffset(lyrics string, fontSize float64, c *canvas.Canvas) {
	chord.pixelOffset = getTextBoxBounds(fontSize, lyrics[0:chord.charOffset], c).W
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

//var canva *canvas.Canvas
//var context *canvas.Context

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

	renderSections(sections)
}

func parseTag(byteLine string) (tag Tag) {
	raw := strings.Split(byteLine, ": ")
	return Tag{raw[0][1:], raw[1][0 : len(raw[1])-1]}
}

func parseLine(byteLine string) (line Line) {
	var lyricRaw []byte
	for i, k := 0, 0; i < len(byteLine); i++ {
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
			chord := Chord{string(chordName), k, 0.0}
			line.chords = append(line.chords, new(Chord))
			line.chords[len(line.chords)-1] = &chord
		} else {
			if byteLine[i] != '\r' {
				lyricRaw = append(lyricRaw, byteLine[i])
			}
			k++
		}
	}

	line.lyrics = string(lyricRaw)

	return line
}

//
// Grapphical Stuff and things
//

func initCanvas(width float64, height float64) (c *canvas.Canvas, context *canvas.Context) {
	fontFamily = canvas.NewFontFamily("Ubuntu")
	fontFamily.Use(canvas.CommonLigatures)
	if err := fontFamily.LoadFontFile("C:\\Windows\\Fonts\\Ubuntu-M.ttf", canvas.FontRegular); err != nil {
		panic(err)
	}

	c = canvas.New(width, height)
	context = canvas.NewContext(c)

	return c, context
}

func renderSections(sections []Section) {
	for _, section := range sections {
		if len(section.tags) > 0 && section.tags[0].name == "comment" {
			renderSection(section)
		}
	}
}

func renderSection(section Section) {

	c, ctx := initCanvas(3840, 1770)

	//setUp canvas
	ctx.SetFillColor(canvas.White)
	fontSize, hMax, wMax := calcFontSize(section, c)

	calcPixelOffset(&section, fontSize, c)

	lineOffset := hMax
	yOffset := math.Max((ctx.Height()-float64(len(section.lines)*2)*hMax)/2, 0)
	xOffset := math.Max((ctx.Width()-wMax)/2, 0)

	face := fontFamily.Face(fontSize, canvas.Black, canvas.FontRegular, canvas.FontNormal)

	i := 0
	for _, line := range section.lines {
		for _, chord := range line.chords {
			chordLine := canvas.NewTextLine(face, chord.name, canvas.Left)
			y := (c.H - yOffset) - (lineOffset * float64(i))
			ctx.DrawText(xOffset+chord.pixelOffset, y, chordLine)
		}
		y := (c.H - yOffset) - (lineOffset * float64(i+1))
		lineBox := canvas.NewTextLine(face, line.lyrics, canvas.Left)
		ctx.DrawText(xOffset, y, lineBox)
		i += 2
	}

	err := c.SavePNG(section.tags[0].value+".png", 1.0)

	handle(err)
}

func getTextBoxBounds(fontSize float64, str string, c *canvas.Canvas) canvas.Rect {
	face := fontFamily.Face(fontSize, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	var box = canvas.NewTextLine(face, str, canvas.Left)

	return box.Bounds()
}

func calcFontSize(section Section, c *canvas.Canvas) (pnt, hMax, wMax float64) {

	fontSize := 100.0
	fontHeight := 0.0
	fontWidth := 0.0

	lines := section.lines

	longestLine := ""
	for _, line := range lines {
		if len(line.lyrics) > len(longestLine) {
			longestLine = line.lyrics
		}
	}

	for fontWidth < c.W && (fontHeight*2.0*float64(len(lines)) < c.H) {
		size := getTextBoxBounds(fontSize, longestLine, c)

		fontHeight = size.H
		fontWidth = size.W

		fontSize += 1
		//fmt.Printf("Testing font %f \n", fontSize)
	}
	return fontSize - 1, fontHeight, fontWidth
}

func calcPixelOffset(section *Section, fontSize float64, c *canvas.Canvas) {
	for _, line := range section.lines {
		for _, chord := range line.chords {
			chord.calcPixelOffset(line.lyrics, fontSize, c)
		}
	}
}
