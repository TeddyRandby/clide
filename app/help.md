# *C L I D E*
Clide is the Command Line Interface for Developer Experience!

## Problem
- Is your codebase littered with useful bash scripts and yarn commands?
- Do they all require passing parameters in a different way?
- Do some contributers find them difficult to find, use, or update?

## Why Clide?
- Organization - The convenience of file based routing .. for the CLI!
- Discoverability - Navigate and search through your scripts with ease!

## Examples
At the top level of your git repository, add a new folder:
```
─┐
 ├─.git
 └─.clide
```
The `.clide` folder is the root of all your command routes.

### Hello World
Lets do a traditional hello world! Simply add a script in the `.clide` directory. 
```
─┐
 ├─.git
 └┬.clide
  └─hello.sh
```
And `hello.sh`:
```bash
#!/usr/bin/env sh

echo "Hello World!"
```
Now, run `clide`. You should be prompted with a menu where you can select `hello`. Press enter!

### Passing Paramters
Often in little bash scripts its useful to have the user provide some input. Clide has this built in!

If you've ever built a backend with file system based routing, this will look very familiar.
```
─┐
 ├─.git
 └┬.clide
  ├─hello.sh
  └┬[Person]
   └─say_hello.sh
```
We've added a directory `[Person]` (brackets included), and thrown a new little script under it:
```bash
#!/usr/bin/env sh

echo "Hello $person!"
```
Similar to tools like *NEXT* or *Sveltekit*, the brackets in the path let Clide know that the scripts under it need an argument with that name. No matter how you capitalize the argument, clide will always supply the script with a variable in all lowercase. Now run `clide`. You should be prompted to provide a value for the argument `person`. Type in a value and press enter!

It can also be useful to have default values for these sorts of text inputs. To do this, add a bash script with the same name as that parameter (no extension). The output of this script will be piped directly into the paramter.


### Passing Select Arguments
Sometimes instead of having the user type in a string, you want to give them an option from a list. Clide has that too! Lets see it in action.
```
─┐
 ├─.git
 └┬.clide
  ├─hello.sh
  ├┬[Person]
  │└─say_hello.sh
  └┬{Person}
   ├─person
   └─say_hola.sh
```
The script `say_hola.sh` is almost identical (because the argument has the same name, but our `[]` around the argument in the path are now `{}`. You could even just use a symlink! We do have a new file adjacent to `say_hola.sh`. Its called `person` and it has no extension. This is the script that Clide will use to find the choices for the list. Ours will look simple.
```bash
#!/usr/bin/env sh

printf "Alice:A friend\nBob:Another friend"
```
Each line that this script outputs will be an option in the list - everything before the `:` is the name, and everything after is the description.

Now, run Clide. Select `say_hola` and pick a friend!

#### Multi-select
To turn a select into a multi-select, add a '+' at the end of the argument name.
```
─┐
 ├─.git
 └┬.clide
  ├─hello.sh
  ├┬[Person]
  │└─say_hello.sh
  ├┬{Person}
  │├─person
  │└─say_hola.sh
  └┬{People+}
   ├─people
   └─greet.sh
```
```bash
#!/usr/bin/env sh
# people

printf "Alice:A friend:Alice Smith\nBob:Another friend:Bob Smith\n"
```
```bash
#!/usr/bin/bash
# greet.sh

printf "Greetings to:\n%s" "$people"
```
Now, the user can select multiple values with `<space>`, and proceed with `<enter>`.

### Shortcuts
For the following shortcuts, we've extended our example file structure:
```
─┐
 ├─.git
 └┬.clide
  ├─hello.sh
  ├┬AnimalFriends
  │└─Pet_Dog.sh
  ├┬[Person]
  │└─say_hello.sh
  └┬{Person}
   ├─person
   └─say_hola.sh
```
Navigating around the menus is great, but sometimes you know the command you want to run ahead of time. Clide provies several ways to shortcut your commands:
 - To shortcut commands, you can use any prefix of the command/module. Eg: `clide a p`
 - To shortcut commands, you may also use any uppercase letters in the command/module as a shortcut. Eg: `clide af pd`
 - To shortcut arguments, use a `-` followed by the uppercase letters in the argument name. Eg: `clide -p=Alice say_hello`

Its important to note that although you define these shortcuts by using uppercase letters, clide only ever shortcuts or passes arguments via lowercase letters.

### Builtins
Clide has support for several builtin commands. They are all prefixed with `@`.
- `clide @ls`: Print a list of all available commands to stdout.
- `clide @help`: Print this help text to stdout.

Eg: List all clide commands, filter for commands with 'hello', and execute the last one.
`clide $(clide @ls | grep hello | awk -F\t 'END { print $3 }')`

### Dependencies
None!

### Installation
Via `go install`:
```
go install github.com/TeddyRandby/clide
```

### Inspiration and references
 - [Charm](https://charm.sh/)
 - [Bubbletea](https://github.com/charmbracelet/bubbletea)
 - [Glow](https://github.com/charmbracelet/glow)
