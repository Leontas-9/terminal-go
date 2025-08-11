/*
Nota: Al borrar, el cursor no se moverá,
lo que significa que el cursor permanecerá en la última posición en la que se
encontraba antes de borrar la línea. Puede usar después de borrar la línea,
para devolver el cursor al inicio de la línea actual.\r
*/
package ansi

import (
	"bytes"
	"image"
	"strings"
)

// Constante de escape ANSI para secuencias de control
const Esc = "\033["

// EraseScreen_FromCursor borra desde la posición actual del cursor hasta el final de la pantalla
func EraseScreen_FromCursor() string {
	return Esc + "0J"
}

// EraseScreen_ToCursor borra desde el inicio de la pantalla hasta la posición actual del cursor
func EraseScreen_ToCursor() string {
	return Esc + "1J"
}

// EraseLastChar borra el último carácter escrito en la terminal
// moviendo el cursor una posición a la izquierda.
func EraseLastChar() string {
	return "\b " + MoveLeft(1)
}

// EraseBlock borra un bloque horizontal desde la columna 'from' hasta la columna 'to'
// en la fila actual del cursor. Si 'to' es menor que 'from', intercambia los valores.
func EraseBlock(from, to int) string {
	if to < from {
		from, to = to, from
	}
	return MoveToColumn(from) + string(bytes.Repeat([]byte(" "), to - from)) + MoveToColumn(from)
}

// EraseRect borra un rectángulo definido por las coordenadas (x1, y1) y (x2, y2).
// Las coordenadas son inclusivas y basadas en 0. El cursor se mueve al inicio del rectángulo.
func EraseRectangle(rectangle image.Rectangle) string {
	var sb strings.Builder
	minX, maxX := rectangle.Min.X+1, rectangle.Max.X+1
	minY := rectangle.Min.Y+1
	
	sb.WriteString(MoveTo(minX, minY))
	for range rectangle.Dy() {
		sb.WriteString(EraseBlock(minX, maxX) + MoveDown(1))
	}
	sb.WriteString(MoveTo(minX, minY))

	return sb.String()
}

// EraseScreen borra toda la pantalla y mueve el cursor a la posición (1,1)
func EraseScreen() string {
	return Esc + "2J"
}

// EraseLines_Saved borra todas las líneas guardadas en el búfer de desplazamiento
// y mueve el cursor a la posición (1,1)
func EraseLines_Saved() string {
	return Esc + "3J"
}

// EraseLine_FromCursor borra desde la posición actual del cursor hasta el final de la línea
func EraseLine_FromCursor() string {
	return Esc + "0K"
}

// EraseLine_ToCursor borra desde el inicio de la línea hasta la posición actual del cursor
func EraseLine_ToCursor() string {
	return Esc + "1K"
}

// EraseLine borra toda la línea actual donde se encuentra el cursor
func EraseLine() string {
	return Esc + "2K"
}