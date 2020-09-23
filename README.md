# WIP PixelUI
**This repository is a proof of concept right now and is not very usable right now.**
Adds GUIs to the [Pixel](https://github.com/faiface/pixel) rendering engine by interpreting [imgui-go](https://github.com/inkyblackness/imgui-go) render data into Pixel render calls.

# Temporary Installation Instructions
There are a couple changes to Pixel that are required for PixelUI to work well; until they're accepted into the project, to use PixelUI you'll have to follow these steps.
1. `go get github.com/faiface/pixel` and `go get github.com/dusk125/pixelui`.
2. cd into the pixelui directory and run `git checkout -b gltriangles origin/gltriangles`; this will switch your PixelUI installation to use the gltriangles branch which relies on the following to build.
3. cd into the pixel directory and run `git remote add pixelui git@github.com:dusk125/pixel.git`; this will add my fork of the Pixel project that contains the in-progress changes to allow PixelUI to work efficiently.
4. Still in the pixel directory, run `git fetch` then `git checkout -b pixelui pixelui/pixelui`; these commands to fetch my forks content and change the Pixel content to pull from my fork (while still allowing importing from github.com/faiface/pixel in GoLang).

After those steps, you should be able to build a project using PixelUI, see the examples repository (look at the example on the gltriangles branch!).

# API
Since this is a work in progress, the API is likely to change at any time until we can fix the issues in Broken Things.

# [Examples](https://github.com/dusk125/pixelui_examples)
The [examples](https://github.com/dusk125/pixelui_examples) repository contains an example demonstrating PixelUI's functionality.

To run an example, navigate to it's directory, then go run the test.go file. For example:

```
$ cd pixelui_examples
$ go run test.go
```

## Current Expected state
This is where we are currently...
![Current State](https://github.com/dusk125/pixelui_examples/blob/master/screenshots/current_state.png)
