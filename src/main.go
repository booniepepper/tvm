package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("tvm")

	hello := widget.NewLabel("Hello Fyne!")

	versions := widget.NewRadioGroup(nil, func(value string) {
		fmt.Println("version set to", value)
	})

	stdout, stderr, err := rtx("plugin", "ls")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: ", stderr)
		return
	}
	options := strings.Split(strings.Trim(stdout, "\n"), "\n")
	plugins := widget.NewRadioGroup(options, selectPlugin(versions))

	tabs := container.NewAppTabs(
		container.NewTabItem("Plugins", container.NewHBox(plugins, container.NewVScroll(versions), container.NewVBox(
			widget.NewButton("set local", func() {
				fmt.Printf("versions.Selected: %v\n", versions.Selected)
				use := strings.Join(strings.Split(versions.Selected, " ")[:2], "@")
				_, stderr, err := rtx("use", use)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: ", stderr)
					return
				}
				selectPlugin(versions)(plugins.Selected)
			}),
			widget.NewButton("set global", func() {
				fmt.Printf("versions.Selected: %v\n", versions.Selected)
				use := strings.Join(strings.Split(versions.Selected, " ")[:2], "@")
				_, stderr, err := rtx("use", "--global", use)
				if err != nil {
					fmt.Fprintf(os.Stderr, "ERROR: ", stderr)
					return
				}
				selectPlugin(versions)(plugins.Selected)
			}),
		))),
		container.NewTabItem("Languages", widget.NewLabel("World!")),
		container.NewTabItem("Misc", container.NewVBox(
			hello,
			widget.NewButton("Hi!", func() {
				hello.SetText("Welcome :)")
			}),
		)),
	)

	w.SetContent(tabs)

	w.ShowAndRun()
}

func selectPlugin(versions *widget.RadioGroup) func(string) {
	return func(plugin string) {
		var stdout string
		var stderr string
		var err error
		if plugin != "" {
			stdout, stderr, err = rtx("ls", plugin)
		} else {
			stdout, stderr, err = rtx("ls")
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: ", stderr)
			return
		}
		versions.Options = strings.Split(strings.Trim(stdout, "\n"), "\n")
		versions.SetSelected("")
		versions.Refresh()
	}
}

func rtx(arg ...string) (string, string, error) {
	cmd := exec.Command("rtx", arg...)
	var stdout strings.Builder
	var stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}
