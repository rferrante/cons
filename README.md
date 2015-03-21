cons
====

Description
-----------

Cons is an interactive console command processing tool. That is, it will help you build console-based applications of the sort that give you a prompt for repeatedly typing commands into until you're done working with whatever your app does. Then you type "Quit" and you're back at the command line.

Its a simple and small package, and there are no dependencies other than the Go standard packages.

It comes along with some simple support for console colors, which are not necessary for building a working console app.
Installation
------------

This package can be installed with the go get command:

    go get github.com/rferrante/cons

Documentation
-------------

##NewRoot

    func NewRoot(usage string) *Word

The NewRoot function returns the root word of your commands. All of your top-level commands will be children of this root Word. This will usually be the first thing you do with cons, unless you have defines a custom PromptFn.

The *usage* string is a short description of the console program.

##AddCommand

    func (parent *Word) AddCommand(match string,
        shortHelp string, longHelp string, fn func([]string)) *Word

The root Word (or any Word) has command words descending from it. For instance, if your console program had a command called "list", the next word you type after "list" might be "connections". Then "list would be added under your root word, and "connections would be added under the "list" word. The AddCommand function returns the new Word, so that you can use it to add more commands under. For top-level commands, you would call this method on the Word returned by the NewRoot function.

The *match* string is a regexp matching whatever the user types. If the user types something the matches the regexp, this Word will be matched.

The two string arguments *shortHelp* and *longHelp* are a simple syntax statement and, if desired, some longer help text for the command.

The function argument *fn* is any function that takes a string slice and returns nothing. The string slice is any remaining argument the user types after this Word gets matched. Note that if the command requires further sub-commands, this argument should be nil. As the sequence of Words the user types is processed, once a Word is matched with a non-nil *fn*, the processing stops, and *fn* is called with the remainder of the command line coming along as the string slice.

##RunLoop

    func RunLoop()

When all command words have been set up, *RunLoop* starts the command processing loop.

##PromptFn

    func PromptFn() string

The default prompt is ">>", but you can change it to anything you want by setting PromptFn to a function of your own.

##Built-In Command Words

In addition to whatever you add, your program will also have the commands *Quit*, and *Help*, which do what you'd expect: exit your program, and print out top-level help. For now, you can't get help on individual commands.

##Example

    func main() {
        root := cons.NewRoot("Example program")
        listCmd := root.AddCommand("[Ll]ist","List things, "", nil)
        _ = listCmd.AddCommand("[Cc]onn.*", "List Connections", "", myFunc)

        root.RunLoop()
    }









