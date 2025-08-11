package ansi

/*
Significado específico:
	38: Aplica color al texto (primer plano).
	48: Aplica color al fondo.
	2: Indica que se usará color en formato RGB de 24 bits (true color).
	{r}, {g}, {b}: Son los valores de Red, Green y Blue de 0 a 255).
	m: Final de una secuencia ANSI significa "modo" o "modo gráfico", activa el modo SGR

puede poner texto rosado sobre un fondo marrón usando:
"\033[38;2;255;82;197;48;2;155;106;0mHello"
*/

import (
	"image/color"
	"math"
	"unicode/utf8"
	"unsafe"
)

// init precomputa la tabla de dígitos para optimizar la conversión a ANSI
func init() {
	// Prellenar la tabla de dígitos
	for n := range 256 {
		digitLookup[n] = byteToDigits(byte(n))
	}
}

// Cache de los codigos ANSI para colorear texto y fondo 
// precomputado en init() para evitar cálculos repetidos
type AnsiCache struct {

	// RGB contiene los códigos ANSI para colores RGB
	// Estructura: R -> G -> B -> Código ANSI
	RGB		[][][][]byte

	// PrefixCode contiene los prefijos ANSI para texto y fondo
	// Estructura: [0] -> Texto, [1] -> Fondo
	PrefixCode 	[][]byte

	// SufixCode es el código ANSI que finaliza la secuencia
	// Estructura: 'm'
	SufixCode 		byte
}

// digitLookup es una tabla optimizada para conversión de 3 dígitos (8 bits)
var digitLookup = [256][]byte{} // Precomputado en init()

// byteToDigits convierte un byte (0-255) en su representación ASCII de dígitos
func byteToDigits(n byte) []byte {
	if n < 10	{ return []byte{'0' + n} }
	if n < 100 	{ return []byte{'0' + n/10, '0' + n%10} }
	return 				 []byte{'0'+ n/100, '0' + (n/10)%10, '0' + n%10}
}

// ANSI escape codes
var (
    ansiPrefixFg  = []byte("\033[38;2;")
    ansiPrefixBg  = []byte("\033[48;2;")
)

// GetANSI_Color genera el código ANSI para un color RGB específico.
// El parámetro isText determina si el color es para texto (true) o fondo (false)
func GetANSI_Color(buf *[]byte, r, g, b uint8, isText bool) {
	code := ansiPrefixBg
	if isText {code = ansiPrefixFg}

	*buf = append(*buf, code...)

	// Acceso directo a dígitos precomputados
	d := digitLookup[r]
	AppendBytes(buf, ';', d)

	d = digitLookup[g]
	AppendBytes(buf, ';', d)

	d = digitLookup[b]
	AppendBytes(buf, 'm', d)
}

// GetANSI_DoubleColor genera el código ANSI para colores RGB específicos
// tanto para texto (rF,gF,bF) como para fondo (rB,gB,bB).
// se usa con el formato: "\033[38;2;R;G;Bm\033[48;2;R;G;Bm"
// y para optimizar el uso de memoria conmenores bytes
func GetANSI_DoubleColor(buf *[]byte, rF,gF,bF,rB,gB,bB byte) {
	// Acceso directo a dígitos precomputados
	*buf = append(*buf, ansiPrefixFg...)

	d := digitLookup[rF]
	AppendBytes(buf, ';', d)
	
	d = digitLookup[gF]
	AppendBytes(buf, ';', d)
	
	d = digitLookup[bF]
	AppendBytes(buf, ';', d)

	*buf = append(*buf, "48;2;"...)

	d = digitLookup[rB]
	AppendBytes(buf, ';', d)
	
	d = digitLookup[gB]
	AppendBytes(buf, ';', d)
	
	d = digitLookup[bB]
	AppendBytes(buf, 'm', d)
}

// AppendBytes agrega una secuencia de bytes seguida de un carácter específico al búfer.
func AppendBytes(buffer *[]byte, char byte, items []byte) {
	*buffer = append(*buffer, items...)
	*buffer = append(*buffer, char)
}

// DefaultColor es un color transparente (usado para resetear colores)
var DefaultColor = color.RGBA{A: 0}

var (
	UpperHalfBlock = '▀' // Bloque superior Unicode '▀'
	LowerHalfBlock = '▄' // Bloque inferior Unicode '▄'
)

// Niveles de transparencia para bloques Unicode
const (
	ALPHA_1 uint8 = (255 * 1) / 5 // ' '
	ALPHA_2 uint8 = (255 * 2) / 5 // '░'
	ALPHA_3 uint8 = (255 * 3) / 5 // '▒'
	ALPHA_4 uint8 = (255 * 4) / 5 // '▓'
	ALPHA_5 uint8 = (255 * 5) / 5 // '█'
)

