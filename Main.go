package main

import (
	"bufio"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/rasterizer"
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

var outPath = ""

func main() {
	renderSong("overcome-A.cho")
}

func getName(stringPath string) (name string) {
	lastSlash := int(math.Max(0, float64(strings.LastIndex(stringPath, "/"))))
	lastDot := strings.LastIndex(stringPath, ".cho")

	return stringPath[lastSlash+1 : lastDot]
}

func renderSong(filePath string) {
	file, err := os.Open(filePath)
	outPath = filePath
	handle(err)
	scanner := bufio.NewScanner(file)

	var sections []Section

	var section = Section{}

	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "CCLI") {
			break
		} else if strings.HasPrefix(text, "{") {
			section.tags = append(section.tags, parseTag(text))
		} else if len(text) > 0 {
			section.lines = append(section.lines, parseLine(text))
		} else {
			sections = append(sections, section)
			section = Section{}
		}

	}
	sections = append(sections, section)
	section = Section{}

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
	if err := fontFamily.LoadFontFile("/usr/share/fonts/truetype/ubuntu/Ubuntu-M.ttf", canvas.FontRegular); err != nil {
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

	c, ctx := initCanvas(1920, 1080)

	//setUp canvas
	ctx.SetFillColor(canvas.Black)
	fontSize, hMax, wMax := calcFontSize(section, c)

	calcPixelOffset(&section, fontSize, c)

	lineOffset := math.Max(hMax, 0)
	yOffset := math.Max((ctx.Height()-float64(len(section.lines)*2)*hMax)/2, 0)
	xOffset := math.Max((ctx.Width()-wMax)/2, 0)

	face := fontFamily.Face(fontSize, canvas.White, canvas.FontRegular, canvas.FontNormal)
	chordFace := fontFamily.Face(fontSize-40, canvas.White, canvas.FontRegular, canvas.FontNormal)

	i := 0
	for _, line := range section.lines {
		for _, chord := range line.chords {
			chordLine := canvas.NewTextLine(chordFace, chord.name, canvas.Left)
			y := (c.H - yOffset) - (lineOffset * float64(i))
			ctx.DrawText(xOffset+chord.pixelOffset, y, chordLine)
		}
		y := (c.H - yOffset) - (lineOffset * float64(i+1))
		lineBox := canvas.NewTextLine(face, line.lyrics, canvas.Left)
		ctx.DrawText(xOffset, y, lineBox)
		i += 2
	}

	name := section.tags[0].value

	err := c.WriteFile(name+".png", rasterizer.PNGWriter(1.0))

	handle(err)
}

func getTextBoxBounds(fontSize float64, str string, c *canvas.Canvas) canvas.Rect {
	face := fontFamily.Face(fontSize, canvas.Black, canvas.FontRegular, canvas.FontNormal)
	var box = canvas.NewTextLine(face, str, canvas.Left)

	return box.Bounds()
}

func calcFontSize(section Section, c *canvas.Canvas) (pnt, hMax, wMax float64) {

	fontSize := 12.0
	fontHeight := 0.0
	fontWidth := 0.0

	lines := section.lines

	longestLine := ""
	for _, line := range lines {
		if len(line.lyrics) > len(longestLine) {
			longestLine = line.lyrics
		}
	}

	if !strings.ContainsAny(longestLine, "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM") {
		return fontSize, 0.0, 0.0
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
