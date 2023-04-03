package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	var progressBar *walk.ProgressBar
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
			ProgressBar{
				AssignTo: &progressBar,
				MaxValue: 100,
			},
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

					filename := selectFile(mw)

					if filename != "" {

						resultLabel.SetText("Scanning " + filename + "...")
						// resultLabel.SetText("Scanning " + filename + "...")
						progressBar.SetValue(0)
						go func() {
							scanFile(filename, resultLabel, progressBar)
							progressBar.SetMarqueeMode(true)
						}()
					}
				},
			},
		},
	}.Run()
}

func selectFile(parent walk.Form) string {
	dlg := new(walk.FileDialog)
	dlg.Title = "Select Folder"
	dlg.Filter = "All Files (*.*)|*.*"
	dlg.ShowReadOnlyCB = true
	// dlg.Multiselect = false
	// dlg.Parent = parent
	if ok, err := dlg.ShowBrowseFolder(parent); err != nil {
		fmt.Println(err)
	} else if !ok {
		return ""
	} else {
		return dlg.FilePath
	}
	return ""
}
func scanFile(filename string, resultLabel *walk.Label, progressBar *walk.ProgressBar) {
	// scannerPath := filepath.Join(os.Getenv("Program Files (x86)"), "ClamAV", "bin", "clamscan.exe")

	cmd := exec.Command("C:\\Program Files\\ClamAV\\clamscan", "-r", filename)
	var out bytes.Buffer
	cmd.Stdout = &out

	go func() {
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
			fmt.Println("error gak tau")
			return
		}
		progressValue := 0.5
		fmt.Println(out.String())
		if strings.Contains(out.String(), "OK") {
			// resultLabel.SetText(filename + " is clean!")
			walk.MsgBox(nil, "Result", "is clean", walk.MsgBoxIconInformation)
			progressBar.SetValue(int(progressValue))
			progressBar.SetMarqueeMode(false)
		} else {
			walk.MsgBox(nil, "Result", "is infected!", walk.MsgBoxIconError)
			resultLabel.SetText(filename + " is infected!")
			progressBar.SetMarqueeMode(false)
		}
		resultLabel.Invalidate()
		progressBar.SetValue(100)
	}()

}
func showMessage(owner walk.Form, title, message string) {
	walk.MsgBox(owner, title, message, walk.MsgBoxIconError)
}
