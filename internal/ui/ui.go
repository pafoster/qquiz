package ui

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"unicode"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/pafoster/qquiz/internal/card"
)

type UI struct {
	cards      []*card.Card
	currentPos int
	showAnswer bool
	app        *tview.Application
}

func New(cards []*card.Card) (ui UI) {
	ui = UI{
		cards: cards,
		app:   tview.NewApplication(),
	}
	return ui
}

func (ui *UI) Run() {
	tv := tview.NewTextView()
	frame := tview.NewFrame(tv).SetBorders(0, 0, 1, 1, 0, 0)

	update := func() {
		// Conditional return required, because after calling ui.app.Stop(), input
		// capturing logic will still call update()
		if ui.currentPos == len(ui.cards) {
			return
		}

		currentCard := ui.cards[ui.currentPos]
		currentKey := card.KEY_Q
		if ui.showAnswer {
			currentKey = card.KEY_A
		}

		qa := "Question"
		if ui.showAnswer {
			qa = "Answer"
		}

		newInfix := ""
		if currentCard.IsNew() {
			newInfix = "(NEW) "
		}
		title := fmt.Sprintf("Card %d of %d %s- %s ", ui.currentPos+1, len(ui.cards), newInfix, qa)

		toolText := "f: flip r: upgrade w: downgrade s: skip e: edit q: quit"
		if !ui.showAnswer {
			toolText = "f: flip s: skip e: edit q: quit"
		}
		frame.Clear().
			AddText(title, true, tview.AlignLeft, tcell.ColorGrey).
			AddText(toolText, false, tview.AlignLeft, tcell.ColorGrey)

		text := currentCard.GetFormatted(currentKey)
		tv.SetText(text)
	}

	update()

	ui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch keyPress := unicode.ToLower(event.Rune()); keyPress {
		case 'f':
			ui.showAnswer = !ui.showAnswer
			update()
			return nil
		case 'e':
			ui.editCurrentQuestion()
			update()
			return nil
		case 's':
			ui.nextQuestion(false, true)
			update()
			return nil
		case 'r', 'w':
			if !ui.showAnswer {
				return nil
			}
			ui.nextQuestion(keyPress == 'r', false)
			update()
			return nil
		case 'q':
			ui.app.Stop()
			return nil
		}
		if unicode.ToLower(event.Rune()) == ' ' { //tcell.KeyCtrlZ {
			// TODO Remapping CtrlZ appears to be insufficient
			// Note that the app doesn't respond as expected when we send
			// the signal manually.
			syscall.Kill(syscall.Getppid(), syscall.SIGTSTP)
			//ui.app.Suspend(func() {
			//})
		}
		return event
	})

	if err := ui.app.SetRoot(frame, true).Run(); err != nil {
		panic(err)
	}
}

func (ui *UI) editCurrentQuestion() {
	currentCard := ui.cards[ui.currentPos]
	editor := "vi"
	if v, ok := os.LookupEnv("EDITOR"); ok {
		editor = v
	}
	ui.app.Suspend(func() {
		cmd := exec.Command(editor, currentCard.Path)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			ui.app.Stop()
		}
	})
	if c, err := card.New(currentCard.Path); err == nil {
		currentCard = c
		ui.cards[ui.currentPos] = c
	}
}

func (ui *UI) nextQuestion(prevAnsIsCorrect, skip bool) {
	currentCard := ui.cards[ui.currentPos]
	if !skip || currentCard.Due == nil || currentCard.Reviewed == nil {
		if err := currentCard.Save(prevAnsIsCorrect); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			ui.app.Stop()
		}
	}

	ui.currentPos += 1
	if ui.showAnswer {
		ui.showAnswer = false
	}
	if ui.currentPos == len(ui.cards) {
		ui.app.Stop()
	}

}
