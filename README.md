# WIP PixelUI
**This repository is a proof of concept right now and is not very usable right now.**
Adds GUIs to the [Pixel](https://github.com/faiface/pixel) rendering engine by interpreting [imgui-go](https://github.com/inkyblackness/imgui-go) render data into Pixel render calls.

# Broken Things
* Clipping rectangles; the command has a clipRect that we're currently ignoring.
* Fix colors; I think this has to do with the imgui shader changing the Frag_Color alpha channel based on the uv of the texture.
* Fix text rendering; no idea what's going on here.
* Scrolling isn't working even though I'm passing the scroll through; it might just need to be scaled.
* Text input/general key handling

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
