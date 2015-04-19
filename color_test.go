package cons

import (
	"fmt"
	"testing"
)

const _GOOD []string = []string{"r", "b+b", "+u", "y+u", ":b", ":b+u", "w:b+u", "w:b", "w:m+v", "c+u:m+b"}
const _BAD []string = []string{"x", "r,u", "r-u", "v", "v:r+u", "r:v", "y:y+y"}

func TestEmpty(t *testing.T) {
	fmt.Printf("%4s %4s %4s %4s\n%4s %4s %4s %4s\n",
		Color("Test", "w"), Color("Red ", "r"), Color("Yell", "y"), Color("Gree", "g"),
		Color("Blue", "b"), Color("Mage", "m"), Color("Cyan", "c"), Color("Blac", "k"))
	fmt.Printf("%4s %4s %4s %4s\n%4s %4s %4s %4s\n",
		Color("Test", "w+u"), Color("Red ", "r+u"), Color("Yell", "y+u"), Color("Gree", "g+u"),
		Color("Blue", "b+u"), Color("Mage", "m+u"), Color("Cyan", "c+u"), Color("Blac", "k+u"))
	fmt.Printf("%4s %4s %4s %4s\n%4s %4s %4s %4s\n",
		Color("Test", "w+b"), Color("Red ", "r+b"), Color("Yell", "y+b"), Color("Gree", "g+b"),
		Color("Blue", "b+b"), Color("Mage", "m+b"), Color("Cyan", "c+b"), Color("Blac", "k+b"))
	fmt.Printf("%4s %4s %4s %4s\n%4s %4s %4s %4s\n",
		Color("Test", "w+i"), Color("Red ", "r+i"), Color("Yell", "y+i"), Color("Gree", "g+i"),
		Color("Blue", "b+i"), Color("Mage", "m+i"), Color("Cyan", "c+i"), Color("Blac", "k+i"))
	fmt.Printf("%4s %4s %4s %4s\n%4s %4s %4s %4s\n",
		Color("Test", "w+k"), Color("Red ", "r+k"), Color("Yell", "y+k"), Color("Gree", "g+k"),
		Color("Blue", "b+k"), Color("Mage", "m+k"), Color("Cyan", "c+k"), Color("Blac", "k+k"))
	fmt.Printf("%4s %4s %4s %4s\n%4s %4s %4s %4s\n",
		Color("Test", "w+v"), Color("Red ", "r+v"), Color("Yell", "y+v"), Color("Gree", "g+v"),
		Color("Blue", "b+v"), Color("Mage", "m+v"), Color("Cyan", "c+v"), Color("Blac", "k+v"))
	fmt.Printf("%4s %4s %4s %4s\n%4s %4s %4s %4s\n",
		Color("Test", "w+f"), Color("Red ", "r+f"), Color("Yell", "y+f"), Color("Gree", "g+f"),
		Color("Blue", "b+f"), Color("Mage", "m+f"), Color("Cyan", "c+f"), Color("Blac", "k+f"))
	fmt.Printf("%4s %4s %4s %4s\n%4s %4s %4s %4s\n",
		Color("Test", "w:w"), Color("Red ", "r:w"), Color("Yell", "y:w"), Color("Gree", "g:w"),
		Color("Blue", "b:w"), Color("Mage", "m:w"), Color("Cyan", "c:w"), Color("Blac", "k:w"))
	t.Log("Example")
}

func TestBadEdgeCases(t *testing.T) {
	bad := []string{"x", "r,u", "r-u", "v", "v:r+u", "r:v", "y:y+y"}
	for _, b := range _BAD {
		fmt.Printf("%4s (%s)\n", Color("Bad ", b), b)
	}
	t.Log("Example")
}
func TestGoodEdgeCases(t *testing.T) {
	good := []string{"r", "b+b", "+u", "y+u", ":b", ":b+u", "w:b+u", "w:b", "w:m+v", "c+u:m+b"}
	for _, g := range _GOOD {
		fmt.Printf("%4s (%s)\n", Color("Good", g), g)
	}
	t.Log("Example")
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
