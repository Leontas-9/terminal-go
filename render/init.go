package terminal

import (
	"log"
)

// initializa el tamaño del terminal y los pools de imágenes reutilizables
func init() {
	// determina el tamaño actual del terminal
	var err error
	terminalSize, err = GetTerminalPixelSize()
	if err != nil {log.Fatal(err)}

	asignRGBA_Pools()
}
