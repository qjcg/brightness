package main

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVal(t *testing.T) {
	testCases := []struct {
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

	for _, tc := range testCases {
		got, err := getVal(tc.reader)

		assert.Equal(t, got, tc.want)
		assert.Equal(t, err, tc.wantErr)
	}
}

func TestSet(t *testing.T) {
	t.Error("NOT IMPLEMENTED")
}
