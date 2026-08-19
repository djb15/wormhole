package sui

// Deliberate violation of the wormhole/go/run-with-scissors-error-return
// CodeQL query, used to verify that alerts surface in the PR view.
// Do not merge.

import (
	"context"
	"errors"

	"github.com/certusone/wormhole/node/pkg/common"
)

var _ = codeqlScissorsDemo

func codeqlScissorsDemo(ctx context.Context, errC chan error) {
	common.RunWithScissors(ctx, errC, "codeql_scissors_demo", func(ctx context.Context) error {
		err := errors.New("demo failure")
		errC <- err
		return nil
	})
}
