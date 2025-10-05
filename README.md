# gradectl

*A lightweight CLI for tracking your self-taught learning progress.*

---

## 📘 Overview

**gradectl** is a simple, interactive command-line tool written in Go that helps
you track your grades while learning independently. It was born out of my own
self-taught learning journey. Working through books, exercises, tests, and
wanting an easy way to see how I’m really doing over time.

With `gradectl`, you can:

- Create and manage **classes** (e.g. Algebra, Calculus, etc.)
- Add **homework**, **quizzes**, and **tests** for each class
- View and update your **grades**
- Adjust **grade weights** (e.g. Tests – 50%, Homework – 30%, Quizzes – 20%)
- See your **overall grade** for each class
- Drill down to view grades for individual assignments

All of this happens in a clean, interactive CLI powered by the [`promptui`](https://github.com/manifoldco/promptui)
package, allowing you to use arrow keys to navigate options and enter grades
intuitively without having to memorize commands.

---

## 🚀 Installation

To install `gradectl` using Go:

```bash
go install github.com/Chance093/gradectl@latest
```

This will place the compiled binary in your `$GOPATH/bin` (or in `~/go/bin` if
you’re using default settings). Make sure it’s included in your system `PATH`.

Then simply run:

```bash
gradectl
```

and you’ll see the interactive menu.

---

## 💻 Running Locally

If you’d rather clone and run it locally (for development or inspection):

```bash
git clone https://github.com/Chance093/gradectl.git
cd gradectl
go run main.go
```

This will launch the same interactive CLI as the installed binary.

---

## 🧠 Why I Built It

While teaching myself math through textbooks and problem sets, I wanted a simple
way to measure my performance, without spreadsheets or heavy apps. `gradectl`
keeps things focused and minimal: no sign-ins, no forgotten spreadsheets, just a CLI 
to log your progress and see how you’re improving.

---

## 🛠️ Tech Stack

- **Language:** Go (Golang)
- **Database**: Sqlite
- **CLI Prompt Library:** [promptui](https://github.com/manifoldco/promptui)

---

## 📄 Example Use

Launch the program:

```bash
gradectl
```

You’ll be greeted with a menu that lets you:

- **Create** or **view** a class  
- Add **assignments** to that class  
- **Edit** or **delete** assignments
- View your **weighted average** grade  

It’s fully interactive, just use your arrow keys to navigate between options
and enter your responses when prompted.

---

## 📬 Feedback

If you have ideas or run into any issues, please open a GitHub issue on the
[gradectl repository](https://github.com/Chance093/gradectl).

---

### Keep learning and keep grading yourself!
