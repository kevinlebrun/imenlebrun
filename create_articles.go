package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gosimple/slug"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Header struct {
	Title string `toml:"title"`
	Date  string `toml:"date"`
	Link  string `toml:"link"`
}

func main() {
	f, err := os.Open("./scartch_mum2be")
	check(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var header []string
	var content []string

	var headerDone bool
	var headerStarted bool

	for scanner.Scan() {
		line := scanner.Text()
		if line == "+++" && headerDone {
			data := strings.Join(header, "\n")
			var h Header
			_, err := toml.Decode(data, &h)
			check(err)
			h.Title = strings.TrimSpace(h.Title)
			h.Date = strings.TrimSpace(h.Date)
			h.Link = strings.TrimSpace(h.Link)
			name := slug.Make(h.Title)
			f, err := os.Create(name + ".md")
			f.Write([]byte("+++\n"))
			err = toml.NewEncoder(f).Encode(h)
			check(err)
			f.Write([]byte("+++\n"))
			f.Write([]byte(strings.Join(content, "\n")))
			f.Close()
			header = make([]string, 0)
			content = make([]string, 0)
			headerDone = false
			headerStarted = true
			continue
		}
		if line == "+++" && headerStarted {
			headerDone = true
			continue
		}
		if line == "+++" && !headerStarted {
			headerStarted = true
			continue
		}
		if headerStarted && !headerDone {
			header = append(header, line)
		} else {
			content = append(content, line)
		}
	}

	check(scanner.Err())
}
