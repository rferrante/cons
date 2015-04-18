package cons

import (
	"fmt"
	"strings"
)

const (
	Start = "\033["
	Reset = "\033[0m"
)

var (
	plain  = false
	colors = map[string][2]string{
		"k": {"30", "40"}, // black
		"r": {"31", "41"}, // red
		"g": {"32", "42"}, // green
		"y": {"33", "43"}, // yellow
		"b": {"34", "44"}, // blue
		"m": {"35", "45"}, // magenta
		"c": {"36", "46"}, // cyan
		"w": {"37", "47"}, // white
	}
	styles = map[string]string{
		"b": "1", // bold
		"B": "5", // blink
		"u": "4", // underline
		"i": "7", // italic
	}
	// Three diagnostics set after each call to Color()
	// Tracers[0] is three indexes used in parsing the arg to Color():
	//     the first +, the :, and the second +
	// Tracers[1] gives the parsed colors and styles as foreground color, style, background color, style
	// Tracers[2] gives the final ansi control string
	Tracers [3]string
)

func append_code(base string, code string) string {
	if len(code) == 0 {
		return base
	}
	appended := base
	if len(appended) > 0 {
		appended += ";"
	}
	appended += code
	return appended
}

type style_spec struct {
	fg_color string
	fg_style string
	bg_color string
	bg_style string
}

func (spec style_spec) String() string {
	var result string
	result = append_code(result, spec.fg_color)
	result = append_code(result, spec.fg_style)
	result = append_code(result, spec.bg_color)
	result = append_code(result, spec.bg_style)
	return result
}

// ColorCode returns the ansi control chars for coloring and styling console text
func ColorCode(code string) string {
	if plain || code == "" {
		return ""
	}
	if code == "reset" {
		return Reset
	}
	var spec style_spec
	ix_plus1 := strings.Index(code, "+")
	ix_plus2 := strings.LastIndex(code, "+")
	ix_colon := strings.Index(code, ":")
	Tracers[0] = fmt.Sprintf("%d %d %d", ix_plus1, ix_colon, ix_plus2)

	if ix_colon == -1 && ix_plus1 == -1 {
		spec.fg_color = code
	} else if ix_colon == -1 {
		spec.fg_color = code[:ix_plus1]
		spec.fg_style = code[ix_plus1+1 : len(code)]
	} else if ix_plus1 == -1 {
		spec.fg_color = code[:ix_colon]
		spec.bg_color = code[ix_colon+1 : len(code)]
	} else if ix_plus1 == ix_plus2 {
		if ix_plus1 < ix_colon {
			spec.fg_color = code[:ix_plus1]
			spec.fg_style = code[ix_plus1+1 : ix_colon]
			spec.bg_color = code[ix_colon+1 : len(code)]
		} else {
			spec.fg_color = code[:ix_colon]
			spec.bg_color = code[ix_colon+1 : ix_plus1]
			spec.bg_style = code[ix_plus1+1 : len(code)]
		}
	} else { // there is at least one colon and at least 2 plusses
		if ix_plus1 < ix_colon && ix_colon < ix_plus2 {
			spec.fg_color = code[:ix_plus1]
			spec.fg_style = code[ix_plus1+1 : ix_colon]
			spec.bg_color = code[ix_colon+1 : ix_plus2]
			spec.bg_style = code[ix_plus2+1 : len(code)]
		}
	}
	Tracers[1] = fmt.Sprintf("%s %s %s %s",
		spec.fg_color, spec.fg_style, spec.bg_color, spec.bg_style)

	// now map to the real styles
	spec.fg_color = colors[spec.fg_color][0]
	spec.fg_style = styles[spec.fg_style]
	spec.bg_color = colors[spec.bg_color][1]
	spec.bg_style = styles[spec.bg_style]
	// get the ansi codes without the start and end escapes
	// (useful for printing the codes without sending escapes)
	Tracers[2] = fmt.Sprintf("%s", spec)

	if len(spec.String()) == 0 {
		return ""
	}
	return fmt.Sprintf("%s%sm", Start, spec)
}

func resetIfNeeded(code string) string {
	if len(code) > 0 {
		return Reset
	}
	return ""
}

// Color(s, style) Surrounds `s` with ANSI color and reset code.
func Color(s, style string) string {
	code := ColorCode(style)
	return code + s + resetIfNeeded(code)
}

// ColorFunc Creates a fast closure.
//
// Prefer ColorFunc over Color as it does not recompute ANSI codes.
func ColorFunc(style string) func(string) string {
	if style == "" {
		return func(s string) string {
			return s
		}
	} else {
		code := ColorCode(style)
		return func(s string) string {
			return code + s + resetIfNeeded(code)
		}
	}
}

var ShowRed = ColorFunc("r")
var ShowMagenta = ColorFunc("m")
var ShowGreen = ColorFunc("g")

// DisableColors disables ANSI color codes. On by default.
func DisableColors(disable bool) {
	plain = disable
}
