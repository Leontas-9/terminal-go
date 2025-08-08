# 💻 terminal-go
Proporciona un conjunto de herramientas para el manejo avanzado de códigos de escape ANSI y la renderización de imágenes directamente en la terminal, facilitando el desarrollo de aplicaciones de línea de comandos interactivas y visualmente atractivas en Go.

### 🖼️: Ejemplo de Imagen:

| Imagen original            | Imagen Renderizada         |
|----------------------------|----------------------------|
| <img width="315" height="216" alt="paisaje" src="https://github.com/user-attachments/assets/d5c13bf8-26db-46a9-9992-b51831bdb393" />|<img width="1360" height="654" alt="paisaje renderizado" src="https://github.com/user-attachments/assets/7ab8c030-fc71-4f62-8a00-6cfcad4b8fa8" />|

----

## Instalacion 

Instale y actualice este paquete go con `go get -u github.com/Leontas-9/terminal-go.git`

----

## Aplicacion

El proyecto incluye diversas funciones que permiten el manejo de codigos ANSI Escape los cuales facilitan la interaccion con el __CLI__ del desarrollador,    
tales como:

``` go
func MoveDown_Start(lines int) string {
	return Esc + fmt.Sprintf("%dE", lines)
}
```

Que permite mover el cursor del terminal hacia abajo una determinada cantidad de lineas y al inicio de la fila.


__Ejemplo simple para renderizar una imagen al terminal:__ `LoadImage()` `PutReusableRGBA`

```go
package main

import (
	"github.com/Leontas-9/terminal/render"
	"fmt"
)

func main() {
	filepath := "ruta/a/tu/imagen.png"

	imgRGBA, err := terminal.LoadImage(filepath)
	if err != nil { fmt.Print(err) }
	defer terminal.PutReusableRGBA(imgRGBA)

	src := terminal.NewImage(img)

	_,err = src.Print()
	if err != nil { fmt.Print(err) }
}
```

## 🎯 Objetivo

- Permitir al usuario un mayor manejo y facilidad de las herraminetas que ofrecen los [__codigos de escape ANSI__](https://gist.github.com/ConnerWill/d4b6c776b509add763e17f9f113fd25b)
- Creacion de imagenes directamente en el terminal.

## ⚙️ ¿Como Funciona?

#### 🎨 Comandos de Entindado
La renderizacion de imagenes dentro del terminal se realiza con una combinacion de comandos, que permiten manejar el color con el que se imprime el texto al igual que su fondo.   

ej.  `\033[38;2;100;30;200m`

<img width="136" height="48" alt="ejemplo Hello World" src="https://github.com/user-attachments/assets/d183bea1-e3ba-4028-ab1a-6b6b822c0792" />


#### 📦 Bloques Unicode
Esto combinado con el uso de caracteres [__Unicode__](https://symbl.cc/es/unicode-table/).   
Proporciona al usuario una variedad de combinaciones para la creacion de imagenes dentro del terminal, tales como:

- `▀` Mitad de Bloque Superior ([U+2580](https://symbl.cc/es/2580/)"`▀`")
- `▄` Mitad de Bloque Inferior ([U+2584](https://symbl.cc/es/2584/)"`▄`")
- `█` Mitad de Bloque Superior ([U+2588](https://symbl.cc/es/2588/)"`█`")

> [!NOTE]
> Estos mismos son los bloques se utilizan dentro del proyecto para la renderizacion de imagenesalestilo Unicode/ANSI


### 🖼️ Otros Ejemplos:

<img width="1355" height="657" alt="chica y fondo" src="https://github.com/user-attachments/assets/26f4af02-60fa-4b56-9877-85d3bca9ac00" />

----

<img width="1365" height="650" alt="tipo" src="https://github.com/user-attachments/assets/aab2bc49-7c86-4678-99fb-ffa330a22377" />

----

<img width="1363" height="655" alt="coloso" src="https://github.com/user-attachments/assets/c2189572-5655-4b0e-9b4f-7a7a7ff035c3" />

----

> [!CAUTION]
> Precaucion a la hora de reajustar los bordes del CLI luego de la impresion de la imagen, ya que este podria causar distoriones en la misma, al no actualizar el estado de `GetTerminalSize()`
