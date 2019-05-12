package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	os.Chdir("./test")
	out, err := exec.Command("sh", "./default.sh").Output()
	if err != nil {
		t.Errorf("FAIL: " + err.Error())
	} else {
		outFile := "default.txt"
		expected, err := ioutil.ReadFile(outFile)
		if err != nil {
			t.Errorf("FAIL: Reading on output file: " + outFile)
		} else if strings.HasPrefix(string(out), strings.TrimSuffix(string(expected), "\n")) {
			t.Logf("PASS: " + "default.sh" + "\n")
		} else {
			t.Errorf("FAIL: output differs: " + "default.sh" + "\n")
		}
	}
}
