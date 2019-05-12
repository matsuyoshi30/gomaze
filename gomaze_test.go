package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	filepath.Walk("./test", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".sh") {
			cmd := exec.Command("bash", filepath.Base(path))
			cmd.Dir = filepath.Dir(path)
			out, err := cmd.Output()
			if err != nil {
				t.Errorf("FAIL: " + err.Error())
			} else {
				outFile := strings.TrimSuffix(path, filepath.Ext(path)) + ".txt"
				expected, err := ioutil.ReadFile(outFile)
				if err != nil {
					t.Errorf("FAIL: Reading on output file: " + outFile)
				} else if strings.HasPrefix(string(out), strings.TrimSuffix(string(expected), "\n")) {
					t.Logf("PASS: " + path + "\n")
				} else {
					t.Errorf("FAIL: output differs: " + path + "\n")
				}
			}
		}
		return nil
	})
}
