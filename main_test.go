package main

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func TestMainBadCmd(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	os.Args = []string{"cloudstate", "bad"}
	main()
	log.SetOutput(os.Stdout)
	if !strings.Contains(buf.String(), "Unknown command: bad") {
		t.Errorf("Expected 'Unknown command: bad' but got '%s'", buf.String())
	}
}
