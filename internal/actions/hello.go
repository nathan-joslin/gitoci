package actions

import (
	"context"
	"fmt"
	"io"

	"github.com/act3-ai/go-common/pkg/logger"
)

// Represents the Hello action
type Hello struct {
	*Tool
}

// Runs the Hello action
func (action *Hello) Run(ctx context.Context, out io.Writer) error {
	log := logger.FromContext(ctx)

	conf, err := action.GetConfig(ctx)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(out, "Hello %s\n", conf.Name)
	if err != nil {
		log.Info("couldn't say hello")
		return err
	}

	return nil
}
