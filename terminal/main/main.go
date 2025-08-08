package main

import (
	"terminal/render"
	"fmt"
	"image"
	"log"
	"golang.org/x/image/draw"
)

func main() {
	filepath := "image_test/man.jpg"

	imgRGBA, err := terminal.LoadImage(filepath)
	if err != nil { fmt.Print(err) }
	defer terminal.PutReusableRGBA(imgRGBA)

	src := terminal.NewCustomImage(
		imgRGBA, 
		imgRGBA.Rect,
		image.Pt(0,0),
		draw.NearestNeighbor, 
		&terminal.UI_Settings{
			ShowCursor: false,
			AlternativeScreen: false,
			EraseScreen: false,
			Auto_Wrap: true,
		},
	)

	_,err = src.Print()
	if err != nil { log.Fatal(err) }

	err = src.Displacement()
	if err != nil { log.Fatal(err) }
}