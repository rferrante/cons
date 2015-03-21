package cons

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Word struct {
	match     string // a regular expression for "text" to match
	kids      []*Word
	fn        func([]string)
	shortHelp string
	longHelp  string
}

func (parent *Word) AddCommand(match string, shortHelp string, longHelp string, fn func([]string)) *Word {
	var newWord Word
	newWord.match = match
	newWord.shortHelp = shortHelp
	newWord.longHelp = longHelp
	newWord.fn = fn
	(*parent).kids = append((*parent).kids, &newWord)
	return &newWord
}
func (parent *Word) RunLoop() {
	consolein := bufio.NewReader(os.Stdin)
	for {
		cmdline := strings.ToLower(get_cmd(consolein))
		err := do_command(cmdline)
		if err != nil {
			fmt.Print(err.Error())
		}
	}
}

type Token string

// when called on a Token, will return it's Command Word (if any), from words
func (tk Token) GetWord(words []*Word) (*Word, error) {
	// fmt.Printf("%d words\n", len(words))
	for _, word_ptr := range words {
		word := *word_ptr
		if Flags.Trace {
			fmt.Printf("testing %s against %s\n", string(tk), word.match)
		}
		if isMatch, err := regexp.MatchString(word.match, string(tk)); err != nil {
			return nil, err
		} else if isMatch == true {
			return word_ptr, nil
		}
	}
	return nil, nil
}

type TokenStack []*Token

// methods of TokenStack
func (q *TokenStack) Push(n *Token) {
	*q = append(*q, n)
}
func (q *TokenStack) Pop() (n *Token) {
	last_ix := q.Len() - 1
	n = (*q)[last_ix]
	*q = (*q)[:last_ix]
	return
}
func (q *TokenStack) Len() int {
	return len(*q)
}

type CommandSet struct {
	Top   Word
	Usage string
}

type FlagSet struct {
	Trace bool
}

// module members
var _tokens TokenStack
var _cmds CommandSet
var Flags FlagSet
var PromptFn func() string

//----

func NewRoot(usage string) *Word {
	_cmds = *(new(CommandSet))
	_cmds.Usage = usage
	return &_cmds.Top
}

func run(prev_word_ptr *Word) (func([]string), error) {
	prev_word := *prev_word_ptr
	if prev_word.fn != nil { // we're at the bottom
		return prev_word.fn, nil
	}
	if _tokens.Len() == 0 {
		return nil, fmt.Errorf("Incomplete command\n")
	}
	// check the next level down
	tk_ptr := _tokens.Pop()
	word_ptr, err := (*tk_ptr).GetWord(prev_word.kids)
	if err != nil {
		fmt.Println("error getting command")
		return nil, err
	}
	if word_ptr == nil { // token does not match
		return nil, fmt.Errorf("Bad token: %s\n", string(*tk_ptr))
	} else {
		return run(word_ptr)
	}
	return nil, err
}

func quit() {
	os.Exit(0)
}

func loadTokens(line string) {
	for ix := 0; ix < _tokens.Len(); ix++ { // clear _tokens
		_tokens.Pop()
	}
	fields := strings.Fields(line)
	for ix := len(fields) - 1; ix >= 0; ix-- {
		if strings.HasPrefix(fields[ix], "/") {
			setFlag(fields[ix])
			continue // discard. do not add to _tokens stack
		}
		tk := Token(fields[ix])
		_tokens.Push(&tk)
	}
}

func setFlag(flag string) {
	switch strings.ToLower(flag) {
	case "/t", "/trace":
		Flags.Trace = !Flags.Trace
	}
}

func ensureBaseCommands() {
	root := &_cmds.Top
	if !isTopCommandMatching("[Qq]|[Qq]uit") {
		root.AddCommand("[Qq]|[Qq]uit", "Exit the program", "", func([]string) { os.Exit(0) })
	}
	if !isTopCommandMatching("[Hh]elp") {
		root.AddCommand("[Hh]elp", "Get help", "", Help)
	}
}

func DisplayTokens() {
	var savedStack TokenStack
	for true {
		if _tokens.Len() <= 0 {
			_tokens = savedStack
			return
		}
		tk_ptr := _tokens.Pop()
		savedStack.Push(tk_ptr)
		fmt.Printf("Token = %s\n", *tk_ptr)
	}
}
func displayWords(w *Word, indent int) {
	fmt.Printf("%s%s\n", strings.Repeat("  ", indent), w.shortHelp)
	if len(w.longHelp) > 0 {
		fmt.Printf("%s%s\n", strings.Repeat("  ", indent), w.longHelp)
	}
	indent++
	for _, kid := range w.kids {
		displayWords(kid, indent)
	}
}

func isTopCommandMatching(match string) bool {
	for _, word := range _cmds.Top.kids {
		if word.match == match {
			return true
		}
	}
	return false
}

func Help([]string) {
	fmt.Printf("Usage: %s\n", _cmds.Usage)
	displayWords(&_cmds.Top, 1)
}

func NYI([]string) {
	fmt.Printf("This function is not yet implemented.\n")
}

func do_command(commandLine string) error {
	ensureBaseCommands()
	loadTokens(commandLine)
	if _tokens.Len() == 0 { // in case it was a stand-alone flag or something
		return nil
	}
	fn, err := run(&(_cmds.Top))
	if err != nil {
		return err
	}
	var args []string
	for _tokens.Len() > 0 {
		args = append(args, string(*_tokens.Pop()))
	}
	if Flags.Trace {
		fmt.Println(args)
	}
	fn(args)
	return nil
}

func default_prompter() string {
	return ">>"
}

func get_cmd(rdr *bufio.Reader) string {
	if PromptFn == nil {
		PromptFn = default_prompter
	}
	fmt.Println()
	fmt.Print(PromptFn())
	input, err := rdr.ReadString('\n')
	if err != nil {
		fmt.Print("error reading input\n")
		os.Exit(1)
	}
	return input
}
