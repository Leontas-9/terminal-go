package terminal

import (
	"bytes"
	"sync"
	"github.com/Leontas-9/terminal-go/ansi"
)

// Precalculo de codigos ANSI
// Estos codigos son utilizados para controlar el cursor, limpiar la pantalla, etc.
// Se utilizan para optimizar el rendimiento y evitar la necesidad de crear nuevos bytes.Buffer cada vez
// que se necesita enviar un comando ANSI al terminal.
var (
	moveToStart = []byte (ansi.MoveToStart())
	eraseScreen_FromCursor = []byte (ansi.EraseScreen_FromCursor())
	alternativeScreen_On = []byte (ansi.AlternativeScreen(true))
	alternativeScreen_Off = []byte (ansi.AlternativeScreen(false))
	resetColor = []byte (ansi.ResetAllColors())
	moveDown = []byte (ansi.MoveDown_Start(1))
	upperBlock = ansi.UpperHalfBlock
	lowerBlock = ansi.LowerHalfBlock
)

// Precalculo de ALPHA_4
// Estos son los valores de ALPHA_4 y ALPHA_1 que se utilizan para optimizar el rendimiento
const (
	ALPHA_4 = ansi.ALPHA_4
	ALPHA_1 = ansi.ALPHA_1
)

// Precalculo de buffers reutilizables
// se utilizan pools para no crear constantes bytes.Buffer
// y evitar la sobrecarga de memoria y mejorar el rendimiento.
var bufferPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

const (
	BPP = 4	// Bytes por Pixel (RGBA)
	PPB = 2 // Pixeles por Bloque, Unicode(inferior  '▄' y superior '▀')
)

// Precañculo de puntos impares
// determina si el alto de la imagen es impar o par
var (
	isYOdd bool
)