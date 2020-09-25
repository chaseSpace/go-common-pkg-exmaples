package exec

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestX(t *testing.T) {
	// win不需要加 ./
	winScriptPath := "windows_script.bat"
	unixScriptPath := "./unix_script.sh"

	log.Println(filepath.Abs(winScriptPath))

	arr := []string{winScriptPath, unixScriptPath}
	for _, st := range arr {
		log.Println(st)
		ext := strings.Split(st, ".")

		var cmd *exec.Cmd

		switch ext[len(ext)-1] {
		case "sh":
			bashPath := "D:\\Program Files\\Git\\bin\\bash.exe"
			cmd = exec.Command(bashPath, st)
		default:
			cmd = exec.Command(st)
		}

		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			t.Error(err)
		}
	}
}
