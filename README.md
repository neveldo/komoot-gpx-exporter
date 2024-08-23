# Go Komoot GPX exporter

## Description

At the moment, Komoot doesn't provide a way to download all the GPX tracks of a specific account.
This Go app can be used to export all the GPX files from your Komoot account into a directory.

## Requirements

- go >=1.22.4

## How to run

```
    go run cmd/export/main.go --user="123" --cookie="komoot_xhr_cookie_value" --dir="." --sport="hike"
    
    go run cmd/export/main.go --help
```

Parameters :
- user : komoot user ID (appears in the app URL)
- cookie : You can retrieve the required cookie value by browsing Komoot app via your web browser, open the debugger, and get the cookie value from a XHR request.
- dir : GPX destination directory

The app will download all your GPX tracks into the specified directory.
It will use up to 10 go routines in parallel but feel free to change this value within cmd/gpx_downloader/main.go file.

## Warning

This is just quick-made personal project to fulfil a personal usecase. Use it at your disposal.

