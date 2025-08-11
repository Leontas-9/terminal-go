package terminal

import (
	"image"
	"log"
	"sync"

	"golang.org/x/sys/windows"
)

// Precalculo del tamaño del terminal para pixeles
// cuenta los pixeles que se pueden usar dentro de un bloque unicode 
// bloque superior ('▀') y bloque inferior ('▄')
var terminalSize image.Point

// initializa el tamaño del terminal y los pools de imágenes reutilizables
func init() {
	// determina el tamaño actual del terminal
	var err error
	terminalSize, err = GetTerminalPixelSize()
	if err != nil {log.Fatal(err)}

	asignRGBA_Pools()
}


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