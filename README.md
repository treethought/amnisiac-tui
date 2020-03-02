# amnisiac

Stream music from reddit and souncloud

## Getting started

Requires [mpv](https://mpv.io/installation/)

Start mpv server

```
mpv --idle --input-ipc-server=/tmp/mpvsocket
```

Start the app

```

go run main.go

```

Keybindings:

- TAB - toggle focused pane
- Ctl-J - move selection down
- Ctl-K - move selection up
- ENTER - select item (fetch subreddit submissions or play result item)

