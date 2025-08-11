package terminal

import (
	"image"
	"os"
	"time"

	"github.com/eiannone/keyboard"
)

//
func (src *RenderImage) Displacement() error {
	keyboard.Open()
	defer keyboard.Close()

	os.Stdout.Write(alternativeScreen_On)
	defer os.Stdout.Write(alternativeScreen_Off)

	_,err := src.Print()
	if err != nil {return err}

	lastPosition 	:= src.InitialPoint
	lastScreen		:= terminalSize
	
	for {
		_, actualKey, err := keyboard.GetKey()
		if err != nil { return err }

		src.moviment(actualKey)

		if actualKey == keyboard.KeyEsc || actualKey == keyboard.KeyCtrlC {
			os.Stdout.Write(moveToStart)
			os.Stdout.Write(eraseScreen_FromCursor)
			return nil
		}

		actualScreen, err := GetTerminalSize()
		if err != nil { return err }
		
		if !lastPosition.Eq(src.InitialPoint) {
			lastPosition = src.InitialPoint
			src.printNewPosition()
		}
		if !lastScreen.Eq(actualScreen) {
			src.InitialPoint = ClampToPoint(src.InitialPoint, actualScreen)
			lastScreen = actualScreen
			src.printNewPosition()
		}
	}
}

func (src *RenderImage) moviment(actualKey keyboard.Key) {
	switch actualKey {
	case keyboard.KeyArrowUp:
		src.MoveUp(StepsDistance)
	case keyboard.KeyArrowDown:
		src.MoveDown(StepsDistance)
	case keyboard.KeyArrowRight:
		src.MoveRight(StepsDistance)
	case keyboard.KeyArrowLeft:
		src.MoveLeft(StepsDistance)
	}
}

func (src *RenderImage) printNewPosition() error {
	os.Stdout.Write(moveToStart)
	os.Stdout.Write(eraseScreen_FromCursor)

	_, err := src.Print()
	if err != nil { return err }

	return nil
}

var( 
	StepsDistance 	= 2
	StepSpeed 		= 100 * time.Millisecond
)

func (src *RenderImage) MoveUp(step int) {
	src.SetInitialPoint(image.Pt(src.InitialPoint.X, src.InitialPoint.Y- step))
}

func (src *RenderImage) MoveDown(step int) {
	src.SetInitialPoint(image.Pt(src.InitialPoint.X, src.InitialPoint.Y+ step))
}

func (src *RenderImage) MoveLeft(step int) {
	src.SetInitialPoint(image.Pt(src.InitialPoint.X- step, src.InitialPoint.Y))
}

func (src *RenderImage) MoveRight(step int) {
	src.SetInitialPoint(image.Pt(src.InitialPoint.X+ step, src.InitialPoint.Y))
}