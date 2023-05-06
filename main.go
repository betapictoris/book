package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/log"
	"github.com/taylorskalyo/goreader/epub"

	"github.com/charmbracelet/glamour"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

var text string

func main() {
	log.Info("Reading file...")

	reader, err := epub.OpenReader("books/aliceInWonderland.de.epub")
	if err != nil {
		log.Fatal("Failed to read", "error", err)
	}
	defer reader.Close() // Close the file when we are done with it
	book := reader.Rootfiles[0]

	log.Info("File read!", "title", book.Title)
	log.Info("Parsing file...")

	for _, item := range book.Manifest.Items {
		if item.MediaType == "application/xhtml+xml" {
			r, err := item.Open()
			if err != nil {
				log.Fatal("Failed to create item reader", "error", err)
			}
			cont, err := io.ReadAll(r)

			converter := md.NewConverter("", true, nil)
			content, err := converter.ConvertString(string(cont))
			if err != nil {
				log.Fatal("Failed to parse book content", "error", err)
			}

			text += content
		}
	}

	out, err := glamour.Render(text, "dark")
	if err != nil {
		log.Fatal("Failed to create glamour reader", "error", err)
	}

	fmt.Print(out)
}
