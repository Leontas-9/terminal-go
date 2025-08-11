package terminal

import (
	"image"
	"sync"

	"golang.org/x/image/draw"
	"golang.org/x/sys/windows"
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



// Precalculo del tamaño del terminal para pixeles
// cuenta los pixeles que se pueden usar dentro de un bloque unicode
// bloque superior ('▀') y bloque inferior ('▄')
var terminalSize image.Point


// Obtiene el tamaño en terminal para pixeles
// cuenta los pixeles que se pueden usar dentro de un bloque unicode 
// bloque superior ('▀') y bloque inferior ('▄')
func GetTerminalPixelSize() (terminalSize image.Point, err error) {
	size, err := GetTerminalSize()

	terminalSize = image.Pt(size.X, size.Y*2)
	return terminalSize, err
}


// Retorna el punto maximo del tamaño actual del terminal
// La cantidad de caracteres de ancho y largo que se pueden utilizar
func GetTerminalSize() (size image.Point, err error) {
	var info windows.ConsoleScreenBufferInfo
	handle := windows.Handle(windows.Stdout)
	err = windows.GetConsoleScreenBufferInfo(handle, &info)
	if err != nil {return image.Point{}, err}
	
	ancho := int(info.Window.Right - info.Window.Left + 1)
	alto := int(info.Window.Bottom - info.Window.Top + 1)


	size = image.Pt(ancho, alto)
	return size, nil
}

// Obtiene una imagen de tipo RGBA reusable
// Es recomendable que al finalizar su uso se devuelva al pool 
// usando PutReusableRGBA para evitar fugas de memoria.
func GetReusableRGBA(bounds image.Rectangle) *image.RGBA {
	pool := getRGBA_Pool(bounds)
	img := pool.Get().(*image.RGBA)
	img.Rect = bounds

	return img
}

// Resoluciones predefinidas para las imágenes reutilizables
// Estas resoluciones son utilizadas para optimizar el uso de memoria y rendimiento
// al reutilizar imágenes de diferentes tamaños según las necesidades del terminal.
var resolutions = []image.Point{
	{141, 58},			// valor por default del terminal
	{160, 120},
	{320, 240},
	{640, 480},
	{1280, 720},
	{1920, 1080},
	{2560, 1440},
	{3840, 2160},
}

// rgbaPools almacena pools de imágenes RGBA reutilizables
var rgbaPools = make(map[image.Point]*sync.Pool)

// asignRGBA_Pools inicializa los pools de imágenes RGBA reutilizables
func asignRGBA_Pools() {
	for _, resolution := range resolutions {
		rgbaPools[resolution] = &sync.Pool{
			New: func() any { return image.NewRGBA(image.Rect(0,0, resolution.X, resolution.Y)) },
		}
	}
}

// getRGBA_Pool obtiene el pool de imágenes RGBA reutilizables adecuado
// según el tamaño de los bounds proporcionados.
// Si no hay un pool adecuado, devuelve un nuevo pool con el tamaño máximo.
// Esto permite reutilizar imágenes de diferentes tamaños sin necesidad de crear nuevas instancias.
func getRGBA_Pool(bounds image.Rectangle) *sync.Pool {
	max := bounds.Max
	for _, resolution := range resolutions {
		if resolution.X >= max.X && resolution.Y >= max.Y {
			return rgbaPools[resolution] 
		}
	}

	return &sync.Pool{
		New: func() any { return image.NewRGBA(image.Rect(0,0, max.X, max.Y)) },
	}
}

// PutReusableRGBA limpia y devuelve una imagen RGBA al pool de reutilización.
// Esto permite que la imagen sea reutilizada en futuras operaciones,
func PutReusableRGBA(img *image.RGBA) {
	if img == nil { return }

	pixels := img.Pix
	for i := range pixels {
		pixels[i] = 0
	}

	pool, ok := rgbaPools[img.Rect.Max]
	if ok { pool.Put(img) }
}