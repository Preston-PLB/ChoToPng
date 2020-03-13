package main

import (
	"bufio"
	"fmt"
	"github.com/tdewolff/canvas"
	"math"
	"os"
	"strings"
)

type Line struct {
	lyrics string
	chords []Chord
	tags   []Tag
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

	var lines []Line

	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "CCLI") {
			break
		}
		lines = append(lines, parseLine(scanner.Text()))
	}

	sections := splitSections(lines)

	initCanvas(3840, 1770)
	renderSections(sections, context)
}

func splitSections(lines []Line) (sections [][]Line) {
	var section []Line
	for i := range lines {

		if len(lines[i].tags) > 0 {
			if lines[i].tags[0].name == "comment" {
				sections = append(sections, section)
				section = nil
			}
		}

		section = append(section, lines[i])
	}
	sections = append(sections, section)

	return sections
}

func parseLine(byteLine string) (line Line) {
	var lyricRaw []byte
	for i := 0; i < len(byteLine); i++ {
		if byteLine[i] == '{' {
			str := string(byteLine[i+1 : len(byteLine)-1])
			raw := strings.Split(str, ": ")

			tag := Tag{raw[0], raw[1]}
			line.tags = append(line.tags, tag)
			return line
		}

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
		renderSection(sections[i], c)
	}
}

func renderSection(section []Line, c *canvas.Context) {
	//setUp canvas
	c.SetFillColor(canvas.White)
	fmt.Println(calcFontSize(section), " in ", section[0].tags)
	//find font size

}

func calcFontSize(section []Line) (pnt float64) {

}
