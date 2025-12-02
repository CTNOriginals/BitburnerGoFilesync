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

- The only thing you are required to install is (TODO: Insert lates release link) [the latest release]() of this project.
- The only files that are required to be present on your system for this tool to work are: the executable and optionally a config file to change the default config fields.
- The executable can be anywhere you like, as long as you can execute it from the commandline to start it.
- Because this is written in golang, the tool is way ligher for your system to run, and way faster then javascript could ever be executed. The speed at which this tool can operate is only limited to how fast bitburner is able to respond back to it.

## More documentation coming very soon