package terminal

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"sync"

	"github.com/Leontas-9/terminal-go/ansi"

	"golang.org/x/image/draw"
	"golang.org/x/sys/windows"
)

// DefaultTerminalSize es el tamaño por defecto del terminal, usado para ajustar los bordes de la imagen
var DefaultTerminalSize = image.Point{ X: 141, Y: 58 }

// Print imprime directamente la imagenen el terminal
func (src *RenderImage) Print() (lenght int, err error) {
	bytes,_,err := src.GetPNG()
	if err != nil { return 0, err }

	return os.Stdout.Write(bytes)
}

// GetPNG obtiene la imagen en formato ANSI/ASCII, la imagen reajustada y error si que hay alguno
func (src *RenderImage) GetPNG() (ASCII_Image []byte, image image.Image, err error) {
	dst := *src
	err = dst.validateInputs()
	if err != nil {return nil, nil, err}

	dst.Margin, err = dst.AdjustLimitsToTerminal()
	if err != nil {return nil, nil, err}

	dst.Image, err = dst.AdjustImage()
	if err != nil {return nil, nil, err}

	dst.InitialPoint = ClampToPoint(dst.InitialPoint, terminalSize.Sub(dst.Image.Rect.Max))

	ASCII_Image, err = dst.RenderImage()
	if err != nil {return nil, nil, err}

	return ASCII_Image, dst.Image, err
}

// validateInputs Valida los inputs de RenderImage
func (src *RenderImage) validateInputs() error {
	if src.Interpolator == nil  { return errors.New("interpolator cannot be nil") }
	if src.Image == nil 		{ return errors.New("image cannot be nil") }
	if src.Margin.Dx() == 0 	{ return errors.New("width cannot be 0") }
	if src.Margin.Dy() == 0 	{ return errors.New("height cannot be 0") }

	return nil
}

// AdjustLimitsToTerminal Ajusta los bordes de imagen a los bordes del terminal
func (src *RenderImage) AdjustLimitsToTerminal() (newBounds image.Rectangle, err error) {
	terminalSize, err := GetTerminalPixelSize()
	if err != nil { return image.Rectangle{}, err }

	newBounds	= src.clampToBounds(image.Rect(0,0, terminalSize.X, terminalSize.Y))
	return newBounds, nil
}


// ClampToPoint Ajusta un punto para que este dentro de otro sin ser mayor 
func ClampToPoint(lastPoint, actualPoint image.Point) image.Point {
	return image.Pt(
		Clamp(lastPoint.X, 0, actualPoint.X), 
		Clamp(lastPoint.Y, 0, actualPoint.Y),
	)
}

// clampToBounds Ajusta los bordes del rectagnulo dentro de los bordes delimitados
func (src *RenderImage) clampToBounds(bounds image.Rectangle) (image.Rectangle) {
	return image.Rect(
		Clamp(src.Margin.Min.X, 0, bounds.Max.X),
		Clamp(src.Margin.Min.Y, 0, bounds.Max.Y),
		Clamp(src.Margin.Max.X, 0, bounds.Max.X),
		Clamp(src.Margin.Max.Y, 0, bounds.Max.Y),
	)
}

// Clamp Clava el valor dentro de los margenes
func Clamp(value, min, max int) int {
	if value < min        { return min
	} else if value > max { return max }

	return value
}

// AdjustImage Ajusta la imagen a una escala acorde a los bordes
func (src *RenderImage) AdjustImage() (*image.RGBA, error) {
    scale := src.CalculateScale()
    dst := src.ScaleImage(scale)

	return dst, nil
}

// CalculateScale Calcula la escala minima para estar dentro de los limites
func (src *RenderImage) CalculateScale() (scale float64) {
	srcWidth	:= float64(src.Image.Bounds().Dx())
	srcHeight	:= float64(src.Image.Bounds().Dy())

	scale = math.Min(
		float64(src.Margin.Dx())/srcWidth,
        float64(src.Margin.Dy())/srcHeight,
	)

	return scale
}

