package collection

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/pafoster/qquiz/pkg/qquiz/card"
)

const extension = ".qq"

type Collection struct {
	// nu is the subset of cards missing a due date
	// nonDue is not currently used during Review()
	due, nu, nonDue []*card.Card
}

func New(dirs []string) (c *Collection, err error) {
	var files []string

	for _, d := range dirs {
		fs, err := ioutil.ReadDir(d)
		if err != nil {
			return nil, err
		}
		for _, f := range fs {
			if filepath.Ext(f.Name()) == extension {
				files = append(files, d+string(filepath.Separator)+f.Name())
			}
		}
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no %s files found", extension)
	}

	var cards []*card.Card
	for _, f := range files {
		card, err := card.New(f)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	due, nu, nonDue := partition(cards)

	rand.Seed(time.Now().UnixNano())
	for _, c := range [][]*card.Card{due, nu, nonDue} {
		shuf(c)
	}

	return &Collection{
		due:    due,
		nu:     nu,
		nonDue: nonDue,
	}, nil
}

func (c *Collection) Review(nMaxDue, nMaxNew int) (cards []*card.Card) {
	filter := func(c []*card.Card, nMax int) []*card.Card {
		if nMax < 0 || nMax > len(c) {
			nMax = len(c)
		}
		return c[:nMax]
	}

	return append(filter(c.due, nMaxDue), filter(c.nu, nMaxNew)...)
}

func partition(cards []*card.Card) (due, nu, nonDue []*card.Card) {

	now := time.Now()
	for _, c := range cards {
		if c.Due == nil {
			nu = append(nu, c)
		} else if c.Due.Before(now) {
			due = append(due, c)
		} else {
			nonDue = append(nonDue, c)
		}
	}

	return

}

func shuf(a []*card.Card) {
	rand.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
}
