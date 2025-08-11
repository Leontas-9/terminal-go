package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/Leontas-9/terminal-go/ansi"
	"github.com/Leontas-9/terminal-go/render"
)

func main() {
	filepath := "Leontas-9/terminal-go/image_test/hollow-knight.png"

	img, err := terminal.LoadImage(filepath)
	if err != nil { fmt.Print(err) }
	defer terminal.PutReusableRGBA(img)

	src := terminal.NewImage(img)

	_,err = src.Print()
	if err != nil { fmt.Print(err) }

	os.Stdout.Write(ansi.PaintRune(' ', 
								   color.RGBA{R: 100, A: 255}, 
								   color.RGBA{R: 255, G: 82, B: 197, A: 255}, 
								   true),
				   )
}
