package ansi

/*
Nota: Algunas secuencias, como guardar y restaurar cursores,
son secuencias privadas y no están estandarizadas.
Si bien algunos emuladores de terminal (es decir, xterm y derivados)
admiten secuencias SCO y DEC, es probable que tengan una funcionalidad diferente.
Por lo tanto, se recomienda utilizar secuencias DEC.
*/

import (
	"fmt"
)

// MoveToStart mueve el cursor a la posición (1,1) en la terminal
func MoveToStart() string {
	return Esc + "H"
}

// MoveUp mueve el cursor hacia arriba un número específico de líneas
func MoveUp(lines int) string {
	return Esc + fmt.Sprintf("%dA", lines)
}

// MoveDown mueve el cursor hacia abajo un número específico de líneas
func MoveDown(lines int) string {
	return Esc + fmt.Sprintf("%dB", lines)
}

// MoveRight mueve el cursor hacia la derecha un número específico de columnas
func MoveRight(columns int) string {
	return Esc + fmt.Sprintf("%dC", columns)
}

// MoveLeft mueve el cursor hacia la izquierda un número específico de columnas
func MoveLeft(columns int) string {
	return Esc + fmt.Sprintf("%dD", columns)
}

func MoveDown_Start(lines int) string {
	return Esc + fmt.Sprintf("%dE", lines)
}

// MoveUp_Start mueve el cursor hacia arriba un número específico de líneas
// y lo posiciona al inicio de la línea.
func MoveTo(column, lines int) string {
	return Esc + fmt.Sprintf("%d;%dH", lines, column)
}

// MoveToColumn mueve el cursor a una columna específica en la línea actual
func MoveToColumn(column int) string {
	return Esc + fmt.Sprintf("%dG", column)
}
