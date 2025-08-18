package main

import "fmt"

func main() {
	//var m map[string]int
	//emptyMap := make(map[string]int)

	colors := map[string]string{
		"red":  "ff0000",
		"blue": "0000ff",
	}

	colors["white"] = "ffffff"

	delete(colors, "red")

	fmt.Println(colors)
}
