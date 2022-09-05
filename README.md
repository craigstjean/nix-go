# nix-go

---

## What is `nix-go`?

`nix-go` is a shell manager for setting up recreatable shell environments using [Nix](https://nix.dev/tutorials/install-nix)

## Why does `nix-go` exist?

To be honest, *you probably shouldn't use this*. This is the start of my journy into [NixOS](https://nixos.org/), which can be used to keep your dev environments pure and your builds reproducible.

To best use Nix, you need to learn the Nix language (e.g. from [nix.dev](https://nix.dev/tutorials/declarative-and-reproducible-developer-environments). There are also tools out there like:

- [home-manager](https://github.com/nix-community/home-manager)
- [lorri](https://github.com/nix-community/lorri/)
- probably [plenty of others](https://nix-community.github.io/awesome-nix/)

This project is essentially a stop-gap until I can consume more documentation and further myself into Nix.

## Using `nix-go`

```
USAGE:
   nix-go [global options] command [command options] [arguments...]

COMMANDS:
   list, ls, l              list existing projects/environments
   new                      create new project/environment
   list-packages, lp        list package in project/environment
   add-package, ap          add package to project/environment
   remove-package, rp       removes package from project/environment
   delete, del, remove, rm  delete project/environment
   shell, run, go           start project/environment
   help, h                  Shows a list of commands or help for one command
```

### Example - C++
```zsh C++
nix-go new cpp
nix-go ap cpp clang_14 clang-tools_14 clang-analyzer ninja autogen autoconf automake
nix-go shell cpp
```

### Example - D
```zsh D
nix-go new dlang
nix-go ap dlang dmd ldc dub
nix-go shell dlang
```

### Example - .NET
```zsh .NET
nix-go new dotnet
nix-go ap dotnet dotnet-sdk
nix-go shell dotnet
```

### Example - Elixir
```zsh Elixir
nix-go new elixir
nix-go ap elixir erlang elixir inotify-tools
nix-go shell elixir
```

### Example - LISP
```zsh LISP
nix-go new lisp
nix-go ap lisp sbcl lispPackages.quicklisp
nix-go shell lisp
```

### Example - Ruby
```zsh Ruby
nix-go new ruby
nix-go ap ruby ruby_3_1
nix-go shell ruby
```

### Example - `nix-go` dev 
```zsh
nix-go new --path ~/Projects/nix-go nix-go
nix-go ap nix-go git gcc go_1_18
nix-go shell nix-go
```

### `NIX_ENV`

`nix-go` starts the Nix shell with an environment variable `NIX_ENV` set, to allow echoing the current shell name in bash/zsh/etc.

An example of how to use this is to add to .zshrc (if you use [oh-my-zsh](https://ohmyz.sh/) with the `candy` theme):

```zsh
if [[ -n "$NIX_ENV" ]]; then
    OPROMPT="$PROMPT"
    PROMPT=$(echo $OPROMPT | sed 's/\(->%{\$fg_bold\[blue\]%} \)%/\1nix-'$NIX_ENV' %/')
fi
```

## Building `nix-go`

```
nix-shell -p git gcc
git clone https://github.com/craigstjean/nix-go.git
cd nix-go
go build
```

