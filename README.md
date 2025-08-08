# 💻 terminal-go
Proporciona un conjunto de herramientas para el manejo avanzado de códigos de escape ANSI y la renderización de imágenes directamente en la terminal, facilitando el desarrollo de aplicaciones de línea de comandos interactivas y visualmente atractivas en Go.

### 🖼️: Ejemplo de Imagen:

[imagen]

----

El proyecto incluye diversas funciones que permiten el manejo de codigos ANSI Escape los cuales facilitan la interaccion con el __CLI__ del desarrollador, tales como:

``` go
func MoveDown_Start(lines int) string {
	return Esc + fmt.Sprintf("%dE", lines)
}
```

Que permite mover el cursor del terminal hacia abajo una determinada cantidad de lineas y en la primera columna.

## Objetivos

- Permitir al usuario un mayor manejo y facilidad de las herraminetas que ofrecen los [__codigos de escape ANSI__](https://gist.github.com/ConnerWill/d4b6c776b509add763e17f9f113fd25b)
- Creacion de imagenes directamente en el terminal.

## ⚙️ ¿Como Funciona?

#### 🎨 Comandos de Entindado
La renderizacion de imagenes dentro del terminal se realiza con una combinacion de comandos, que permiten manejar el color con el que se imprime el texto al igual que su fondo.
ej.  `\033[38;2;100;30;200m`

[imagen de ejemplo]

#### 📦 Bloques Unicode
Esto combinado con el uso de caracteres [__Unicode__](https://symbl.cc/es/unicode-table/), tales como:

- `▀` Mitad de Bloque Superior ([U+2580](https://symbl.cc/es/2580/)"`▀`")
- `▄` Mitad de Bloque Inferior ([U+2580](https://symbl.cc/es/2584/)"`▄`")
- `█` Mitad de Bloque Superior ([U+2588](https://symbl.cc/es/2588/)"`█`")

> [!NOTE]
> Estos mismos son los bloques se utilizan dentro del proyecto para la creacion de imagenes

Proporciona al usuario una variedad de combinaciones para la creacion de imagenes dentro del terminal

### 🖼️ Otros Ejemplos:

[imagen 1]
----
[imagen 2]
----
[imagen 3]
----

> [!CAUTION]
> Precaucion a la hora de reajustar los bordes del CLI luego de la impresion de la imagen, ya que este podria causar distoriones en la misma, al no actualizar el estado de `GetTerminalSize()`
