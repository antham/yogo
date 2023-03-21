package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	version = "v2.7.0"

	var output bytes.Buffer
	var outputErr bytes.Buffer
	cmd := &cobra.Command{}
	cmd.SetOut(&output)
	cmd.SetErr(&outputErr)
	versionCmd.Run(cmd, []string{})
	assert.Equal(t, "v2.7.0\n", output.String())
	assert.Empty(t, outputErr.String())
}
