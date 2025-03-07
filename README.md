# WarCraft III CustomKeys.txt Generator

## Prerequisites

- [Go 1.18](https://go.dev/doc/install)

## How to generate CustomKeys.txt

1. Change the `Hotkey=` entries to the desired hotkeys in
   [`data/CustomKeysDefault.txt`](data/CustomKeysDefault.txt)
1. Generate the keys:
   ```shell-session
   go run . data/CustomKeysDefault.txt data/CustomKeys.txt
   ```
1. Copy the generated [`data/CustomKeys.txt`](data/CustomKeys.txt) to the
   location specified in the game
   (`%UserProfile%\Documents\Warcraft III\CustomKeyBindings` on Windows)

## Usage

```shell-session
./wc3-custom-keys <input-file|default:stdin> <output-file|default:stdout>
```

