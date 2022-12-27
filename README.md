# Mondu

"mon-DI-you", my `du`-alike written in Go.

## Build

To build, run `make`. This will generate the `mondu` binary.

## Install

Install `mondu` by copying the compiled binary to your preferred `bin` directory (maybe `/usr/local/bin` is right for you).

## Run

Calculate the recursive size of the directory `path/to/my/stuff` by running `mondu path/to/my/stuff`. It will output a single number indicating the total number of bytes.

If you did not install `mondu` to a `bin` directory, then you can run `./mondu path/to/my/stuff`.

If you are seeing a lot of unwanted error output, then try running in quiet mode with `mondu -q path/to/my/stuff`.
