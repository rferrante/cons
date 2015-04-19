package cons

import (
	"fmt"
	"testing"
)

var _GOOD []string = []string{"r", "bb", ".u", "yu", "--b", "-bb", "-ub", "wub", "w-b", "wvm", "cum"}
var _BAD []string = []string{"x", "r,u", "r-u", "v", "v:r+u", "r:v", "y:y+y"}

func TestEmpty(t *testing.T) {
	fmt.Printf("%4s %4s %4s %4s\n%4s %4s %4s %4s\n",
		Color("Test", "w"), Color("Red ", "r"), Color("Yell", "y"), Color("Gree", "g"),
		Color("Blue", "b"), Color("Mage", "m"), Color("Cyan", "c"), Color("Blac", "k"))
	fmt.Printf("%4s %4s %4s %4s\n%4s %4s %4s %4s\n",
		Color("Test", "wu"), Color("Red ", "ru"), Color("Yell", "yu"), Color("Gree", "gu"),
		Color("Blue", "bu"), Color("Mage", "mu"), Color("Cyan", "cu"), Color("Blac", "ku"))
	fmt.Printf("%4s %4s %4s %4s\n%4s %4s %4s %4s\n",
		Color("Test", "wb"), Color("Red ", "rb"), Color("Yell", "yb"), Color("Gree", "gb"),
		Color("Blue", "bb"), Color("Mage", "mb"), Color("Cyan", "cb"), Color("Blac", "kb"))
	fmt.Printf("%4s %4s %4s %4s\n%4s %4s %4s %4s\n",
		Color("Test", "wv"), Color("Red ", "rv"), Color("Yell", "yv"), Color("Gree", "gv"),
		Color("Blue", "bbw"), Color("Mage", "mbw"), Color("Cyan", "cbw"), Color("Blac", "kbw"))
	fmt.Printf("%4s %4s %4s %4s\n%4s %4s %4s %4s\n",
		Color("Test", "wur"), Color("Red ", "ruw"), Color("Yell", "yuc"), Color("Gree", "gum"),
		Color("Blue", "b_w"), Color("Mage", "m_w"), Color("Cyan", "c_w"), Color("Blac", "k_w"))
	t.Log("Example")
}

func TestBadEdgeCases(t *testing.T) {
	fmt.Println("Bad cases...")
	for _, b := range _BAD {
		fmt.Printf("%4s (%s)\n", Color("Bad ", b), b)
	}
	t.Log("Example")
}
func TestGoodEdgeCases(t *testing.T) {
	fmt.Println("Good cases...")
	for _, g := range _GOOD {
		fmt.Printf("%4s (%s)\n", Color("Good", g), g)
	}
	t.Log("Example")
}

func TestFormattedPrinting(t *testing.T) {
	Printfs("ru", "%s and %s\n", "Some", "Some more")
	Printfs("ru", "%d and %d\n", 27, 29)
}

func TestFlagged(t *testing.T) {
	fmt.Println(StyleIf("bu", "Should be styled", true))
	fmt.Println(StyleIf("bu", "Should NOT be styled", false))
}

func inspect(code string) string {
	var s string = "inspector: "
	for _, roon := range code {
		s += fmt.Sprintf("% q", roon)
	}
	return s
}

func TestCodes(t *testing.T) {
	if ColorCode("r") != Start+"31m" {
		t.Logf("r failed, is %s\n", inspect(ColorCode("r")))
		t.Fail()
	}
}

func TestValidCodes(t *testing.T) {
	for _, g := range _GOOD {
		if !IsValid(g) {
			t.Fatalf("IsValid failed a good pattern: %s", g)
		}
	}
	for _, b := range _BAD {
		if IsValid(b) {
			t.Fatalf("IsValid passed a bad pattern: %s", b)
		}
	}
}
