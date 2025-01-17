# Browser Images

[![Build Status](https://github.com/aerokube/images/workflows/build/badge.svg)](https://github.com/aerokube/images/actions?query=workflow%3Abuild)
[![Release](https://img.shields.io/github/release/aerokube/images.svg)](https://github.com/aerokube/images/releases/latest)

This repository contains [Docker](http://docker.com/) build files to be used for [Selenoid](http://github.com/aerokube/selenoid) and [Moon](http://github.com/aerokube/moon) projects. You can find prebuilt images [here](https://hub.docker.com/u/selenoid/).

## Download Statistics

### Firefox: [![Firefox Docker Pulls](https://img.shields.io/docker/pulls/selenoid/firefox.svg)](https://hub.docker.com/r/selenoid/firefox)

### Chrome: [![Chrome Docker Pulls](https://img.shields.io/docker/pulls/selenoid/chrome.svg)](https://hub.docker.com/r/selenoid/chrome)

### Opera: [![Opera Docker Pulls](https://img.shields.io/docker/pulls/selenoid/opera.svg)](https://hub.docker.com/r/selenoid/opera)

### Android: [![Android Docker Pulls](https://img.shields.io/docker/pulls/selenoid/android.svg)](https://hub.docker.com/r/selenoid/android)

## Building Images

Moved to: http://aerokube.com/images/latest/#_building_images

### Chrome

```bash
$ go install github.com/markbates/pkger/cmd/pkger@latest
$ go generate github.com/aerokube/images
$ go build
$ ./images chrome -b 132.0.6834.83 -d 132.0.6834.83 -t selenoid/chrome:132.0
```

Use `--platform=linux/amd64` to build images on M1 Mac in the
`static/chrome/Dockerfile` 

## Image information

Moved to: http://aerokube.com/images/latest/#_browser_image_information