//  ScaleImage  hace el escalamiento de la imagen
func (src *RenderImage) ScaleImage(scale float64) (*image.RGBA) {
	if scale != 1.0 {
		dst := GetReusableRGBA(image.Rect(0,0, 
			int(math.Round(float64(src.Image.Rect.Dx())* scale)), 
			int(math.Round(float64(src.Image.Rect.Dy())* scale)),
		))

		src.Interpolator.Scale(
			dst, dst.Bounds(), src.Image, src.Image.Bounds(),draw.Over, nil)

		return dst
	}

	return src.Image
}

// RenderImage realiza la transformacion de imagenes a texto ([]byte) Unicode/ANSI
func (src *RenderImage) RenderImage() (ASCII_Image []byte, err error) {
	blocks := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(blocks)
	blocks.Reset()

	err = src.initializeRender(blocks)
	if  err != nil { return nil, err }
	
	err = src.renderBlocks(blocks)
	if  err != nil { return nil, err }

	err = src.finalizeRender(blocks)
	if  err != nil { return nil, err }

	ASCII_Image = blocks.Bytes()
	if len(ASCII_Image) == 0 {
		return nil, fmt.Errorf("renderizado vacio, verifique la imagen y los bordes")
	}

	return ASCII_Image, err
}

// validateUI_Settings guarda en un buffer las configuraciones de UI si es true, 
// si es falso este dejara las configuraciones por defecto
func (src *UI_Settings) validateUI_Settings(buf *bytes.Buffer, isDefaultSettings bool) (err error) {
isCursor := src.ShowCursor
	isScreen := src.AlternativeScreen
	isErase  := src.EraseScreen
	isWrap   := src.Auto_Wrap

	if isDefaultSettings {
		settings := UI_Settings{}.Default()
		isCursor  = settings.ShowCursor
		isScreen  = settings.AlternativeScreen
		isErase   = settings.EraseScreen
		isWrap	  = settings.Auto_Wrap
	}

	buf.Grow(6+ 8+ 5)		// espacio para 3 codigos ANSI 6 -> ShowCursor, 8 -> AlternativeScreen, 5 -> Auto_Wrap
	_,err = buf.WriteString(ansi.ShowCursor(isCursor))
	if err != nil { return err }
	
	_,err = buf.WriteString(ansi.AlternativeScreen(isScreen))
	if err != nil { return err }


	_,err = buf.WriteString(ansi.Auto_Wrap(isWrap))
	if err != nil { return err }
	
	if isErase {
		buf.Grow(3+ 4)		// espacio para 2 codigos ANSI 3 -> MoveToStart, 4 -> EraseScreen_FromCursor
		_,err = buf.Write(moveToStart)
		if err != nil { return err }
	
		_,err = buf.Write(eraseScreen_FromCursor)
		if err != nil { return err }
	}

	return
}

// Renderiza los bloques dentro de una imagen a un formato Unicode/ANSI y los guarda en un buffer
func (src *RenderImage) renderBlocks(buf *bytes.Buffer) (err error) {
	isYOdd = src.isYOdd()
	for y := 0; y < src.Image.Rect.Dy(); y += PPB {
		line := y * src.Image.Stride

		for x := 0; x < src.Image.Stride; x += BPP {
			err = src.renderBlock(buf, line + x)
			if err != nil { return err }
		}
		err = src.endLine(buf)
		if err != nil { return err }

		if isYOdd && y == 0 { y-- }
	}
	
	return
}


// Renderiza un unico bloque Unicode
func (src *RenderImage) renderBlock(buf *bytes.Buffer, index int) (err error) {
	src.validateIndex(buf, index)

	fgColor, bgColor := src.getPixels(index) // frente y fondo
	
	block := src.determineBlockType(index, fgColor, bgColor)
	if block == 0 { return }

	if src.sameColor(buf, index, block, fgColor, bgColor) { return }

	buf.Grow(39)
	buf.Write(ansi.PaintRune(block, fgColor, bgColor, false))
	return
}

