package main

import "testing"

func TestGetIntCondition(t *testing.T) {
	in := "(0<N<100)"
	var got string
	var wants string

	wants = "0<N<100"
	got = GetCondition(in).String()
	if got != wants {
		t.Errorf("got %s, wants %s", got, wants)
	}
}

func TestProcessValueParam(t *testing.T) {
	in := "int(0<N<100)"

	wants1 := "int"
	wants2 := "0<N<100"
	got1, got2 := ProcessValueParam(in)
	if got1 != wants1 {
		t.Errorf("got %s, wants %s", got1, wants1)
	}
	if got2.String() != wants2 {
		t.Errorf("got %s, wants %s", got2, wants2)
	}
}
