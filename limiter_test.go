package otto

import (
	"context"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestLimiter(t *testing.T) {
	limit := 5

	vm := New()
	vm.Limiter = nil

	script := `(function t() {
	           for (var i = 0; i < 2; i++) {
	           	var a = 1+1;
	           }
	       }())`
	start := time.Now()
	_, err := vm.Eval(script)

	is(err, nil)
	is(vm.Ticks(), 29)
	if time.Since(start) > time.Second*time.Duration(limit) {
		t.Fatalf("expected test to take less than %d seconds", limit)
	}

	vm = New()
	vm.Limiter = rate.NewLimiter(rate.Limit(limit), 1)
	start = time.Now()
	_, err = vm.Eval(script)

	is(err, nil)
	is(vm.Ticks(), 29)
	if time.Since(start) < time.Second*time.Duration(limit) {
		t.Fatalf("expected test to take more than %d seconds", limit)
	}

	vm = New()
	vm.Limiter = rate.NewLimiter(rate.Inf, 1)
	start = time.Now()
	_, err = vm.Eval(script)

	is(err, nil)
	is(vm.Ticks(), 29)
	if time.Since(start) > time.Second*time.Duration(limit) {
		t.Fatalf("expected test to take less than %d seconds", limit)
	}

	script = `(function t() {
		try {
			for (var i = 0; i < 100; i++) {
	           	var a = 1+1;
	           }
		} catch(e) {

		}

	       }())`
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	vm = NewWithContext(ctx)
	vm.Limiter = rate.NewLimiter(rate.Limit(limit), 1)
	start = time.Now()
	evalErr := func() (err interface{}) {
		defer func() {
			err = recover()
		}()
		vm.EvalWithContext(ctx, script)
		return nil
	}()
	is(evalErr, context.DeadlineExceeded)
	is(vm.Ticks(), 11)
	if time.Since(start) > time.Second*time.Duration(limit) {
		t.Fatalf("expected test to take less than %d seconds", limit)
	}

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	vm = NewWithContext(ctx)
	vm.Limiter = rate.NewLimiter(rate.Limit(limit), 1)
	start = time.Now()
	go func() {
		time.Sleep(time.Second * 2)
		cancel()
	}()
	evalErr = func() (err interface{}) {
		defer func() {
			err = recover()
		}()
		vm.EvalWithContext(ctx, script)
		return nil
	}()
	is(evalErr, context.Canceled)
	if vm.Ticks() > 50 {
		t.Fatal("expected to not process many ticks")
	}
	if time.Since(start) > time.Second*time.Duration(limit) {
		t.Fatalf("expected test to take less than %d seconds", limit)
	}
}
