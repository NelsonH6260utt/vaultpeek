package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/vaultpeek/internal/output"
)

func TestPrintSecrets_SortedKeys(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatText)

	data := map[string]interface{}{
		"zebra": "last",
		"alpha": "first",
		"middle": "second",
	}
	f.PrintSecrets("secret/test", data)

	result := buf.String()
	alphaIdx := strings.Index(result, "alpha")
	middleIdx := strings.Index(result, "middle")
	zebraIdx := strings.Index(result, "zebra")

	if alphaIdx == -1 || middleIdx == -1 || zebraIdx == -1 {
		t.Fatal("expected all keys to appear in output")
	}
	if !(alphaIdx < middleIdx && middleIdx < zebraIdx) {
		t.Errorf("expected keys in sorted order, got positions: alpha=%d middle=%d zebra=%d",
			alphaIdx, middleIdx, zebraIdx)
	}
}

func TestPrintSecrets_ContainsPath(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatText)
	f.PrintSecrets("secret/myapp/prod", map[string]interface{}{"key": "val"})

	if !strings.Contains(buf.String(), "secret/myapp/prod") {
		t.Error("expected output to contain the path")
	}
}

func TestPrintDiffLine_Symbols(t *testing.T) {
	tests := []struct {
		symbol string
		key    string
		value  interface{}
	}{
		{"+", "newkey", "newval"},
		{"-", "oldkey", "oldval"},
		{"~", "changedkey", "changedval"},
		{" ", "samekey", "sameval"},
	}

	for _, tt := range tests {
		var buf bytes.Buffer
		f := output.New(&buf, output.FormatText)
		f.PrintDiffLine(tt.symbol, tt.key, tt.value)
		if !strings.Contains(buf.String(), tt.key) {
			t.Errorf("expected key %q in output", tt.key)
		}
	}
}

func TestPrintError_ContainsMessage(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatText)
	f.PrintError(fmt.Errorf("something went wrong"))

	if !strings.Contains(buf.String(), "something went wrong") {
		t.Error("expected error message in output")
	}
}
