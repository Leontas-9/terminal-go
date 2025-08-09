package main

import (
	"terminal/render"
	"fmt"
)

func main() {
	filepath := "C:/Users/lenna/OneDrive/Programacion/leontas-9/terminal/image_test/hollow-knight.png"

	img, err := terminal.LoadImage(filepath)
	if err != nil { fmt.Print(err) }
	defer terminal.PutReusableRGBA(img)

	src := terminal.NewImage(img)

	_,err = src.Print()
	if err != nil { fmt.Print(err) }
	// fmt.Print("\n\n\n\033[38;2;255;82;197;48;2;155;106;0mHello, World!\033[0m\n\n\n")

}