package godexrcvr

import "testing"

func TestMmmolToMg(t *testing.T) {
	tt := []struct {
		name     string
		mmol     float32
		expected int
	}{
		{
			name:     "test 1",
			mmol:     1.0,
			expected: 18,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual := MmolToMg(tc.mmol)
			if actual != tc.expected {
				t.Errorf("expected %d, got %d", tc.expected, actual)
			}
		})
	}
}

func TestMgToMmol(t *testing.T) {
	tt := []struct {
		name     string
		mg       int
		expected float32
	}{
		{
			name:     "test 18",
			mg:       18,
			expected: 1.0,
		},
		{
			name:     "test 123",
			mg:       123,
			expected: 6.833333,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual := MgToMmol(tc.mg)
			if actual != tc.expected {
				t.Errorf("expected %f, got %f, diff: %f", tc.expected, actual, actual-tc.expected)
			}
		})
	}
}
