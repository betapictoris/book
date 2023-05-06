# Book CLI [![Build](https://github.com/BetaPictoris/book/actions/workflows/build.yml/badge.svg)](https://github.com/BetaPictoris/book/actions/workflows/build.yml)

Read epubs through the command line.

[![Book CLI](https://cdn.ozx.me/betapictoris/book.svg)](https://github.com/BetaPictoris/book)

## Installation

### From release

```bash
curl -LO https://github.com/BetaPictoris/book/releases/latest/download/book    # Download the latest binary.
sudo install -Dt /usr/local/bin -m 755 book                                    # Install Book CLI to "/usr/local/bin" with the mode "755"
```

### Build from source

#### Dependencies

You need Go (1.19+) installed to build this program. You can install it from your distro's repository using one of the following commands:

```bash
# Arch/Manjaro (and derivatives)
sudo pacman -Syu go

# Debian/Ubuntu (and derivatives)
sudo apt install golang-go
```

Alternatively, you can install it from [Go's official website](https://go.dev/doc/install).

Then, to build & install Book CLI run:

```bash
git clone git@github.com:BetaPictoris/book.git      # Clone the repository
cd book                                             # Change into the repository's directory
make                                                # Build Book CLI
sudo make install                                   # Install Book CLI to "/usr/local/bin" with the mode "755"
```

### User install

If you don't have access to `sudo` on your system you can install to your user's `~/.local/bin` directory with this command:

```bash
install -Dt ~/.local/bin -m 755 book
```

## Usage

```
book <path to file>
```

---

[![Beta Pictoris](https://cdn.ozx.me/betapictoris/header.svg)](https://github.com/BetaPictoris)
