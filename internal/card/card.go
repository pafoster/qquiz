package card

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	minReviewInterval      = time.Hour * 6
	spacingFactor          = 2.0
	keySeparator      rune = ':'
	KEY_Q             rune = 'q'
	KEY_A             rune = 'a'
	key_r             rune = 'r'
	key_d             rune = 'd'
)

type Card struct {
	Path          string
	lines         []string
	keys          map[rune]int
	Reviewed, Due *time.Time
}

func New(path string) (c *Card, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read lines in file
	sc := bufio.NewScanner(f)
	lines := make([]string, 0)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	if err := sc.Err(); err != nil {
		return nil, err
	}

	// Determine first line of occurrence of each key
	keys := map[rune]int{}
	for n, l := range lines {
		if len(l) > 1 && rune(l[1]) == keySeparator {
			if _, ok := keys[rune(l[0])]; !ok {
				keys[rune(l[0])] = n
			}
		}
	}

	// File must at least have 'q' and 'a' keys
	for _, k := range []rune{KEY_Q, KEY_A} {
		if _, ok := keys[k]; !ok {
			return nil, fmt.Errorf("reading %s: missing required q/a keys", path)
		}
	}

	card := &Card{
		Path:  path,
		lines: lines,
		keys:  keys,
	}

	for _, k := range []rune{key_r, key_d} {
		if lineNo, ok := keys[k]; ok {
			s := strings.TrimSpace(lines[lineNo][2:])
			t, err := time.Parse(time.RFC3339, s)
			if err != nil {
				return nil, fmt.Errorf("reading %s: could not parse date for %c key", path, k)
			}
			if k == key_r {
				card.Reviewed = &t
			} else if k == key_d {
				card.Due = &t
			}
		}
	}

	return card, nil
}

func (c *Card) GetFormatted(key rune) string {
	var field []string

	for _, l := range c.lines {
		if len(l) > 1 && l[:2] == string(key)+string(keySeparator) {
			field = append(field, strings.TrimSpace(l[2:]))
		}
	}

	return strings.Join(field, "\n")
}

func (c *Card) IsNew() bool {
	return c.Reviewed == nil
}

func (c *Card) Save(answerIsCorrect bool) error {
	now := time.Now()
	due := now.Add(minReviewInterval)
	if c.Reviewed != nil && c.Due != nil && answerIsCorrect {
		due = now.Add(c.Due.Sub(*c.Reviewed) * spacingFactor)
	}
	c.Reviewed = &now
	c.Due = &due

	f, err := os.OpenFile(c.Path, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()

	for _, l := range c.lines {
		if len(l) > 1 && rune(l[1]) == keySeparator {
			if rune(l[0]) == key_r || rune(l[0]) == key_d {
				continue
			}
			fmt.Fprintf(w, "%s\n", l)
		}
	}

	fmt.Fprintf(w, "%s: %s\n", "r", c.Reviewed.Format(time.RFC3339))
	fmt.Fprintf(w, "%s: %s\n", "d", c.Due.Format(time.RFC3339))

	return nil
}
