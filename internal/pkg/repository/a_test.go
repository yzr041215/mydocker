package repository

import "testing"

func TestA(t *testing.T) {
	a, err := GetImage("nginx", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(a)
}
func TestB(t *testing.T) {
	b, err := GetImagesList()
	if err != nil {
		t.Fatal(err)
	}
	for _, img := range b {
		t.Log(img)
	}
}
func TestC(t *testing.T) {
	SaveImage("A", "latest", "----------")
}
