package cmd

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func Test_CheckCMD(t *testing.T) {
	expected := "check       Run checks"
	b := bytes.NewBufferString("")
	checkCmd.SetOut(b)
	checkCmd.Execute()
	out, err := io.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), expected) {
		t.Fatalf("expected \"%s\" got \"%s\"", expected, string(out))
	}
}
