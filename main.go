package main 

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	//Define the commands
	dirPtr := flag.String("dir", ".", "The directory to scan for code auditing")
	keyPtr := flag.String("key", "", "Your Gemini API Key (can also be set via GEMINI_API_KEY env variable)")

	//parse flag to the terminal 
	flag.Parse()


	// Retrieve values (deferencing pointers)
	dir := *dirPtr
	apiKey := *keyPtr

	//Validate flag values 

	if apiKey == "" {
		//If flag is empty, check the enviroment variables
		apiKey = os.Getenv("GEMINI_API_KEY") 
	}


	if apiKey == "" {
		fmt.Println("Error: Gemini API Key is required. Set it using the -key flag or GEMINI_API_KEY environment variable.") 
		flag.Usage() // Prints the default help message listing all flags
		os.Exit(1) //Exists the program with status code 1 indicating the error 
 	}



	//OutPut values to confirm they work
	fmt.Printf("Scanning directory: %s\n", dir)

	files, err := scanDirectory(dir)
	if err != nil {
		fmt.Printf("Error scanning directory: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Found %d source files to audit.\n", len(files))

		// Loop through files and print their paths (using %s and len)
	for i, file := range files {
		fmt.Printf("  [%d] File: %s (%d bytes)\n", i+1, file.Path, len(file.Content))
	}
	
	
}