// Escala de bloques Unicode para diferentes niveles de transparencia
var transparent_BlockScale = []rune{' ', '░', '▒', '▓', '█'}

// Reset_ColorForeground devuelve el código ANSI para resetear el color del texto (foreground)
func Reset_ColorForeground() []byte {
	return []byte("\033[39m")
}

// Reset_ColorBackground devuelve el código ANSI para resetear el color del fondo (background)
func Reset_ColorBackground() []byte {
	return []byte("\033[49m")
}

// ResetAllColors devuelve el código ANSI para resetear ambos colores (texto y fondo)
func ResetAllColors() []byte {
	return []byte("\033[39;49m")
}

// buffer para almacenar códigos ANSI generados
// especialmente paralos colores ANSI
// utiliza bigANSI_Code para ambos colores
// y littleANSI_Code para un solo color
var(
	bigANSI_Code [36]byte
	littleANSI_Code [19]byte
)

// PaintBase genera el código ANSI para aplicar colores de texto y fondo
// Si ambos colores son opacos, usa el código extendido para ambos
func PaintBase(fgColor, bgColor color.RGBA) (block []byte) {
	var buf []byte

	bgAlpha := bgColor.A > ALPHA_1
	fgAlpha := fgColor.A > ALPHA_1

	if bgAlpha && fgAlpha {
		buf = bigANSI_Code[:0]
		GetANSI_DoubleColor(
			&buf, fgColor.R, fgColor.G, fgColor.B, bgColor.R, bgColor.G, bgColor.B)
		
		return buf
	}

	buf = littleANSI_Code[:0]
	if bgAlpha {
		GetANSI_Color(&buf, bgColor.R, bgColor.G, bgColor.B, false)
	} else if fgAlpha {
		GetANSI_Color(&buf, fgColor.R, fgColor.G, fgColor.B, true)
	}

	return buf
}

// ansiBlock es un búfer reutilizable para construir cadenas ANSI 
// y evitar asignaciones repetidas
var ansiBlock 	[47]byte                                                              

// PaintRune genera el código ANSI para un solo carácter con colores específicos.
// Si resetColor es true, agrega el código para resetear los colores al final.
func PaintRune(character rune, textColor, backgroundColor color.RGBA, resetColor bool) (block []byte) {
	buf := ansiBlock[:0]

	// Conversión directa del rune
	buf = append(buf, PaintBase(textColor, backgroundColor)...)
	utf8.EncodeRune(buf[len(buf):len(buf)+3], character)
	buf = buf[:len(buf)+utf8.RuneLen(character)]

	if resetColor {
		buf = append(buf, ResetAllColors()...)
	}
	return buf
}

// PaintString genera el código ANSI para una cadena de texto con colores específicos.
// Si resetColor es true, agrega el código para resetear los colores al final.
func PaintString(text string, textColor, backgroundColor color.RGBA, resetColor bool) (block string) {
	buf := ansiBlock[:0]
	buf = append(buf, PaintBase(textColor, backgroundColor)...)

	for _, r := range text {
        var charBuf [4]byte // Máximo tamaño de un rune en UTF-8
        n := utf8.EncodeRune(charBuf[:], r)
        buf = append(buf, charBuf[:n]...)
    }

	if resetColor {
		buf = append(buf, ResetAllColors()...)
	}
	return unsafe.String(unsafe.SliceData(buf), len(buf))
}

// BlockShade devuelve un carácter Unicode que representa el nivel de 
// transparencia del color dado.
// Usa una escala de bloques Unicode para diferentes niveles de transparencia.
func BlockShade(color color.RGBA) (block rune) {
	transparencyLevel := int(math.Round((float64(color.A) / 255.0) * float64(len(transparent_BlockScale)-1)))
	return transparent_BlockScale[transparencyLevel]
}

// AverageColor calcula el color promedio entre dos colores RGBA.(incluyendo alfa)
func AverageColor(color1, color2 color.RGBA) (color3 color.RGBA) {
	R1, G1, B1, A1 := color1.RGBA()
	R2, G2, B2, A2 := color2.RGBA()
	return color.RGBA{
		R: uint8(R1>>8 + R2>>8) / 2,
		G: uint8(G1>>8 + G2>>8) / 2,
		B: uint8(B1>>8 + B2>>8) / 2,
		A: uint8(A1>>8 + A2>>8) / 2,
	}
}

// AverageAlpha calcula el valor alfa promedio entre dos colores RGBA. (unicamente alfa)
func AverageAlpha(color1, color2 color.RGBA) (color3 color.RGBA) {
	_, _, _, A1 := color1.RGBA()
	_, _, _, A2 := color2.RGBA()
	return color.RGBA{
		A: uint8(A1>>8 + A2>>8) / 2,
	}
}