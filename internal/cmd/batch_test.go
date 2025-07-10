package cmd

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_batcher_Read(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name       string
		mockGitOut []string
		want       Git
		wantErr    bool
	}{
		{
			name: "Capabilities",
			mockGitOut: []string{
				"capabilities",
			},
			want: Git{
				Cmd:    Capabilities,
				SubCmd: "",
				Data:   []string{},
			},
			wantErr: false,
		},
		{
			name: "Option Verbosity",
			mockGitOut: []string{
				"option verbosity 4",
			},
			want: Git{
				Cmd:    Option,
				SubCmd: OptionVerbosity,
				Data:   []string{"4"},
			},
			wantErr: false,
		},
		{
			name: "Empty/Done",
			mockGitOut: []string{
				"\n",
			},
			want: Git{
				Cmd:    Empty,
				SubCmd: "",
				Data:   []string{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitOut := new(bytes.Buffer)
			// mock Git input, the output of batcher
			gitIn := new(bytes.Buffer)

			batcher := NewBatcher(gitOut, gitIn)
			for _, line := range tt.mockGitOut {
				_, err := gitOut.WriteString(line)
				if err != nil {
					t.Fatalf("failed to mock Git output error =  %v", err)
				}
			}
			got, err := batcher.Read(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("batcher.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want.Cmd, got.Cmd)
			assert.Equal(t, tt.want.SubCmd, got.SubCmd)
			assert.ElementsMatch(t, tt.want.Data, got.Data)
		})
	}
}
