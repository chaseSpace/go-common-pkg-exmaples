package color

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"testing"
)

func TestColorStr(t *testing.T) {
	color.Red("Prints %s in red.", "text")
	color.Blue("Prints %s in blue.", "text")

	// Create a new color object
	c := color.New(color.FgCyan).Add(color.Underline)
	c.Println("Prints cyan text with an underline.")

	// Or just add them to New()
	d := color.New(color.FgCyan, color.Bold)
	d.Printf("This prints bold cyan %s\n", "too!.")

	// Mix up foreground and background colors, create new mixes!
	red := color.New(color.FgRed)

	boldRed := red.Add(color.Bold)
	boldRed.Println("This will print text in bold red.")

	whiteBackground := red.Add(color.BgWhite)
	whiteBackground.Println("Red text with white background.")

	// Use your own io.Writer output
	// NOTE: Fprintln won't work on windows terminal
	color.New(color.FgBlue).Fprintln(os.Stdout, "blue color!")
	blue := color.New(color.FgBlue)
	fmt.Fprint(os.Stdout, "")
	blue.Fprint(os.Stdout, "This will print text in blue.\n")

	// Create a custom print function for convenience
	redfn := color.New(color.FgRed).PrintfFunc()
	redfn("\nWarning: ")
	redfn("Error: %s\n", fmt.Errorf("x"))

	// Mix up multiple attributes
	notice := color.New(color.Bold, color.FgGreen).PrintlnFunc()
	notice("Don't forget this...")

	// ...
}
