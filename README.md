# WarCraft III CustomKeys.txt Generator

## How to generate CustomKeys.txt

1. Modify `CustomKeys.txt`: change the `Hotkey=?` entries to desired hotkeys
1. Run `go build && ./wc3-custom-keys CustomKeys.txt CustomKeysGen.txt`
1. `CustomKeysGen` will now have properly formatted hokeys, tooltips, etc.

## Usage

```shell-session
./wc3-custom-keys <input-file|default:stdin> <output-file|default:stdout>
```

