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
	chords *[]Chord
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
			*line.chords = append(*line.chords, chord)
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
	for i := range sections {
		if len(sections[i].tags) > 0 && sections[i].tags[0].name == "comment" {
			renderSection(sections[i])
		}
	}
}

func renderSection(section Section) {

	c, ctx := initCanvas(3840, 1770)

	//setUp canvas
	ctx.SetFillColor(canvas.White)
	fontSize, hMax, wMax := calcFontSize(section, c)

	calcPixelOffset(&section, fontSize, c)

	lineOffset := hMax / float64(len(section.lines)*2)
	yOffset := (ctx.Height() - hMax) / 2
	xOffset := (ctx.Width() - wMax) / 2

	face := fontFamily.Face(fontSize, canvas.Black, canvas.FontRegular, canvas.FontNormal)

	for i := range section.lines {
		line := section.lines[i]
		for j := range *line.chords {
			chord := *line.chords
			chordLine := canvas.NewTextLine(face, chord[j].name, canvas.Left)
			ctx.DrawText(xOffset+chord[j].pixelOffset, (lineOffset*float64(i))+yOffset, chordLine)
		}
		lineBox := canvas.NewTextLine(face, line.lyrics, canvas.Left)
		ctx.DrawText(xOffset, (lineOffset*float64(i))+yOffset+lineOffset, lineBox)
	}

	err := c.SavePNG(section.tags[0].name+".png", 5.0)

	handle(err)
}

func getTextBoxBounds(fontSize float64, str string, c *canvas.Canvas) canvas.Rect {
	face := fontFamily.Face(fontSize, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	box := canvas.NewTextBox(face, str, c.H, c.W, canvas.Left, canvas.Top, 0.0, 0.0)

	return box.Bounds()
}

func calcFontSize(section Section, c *canvas.Canvas) (pnt, hMax, wMax float64) {

	fontSize := 100.0
	fontHeight := 0.0
	fontWidth := 0.0

	lines := section.lines

	longestLine := ""
	for i := range lines {
		if len(lines[i].lyrics) > len(longestLine) {
			longestLine = lines[i].lyrics
		}
	}

	if !strings.ContainsAny(longestLine, "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM") {
		return fontSize, 0.0, 0.0
	}

	for fontWidth < c.W && (fontHeight*2.0*float64(len(lines)-1) < c.H) {
		size := getTextBoxBounds(fontSize, longestLine, c)

		fontHeight = size.H
		fontWidth = size.W

		fontSize += 3
		//fmt.Printf("Testing font %f \n", fontSize)
	}
	return fontSize, fontHeight, fontWidth
}

func calcPixelOffset(section *Section, fontSize float64, c *canvas.Canvas) {
	lines := section.lines
	for i := range lines {
		line := lines[i]
		for j := range *line.chords {
			chord := *line.chords
			chord[j].pixelOffset = getTextBoxBounds(fontSize, lines[j].lyrics[0:chord[j].charOffset], c).W
		}
	}
}
