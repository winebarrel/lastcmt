#!/bin/sh
ffmpeg -i "$1" -filter_complex "[0:v] fps=10,scale=1024:-1,split [a][b];[a] palettegen [p];[b][p] paletteuse=dither=none" -loop 0 demo.gif
