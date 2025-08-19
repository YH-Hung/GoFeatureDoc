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

	//delete(colors, "red")

	//fmt.Println(colors)
	printMap(colors)
}

func printMap(m map[string]string) {
	for k, v := range m {
		fmt.Println("Hex code for", k, "is", v)
	}
}
