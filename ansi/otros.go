package ansi
/*
Modos que funcionan en PowerShell moderno:
	ESC[=7h → Habilita el ajuste de línea (auto-wrap) Este sí funciona.
	El texto salta automáticamente a la siguiente línea al llegar al borde. Puedes desactivarlo con ESC[=7l.
*/

// ShowCursor muestra u oculta el cursor en la terminal.
// Parámetro 'show': true para mostrar, false para ocultar.
func ShowCursor(show bool) string {
	if show {return Esc + "?25h"
	} else  {return Esc + "?25l"}
}

// Auto_Wrap habilita o deshabilita el ajuste automático de línea en la terminal.
// Parámetro 'isActive': true para habilitar, false para deshabilitar.
func Auto_Wrap(isActive bool) string {
	if isActive {return Esc + "?7h"
	} else 		{return Esc + "?7l"}
}

// AlternativeScreen activa o desactiva el modo de pantalla alternativa.
// Parámetro 'isActive': true para activar, false para desactivar.
func AlternativeScreen(isActive bool) string {
	if isActive {return Esc + "?1049h"
	} else 		{return Esc + "?1049l"}	
}