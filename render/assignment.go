package terminal

import (
	"image"
	"golang.org/x/image/draw"
)

type RenderImage struct {
	// Imagen original (sin cambios de ajuste de bordes)
	Image 		*image.RGBA

	// Bordes en los que se desea ajustar la imagen dentro del terminal
	Margin 		image.Rectangle

	// Punto inicial en donde se posicionara la imagen (el punto superior izquierdo)
	InitialPoint	image.Point

	// Tipo de interpolador
	Interpolator draw.Interpolator

	// Configuracion UI
	// Permite cambiar la configuracion de la imagen
	// como el cursor, pantalla alternativa, borrar pantalla, auto ajuste de imagen
	opts        UI_Settings
}

type UI_Settings struct {
	// Al final del renderizado muestra el cursor?
	ShowCursor 			bool
	
	// Cambia a la pantalla alternativa en el renderizado?
	// Si es verdadero, la imagen final no se vera reflejada en la pantalla principal
	// y se conservara todo dato impreso anteriormente a la imagen
	AlternativeScreen	bool

	// Borrar la imagen al finalizar?
	// Si es verdadero, no conservara los datos anteriores
	EraseScreen			bool
	
	// Permite el auto ajuste de la imagen al final del renderizado?
	// Si es verdadero, la imagen se ajustara a los bordes de la terminal
	Auto_Wrap			bool
}

// NewRenderImage crea una nueva imagen renderizable con los parametros especificados
func newRenderImage(img *image.RGBA, bounds image.Rectangle, initialPoint image.Point, interpolator draw.Interpolator, opts UI_Settings) *RenderImage {
	return &RenderImage{
		Image:        img,
		Margin:       bounds,
		InitialPoint: initialPoint,
		Interpolator: interpolator,
		opts:         opts,
	}
}

// NewImage crea una nueva imagen renderizable con el tamaño de la imagen original
// y los valores por defecto de UI_Settings
func NewImage(img *image.RGBA) *RenderImage {
	return newRenderImage(
		img,
		DefaultSize(img.Rect),
		image.Point{},
		draw.NearestNeighbor,
		*UI_Settings{}.Default(),
	)
}

// Default devuelve una configuracion por defecto de UI_Settings
// con los valores establecidos en el constructor
func (UI_Settings) Default() *UI_Settings {
	return &UI_Settings{
			ShowCursor: true,
			AlternativeScreen: false,
			EraseScreen: false,
			Auto_Wrap: true,
		}
}

// NewCustomImage crea una nueva imagen renderizable con los parametros especificados
func NewCustomImage(img *image.RGBA, bounds image.Rectangle, initialPoint image.Point, interpolator draw.Interpolator, opts *UI_Settings) *RenderImage {
	return newRenderImage(img, bounds, initialPoint, interpolator, *opts)
}

// GetUI_Settings cambia la configuracion de UI_Settings de la imagen
func (img *RenderImage) SetUI_Settings(new *UI_Settings) {
	img.opts = *new
}

// GetImage cambia los margenes en el que se imprimira la imagen dentreo del terminal
func (img *RenderImage) SetMargins(new image.Rectangle) {
	img.Margin = new
}

// SetInterpolator cambia el interpolador de la imagen
// Permite cambiar el tipo de interpolacion que se usara al renderizar la imagen
func (img *RenderImage) SetInterpolator(new draw.Interpolator) {
	img.Interpolator = new
}

// SetImage cambia la imagen que se renderizara
func (img *RenderImage) SetImage(new *image.RGBA) {
	img.Image = new
}

// SetInitialPoint cambia el punto inicial en donde se posicionara la imagen
// Este punto es el punto superior izquierdo de la imagen dentro del terminal
func (img *RenderImage) SetInitialPoint(new image.Point) {
	img.InitialPoint = new
}

// DefaultSize devuelve el tamaño por defecto de la imagen
// Si el area de la imagen es menor o igual al area del terminal, devuelve el tamaño
// de la imagen, de lo contrario devuelve el tamaño del terminal
// Esto se usa para ajustar la imagen al tamaño del terminal
func DefaultSize(src image.Rectangle) image.Rectangle {
	terminalRect := image.Rect(0, 0, terminalSize.X, terminalSize.Y)
	
	if GetAreaRect(src) <= GetAreaRect(terminalRect){
			return src
	} else {return terminalRect}
}

// GetAreaRect calcula el area de un rectangulo
// Se usa para comparar el area de la imagen con el area del terminal
func GetAreaRect(rect image.Rectangle) int {
	return rect.Dx()*rect.Dy()
}