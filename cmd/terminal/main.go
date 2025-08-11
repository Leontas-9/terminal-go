package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/Leontas-9/terminal-go/ansi"
	"github.com/Leontas-9/terminal-go/render"
)

func main() {
	filepath := "Leontas-9/terminal-go/image/burger.gif"

	img, err := terminal.LoadImage(filepath)
	if err != nil { fmt.Print(err) }
	defer terminal.PutReusableRGBA(img)

	src := terminal.NewImage(img)

	_,err = src.Print()
	if err != nil { fmt.Print(err) }
	os.Stdout.Write(ansi.PaintRune(' ', color.RGBA{R: 100, A: 255}, color.RGBA{R: 255, G: 82, B: 197, A: 255}, true))
	// fmt.Print("\n\n\n\033[38;2;255;82;197;48;2;155;106;0mHello, World!\033[0m\n\n\n")

}