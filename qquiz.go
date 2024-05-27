package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pafoster/qquiz/internal/collection"
	"github.com/pafoster/qquiz/internal/ui"
)

func main() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "usage: %s [-c] [-d int] [-n int] directory ...\n", os.Args[0])
		flag.PrintDefaults()
	}

	nMaxNew := flag.Int("n", -1, "maximum number of new cards")
	nMaxDue := flag.Int("d", -1, "maximum number of due cards")
	f := flag.Lookup("n")
	f.DefValue = ": all new"
	f = flag.Lookup("d")
	f.DefValue = ": all due"

	nonInteractive := flag.Bool("c", false, "print number of cards we would review (if any) and exit")

	flag.Parse()
	dirs := flag.Args()

	os.Exit(run(dirs, *nMaxDue, *nMaxNew, *nonInteractive))
}

func run(dirs []string, nMaxDue, nMaxNew int, nonInteractive bool) int {
	if len(dirs) < 1 {
		flag.Usage()
		return 1
	}

	coll, err := collection.New(dirs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)

		return 1
	}

	cards := coll.Review(nMaxDue, nMaxNew)
	if nonInteractive {
		if len(cards) > 0 {
			fmt.Fprintf(os.Stdout, "%d cards for review\n", len(cards))
		}
		return 0
	}
	if len(cards) == 0 {
		fmt.Fprintf(os.Stderr, "No cards for review\n")
		return 0
	}

	ui := ui.New(cards)
	ui.Run()

	return 0
}
