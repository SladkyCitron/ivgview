# ivgview

[![CI (Go)](https://github.com/SladkyCitron/ivgview/actions/workflows/ci.yml/badge.svg)](https://github.com/SladkyCitron/ivgview/actions/workflows/ci.yml) [![License](https://img.shields.io/github/license/SladkyCitron/ivgview)](LICENSE) ![Made in Slovakia](https://raw.githubusercontent.com/pedromxavier/flag-badges/refs/heads/main/badges/SK.svg)

**ivgview** is a small IconVG viewer and disassembler written in Go.
Open an IconVG file to inspect its rendered image or the instructions that make it up.

## ✨ Features

- Render `.ivg` and `.iconvg` files
- Open files from the command line, a file dialog, or by dragging them onto the window
- Zoom with the mouse wheel
- Toggle a transparency checkerboard
- View and scroll a textual IconVG disassembly

### Known Limitations

No support for FFV1 and FFV2 yet.

## 🤔 Usage

1. Install [Go](https://go.dev).
2. Run `go install github.com/SladkyCitron/ivgview@latest` to install ivgview.
3. Run `ivgview` to open the viewer. You can pass a file path as an argument to open a specific IconVG file, left click or press `Enter` to open a file, or drag and drop a file onto the window.

## 🎮 Controls

| Input                    | Action                               |
| ------------------------ | ------------------------------------ |
| Left click or `Enter`    | Open a file                          |
| Mouse wheel              | Zoom the image                       |
| `B`                      | Toggle the transparency checkerboard |
| `D`                      | Toggle disassembly                   |
| `↑` / `↓` or mouse wheel | Scroll disassembly                   |
| Drag and drop            | Open an IconVG file                  |

## ⚖️ License

MIT License (see [LICENSE](LICENSE))

## ❤️ Acknowledgements

- [IconVG](https://github.com/google/iconvg)
- [Ebitengine](https://ebitengine.org/)
