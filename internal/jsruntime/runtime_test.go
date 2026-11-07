package jsruntime

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestSynchronousBashAndPersistence(t *testing.T) {
	r := New()
	if result := r.Run(context.Background(), `var result = bash("printf hello", 1); result.stdout`); result != `"hello"` {
		t.Fatal(result)
	}
	if result := r.Run(context.Background(), `result.exitCode`); result != "0" {
		t.Fatal(result)
	}
	if result := r.Run(context.Background(), `typeof console + ":" + typeof require`); result != `"undefined:undefined"` {
		t.Fatal(result)
	}
}

func TestRuntimeCancellationAndRecovery(t *testing.T) {
	r := New()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	if result := r.Run(ctx, `while (true) {}`); !strings.Contains(result, "timed out") {
		t.Fatal(result)
	}
	if result := r.Run(context.Background(), `2 + 2`); result != "4" {
		t.Fatal(result)
	}
}

func TestSerializationErrorsAreToolResults(t *testing.T) {
	r := New()
	if result := r.Run(context.Background(), `var cyclic = {}; cyclic.self = cyclic; cyclic`); !strings.Contains(result, "Cannot serialize") {
		t.Fatal(result)
	}
}

func TestBashIsSynchronous(t *testing.T) {
	result := New().Run(context.Background(), `var r = bash("sleep 0.02; printf done", 1); ({output: r.stdout, promise: typeof r.then})`)
	if result != `{"output":"done","promise":"undefined"}` {
		t.Fatal(result)
	}
}
