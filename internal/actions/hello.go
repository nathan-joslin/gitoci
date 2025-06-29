package actions

import (
	"context"
	"fmt"
	"io"

	"git.act3-ace.com/ace/go-common/pkg/logger"
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
