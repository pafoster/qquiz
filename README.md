# qquiz - Spaced Repetition Learning in the Terminal 
qquiz (quick quiz) is a tool for **spaced repetition learning** (SRL), a method which helps the user memorise facts by showing flashcards to the user. Each flashcard specifies a quiz question and the correct answer, for example the question "On what date (New Style) was J.S. Bach born?" with the correct answer "31st March 1685". In SRL, a flashcard is shown to the user more or less frequently, depending on a) the user's ability to recall the correct answer to the question and b) how recently the flashcard was introduced. In qquiz, performance is evaluated based on self-assessment, i.e. the user may manually 'upgrade' or 'downgrade' a flashcard, after being presented with the correct answer. Upgrading means showing a flashcard less frequently; conversely downgrading means showing a flashcard more frequently.

qquiz has a terminal user interface (TUI) and a minimalist flavour. Each flashcard is a human-readable `.qq` file, with all metadata (reviewed and due dates) contained in the file. No additional files are needed for maintaining state.

# Similar Software
* [Anki](https://apps.ankiweb.net/) (more features, GUI-based)
* [Vocage](https://github.com/proycon/vocage) (self-described minimalistic, TUI-based)

(in the sense of [Anki](https://apps.ankiweb.net/)), but for the **terminal** and with a **minimalist flavour**.

# Requirements
* Go 1.20
* [tview](https://github.com/rivo/tview)

# Installation
* To build and install the `qquiz` executable, you can simply use `go install github.com/pafoster/qquiz`. This will download the source (along with dependencies) and build the executable. You should end up with a `qquiz` executable in `$GOBIN` (defaults to `~/go/bin`).