// sameColor verifica si el bloque actual es del mismo color que el bloque superior e inferior
// si es asi, escribe un espacio en blanco en el buffer y retorna true, de lo contrario,
// escribe el bloque actual en el buffer y retorna false
// esto con el objetivo de optimizar el renderizado y evitar bloques o colores repetidos
func (src *RenderImage) sameColor(blockBuf *bytes.Buffer, index int, block rune, fgColor, bgColor color.RGBA) bool {
	lowerIndex := index + src.Image.Stride
	isX_0 := src.isX_0(index)
	var byteBuf []byte 

	if !isX_0 {
		var sameUpper bool
		var sameLower bool
		var sameBlock bool

		// Verifica que el indice inferior este dentro del rango
		// y que el indice superior no sea el primer pixel de la fila
		if lowerIndex+3 < len(src.Image.Pix) && index-4 >= 0 {
			sameUpper	= src.isSameColor(index, index - 4)
			sameLower 	= src.isSameColor(lowerIndex, lowerIndex - 4)
			sameBlock	= src.isSameColor(index, lowerIndex)
		}

		if sameUpper && sameLower {
			if sameBlock {
				blockBuf.Grow(1)
				blockBuf.WriteRune(' ')

				return true
			}
			blockBuf.Grow(3)
			blockBuf.WriteRune(block)
			
			return true

		} else if sameUpper {
			blockBuf.Grow(22)
			ansi.GetANSI_Color(&byteBuf, bgColor.R, bgColor.G, bgColor.B, false)
			blockBuf.Write(byteBuf)
			blockBuf.WriteRune(block)
			
			return true
		
		} else if sameLower {
			blockBuf.Grow(22)
			ansi.GetANSI_Color(&byteBuf, fgColor.R, fgColor.G, fgColor.B, true)
			blockBuf.Write(byteBuf)
			blockBuf.WriteRune(block)
	
			return true
		}
	}

	return false
}

// isX_0 verifica si el indice es el primer pixel de la fila
func (src *RenderImage)  isX_0(index int) bool {
	return index % src.Image.Stride == 0
}

// isSameColor verifica si dos colores son iguales, teniendo en cuenta el alpha, y los valores RGB
// si los colores son iguales, retorna true, de lo contrario, retorna false
func (src *RenderImage) isSameColor(indexColor1, indexColor2 int) bool {
	 if indexColor1 < 0 ||
	 	indexColor2 < 0 ||
        indexColor1+3 >= len(src.Image.Pix) ||
		indexColor2+3 >= len(src.Image.Pix) {

        return false
    }

	Pixel1 := src.getPixel(indexColor1)
	Pixel2 := src.getPixel(indexColor2)

	return Pixel1 == Pixel2
}


// validateIndex Valida que el indice este dentro del rango
func (src *RenderImage) validateIndex(buf *bytes.Buffer, index int) {
	if index < 0 || index+3 >= len(src.Image.Pix) {
        buf.Write(ansi.PaintRune(' ', color.RGBA{}, color.RGBA{}, false))
    }
}

// determineBlockType Determina el tipo de bloque a usar
func (src *RenderImage) determineBlockType(index int, foreground, background color.RGBA) rune {
    // Caso 1: Transparencia en ambos píxeles
    if foreground.A < ALPHA_4 && background.A < ALPHA_4 {
		if foreground.A < ALPHA_1 && background.A < ALPHA_1 { return 0 }

        return ansi.BlockShade(ansi.AverageAlpha(foreground, background))
    }

    // Caso 2: Primera línea impar
    if isYOdd && isFirstRow(index, src.Image.Stride) { return lowerBlock }

    // Caso por defecto
    return upperBlock
}

// isFirstRow verifica si el indice es el primer pixel de la fila
// esto es necesario para evitar que se dibuje un bloque superior en la primera fila
func isFirstRow(index, width int) bool {
	return index < width
}

