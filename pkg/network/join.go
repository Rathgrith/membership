package network

import (
	// "fmt"
	"os"
)

// GetHostname returns the hostname of the machine.
func GetHostname() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return hostname, nil
}

// func main() {
// 	host, err := GetHostname()
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return
// 	}
// 	fmt.Println("Hostname:", host)
// }
