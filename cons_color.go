package cons

import (
	"fmt"
	"regexp"
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
		"b": "1",  // bold
		"k": "5",  // blink
		"u": "4",  // underline
		"i": "3",  // italic
		"v": "7",  // inverted
		"f": "51", // framed
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
	//result = append_code(result, spec.bg_style)
	return result
}

func IsValid(s string) bool {
	rx := []string{"(?:^[krgbymcw]$)", "(?:^[bkuivf]$)"}
	halves := strings.Split(s, ":")
	for _, h := range halves {
		if len(h) == 0 {
			continue
		}
		quarters := strings.Split(h, "+")
		rxIx := 0
		for _, q := range quarters {
			if len(q) == 0 {
				rxIx++
				continue
			}
			//fmt.Printf("pattern = %s\n", q)
			match, err := regexp.MatchString(rx[rxIx], q)
			if err != nil || !match {
				return false
			}
			rxIx++
		}
	}
	return true
}

// ColorCode returns the ansi control chars for coloring and styling console text
// code should be in UTF-8
func ColorCode(code string) string {
	code = strings.ToLower(code)
	if plain || code == "" || code == "reset" {
		return Reset
	}

	var spec style_spec
	var ptrs = []*string{&spec.fg_color, &spec.fg_style, &spec.bg_color, &spec.bg_style}
	ptrIx := 0
	for _, roon := range code {
		switch roon {
		case ':':
			ptrIx = 2
		case 'v', 'i', 'b', 'u', 'k', 'f', 'y', 'w', 'r', 'g', 'm', 'c':
			*ptrs[ptrIx] = string(roon)
			ptrIx++
		case '+':
			// do nothing, for backwards compatibility
		default: // i.e., +, _, -, .
			ptrIx++
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

func Code2Ansi(code string) string {
	code = strings.ToLower(code)
	if plain || code == "" {
		return ""
	}
	return "" // TODO
}

func resetIfNeeded(code string) string {
	if len(code) > 0 {
		return Reset
	}
	return ""
}

// Color(s, style) Surrounds `s` with ANSI color and reset code.
// deprecated, use Style instead
func Color(s, style string) string {
	code := ColorCode(style)
	return code + s + resetIfNeeded(code)
}

func Style(style string, s string) string {
	code := ColorCode(style)
	return code + s + resetIfNeeded(code)
}

// if flag is false, returns s unchanged
func StyleIf(style string, s string, flag bool) string {
	if flag {
		return Style(style, s)
	}
	return s
}

// ColorFunc Creates a fast closure.
//
// Prefer ColorFunc over Style, Printf etc. if you are concerned about performance
// as it does not recompute ANSI codes.
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

// convenience function vars
var ShowRed = ColorFunc("r")
var ShowBlue = ColorFunc("b")
var ShowYellow = ColorFunc("y")
var ShowMagenta = ColorFunc("m")
var ShowGreen = ColorFunc("g")
var ShowCyan = ColorFunc("c")

// DisableColors disables ANSI color codes. On by default.
func DisableColors(disable bool) {
	plain = disable
}

// convenience function
func Printfs(style, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args)
	fmt.Printf("%s", Style(style, s))
}

// convenience function
func Sprintfs(style, format string, args ...interface{}) string {
	s := fmt.Sprintf(format, args)
	return Style(style, s)
}
