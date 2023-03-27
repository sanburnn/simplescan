package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	// var filePath *walk.LineEdit
	mw, err := walk.NewMainWindow()
	if err != nil {
		fmt.Println(err)
		return
	} // Create the UI controls
	resultLabel, err := walk.NewLabel(mw)
	if err != nil {
		fmt.Println(err)
		return
	}
	resultLabel.SetText("")
	MainWindow{
		Title:   "ClamAV Scanner",
		MinSize: Size{Width: 100, Height: 200},
		Layout:  VBox{},
		Children: []Widget{
			Label{
				Text: "Enter file path:",
			},
			// LineEdit{
			// 	AssignTo: &filePath,
			// },
			// PushButton{
			// 	OnClicked: func() {
			// 		filename := selectFile()
			// 		if filename != "" {
			// 			resultLabel.SetText("Scanning " + filename + "...")
			// 			scanFile(filename, resultLabel)
			// 		}
			// 	},
			// },
			PushButton{
				Text: "Scan",
				OnClicked: func() {
					// path := filePath.Text()
					// cmd := exec.Command("clamscan", path)
					// output, err := cmd.CombinedOutput()
					// if err != nil {
					// 	fmt.Println(err)
					// 	return
					// }
					// walk.MsgBox(nil, "Scan Result", string(output), walk.MsgBoxOK|walk.MsgBoxIconInformation)
					filename := selectFile(mw)
					if filename != "" {
						resultLabel.SetText("Scanning " + filename + "...")
						scanFile(filename, resultLabel)
					}
				},
			},
		},
	}.Run()
}
func selectFile(parent walk.Form) string {
	dlg := new(walk.FileDialog)
	dlg.Title = "Select file"
	dlg.Filter = "All Files (*.*)|*.*"
	dlg.ShowReadOnlyCB = true
	// dlg.Multiselect = false
	// dlg.Parent = parent

	if ok, err := dlg.ShowOpen(parent); err != nil {
		fmt.Println(err)
	} else if !ok {
		return ""
	} else {
		return dlg.FilePath
	}

	return ""
}
func scanFile(filename string, resultLabel *walk.Label) {
	// scannerPath := filepath.Join(os.Getenv("Program Files (x86)"), "ClamAV", "bin", "clamscan.exe")

	cmd := exec.Command("C:\\Program Files (x86)\\ClamAV\\clamscan.exe", "-v", "-i", "-f", filename)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		fmt.Println("error gak tau")
		return
	}

	if strings.Contains(string(output), "OK") {
		resultLabel.SetText(filename + " is clean!")
	} else {
		resultLabel.SetText(filename + " is infected!")
	}
}
