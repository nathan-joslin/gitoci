// Package comms implements utilities for communicating with Git, i.e. reading from and writing to Git.
package comms

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_batcher_Read(t *testing.T) {
	// NOTE: this test setup is not safe to run tests in parallel.
	// mock Git output, the input of batcher
	gitOut := new(bytes.Buffer)
	// mock Git input, the output of batcher
	gitIn := new(bytes.Buffer)

	batcher := NewBatcher(gitOut, gitIn)

	tests := []struct {
		name       string
		mockGitOut []string
		want       Command
		wantErr    bool
	}{
		// TODO: Add test cases.
		{name: "Capabilities",
			mockGitOut: []string{
				"capabilities",
			},
			want: Command{
				CommandType: CmdCapabilities,
				Data:        []string{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, line := range tt.mockGitOut {
				_, err := gitOut.WriteString(line)
				if err != nil {
					t.Fatalf("failed to mock Git output error =  %v", err)
				}
			}
			got, err := batcher.Read()
			if (err != nil) != tt.wantErr {
				t.Errorf("batcher.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want.CommandType, got.CommandType)
			assert.ElementsMatch(t, tt.want.Data, got.Data)
		})
	}
}