// isYOdd verifica si la fila es impar, 
// esto es necesario para evitar que se dibuje un bloque superior en la primera fila
func (src *RenderImage) isYOdd() bool {
	return src.InitialPoint.Y & 1 == 1
}

// getPixels Obtiene el color del pixel superior e inferior, segun el indice del pixel superior 
// y retorna el color del pixel superior y el color del pixel inferior
// si el pixel inferior no existe, retorna un color transparente
func (src *RenderImage) getPixels(upperIndex int) (foreground, background color.RGBA) {
    foreground = src.getPixel(upperIndex)
	background = color.RGBA{}
	
    lowerIndex := upperIndex + src.Image.Stride
    if lowerIndex+3 >= len(src.Image.Pix) {return } 					// No hay píxel de fondo

	if isFirstRow(upperIndex, src.Image.Stride) && isYOdd { return }
    background = src.getPixel(lowerIndex)

    return foreground, background
}

// getPixel Obtiene el color del pixel, segun el indice (RGBA)
func (src *RenderImage) getPixel(index int) (pixel color.RGBA) {
	return color.RGBA{
        R: src.Image.Pix[index],
        G: src.Image.Pix[index+1],
        B: src.Image.Pix[index+2],
        A: src.Image.Pix[index+3],
    }
}

// endLine Coloca en el buffer codigoANSI para retornar el cursor al inicio de los bordes
// y bajar una linea
func (src *RenderImage) endLine(buf *bytes.Buffer) (err error) {
	_,err = buf.Write(resetColor)
	if err != nil { return err }

	_,err = buf.Write(moveDown)
	if err != nil { return err }
	
	_,err = buf.WriteString(ansi.MoveToColumn(src.InitialPoint.X + 1))
	if err != nil { return err }

	return
}

// initializeRender Coloca el cursor en la posision inicial para imprimir la imagen
func (src *RenderImage) initializeRender(buf *bytes.Buffer) (err error) {
	startCol, startRow := src.calculateStartPosition()

	err = src.opts.validateUI_Settings(buf, false); 
	if err != nil {	return err }

	_, err = buf.WriteString(ansi.MoveTo(startCol, startRow)); 
	if err != nil {	return err }

	return
}

// finalizeRender Finaliza el renderizando de imagen, volviendo a los valores default
func (src *RenderImage) finalizeRender(buf *bytes.Buffer) (err error) {
	finalCol, finalRow := src.calculateFinalPosition()

	err = src.opts.validateUI_Settings(buf, true)
	if err != nil { return err }

	err = src.finalPosition(buf, finalCol, finalRow)
	if err != nil { return err }
	return 
}

// calculateStartPosition Calcula la posicion inicial del cursor
// teniendo en cuenta el punto inicial y el tamaño de la imagen
func (src *RenderImage) calculateStartPosition() (col, row int) {
	starterCol := src.InitialPoint.X + 1
    starterRow := (src.InitialPoint.Y / 2) + 1

    return starterCol, starterRow
}

// calculateFinalPosition Calcula la posicion final del cursor
func (src *RenderImage) calculateFinalPosition() (col, row int) {
	finalCol := src.Image.Rect.Dx() + src.InitialPoint.X + 1
    finalRow := int(math.Ceil(float64(src.Image.Rect.Dy()) / 2.0) + float64(src.InitialPoint.Y) / 2)

    return finalCol, finalRow
}

// finalPosition Posiciona el cursor en la ultima columna de la ultima fila
func (src *RenderImage) finalPosition(buf *bytes.Buffer, finalCol, finalRow int) (err error) {
	_, err = buf.Write(resetColor)
	if err != nil { return err }

	_, err = buf.WriteString(ansi.MoveTo(finalCol, finalRow))
	if err != nil { return err }
	
	if src.Image.Rect.Max.X >= terminalSize.X {
		_, err = buf.Write(moveDown)
		if err != nil { return err }	
	}

	return
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