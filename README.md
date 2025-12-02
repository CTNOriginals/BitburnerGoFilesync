# Bitburner Go Filesync

This tool allows you to sync up your local files to Bitburner automatically without needing NodeJS or NPM installed at all.

Current automated triggers:
- On File Modify
- On File Create
- On File Delete

## Why use this tool?

### Resons not to use NPM

- The reason that I started this project is because I have gotten frustraighted with the NPM package manager and how many packages it installs to use one single thing on it.
- [NPM also keeps getting hacked](https://www.youtube.com/watch?v=M_W-dleZXCs) and it doesnt seem to get better anytime soon, which simply makes it unsafe to use.
- But overall, I just like having a tool that does one single thing whithout requiring to install ([at the time of writing](https://github.com/bitburner-official/typescript-template/blob/4136ab7107b75cf381acb21d1da7aa9e5e5b80fa/package-lock.json)) 185 dependencies...

### Reasons to use this instead

- The only thing you are required to install is [the latest release](https://github.com/CTNOriginals/BitburnerGoFilesync/releases) of this project.
- The only files that are required to be present on your system for this tool to work are: the executable and optionally a config file to change the default config fields.
- Because this is written in golang, the tool is way ligher for your system to run, and way faster then javascript could ever be executed. The speed at which this tool can operate is only limited to how fast bitburner is able to respond back to it.
- The executable can be anywhere you like, as long as you can execute it from the commandline to start it.


## How to install

1. Download the [latest release](https://github.com/CTNOriginals/BitburnerGoFilesync/releases) executable
2. Put the executable in same directory where you keep all your bitburner scripts
3. Open a commandline in the directory you put the executable in
4. Move on to **Usage**

You may skip step 2 and put the executable anywhere you like, but for the program to work you will then have to supply it with a `--dir` argument flag, more about that in **Arguments**.

## Usage

Once the executable is installed you are ready to get started with the game, the only thing you would still need to do is start the filesyncer via the commandline that you should have opened in step 3 of the instaliation guide.

Simply enter `BitburnerGoFilesync.exe` (or whatever the executable is called on your system) and pass in any arguments that you like.


## Arguments

This is the text that will pop up if you enter `--full-help` in the commandline arguments
```
Formatting Rules:
    Each new argument always has to start with a double dash '--'.
    If the argument does not start with '--' it is considered a parameter
    for the most recent argument that started with '--'.

    Each argument may have any number of parameters,
    to check what an argument may accept or require,
    you can do --help followed by the name of the argument without the '--'.

    Some arguments may accept a specific amount of parameters where others accept a range.
    If an argument doesnt have its required parameters, it will say say so in the console,
    this argument will not execute anything after that and will be ignored.
    If you pass in more parameters then an argument needs, it simply ignores the overflow.

--help, --wtf:
    Prints a list of arguments and their descriptions.
    Follow it up with another argument (without the -- before it)
    to get a more detailed explination about that argument.
  Parameters:
    command:
      The name (without the -- before it) of a command.
      Print a detailed explination about a specific command.

--full-help, --fhelp:
    The same as --help, but it also includes all of the extra information
    as if you entered --help <command> for each argument.

--dir:
    Specify the directory where this tool should watch
    for file changes to sync up with bitburner
  Parameters:
    dir:
      The path to the directory where you keep your bitburner scripts.
      Make sure to surround this parameter with double quotes (").

--include-ext, --ext:
    Specify which file extensions the file watcher should include.
    Default: js ts txt
  Parameters:
    extensions:
      Any number of file extensions seperated with spaces.
      Example: js ts json

DEBUG ARGUMENTS

--test:
    Runs the test function if it exists

--no-watcher:
    Prevents the program from watching file events.

--no-server:
    Prevents the program from creating a server and connecting to bitburner.
  Parameters:
    keep-alive:
      Accepts: true, false
      Usually when a server is ran, the program wont exit as it keeps evaluating it,
      if this parameter is set to true, the program will still be prevented from exiting.
```