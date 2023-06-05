package hclogr

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/require"
)

func TestLogr(t *testing.T) {
	var output bytes.Buffer

	log := hclog.New(&hclog.LoggerOptions{
		Level:  hclog.Info,
		Output: &output,
	})

	w := Wrap(log)

	r := require.New(t)

	r.True(w.Enabled())

	w.Info("check", "key", 314)

	r.Contains(output.String(), "key=314")

	w.Error(fmt.Errorf("err=889"), "bad stuff", "key", 234)

	r.Contains(output.String(), "key=234")
	r.Contains(output.String(), "err=889")

	w.V(4).Info("check", "key", 345)
	r.NotContains(output.String(), "key=345")

	w.V(3).Info("check", "key", 456)
	r.NotContains(output.String(), "key=456")

	w.V(2).Info("check", "key", 567)
	r.Contains(output.String(), "key=567")

	w.V(1).Info("check", "key", 678)
	r.Contains(output.String(), "key=678")

	w.V(0).Info("check", "key", 789)
	r.Contains(output.String(), "key=789")

	var o2 bytes.Buffer

	l2 := New(&hclog.LoggerOptions{
		Level:           hclog.Info,
		Output:          &o2,
		IncludeLocation: true,
	})

	l2.Info("includes location")

	r.Contains(o2.String(), "hclogr/hclogr_test.go")
}
