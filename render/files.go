
package terminal

import (
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
)

//LoadImage carga una imagen desde un archivo y la convierte a formato RGBA.
// Si el archivo no se puede abrir o el formato no es compatible, devuelve un error.
// Formatos soportados: JPEG, PNG, BMP, TIFF, WebP y GIF (solo la primera imagen del GIF).
func LoadImage(filepath string) (*image.RGBA, error) {
	file, err := OpenFile(filepath)
	if err != nil { return nil, err }
	defer file.Close()

	img, err := DecodeImage(file, filepath)
	if err != nil { return nil, fmt.Errorf("decode: %v", err) }

	dst, err := ConvertToRGBA(img)
	if err != nil { return nil, fmt.Errorf("convertToRGBA: %v", err) }

	return dst, nil
}

// DecodeImage decodifica una imagen desde un lector de archivos
// según su extensión. Soporta formatos comunes como JPEG, PNG, BMP, TIFF, WebP y GIF.
// (En el case de GIF este codifica unuicamente la primera imagen del GIF). 
// Si el formato no es compatible, intenta decodificar como imagen genérica.
func DecodeImage(file io.Reader, fileName string) (img  image.Image, err error) {
	extension := filepath.Ext(fileName)

	switch extension {
	case ".jpeg": 	img, err = jpeg.Decode(file)
	case ".jpg": 	img, err = jpeg.Decode(file)
	case ".png":	img, err = png.Decode(file)
	case ".bmp":	img, err = bmp.Decode(file)
	case ".tiff":	img, err = tiff.Decode(file)
	case ".webp":	img, err = webp.Decode(file)
	case ".gif":	img, err = gif.Decode(file)
	default:		img,_,err = image.Decode(file)
		if err != nil { return nil, fmt.Errorf("formato no soportado: %v", extension) }
	}
	if err != nil { return nil, err }

	return
}

// OpenFile abre un archivo en la ruta especificada y devuelve un puntero al archivo.
// Si hay un error al cambiar el directorio o al abrir el archivo, devuelve un error
func OpenFile(filepath string) (*os.File, error) {
	file, err := os.Open(filepath)
	if err != nil { return nil, err }

	return file, nil
}

// ConvertToRGBA convierte una imagen a formato RGBA.
// Si la imagen ya es de tipo RGBA, la devuelve directamente.
// esto evita la necesidad de crear una nueva imagen innecesariamente.
func ConvertToRGBA(src image.Image) (*image.RGBA, error) {
	dst ,ok := src.(*image.RGBA)
	if ok { return dst, nil	}

	dst = GetReusableRGBA(src.Bounds())
	
	draw.Draw(dst, src.Bounds(), src, image.Point{}, draw.Src)

	return dst, nil
}