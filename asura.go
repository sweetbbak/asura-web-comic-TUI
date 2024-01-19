package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/jessevdk/go-flags"
	fzf "github.com/ktr0731/go-fuzzyfinder"
)

var opts struct {
	Url       string `short:"u" long:"url" description:"url of specific manga/novel from helscans to retrieve"`
	OutputDir string `short:"o" long:"output" description:"directory to download the files to"`
	UserAgent string `short:"U" long:"user-agent" description:"user agent to use"`
	Latest    bool   `short:"l" long:"latest" description:"get the latest updated content"`
	Json      bool   `short:"j" long:"json" description:"output results in json (non-interactive)"`
	Search    string `short:"s" long:"search" description:"search sources for a query"`
	Verbose   bool   `short:"v" long:"verbose" description:"print debugging information and verbose output"`
}

var Debug = func(string, ...interface{}) {}

func docFromFile(file *os.File) (*goquery.Document, error) {
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func searchInput(prompt string) string {
	var search string
	huh.NewInput().
		Title(prompt).
		Value(&search).
		Run()

	return search
}

// this is dope but if the list is longer than the terminal it doesnt render lmaoooo
func list(args []string, title string) string {
	var opt string
	huh.NewSelect[string]().Title(title).Value(&opt).Options(huh.NewOptions(args...)...).Run()
	return opt
}

func listfzf(args []string) (string, error) {
	idx, err := fzf.Find(
		args,
		func(i int) string {
			return args[i]
		},
		fzf.WithPromptString("$ "),
	)
	if err == fzf.ErrAbort {
		return "", fmt.Errorf("No selection")
	}

	return args[idx], nil
}

func Menu() string {
	var opt string
	huh.NewSelect[string]().Title("Menu").Options(
		huh.NewOption("Latest", "latest"),
		huh.NewOption("Search [not implemented]", "search"),
		huh.NewOption("Exit", "exit"),
	).Value(&opt).Run()
	return opt
}

func Interactive() error {
	var mode string
	mode = Menu()
	switch mode {
	case "search":
		input := searchInput("Search: ")
		fmt.Println("searching for: ", input)
		var f fn
		f = func() {
			time.Sleep(time.Second * 2)
		}
		spinnerAction(f)
	case "latest":
		return frontPage()
	case "exit":
		return nil
	}
	return nil
}

type fn func()

func spinnerAction(f fn) error {
	err := spinner.New().
		Title("Getting manga...").
		Action(f).
		Run()

	if err != nil {
		return err
	}
	return nil
}

func frontPage() error {
	latest, err := getLatest()
	if err != nil {
		return err
	}

	var titles []string
	for title := range latest {
		titles = append(titles, title)
	}

	sel, err := listfzf(titles)
	if sel == "" || err != nil {
		return fmt.Errorf("No selection")
	}

	fmt.Printf("title [%v] - link [%v]\n", sel, latest[sel])
	// action := list([]string{"select chapter", "download all", "download range"}, "Choose an action")
	chapters, err := getChapterList(latest[sel])
	if err != nil {
		return err
	}

	var ch []string
	for c := range chapters {
		ch = append(ch, c)
	}

	if len(ch) == 0 {
		return fmt.Errorf("no found chapters")
	}

	ch, _ = sortChapters(ch)

	chsel, err := listfzf(ch)
	if err != nil {
		return err
	}

	images, err := getImages(chapters[chsel])
	if err != nil {
		return err
	}

	if len(images) == 0 {
		return fmt.Errorf("error getting chapter images")
	}

	if opts.OutputDir != "" {
		if !isExists(opts.OutputDir) {
			if err := createDir(opts.OutputDir); err != nil {
				return err
			}
		}
	}

	var dir string
	dir = clearString(sel)
	dir = titleCleanup(dir) // Hello, World! -> hello-world

	if opts.OutputDir != "" {
		dir = filepath.Join(opts.OutputDir, dir)
	}

	err = os.MkdirAll(dir, 0o755)
	if err != nil {
		return err
	}

	var f fn
	var imgs []string

	f = func() {
		imgs, err = Download(images, dir)
	}

	spinnerAction(f)
	if err != nil {
		return err
	}

	if err := makeCbz(chsel, chsel, imgs); err != nil {
		return err
	}

	return nil
}

func Asura(args []string) error {
	return Interactive()
}

func main() {
	args, err := flags.Parse(&opts)
	if flags.WroteHelp(err) {
		os.Exit(0)
	}

	if err != nil {
		log.Fatal(err)
	}

	if opts.Verbose {
		Debug = log.Printf
	}

	if err := Asura(args); err != nil {
		log.Fatal(err)
	}

}
