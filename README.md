# Bitburner Go Filesync

This tool allows you to sync up your local files to Bitburner automatically without needing NodeJS or NPM installed at all.

Current automated triggers:
- On File Modify
- On File Create
- On File Delete

## Why use this tool?

### Reasons not to use NPM

- The reason that I started this project is because I have gotten frustrated with the NPM package manager and how many packages it installs to use one single thing on it.
- [NPM also keeps getting hacked](https://www.youtube.com/watch?v=M_W-dleZXCs) and it doesn't seem to get better anytime soon, which simply makes it unsafe to use.
- But overall, I just like having a tool that does one single thing without requiring to install ([at the time of writing](https://github.com/bitburner-official/typescript-template/blob/4136ab7107b75cf381acb21d1da7aa9e5e5b80fa/package-lock.json)) 185 dependencies...

### Reasons to use this instead

- The only thing you are required to install is [the latest release](https://github.com/CTNOriginals/BitburnerGoFilesync/releases) of this project.
- The only files that are required to be present on your system for this tool to work are: the executable and optionally a config file to change the default config fields.
- Because this is written in golang, the tool is way lighter for your system to run, and way faster than javascript could ever be executed. The speed at which this tool can operate is only limited to how fast bitburner is able to respond back to it.
- The executable can be anywhere you like, as long as you can execute it from the commandline to start it.


## How to install

1. Download the [latest release](https://github.com/CTNOriginals/BitburnerGoFilesync/releases) file
    - windows: `BitburnerGoFilesync_win.exe`
    - linux: `BitburnerGoFilesync_linux`
    - mac: `BitburnerGoFilesync_mac`
2. Put the executable in same directory where you keep all your bitburner scripts
3. Open a commandline in the directory you put the executable in
4. Move on to [Usage](#usage)

You may skip step 2 and put the executable anywhere you like, but for the program target the correct files,
<br>you may need to specify youre preferred file path in the [config](#config),
<br>or pass in the `--dir` argument, more about that in [Arguments](#arguments).

Alternatively if you rather not download an executable from this repository (I would not blame you for being careful)
you may also clone this repository and [run/build it yourself](#how-to-build).

## Usage

Once the executable is installed you are ready to get started with the game, the only thing you would still need to do is start the filesyncer via the commandline that you should have opened in step 3 of the installation guide.

Simply enter `BitburnerGoFilesync.exe` (or whatever the executable is called on your system) and pass in any arguments that you like to run it with.

## Config

A config file will be created for you once you run the tool for the first time.

The file will be located in the same directory that you ran the tool in, you can change that with the `--config` argument (use `--help config` for more info).

The config file by default will be named `config.toml` and the initial content will be all possible fields for the config with their default values assigned to them.

### Config Fields

**Legend**:
- `[field]`: Group Section Header. This row's default column will instead be empty or hold relevant info.
- `- field`: Section child field.

| Field | Description | Default |
|:-----|:------------|:-------:|
| Port | Set the port for the server to connect to. | `"8080"` |
| Directory | Specify the directory where this tool should watch for file changes to sync up with bitburner. | `"./"` |
| FileScanInterval | The amount of miliseconds the file scanner waits each loop. | `1000` |
| [FilePatterns] | Holds include and exclude file pattern matching fields which allow you to define which files to sync and which not to. | [Pattern matching rules](https://github.com/bmatcuk/doublestar?tab=readme-ov-file#patterns) |
| - Include | Which files should be included. | `["**/*.js", "**/*.ts"]` |
| - Exclude | Which files to ignore. This is checked before the include patterns.  | `["**/*.d.ts"]` |
| [Logging] | Control what is and isnt logged |  |
| - NoColor | Prevent the logger from applying ansi color coding to the log prefix. | `false` |

## How to build

### Requirements
- [golang](https://go.dev/doc/install): To be able to run the project.
- (Optional) [wgo](https://github.com/bokwoon95/wgo): This is required for some Makefile targets.
<br>It is a tool that also watches the project for file changes and restarts it once anything is detected.

Once you have a clone of this repository and you installed the "*long list*" of requirements you can run the project in one of two ways:
1. Use `make [target]`.
<br>To see all possible targets and their descriptions, run `make help`.
2. Manually type out what you want to run.
<br>This would be something along the lines of `go run . --arg1 param --arg2 ...`.

If you are looking to contribute, please read the [Contribution Guidelines](https://github.com/CTNOriginals/BitburnerGoFilesync?tab=readme-ov-file).

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
    If an argument doesnt have its required parameters, it will say so in the console,
    this argument will not execute anything after that and will be ignored.
    If you pass in more parameters than an argument needs, it simply ignores the overflow.

--help, --wtf:
    Prints a list of arguments and their descriptions.
    Follow it up with another argument (without the -- before it)
    to get a more detailed explanation about that argument.
  Parameters:
    command:
      The name (without the -- before it) of a command.
      Print a detailed explanation about a specific command.

--full-help, --fhelp:
    The same as --help, but it also includes all of the extra information
    as if you entered --help <command> for each argument.

--config:
    Define the config.toml file path.
    By default, the config file is located in the same directory as the binary.
    If no config file exists at the specified location, one will be created.
  Parameters:
    filepath:
      The path to the config file
      Default: /home/ctn/code/bitburner/gofilesync/config.toml

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
      Any number of file extensions separated with spaces.
      Example: js ts json

--port:
    Set the port for the server to connect to.
    By default, the server will try to connect to 'localhost:8080'.
  Parameters:
    port:
      The port number.
      Default: 8080

--scan-interval, --interval:
    The amount of miliseconds the file scanner waits each loop.
    By default 100, if <= 0 it will skip the sleep function entirely.
  Parameters:
    interval:
      The interval in miliseconds
      Default: 100

--get-definitions:
    Requests the NetscriptDefinitions.d.ts file when a connection is established.
    The definitions file will be created in bitburners root directory.

DEBUG ARGUMENTS

--test, --debug:
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

