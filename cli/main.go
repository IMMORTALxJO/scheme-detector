package main

import (
	"encoding/json"
	"fmt"

	schemedetector "github.com/IMMORTALxJO/scheme-detector"
)

func main() {
	schemas := schemedetector.FromEnv()
	b, _ := json.MarshalIndent(schemas, "", "  ")
	fmt.Println(string(b))
}
