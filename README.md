# Maze written in Go [![circleci](https://circleci.com/gh/matsuyoshi30/gomaze.svg?style=shield&circle-token=99ca77ba57ac5cdd3276aade4335a04c5f68879a)](https://circleci.com/gh/matsuyoshi30/gomaze)

### Usage

```
$ go get -u github.com/matsuyoshi30/gomaze

$ gomaze
||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||
S           ||  ||                      ||  ||          ||  ||
||  ||||||  ||  ||  ||  ||||||||||  ||||||  ||  ||  ||||||  ||
||  ||  ||  ||      ||      ||  ||  ||  ||      ||      ||  ||
||||||  ||  ||  ||||||  ||||||  ||  ||  ||||||||||  ||  ||  ||
||  ||              ||  ||  ||          ||          ||  ||  ||
||  ||  ||  ||||||||||||||  ||||||  ||||||  ||  ||||||  ||  ||
||  ||  ||  ||                  ||          ||      ||  ||  ||
||  ||||||  ||||||  ||||||  ||||||  ||  ||  ||  ||  ||||||  ||
||          ||  ||  ||  ||  ||      ||  ||  ||  ||          ||
||  ||  ||||||  ||  ||  ||  ||  ||||||||||  ||||||||||||||  ||
||  ||          ||  ||  ||  ||  ||                  ||  ||  ||
||||||  ||||||  ||||||  ||  ||  ||||||||||||||||||||||  ||  ||
||      ||  ||          ||  ||  ||          ||  ||          ||
||||||||||  ||  ||||||||||  ||||||  ||||||  ||  ||  ||||||||||
||              ||  ||      ||      ||      ||      ||      ||
||||||  ||||||||||  ||  ||  ||  ||  ||||||||||  ||  ||  ||||||
||      ||  ||  ||      ||  ||  ||  ||          ||      ||  ||
||  ||||||  ||  ||||||  ||  ||  ||  ||  ||||||  ||||||  ||  ||
||  ||              ||  ||  ||  ||      ||  ||  ||  ||      ||
||  ||||||||||  ||||||  ||||||||||||||  ||  ||  ||  ||  ||||||
||      ||          ||              ||  ||  ||      ||  ||  ||
||  ||||||  ||||||||||  ||||||||||  ||  ||  ||  ||  ||||||  ||
||  ||              ||  ||                  ||  ||          ||
||  ||||||||||  ||||||  ||||||  ||||||||||||||||||  ||||||  ||
||      ||  ||  ||  ||      ||  ||  ||              ||      ||
||||||  ||  ||  ||  ||  ||||||  ||  ||||||||||  ||||||||||||||
||          ||          ||              ||  ||              ||
||  ||  ||  ||||||||||  ||||||||||||||||||  ||  ||||||  ||  ||
||  ||  ||      ||          ||                      ||  ||   G
||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||||
```

You can specify the size of maze with `--height` or `--width`.  

```
$ gomaze --height 10 --width 10
||||||||||||||||||||||
S       ||      ||  ||
||||||  ||  ||||||  ||
||      ||          ||
||  ||||||  ||  ||||||
||  ||      ||      ||
||  ||  ||  ||||||  ||
||      ||      ||  ||
||||||  ||  ||||||  ||
||      ||      ||   G
||||||||||||||||||||||
```

You can play the maze with `--screen` and use arrow keys.

![](./_resource/screen.gif)

You can see the maze solution with `--bfs` or `--dfs`.

![](./_resource/bfs.gif)