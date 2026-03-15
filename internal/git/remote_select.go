package git

import (
	"bufio"
	"fmt"
	"os"
)

func SelectRemote() (string, error) {
	remotes, err := GetRemotes()
	if err != nil {
		return "", err
	}
	if len(remotes) == 0 {
		return "", fmt.Errorf("no git remotes found")
	}
	// 1️⃣ Prefer origin
	for _, r := range remotes {
		if r == "origin" {
			return "origin", nil
		}
	}
	// 2️⃣ Then upstream
	for _, r := range remotes {
		if r == "upstream" {
			return "upstream", nil
		}
	}
	// 3️⃣ If only one remote
	if len(remotes) == 1 {
		return remotes[0], nil
	}
	// 4️⃣ Ask user
	fmt.Println("Multiple remotes detected. Please select one:")
	for i, r := range remotes {
		fmt.Printf("%d) %s\n", i+1, r)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter number: ")
		var choice int
		_, err := fmt.Scanln(&choice)
		if err != nil || choice < 1 || choice > len(remotes) {
			fmt.Println("Invalid choice")
			reader.ReadString('\n')
			continue
		}
		return remotes[choice-1], nil
	}
}
