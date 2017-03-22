package gocql

import "testing"

func TestUnmarshalCassVersion(t *testing.T) {
	tests := [...]struct {
		data    string
		version cassVersion
	}{
		{"3.2", cassVersion{3, 2, 0}},
		{"2.10.1-SNAPSHOT", cassVersion{2, 10, 1}},
		{"1.2.3", cassVersion{1, 2, 3}},
	}

	for i, test := range tests {
		v := &cassVersion{}
		if err := v.UnmarshalCQL(nil, []byte(test.data)); err != nil {
			t.Errorf("%d: %v", i, err)
		} else if *v != test.version {
			t.Errorf("%d: expected %#+v got %#+v", i, test.version, *v)
		}
	}
}

func TestCassVersionBefore(t *testing.T) {
	tests := [...]struct {
		version             cassVersion
		major, minor, patch int
	}{
		{cassVersion{1, 0, 0}, 0, 0, 0},
		{cassVersion{0, 1, 0}, 0, 0, 0},
		{cassVersion{0, 0, 1}, 0, 0, 0},

		{cassVersion{1, 0, 0}, 0, 1, 0},
		{cassVersion{0, 1, 0}, 0, 0, 1},
	}

	for i, test := range tests {
		if !test.version.Before(test.major, test.minor, test.patch) {
			t.Errorf("%d: expected v%d.%d.%d to be before %v", i, test.major, test.minor, test.patch, test.version)
		}
	}

}
