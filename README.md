cons
====

Description
-----------

Cons is an interactive console command processing toolkit. That is, it will help you build console-based applications of the sort that give you a prompt for repeatedly typing commands into until you're done working with whatever your app does. Then you type "Quit" and you're back at the command line.

Its a simple and small package, and there are no dependencies other than the Go standard packages.

It comes along with some simple support for console colors, which are not necessary but might be handy.
Installation
------------

This package can be installed with the go get command:

    go get github.com/rferrante/cons

Documentation
-------------

####NewRoot

    func NewRoot(usage string) *Word

The NewRoot function returns the root word of your commands. All of your top-level commands will be children of this root Word. This will usually be the first thing you do with cons, unless you have defines a custom PromptFn.

The *usage* string is a short description of the console program.

####AddCommand

    func (parent *Word) AddCommand(match string,
        shortHelp string, longHelp string, fn func([]string)) *Word

The root Word (or any Word) has command words descending from it. For instance, if your console program had a command called "list", the next word you type after "list" might be "connections". Then "list would be added under your root word, and "connections would be added under the "list" word. The AddCommand function returns the new Word, so that you can use it to add more commands under. For top-level commands, you would call this method on the Word returned by the NewRoot function.

The *match* string is a regexp matching whatever the user types. If the user types something the matches the regexp, this Word will be matched.

The two string arguments *shortHelp* and *longHelp* are a simple syntax statement and, if desired, some longer help text for the command.

The function argument *fn* is any function that takes a string slice and returns nothing. The string slice is any remaining argument the user types after this Word gets matched. Note that if the command requires further sub-commands, this argument should be nil. As the sequence of Words the user types is processed, once a Word is matched with a non-nil *fn*, the processing stops, and *fn* is called with the remainder of the command line coming along as the string slice.

####RunLoop

    func RunLoop()

When all command words have been set up, *RunLoop* starts the command processing loop.

####PromptFn

    func PromptFn() string

The default prompt is ">>", but you can change it to anything you want by setting PromptFn to a function of your own.

####Built-In Command Words

In addition to whatever you add, your program will also have the commands *Quit*, and *Help*, which do what you'd expect: exit your program, and print out top-level help. For now, you can't get help on individual commands.

####Example

    func main() {
        root := cons.NewRoot("Example program")
        listCmd := root.AddCommand("[Ll]ist","List things, "", nil)
        _ = listCmd.AddCommand("[Cc]onn.*", "List Connections", "", myFunc)

        root.RunLoop()
    }

Color Support
-------------

Ansi terminal color support is provided by using a 'style' code. The code is a three-letter string, with the letters coding for foreground color, attribute, and background color in order. Colors can be one of the letters below:

 - w: white
 - r: red
 - g: green
 - b: blue
 - y: yellow
 - c: cyan
 - m: magenta
 - k: black

 Styles can be as follows:

 - b: bold
 - u: underline
 - v: inverted
 - i: italic (untested)
 - k: blinking (untested)
 - f: framed (untested)

 For example, the style code "bum" would be for blue text on a magenta background, and underlined.

 For any of the three positions you can skip setting them by using an underline, hyphen, or period instead of a code. For example, "b.." would be simply blue text on the default background and with no attribute.

In addition to the functions covered below, there are some very simple helpers for quickly specifying colored forground text. They are function variables so they can be used just like functions, and have the added benefit of being efficient since no style code decoding need be done. They are ShowRed, ShowBlue, ShowGreen, ShowYellow, ShowMagenta, and ShowCyan. Just use them like this: `ShowRed("Something Alarming!")`.

####Style

    func Style(style string, someText string) string

The Style function applies a style to some text. Styles are three-letter strings as described above. It returns the text with some ansi escape sequences around it, for printing on a terminal. Reset codes are included at the end so the terminal is left in the normal state.

    func StyleIf(style string, someText string, flag bool) string

The StyleIf function works just like Style, except that it will not style the text if flag is false.

####Example

    fmt.Printf("Here is an %s\n", cons.Style("rb", "Important Phrase"))










