package id

import "testing"

func TestGenerateID(t *testing.T) {
	cases := map[string]struct {
		counter int
		want    string
	}{
		"two":         {2, "2"},
		"one":         {1, "1"},
		"ninety-nine": {99, "99"},
		"seven":       {7, "7"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := GenerateID(tc.counter); got != tc.want {
				t.Fatalf("GenerateID = %q, want %q", got, tc.want)
			}
		})
	}
}
