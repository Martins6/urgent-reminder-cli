package display

import (
	"fmt"
	"os"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

type Display struct {
	noColor bool
}

func NewDisplay(noColor bool) *Display {
	if noColor {
		color.NoColor = true
	}
	return &Display{noColor: noColor}
}

func (d *Display) PrintBanner() {
	myFigure := figure.NewFigure("URGENT", "", true)
	myFigure.Print()
	figure2 := figure.NewFigure("REMINDERS", "", true)
	figure2.Print()
}

func (d *Display) PrintSeparator() {
	width := 60
	separator := strings.Repeat("=", width)
	fmt.Println(separator)
}

func (d *Display) PrintReminder(description string, daysRemaining int, alertEnabled bool) {
	var status string
	var dayText string

	if alertEnabled {
		status = color.New(color.FgCyan, color.Bold).Sprint("[ACTIVE]")
	} else {
		status = color.New(color.FgHiBlack).Sprint("[DONE]")
	}

	if daysRemaining < 0 {
		dayText = color.New(color.FgRed).SprintFunc()(fmt.Sprintf("Overdue by %d days", -daysRemaining))
	} else if daysRemaining == 0 {
		dayText = color.New(color.FgYellow).SprintFunc()("Due today!")
	} else if daysRemaining == 1 {
		dayText = color.New(color.FgYellow).SprintFunc()("Due in 1 day")
	} else {
		dayText = fmt.Sprintf("Due in %d days", daysRemaining)
	}

	fmt.Printf("%s %s -- %s\n", status, description, dayText)
}

func (d *Display) PrintSuccess(message string) {
	fmt.Println(color.GreenString(message))
}

func (d *Display) PrintError(message string) {
	fmt.Fprintln(os.Stderr, color.RedString(message))
}

func (d *Display) PrintInfo(message string) {
	fmt.Println(color.BlueString(message))
}

func (d *Display) PrintWarning(message string) {
	fmt.Println(color.YellowString(message))
}

func (d *Display) PrintEmpty() {
	fmt.Println()
}

func (d *Display) PrintHeader(text string) {
	d.PrintSeparator()
	fmt.Printf("%s\n", text)
	d.PrintSeparator()
}

func (d *Display) PrintSimpleReminder(id int, title, date, time string) {
	if time != "" {
		fmt.Printf("[%d] %s -- %s -- %s\n", id, title, date, time)
	} else {
		fmt.Printf("[%d] %s -- %s\n", id, title, date)
	}
}
