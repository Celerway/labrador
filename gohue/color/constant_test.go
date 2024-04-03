package color

import (
	"testing"
)

func TestFindColorByName(t *testing.T) {

	color, err := FindColorByName("red")
	if err != nil {
		t.Fatalf("Color red should exits")
	}
	if !AlmostEqual32(color.X, 0.640075, 1e-6) {
		t.Fatalf("X for red should be %f", color.X)
	}
	if !AlmostEqual32(color.Y, 0.329971, 1e-6) {
		t.Fatalf("Y for red should be %f", color.Y)
	}
}

func TestCannotFindColorByName(t *testing.T) {
	_, err := FindColorByName("foobar")
	if err == nil {
		t.Fatalf("foobar color should not exist")
	}
}

func AlmostEqual32(a, b, epsilon float32) bool {
	if a == b {
		return true
	}
	if a == 0 || b == 0 {
		return false
	}
	diff := (a - b) / (a + b)
	return diff < epsilon
}
