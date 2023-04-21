package main

/*
#cgo LDFLAGS: -lclamav
#include <stdlib.h>
#include <clamav.h>
*/

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	var progressBar *walk.ProgressBar
	// var filePath *walk.LineEdit
	// Initialize ClamAV
	// C.cl_init(C.CL_INIT_DEFAULT)

	// Load the virus signature database
	// C.cl_load(C.CString("/var/lib/clamav"))
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
					checkClamAVVersion()
					output, err := checkCurrentVersion("C:\\Program Files\\ClamAV\\clamscan", "-V")
					if err != nil {
						fmt.Println("Failed to execute command:", err)
						return
					}

					fmt.Println("Output of 'date' command:", output)

					// resultLabel.SetText("Scanning " + filename + "...")
					progressBar.SetValue(0)
					go func() {
						updateDatabase(resultLabel, progressBar)
						scanProcess(resultLabel, progressBar)
						// scanFile(filename, resultLabel, progressBar)
						progressBar.SetMarqueeMode(true)
					}()

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
func updateDatabase(resultLabel *walk.Label, progressBar *walk.ProgressBar) {
	cmd := exec.Command("C:\\Program Files\\ClamAV\\freshclam", "--show-progress")
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
		if strings.Contains(out.String(), "up-to-date") {
			// resultLabel.SetText(filename + " is clean!")
			walk.MsgBox(nil, "Result", "Database is Up to date", walk.MsgBoxIconInformation)
			progressBar.SetValue(int(progressValue))
			progressBar.SetMarqueeMode(false)
		} else {
			walk.MsgBox(nil, "Result", "Database is Updated", walk.MsgBoxIconError)
			// resultLabel.SetText(filename + " is infected!")
			progressBar.SetMarqueeMode(false)
		}
		resultLabel.Invalidate()
		progressBar.SetValue(100)
	}()
}
func scanProcess(resultLabel *walk.Label, progressBar *walk.ProgressBar) {
	cmd := exec.Command("C:\\Program Files\\ClamAV\\clamscan", "--memory=yes")
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
			// resultLabel.SetText(filename + " is infected!")
			progressBar.SetMarqueeMode(false)
		}
		resultLabel.Invalidate()
		progressBar.SetValue(100)
	}()
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

func checkCurrentVersion(command string, args ...string) (string, error) {
	// cmd := exec.Command("C:\\Program Files\\ClamAV\\clamscan", "-V")
	cmd := exec.Command(command, args...)
	outputBytes, err := cmd.Output()
	if err != nil {
		return "", err
	}

	output := string(outputBytes)
	fmt.Println(output)
	return output, nil
}

func checkClamAVVersion() {
	// Fetch the ClamAV download page HTML
	resp, err := http.Get("https://www.clamav.net/downloads")
	if err != nil {
		fmt.Println("Failed to fetch ClamAV download page:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read ClamAV download page response:", err)
		return
	}
	fmt.Println(body)
	// Search for the version number in the HTML using a regular expression
	value := body[0]
	parts := []string{
		strconv.Itoa(int(value / 100)),
		".",
		strconv.Itoa(int((value / 10) % 10)),
		".",
		strconv.Itoa(int(value % 10)),
	}
	result := strings.Join(parts, "")
	// re := regexp.MustCompile(`(\d)(\d)(\d)`)
	// matches := re.FindStringSubmatch(string(body))
	fmt.Println(result)
	if len(result) < 2 {
		fmt.Println("Failed to find ClamAV version number in download page")
		return
	}

	latestVersion := result
	//* PRINT THE RESULT OF VERSION
	// fmt.Println(latestVersion)
	// fmt.Println("Latest ClamAV version is", latestVersion)

	// Check if the installed version is out of date
	if isOutOfDate(installedVersion(), latestVersion) {
		fmt.Println("Your ClamAV version is out of date. Please update to the latest version.")
	} else {
		fmt.Println("Your ClamAV version is up to date.")
	}
}

func installedVersion() string {
	// TODO: implement a function that retrieves the currently installed version of ClamAV
	return "1.0.2" // replace with actual version number
}

func isOutOfDate(currentVersion string, latestVersion string) bool {
	// Extract version numbers from strings
	currentRe := regexp.MustCompile(`(\d+\.\d+\.\d+)`)
	currentMatch := currentRe.FindStringSubmatch(currentVersion)
	latestMatch := currentRe.FindStringSubmatch(latestVersion)

	if len(currentMatch) < 2 || len(latestMatch) < 2 {
		fmt.Println("Failed to extract version numbers")
		return false
	}

	// Compare version numbers
	for i := 1; i < 4; i++ {
		currentNum := strings.Split(currentMatch[1], ".")[i-1]
		latestNum := strings.Split(latestMatch[1], ".")[i-1]
		currentInt := 0
		latestInt := 0
		fmt.Sscanf(currentNum, "%d", &currentInt)
		fmt.Sscanf(latestNum, "%d", &latestInt)
		if latestInt > currentInt {
			return true
		} else if latestInt < currentInt {
			return false
		}
	}

	return false
}

//	func showMessage(owner walk.Form, title, message string) {
//		walk.MsgBox(owner, title, message, walk.MsgBoxIconError)
//	}
// func checkClamAVVersion() {
// 	// Fetch the latest version number from the ClamAV website
// 	resp, err := http.Get("https://www.clamav.net/downloads")
// 	if err != nil {
// 		fmt.Println("Failed to fetch ClamAV download page:", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println("Failed to read ClamAV download page response:", err)
// 		return
// 	}

// 	// Search for the version number in the HTML response
// 	versionStart := strings.Index(string(body), "ClamAV ")
// 	if versionStart < 0 {
// 		fmt.Println("Failed to find ClamAV version number in download page")
// 		return
// 	}
// 	versionStart += 7 // move past "ClamAV "

// 	versionEnd := strings.Index(string(body[versionStart:]), "<")
// 	if versionEnd < 0 {
// 		fmt.Println("Failed to find end of ClamAV version number in download page")
// 		return
// 	}
// 	versionEnd += versionStart

// 	version := string(body[versionStart:versionEnd])
// 	fmt.Println("Latest ClamAV version is", version)

// 	// Check if the installed version is out of date
// 	if installedVersion() < version {
// 		fmt.Println("Your ClamAV version is out of date. Please update to the latest version.")
// 	} else {
// 		fmt.Println("Your ClamAV version is up to date.")
// 	}
// }

// func installedVersion() string {
// 	// TODO: implement a function that retrieves the currently installed version of ClamAV
// 	return "1.0.1" // replace with actual version number
// }
