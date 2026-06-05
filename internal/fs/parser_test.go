package fs

import (
	"testing"

	"gotest.tools/v3/assert"
)

type Object struct {
	Working bool
}

func TestConstructFromYAML(t *testing.T) {
	path := "/tmp/dockerOrchestrator/dockerOrchestrator/hardware.yml"

	var obj Object

	err := ConstructFromYAML(obj, path)

	assert.Assert(t, obj.Working != true)
	assert.Error(t, err, "object must be a non-nil pointer")

	err = ConstructFromYAML(&obj, path)

	assert.Assert(t, obj.Working == true)
}
