package tui

import (
	"fmt"

	"github.com/logrusorgru/aurora"
)

// GenerateASCIIArt prints the lighthouse onto the screen
func GenerateASCIIArt() {
	fmt.Print("\033[H\033[2J")

	// |	. _  .     __        __,--'
	// |	 (_)      /__\ __,--'
	// | '  .  '    | o|
	// |	  			 [IIII]`--.__
	// | 				  	|  |       `--.__
	// |	  				| :|             `--.__
	// | 	 		  		|  |                   `--.__
	// | ._,,.-,.__.'__`.___.,.,.-..,_.,.,.,-._..`--..-.,._.,,._,-,..,._..,.,_,,.

	fmt.Printf("  %s     %s        %s\n", aurora.Yellow(". _  ."), aurora.BrightRed("__"), aurora.White("__,--'"))
	fmt.Printf("   %s      %s %s\n", aurora.Yellow("(_)"), aurora.BrightWhite("/__\\"), aurora.White("__,--'"))
	fmt.Printf(" %s    %s\n", aurora.Yellow("'  .  '"), aurora.BrightRed("| o|"))
	fmt.Printf("           %s%s\n", aurora.BrightWhite("[IIII]"), aurora.White("`--.__"))
	fmt.Printf("            %s       %s\n", aurora.BrightRed("|  |"), aurora.White("`--.__"))
	fmt.Printf("            %s             %s\n", aurora.BrightWhite("| :|"), aurora.White("`--.__"))
	fmt.Printf("            %s                   %s\n", aurora.BrightRed("|  |"), aurora.White("`--.__"))
	fmt.Printf("%s%s%s\n", aurora.BrightBlue(" ._,,.-,.__."), aurora.BrightWhite("'__`"), aurora.BrightBlue(".___.,.,.-..,_.,.,.,-._..`--..-.,._.,,._,-,..,._..,.,_,,."))

	fmt.Println("")
	fmt.Printf("%s", aurora.Bold(aurora.White("Lighthouse v0.1")))
	fmt.Println(aurora.White(""))
}
