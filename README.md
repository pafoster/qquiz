# qquiz - Spaced Repetition Learning in the Terminal 
`qquiz` (quick quiz) is a tool for **spaced repetition learning** (SRL), a method which helps you memorise facts using flashcards. Each flashcard specifies a quiz question and its correct answer. For example, a flashcard might specify the question "In which year was J.S. Bach born?" and its correct answer "1685 CE". In SRL, a flashcard is reviewed more frequently or less frequently, depending on **a)** your ability to recall the correct answer and **b)** how recently the flashcard was introduced. In `qquiz`, performance is evaluated based on self-assessment, i.e. you may manually 'upgrade' or 'downgrade' a flashcard after being shown the correct answer. Upgrading means reviewing a flashcard less frequently; conversely downgrading means reviewing a flashcard more frequently.

`qquiz` has a terminal user interface (TUI) and a minimalist flavour. Each flashcard is a human-readable file with a `.qq` extension, with all metadata (reviewed and due dates) stored in the `.qq` file. No additional files are needed for maintaining state.

# Similar Software
* [Anki](https://apps.ankiweb.net/) (more features, GUI-based)
* [Vocage](https://github.com/proycon/vocage) (self-described minimalistic, TUI-based)

# Requirements
* Go 1.20
* [tview](https://github.com/rivo/tview)

# Installation
To build and install the `qquiz` executable, you can simply run `go install github.com/pafoster/qquiz`. This will download the source (along with dependencies) and build the executable. You should end up with a `qquiz` executable in `$GOBIN` (defaults to `~/go/bin`). If desired, `export PATH:$PATH:~/go/bin` in your `.profile`.

# Building
```
git clone 'https://github.com/pafoster/qquiz/'
cd qquiz
make build
```

# The Learning Scheme
`qquiz` implements the following (simple) learning scheme based on timestamps. The collection of flashcards is partitioned into three (disjoint) subsets:
* *new*: Flashcards which were never reviewed
* *due*: Flashcards whose due date is in the past
* *non-due*: Flashcards whose due date is in the future

Those flashcards designated as '*new*' and '*due*' are shuffled and reviewed when `qquiz` is run (subject to any user-specified limits, see below). The new due date for new and downgraded flashcards is 6 hours from the most recent time of review. The new due date for upgraded flashcards is 2.0 times the interval between the previous time of review and existing due date. Upon upgrading or downgrading a flashcard, the new due date and most recent time of review are written to the relevant `.qq` file as lines beginning with `d:` and `r:`, respectively.

# Tutorial
## Creating a Flashcard
Create a directory for storing your flashcards, e.g. `mkdir ~/flashcards`. To create a new flashcard, open a new `.qq` file in your editor, e.g. `vi ~/flashcards/jsbach.qq`. A minimal `.qq` file might look like this:
```
q: In which year was J.S. Bach born?
a: 1685 CE
```
That is, lines beginning with `q:` denote the question and lines beginning `a:` denote the answer. (For multi-line questions or multi-line answers, it's possible to have multiple `q:` and `a:` lines in a `.qq` file.)

## Reviewing Flashcards
Review all '*due*' and '*new*' flashcards in the directory `~/flashcards`:
```
qquiz ~/flashcards
```
Review a maximum of 8 '*due*' flashcards and a maximum of 2 '*new*' flashcards:
```
qquiz -d 8 -n 2 ~/flashcards
```
Review only '*due*' flashcards:
```
qquiz -n 0 ~/flashcards
```
Print how many flashcards we would review (if any) and exit. (Useful for inclusion in a cronjob):
```
qquiz -c ~/flashcards
```
You can invoke `qquiz` using multiple directories as positional arguments. This permits you to organise your flashcards into multiple subdirectories (collections), which you might decide to combine into one review session. Using shell expansion, you could write:
```
qquiz ~/flashcards/{machine_learning,statistics}
```
# Screenshot
![screenshot](screenshots/qquiz.png)

# Key Bindings
* `f` Flip between question and answer
* `r` Upgrade current flashcard, modifying `.qq` file (available when answer is displayed)
* `d` Downgrade current flashcard, modifying `.qq` file (available when answer is displayed)
* `s` Skip current flashcard without modifying `.qq` file
* `e` Open current flashcard in `$EDITOR` (useful for correcting typos or inaccuracies)
* Arrow keys, `PgUp`, `PgDn` Scroll text (useful if you want to include extensive notes in the answer)
* `q` Quit
