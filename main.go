package main

func main() {
	text := "Hello  World!"

	cleanText := cleanInputString(text)

	for _, word := range cleanText {
		println(word)
	}

}
