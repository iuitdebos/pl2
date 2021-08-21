<!-- PROJECT LOGO -->
<h1 align="center">PL2</h1>
<p align="center">
  Package for creating PL2 palette transform files
  <br />
  <br />
  <a href="https://github.com/OpenDiablo2/pl2/issues">Report Bug</a>
  ·
  <a href="https://github.com/OpenDiablo2/pl2/issues">Request Feature</a>
</p>

<!-- ABOUT THE PROJECT -->
## About

This package provides a PL2 palette transformation codec.

This package also contains command-line and graphical applications for working with PL2 files.

## Project Structure
* `pkg/` - This directory contains the core PL2 library.
    ```golang
   import (
        pl2 "github.com/OpenDiablo2/pl2/pkg"
  )
    ```
* `cmd/` - This directory contains command-line and graphical applications, each having their own sub-directory.
* `assets/` - This directory contains files, like the images displayed in this README, or test pl2 file data.

## Getting Started

### Prerequisites
You need to install [Go 1.16][golang], as well as set up your go environment.
In order to install the applications inside of `cmd/`, you will need to
make sure that `$GOBIN` is defined and points to a valid directory,
and this will also need to be added to your `$PATH` environment variable.
```shell
export GOBIN=$HOME/.gobin
mkdir -p $GOBIN
PATH=$PATH:$GOBIN
```

### Installation
As long as `$GOBIN` is defined and on your `$PATH`, you can build and install all apps inside of
`cmd/` by running these commands:

```shell
# clone the repo, enter the dir
git clone http://github.com/OpenDiablo2/pl2
cd pl2

# build and install inside of $GOBIN
go build ./cmd/...
go install ./cmd/...
```

At this point, you should be able to run the apps inside of `cmd/` from the command-line, like `pl2-to-gpl`.

<!-- CONTRIBUTING -->
## Contributing

I've set up all of the repos with a similar project structure. `~/pkg/` is where the actual
transcoder library is, and `~/cmd/` has subdirectories for each CLI/GUI application that can be
compiled.

Any contributions are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<!-- MARKDOWN LINKS & IMAGES -->
[dt1]: https://github.com/OpenDiablo2/dt1
[dc6]: https://github.com/OpenDiablo2/dc6
[dc6]: https://github.com/OpenDiablo2/dcc
[dat_palette]: https://github.com/OpenDiablo2/dat_palette
[ds1]: https://github.com/OpenDiablo2/ds1
[cof]: https://github.com/OpenDiablo2/cof
[golang]: https://golang.org/dl/