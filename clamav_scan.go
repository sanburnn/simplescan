// package main

// import (
// 	"fmt"
// 	"os"
// 	"os/exec"
// 	"path/filepath"
// 	"strings"

// 	"github.com/lxn/walk"
// )

// func main() {
// 	mw, err := walk.NewMainWindow()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	// Create a label to show the result of the scan
// 	resultLabel, err := walk.NewLabel(mw)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	resultLabel.SetText("")

// 	// Create a button to select a file to scan
// 	selectButton, err := walk.NewPushButton(mw)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	selectButton.SetText("Select file")
// 	selectButton.Clicked().Attach(func() {
// 		filename := selectFile(mw)
// 		if filename != "" {
// 			resultLabel.SetText("Scanning " + filename + "...")
// 			scanFile(filename, resultLabel)
// 		}
// 	})

// 	// Create a vertical box layout and add the controls to it
// 	vbox, err := walk.NewVBoxLayout()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	vbox.Add(resultLabel)
// 	vbox.Add(selectButton)

// 	// Set the layout of the main window
// 	mw.SetLayout(vbox)

// 	// Show the main window
// 	mw.Show()

// 	// Run the main event loop
// 	walk.Run()
// }

// func selectFile(parent walk.Form) string {
// 	dlg := new(walk.FileDialog)
// 	dlg.Title = "Select file"
// 	dlg.Filter = "All Files (*.*)|*.*"
// 	dlg.ShowReadOnlyCB = true
// 	// dlg.Multiselect = false
// 	// dlg.Parent = parent

// 	if ok, err := dlg.ShowOpen(parent); err != nil {
// 		fmt.Println(err)
// 	} else if !ok {
// 		return ""
// 	} else {
// 		return dlg.FilePath
// 	}

// 	return ""
// }

// func scanFile(filename string, resultLabel *walk.Label) {
// 	// Get the path to the ClamAV scanner executable
// 	scannerPath := filepath.Join(os.Getenv("ProgramFiles"), "ClamAV", "bin", "clamscan.exe")

// 	// Run the ClamAV scanner on the selected file
// 	cmd := exec.Command(scannerPath, "-v", "-i", "-f", filename)
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

//		// Set the result label text based on the scanner output
//		if strings.Contains(string(output), "OK") {
//			resultLabel.SetText(filename + " is clean!")
//		} else {
//			resultLabel.SetText(filename + " is infected!")
//		}
//	}
package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	// Create the main window
	var mainWindow *walk.MainWindow
	var filePathTE *walk.TextEdit
	var resultTE *walk.TextEdit

	MainWindow{
		AssignTo: &mainWindow,
		Title:    "ClamAV Scan",
		MinSize:  Size{Width: 400, Height: 200},
		Layout:   VBox{},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					TextLabel{
						Text: "File path to scan:",
					},
					TextEdit{
						AssignTo: &filePathTE,
					},
					PushButton{
						Text: "Scan",
						OnClicked: func() {
							go func() {
								filePath := filePathTE.Text()
								if _, err := os.Stat(filePath); err == nil {
									// Run the ClamAV scan command
									cmd := exec.Command("clamscan", "-r", "-i", filePath)
									out, err := cmd.Output()

									// Update the result text edit with the scan result
									resultTE.SetText(fmt.Sprintf("%s\n\n%s", out, err))
								} else {
									resultTE.SetText("File not found.")
								}
							}()
						},
					},
				},
			},
			TextEdit{
				AssignTo: &resultTE,
				ReadOnly: true,
				VScroll:  true,
			},
		},
	}.Run()
}
