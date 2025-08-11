# 🖼️ Terminal-Go

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)
[![Platform](https://img.shields.io/badge/Platform-Windows-lightgrey.svg)](https://www.microsoft.com/windows)

> **Renderizador de imágenes avanzado para terminal CLI** que transforma cualquier imagen en arte ASCII/Unicode utilizando códigos de escape ANSI y bloques Unicode optimizados.

## 📖 Descripción

**Terminal-Go** es una biblioteca Go especializada que permite renderizar imágenes de alta calidad directamente en la terminal CLI. Utiliza una combinación sofisticada de:

- 🎨 **Códigos ANSI RGB de 24 bits** para colores verdaderos
- 📦 **Bloques Unicode** (`▀`, `▄`, `█`) para simulación de píxeles 
- ⚡ **Algoritmos optimizados** para rendimiento en tiempo real
- 🎮 **Control interactivo** con navegación por teclado

## ✨ Características Principales

### 🖼️ Renderización Avanzada
- **Soporte multi-formato**: JPEG, PNG, BMP, TIFF, WebP, GIF
- **Escalado inteligente** que preserva las proporciones
- **Interpolación configurable** (Nearest Neighbor, Linear, etc.)
- **Ajuste automático** al tamaño del terminal

### 🎨 Optimización de Colores
- **True Color RGB** (16.7 millones de colores)
- **Transparencia multinivel** con 5 escalas de opacidad
- **Cache de códigos ANSI** para máximo rendimiento
- **Detección de colores similares** para optimización

### 🎮 Interactividad
- **Navegación con flechas** para mover la imagen
- **Ajuste dinámico** al redimensionar terminal
- **Modo pantalla alternativa** para preservar contenido
- **Control total del cursor** y configuración UI

## 🚀 Instalación

```powershell
go get -u github.com/Leontas-9/terminal-go
```

### Dependencias
```powershell
go get golang.org/x/image
go get golang.org/x/sys/windows
go get github.com/eiannone/keyboard
```

## 🖼️ Ejemplos de Imagenes

<img width="408" height="408" alt="hollow-knight" src="https://github.com/user-attachments/assets/5c0d2f09-b974-4dc7-bdda-6cd3d8d042a2" />


<img width="421" height="402" alt="hollow-knight-render" src="https://github.com/user-attachments/assets/3098d6a0-c5a1-4509-9b7a-c1aa7028c003" />

---
    
<img width="597" height="409" alt="paisaje" src="https://github.com/user-attachments/assets/d5c13bf8-26db-46a9-9992-b51831bdb393" />


<img width="837" height="402" alt="paisaje-render" src="https://github.com/user-attachments/assets/7ab8c030-fc71-4f62-8a00-6cfcad4b8fa8" />

---

## 💻 Uso Básico

### Ejemplo Simple
```go
package main

import (
    "github.com/Leontas-9/terminal-go/terminal/render"
    "fmt"
    "log"
)

func main() {
    // Cargar imagen desde archivo
    imgRGBA, err := terminal.LoadImage("mi_imagen.jpg")
    if err != nil {
        log.Fatal(err)
    }
    defer terminal.PutReusableRGBA(imgRGBA) // Liberar memoria

    // Crear imagen renderizable con configuración básica
    src := terminal.NewImage(imgRGBA)

    // Renderizar en terminal
    _, err = src.Print()
    if err != nil {
        log.Fatal(err)
    }
}
```

### Ejemplo Avanzado con Configuración Personalizada
```go
package main

import (
    "github.com/Leontas-9/terminal-go/terminal/render"
    "golang.org/x/image/draw"
    "image"
    "log"
)

func main() {
    // Cargar imagen
    imgRGBA, err := terminal.LoadImage("imagen_compleja.png")
    if err != nil {
        log.Fatal(err)
    }
    defer terminal.PutReusableRGBA(imgRGBA)

    // Configuración UI personalizada
    uiConfig := &terminal.UI_Settings{
        ShowCursor:        false,           // Ocultar cursor durante renderizado
        AlternativeScreen: true,            // Usar pantalla alternativa
        EraseScreen:       true,            // Limpiar pantalla antes
        Auto_Wrap:         false,           // Desactivar ajuste automático
    }

    // Crear imagen con parámetros personalizados
    src := terminal.NewCustomImage(
        imgRGBA,                    // Imagen a renderizar
        image.Rect(0, 0, 80, 40),   // Área de renderizado (80x40 caracteres)
        image.Pt(10, 5),            // Posición inicial (columna 10, fila 5)
        draw.BiLinear,              // Interpolación suave
        uiConfig,                   // Configuración UI
    )

    // Renderizar imagen
    _, err = src.Print()
    if err != nil {
        log.Fatal(err)
    }

    // Modo interactivo - navegar con flechas del teclado
    err = src.Displacement()
    if err != nil {
        log.Fatal(err)
    }
}
```

## 🎛️ Configuración Avanzada

### UI_Settings - Opciones de Interfaz

| Parámetro | Tipo | Descripción | Valor por Defecto |
|-----------|------|-------------|-------------------|
| `ShowCursor` | `bool` | Muestra/oculta el cursor al finalizar | `true` |
| `AlternativeScreen` | `bool` | Usa pantalla alternativa (preserva contenido previo) | `false` |
| `EraseScreen` | `bool` | Limpia la pantalla antes del renderizado | `false` |
| `Auto_Wrap` | `bool` | Habilita ajuste automático de línea | `true` |

### Interpoladores Disponibles

```go
import "golang.org/x/image/draw"

// Interpoladores soportados:
draw.NearestNeighbor    // Más rápido, pixelado
draw.BiLinear          // Balance calidad/velocidad  
draw.CatmullRom        // Máxima calidad, más lento
```

## 🎮 Controles Interactivos

En el modo `Displacement()`, puedes controlar la imagen con:

| Tecla | Acción |
|-------|---------|
| `↑` | Mover imagen arriba |
| `↓` | Mover imagen abajo |
| `←` | Mover imagen izquierda |
| `→` | Mover imagen derecha |
| `Esc` / `Ctrl+C` | Salir del modo interactivo |

## 📁 Formatos Soportados

| Formato | Extensión | Notas |
|---------|-----------|-------|
| **JPEG** | `.jpg`, `.jpeg` | Compresión con pérdida |
| **PNG** | `.png` | Soporte completo de transparencia |
| **BMP** | `.bmp` | Formato bitmap sin compresión |
| **TIFF** | `.tiff` | Alta calidad, múltiples capas |
| **WebP** | `.webp` | Formato moderno de Google |
| **GIF** | `.gif` | ⚠️ Solo renderiza el primer frame |

## 🔧 API Detallada

### Funciones de Carga
```go
// LoadImage carga cualquier formato soportado y lo convierte a RGBA
func LoadImage(filepath string) (*image.RGBA, error)

// PutReusableRGBA libera la memoria de imagen (usar con defer)
func PutReusableRGBA(img *image.RGBA)
```

### Constructores de Imagen
```go
// NewImage - Configuración automática básica
func NewImage(img *image.RGBA) *RenderImage

// NewCustomImage - Control total de parámetros
func NewCustomImage(
    img *image.RGBA, 
    bounds image.Rectangle, 
    initialPoint image.Point, 
    interpolator draw.Interpolator, 
    opts *UI_Settings
) *RenderImage
```

### Métodos de RenderImage
```go
// Print - Renderiza imagen directamente en terminal
func (src *RenderImage) Print() (int, error)

// GetPNG - Obtiene datos de imagen en formato ANSI
func (src *RenderImage) GetPNG() ([]byte, image.Image, error)

// Displacement - Modo interactivo con controles de teclado
func (src *RenderImage) Displacement() error

// Métodos de configuración
func (src *RenderImage) SetUI_Settings(new *UI_Settings)
func (src *RenderImage) SetMargins(new image.Rectangle)
func (src *RenderImage) SetInterpolator(new draw.Interpolator)
func (src *RenderImage) SetInitialPoint(new image.Point)
```

## 🛠️ Funciones ANSI Utilitarias

### Control de Cursor
```go
import "github.com/Leontas-9/terminal-go/ansi"

ansi.MoveToStart()              // Cursor a posición (1,1)
ansi.MoveTo(column, row)        // Mover a posición específica
ansi.MoveUp(lines)              // Mover hacia arriba
ansi.MoveDown(lines)            // Mover hacia abajo
ansi.MoveLeft(columns)          // Mover hacia izquierda
ansi.MoveRight(columns)         // Mover hacia derecha
```

### Limpieza de Pantalla
```go
ansi.EraseScreen()              // Limpiar pantalla completa
ansi.EraseLine()                // Limpiar línea actual
ansi.EraseRectangle(rect)       // Limpiar área rectangular
```

### Configuración de Terminal
```go
ansi.ShowCursor(show)           // Mostrar/ocultar cursor
ansi.AlternativeScreen(active)  // Pantalla alternativa
ansi.Auto_Wrap(active)          // Ajuste automático de línea
```

## 🔬 Aspectos Técnicos

### Algoritmo de Renderizado

1. **Carga y Validación**: La imagen se carga y convierte a formato RGBA
2. **Ajuste de Bordes**: Se calculan los límites dentro del terminal
3. **Escalado Proporcional**: La imagen se redimensiona preservando aspecto
4. **Renderizado por Bloques**: Cada 2 píxeles verticales se convierten en 1 bloque Unicode
5. **Optimización de Color**: Se detectan colores similares para reducir códigos ANSI
6. **Salida Optimizada**: Se genera la secuencia mínima de caracteres ANSI

### Pool de Memoria
El sistema utiliza pools de memoria reutilizable para optimizar el rendimiento:

```go
// Resoluciones precargadas para diferentes tamaños
var resolutions = []image.Point{
    {141, 58},    // Terminal por defecto
    {320, 240},   // QVGA
    {640, 480},   // VGA
    {1280, 720},  // HD
    {1920, 1080}, // Full HD
    {2560, 1440}, // QHD
    {3840, 2160}, // 4K
}
```

### Códigos ANSI Precomputados
Para máximo rendimiento, los códigos ANSI se precalculan:

```go
// Ejemplo de código ANSI para color RGB
"\033[38;2;255;82;197;48;2;155;106;0m"
//^^^ ^^^^ ^^^ ^^ ^^^ ^^^^ ^^^ ^^^ ^
// |    |   |  |   |   |    |   |  | 
// |    |   R  G   B   |    R   G  B 
// |    |              | 
// | color RGB     color RGB 
// |  (texto)       (fondo)
// |
// escape ANSI
```

## 🎯 Ejemplos de Código

### Renderizado con Márgenes Personalizados
```go
func RenderWithCustomMargins() {
    imgRGBA, _ := terminal.LoadImage("imagen.png")
    defer terminal.PutReusableRGBA(imgRGBA)
    
    src := terminal.NewImage(imgRGBA)
    
    // Establecer área de renderizado de 100x50 caracteres
    src.SetMargins(image.Rect(0, 0, 100, 50))
    
    // Posicionar en el centro-derecha del terminal
    src.SetInitialPoint(image.Pt(50, 10))
    
    src.Print()
}
```

### Modo Presentación (Sin Interferencias)
```go
func PresentationMode() {
    imgRGBA, _ := terminal.LoadImage("presentacion.jpg")
    defer terminal.PutReusableRGBA(imgRGBA)
    
    config := &terminal.UI_Settings{
        ShowCursor:        false,  // Sin cursor visible
        AlternativeScreen: true,   // Pantalla limpia
        EraseScreen:       true,   // Limpiar antes
        Auto_Wrap:         false,  // Control total
    }
    
    src := terminal.NewCustomImage(
        imgRGBA,
        imgRGBA.Rect,
        image.Pt(0, 0),
        draw.BiLinear,
        config,
    )
    
    src.Print()
}
```

### Galería Interactiva
```go
func InteractiveGallery() {
    images := []string{"img1.jpg", "img2.png", "img3.webp"}
    
    for _, imgPath := range images {
        imgRGBA, err := terminal.LoadImage(imgPath)
        if err != nil {
            continue
        }
        
        src := terminal.NewImage(imgRGBA)
        src.Print()
        
        // Permitir navegación interactiva
        src.Displacement() // Usar flechas para mover, Esc para siguiente
        
        terminal.PutReusableRGBA(imgRGBA)
    }
}
```

## 🏗️ Arquitectura del Proyecto

```
terminal-go/
├── terminal/
│   ├── ansi/                       # Módulo de códigos ANSI
│   │   ├── colors.go           # Gestión de colores RGB/ANSI
│   │   ├── cursor.go           # Control de posición del cursor  
│   │   ├── erase.go            # Funciones de limpieza
│   │   └── otros.go            # Configuraciones adicionales
│   ├── render/                     # Motor de renderizado principal
│   │   ├── assignment.go       # Estructuras y constructores
│   │   ├── files.go            # Carga de archivos de imagen
│   │   ├── init.go             # Inicialización y pools de memoria
│   │   ├── moviment.go         # Sistema de navegación interactiva
│   │   ├── render.go           # Algoritmo de renderizado principal
│   │   └── variables.go        # Constantes y variables globales
│   ├── main/
│   │   └── main.go             # Ejemplo de uso
│   └── image_test/
│       └── soldado.webp        # Imagen de prueba
└── README.md
```

## ⚙️ Funcionamiento Interno

### 1. **Procesamiento de Imagen**
```go
// El sistema lee la imagen y la convierte a RGBA
imgRGBA, err := terminal.LoadImage("imagen.jpg")

// Cada pixel contiene 4 bytes: Red, Green, Blue, Alpha
pixel := color.RGBA{R: 255, G: 128, B: 64, A: 255}
```

### 2. **Mapeo de Bloques Unicode**
```go
// Cada carácter del terminal representa 2 píxeles verticales:
// Píxel superior → Color de texto del bloque ▀
// Píxel inferior → Color de fondo del bloque ▀

upperPixel := image.getPixel(y)       // Fila actual  
lowerPixel := image.getPixel(y+1)     // Fila siguiente

// Resultado: Un bloque ▀ con colores específicos para texto y fondo
```

### 3. **Optimización de Códigos ANSI**
```go
// Sistema de cache para evitar regenerar códigos idénticos
var digitLookup = [256][]byte{} // Tabla precomputada 0-255

// Generación eficiente de códigos RGB:
// "\033[38;2;R;G;B;48;2;R;G;Bm" para texto+fondo
func GetANSI_DoubleColor(buf *[]byte, rF,gF,bF,rB,gB,bB byte)
```

### 4. **Gestión de Transparencia**
```go
// 5 niveles de transparencia con bloques Unicode
const (
    ALPHA_1 = (255 * 1) / 5  // ' ' (vacío)
    ALPHA_2 = (255 * 2) / 5  // '░' (puntos ligeros)  
    ALPHA_3 = (255 * 3) / 5  // '▒' (puntos medios)
    ALPHA_4 = (255 * 4) / 5  // '▓' (puntos densos)
    ALPHA_5 = (255 * 5) / 5  // '█' (sólido)
)
```

## 🔍 Casos de Uso

### 📊 Visualización de Datos
- Mostrar gráficos generados como imágenes
- Dashboards en terminal para servidores
- Monitoring visual en tiempo real

### 🎨 Arte ASCII Moderno  
- Conversión de fotografías a arte terminal
- Logos y banners para aplicaciones CLI
- Efectos visuales en herramientas de desarrollo

### 🖥️ Interfaces de Usuario Avanzadas
- Menús gráficos en aplicaciones de consola
- Previsualizadores de imágenes en terminal
- Herramientas de desarrollo con feedback visual

## ⚡ Optimizaciones de Rendimiento

### Pool de Memoria Reutilizable
- **Evita allocaciones**: Reutiliza objetos `image.RGBA`
- **Resoluciones precargadas**: Pools optimizados para tamaños comunes
- **Limpieza automática**: Gestión transparente de memoria

### Cache de Códigos ANSI
- **Tabla de dígitos precomputada**: Conversión instantánea 0-255
- **Buffers reutilizables**: Minimiza creación de strings
- **Detección de similitud**: Reduce códigos ANSI redundantes

### Algoritmos Eficientes
- **Renderizado por bloques**: Procesa 2 píxeles por iteración
- **Salto de transparencias**: Omite píxeles completamente transparentes
- **Clampeo optimizado**: Valida límites sin branches

## 🚨 Limitaciones y Consideraciones

### 🖥️ Compatibilidad
- **Windows únicamente**: Usa APIs específicas de Windows para detección de tamaño
- **Terminal moderno requerido**: Necesita soporte para ANSI True Color
- **PowerShell/CMD**: Funciona mejor en terminales modernos

### 🎯 Limitaciones de Formato
- **GIF animado**: Solo renderiza el primer frame
- **Imágenes muy grandes**: Se redimensionan automáticamente al terminal
- **Colores limitados**: Aunque son 24-bit, la percepción depende del terminal

### ⚠️ Precauciones
```go
// ⚠️ IMPORTANTE: Siempre liberar memoria
imgRGBA, err := terminal.LoadImage("imagen.jpg")
if err != nil { return err }
defer terminal.PutReusableRGBA(imgRGBA) // ← CRÍTICO

// ⚠️ Evitar redimensionar terminal durante renderizado
// Puede causar distorsiones hasta el próximo renderizado
```

## 🔮 Desarrollo Futuro

### Características Planeadas
- [ ] 🐧 **Soporte Linux/macOS**: Detección de tamaño multi-plataforma
- [ ] 🎞️ **GIF animado**: Renderizado de múltiples frames
- [ ] 🎨 **Paletas de color**: Reducción automática para terminals limitados
- [ ] 📱 **Modo responsivo**: Ajuste automático a redimensionamiento
- [ ] 🎮 **Más controles**: Zoom, rotación, filtros en tiempo real

### Optimizaciones Futuras
- [ ] ⚡ **Renderizado paralelo**: Goroutines para bloques independientes
- [ ] 🧠 **AI upscaling**: Mejora de calidad para imágenes pequeñas
- [ ] 💾 **Compresión inteligente**: Reducción de códigos ANSI repetitivos

## 🤝 Contribuciones

¡Las contribuciones son bienvenidas! Por favor:

1. 🍴 Fork el repositorio
2. 🌿 Crea una rama para tu feature (`git checkout -b feature/nueva-funcionalidad`)
3. 💾 Commit tus cambios (`git commit -am 'Agregar nueva funcionalidad'`)
4. 📤 Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. 🔄 Abre un Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia GNU General Public License v3.0. Ver el archivo `LICENSE` para más detalles.

## 👨‍💻 Autor

**Leontas-9** - [GitHub](https://github.com/Leontas-9)

---

⭐ **¡Si este proyecto te resulta útil, considera darle una estrella!** ⭐

> *"Transformando píxeles en arte ASCII, una terminal a la vez."* 🎨✨
>>>>>>> e121ace (Este es mi primer commit)
