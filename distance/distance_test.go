package distance

import (
	"context"
	"math"
	"testing"
)

func TestCalculateDistance(t *testing.T) {
	tests := []struct {
		name    string
		lat1    float64
		lng1    float64
		lat2    float64
		lng2    float64
		units   []string
		wantMin float64
		wantMax float64
	}{
		{
			name:    "same point returns zero",
			lat1:    -6.2088,
			lng1:    106.8456,
			lat2:    -6.2088,
			lng2:    106.8456,
			wantMin: 0,
			wantMax: 0.01, // floating point tolerance
		},
		{
			name:    "Jakarta to Bandung in default unit (miles*1000)",
			lat1:    -6.2088,
			lng1:    106.8456,
			lat2:    -6.9175,
			lng2:    107.6191,
			wantMin: 60000,
			wantMax: 90000,
		},
		{
			name:    "Jakarta to Bandung with K unit (km*1000 = meters)",
			lat1:    -6.2088,
			lng1:    106.8456,
			lat2:    -6.9175,
			lng2:    107.6191,
			units:   []string{"K"},
			wantMin: 100000,
			wantMax: 250000,
		},
		{
			name:    "New York to London (approx 5570km)",
			lat1:    40.7128,
			lng1:    -74.0060,
			lat2:    51.5074,
			lng2:    -0.1278,
			units:   []string{"K"},
			wantMin: 5000000,
			wantMax: 6000000,
		},
		{
			name:    "equator points 1 degree apart",
			lat1:    0,
			lng1:    0,
			lat2:    0,
			lng2:    1,
			units:   []string{"K"},
			wantMin: 100000,
			wantMax: 120000,
		},
		{
			name:    "north pole to south pole",
			lat1:    90,
			lng1:    0,
			lat2:    -90,
			lng2:    0,
			units:   []string{"K"},
			wantMin: 19000000,
			wantMax: 21000000,
		},
		{
			name:    "unknown unit treated as miles*1000",
			lat1:    0,
			lng1:    0,
			lat2:    0,
			lng2:    1,
			units:   []string{"M"},
			wantMin: 60000,
			wantMax: 80000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateDistance(context.Background(), tt.lat1, tt.lng1, tt.lat2, tt.lng2, tt.units...)
			if math.IsNaN(got) {
				t.Errorf("CalculateDistance() returned NaN")
				return
			}
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("CalculateDistance() = %v, want between %v and %v", got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestCalculateDistance_Symmetry(t *testing.T) {
	// Distance from A to B should equal distance from B to A
	lat1, lng1 := -6.2088, 106.8456
	lat2, lng2 := -6.9175, 107.6191

	d1 := CalculateDistance(context.Background(), lat1, lng1, lat2, lng2, "K")
	d2 := CalculateDistance(context.Background(), lat2, lng2, lat1, lng1, "K")

	if math.Abs(d1-d2) > 0.01 {
		t.Errorf("Distance is not symmetric: %v vs %v", d1, d2)
	}
}
