package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestNewBacklight(t *testing.T) {
	want := &Backlight{50, 100}

	got, err := NewBacklight(
		strings.NewReader("50"),
		strings.NewReader("100"),
	)
	if err != nil {
		t.Error(err)
	}
	if *got != *want {
		t.Errorf("Wanted %#v, Got %#v\n", *want, *got)
	}
}

func TestGet(t *testing.T) {
	cases := []struct {
		name    string
		reader  io.Reader
		want    float64
		wantErr error
	}{
		{
			name:    "basic",
			reader:  strings.NewReader("123"),
			want:    123,
			wantErr: nil,
		},
		{
			name:    "extra-whitespace",
			reader:  strings.NewReader("  123  "),
			want:    123,
			wantErr: nil,
		},
	}

	for _, tc := range cases {
		got, err := get(tc.reader)
		if err != tc.wantErr {
			t.Error(err)
		}
		if got != tc.want {
			t.Errorf("Wanted %v, got %v\n", tc.want, got)
		}
	}
}

func TestSet(t *testing.T) {
	cases := []struct {
		pct  float64
		max  float64
		want string
	}{
		{50, 1000, "500"},
	}
	for _, tc := range cases {
		var w bytes.Buffer
		bl := Backlight{MaxBrightness: tc.max}
		err := bl.Set(&w, tc.pct)
		if err != nil {
			t.Fatal(err)
		}
		got := w.String()
		if got != tc.want {
			t.Fatalf("Wanted %v, got %v\n", tc.want, got)
		}
	}
}

func TestSetIncr(t *testing.T) {
	cases := []struct {
		pctIncr float64
		cur     float64
		max     float64
		want    string
	}{
		{5, 50, 100, "55"},
		{-5, 50, 100, "45"},
	}
	for _, tc := range cases {
		var w bytes.Buffer
		bl := Backlight{CurrentBrightness: tc.cur, MaxBrightness: tc.max}
		err := bl.SetIncr(&w, tc.pctIncr)
		if err != nil {
			t.Fatal(err)
		}
		got := w.String()
		if got != tc.want {
			t.Fatalf("Wanted %v, got %v\n", tc.want, got)
		}
	}
}

func TestPercent(t *testing.T) {
	cases := []struct {
		cur  float64
		max  float64
		want float64
	}{
		{50, 100, 50},
		{0, 100, 0},
	}
	for _, tc := range cases {
		bl := Backlight{
			CurrentBrightness: tc.cur,
			MaxBrightness:     tc.max,
		}
		if got := bl.Percent(); got != tc.want {
			t.Errorf("Wanted %v, got %v\n", tc.want, got)
		}
	}
}

func TestPercentToBrightness(t *testing.T) {
	cases := []struct {
		pct  float64
		max  float64
		want float64
	}{
		{10, 500, 50},
		{50, 4, 2},
		{100, 100, 100},
	}
	for _, tc := range cases {
		bl := Backlight{MaxBrightness: tc.max}
		if got := bl.PercentToBrightness(tc.pct); got != tc.want {
			t.Errorf("Wanted %v, got %v\n", tc.want, got)
		}
	}
}
