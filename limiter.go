package otto

import (
	"context"
	"strings"
)

func (self *_runtime) waitOneTick() {
	self.ticks++
	if self.otto.Limiter == nil {
		return
	}
	var ctx context.Context
	if self.scope != nil {
		ctx = self.scope.context
	} else {
		ctx = self.ctx
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if waitErr := self.otto.Limiter.Wait(ctx); waitErr != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			panic(ctxErr)
		}
		if strings.Contains(waitErr.Error(), "would exceed") {
			panic(context.DeadlineExceeded)
		}
		panic(waitErr)
	}
}
