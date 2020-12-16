package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSolveSAT(t *testing.T) {
	for _, tt := range []struct {
		cnf    [][]int
		want   []int
		wantOK bool
	}{
		{
			cnf:    [][]int{{3}, {3, 4}, {3, 4, 5}, {-3, -4}, {-3, -5}, {-4, -5}},
			want:   []int{3},
			wantOK: true,
		},
		{
			cnf:    [][]int{{3, 4}, {-3, -4}, {3, -4}, {-3, 4}},
			wantOK: false,
		},
		{
			cnf:    [][]int{{3, 4}, {4, 5}, {-5}, {4, -3}, {5, 3}},
			want:   []int{3, 4},
			wantOK: true,
		},
	} {
		got, ok := solveSAT(tt.cnf)
		if !ok {
			if tt.wantOK {
				t.Errorf("solveSAT(%v): got !ok, want %v", tt.cnf, tt.want)
			}
			continue
		}
		if !tt.wantOK {
			t.Errorf("solveSAT(%v): got %v, want !ok", tt.cnf, got)
			continue
		}
		if diff := cmp.Diff(got, tt.want); diff != "" {
			t.Errorf("solveSAT(%v) (-got, +want):\n%s", tt.cnf, diff)
		}
	}
}
