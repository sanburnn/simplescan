package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	var filePath *walk.LineEdit

	MainWindow{
		Title:   "ClamAV Scanner",
		MinSize: Size{Width: 300, Height: 200},
		Layout:  VBox{},
		Children: []Widget{
			Label{
				Text: "Enter file path:",
			},
			LineEdit{
				AssignTo: &filePath,
			},
			PushButton{
				Text: "Scan",
				OnClicked: func() {
					path := filePath.Text()
					cmd := exec.Command("clamscan", path)
					output, err := cmd.CombinedOutput()
					if err != nil {
						fmt.Println(err)
						return
					}
					walk.MsgBox(nil, "Scan Result", string(output), walk.MsgBoxOK|walk.MsgBoxIconInformation)
				},
			},
		},
	}.Run()
}
func scanFile(filename string, resultLabel *walk.Label) {
	scannerPath := filepath.Join(os.Getenv("ProgramFiles"), "ClamAV", "bin", "clamscan.exe")

	cmd := exec.Command(scannerPath, "-v", "-i", "-f", filename)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return
	}

	if strings.Contains(string(output), "OK") {
		resultLabel.SetText(filename + " is clean!")
	} else {
		resultLabel.SetText(filename + " is infected!")
	}
}
