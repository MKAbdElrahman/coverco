package finder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchPattern(t *testing.T) {
	tests := []struct {
		packageName string
		patterns    []string
		expected    bool
		errExpected bool
	}{
		{"github.com/example/project/pkg", []string{"github.com/example/*"}, true, false},
		{"github.com/example/project/pkg", []string{"github.com/test/*"}, false, false},
		{"github.com/example/project/pkg", []string{"github.com/example/project/pkg"}, true, false},
		{"github.com/example/project/pkg", []string{"github.com/example/project/pkg*", "github.com/another/*"}, true, false},
		{"github.com/example/project/pkg", []string{"github.com/example/project/pkg[invalid"}, false, true}, // Invalid regex pattern
	}

	for _, tt := range tests {
		matched, err := matchPattern(tt.packageName, tt.patterns)
		if tt.errExpected {
			assert.Error(t, err)
			var patternMatchErr *PatternMatchError
			assert.ErrorAs(t, err, &patternMatchErr)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, tt.expected, matched)
	}
}
