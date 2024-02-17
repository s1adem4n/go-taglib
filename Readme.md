# go-taglib

A taglib wrapper for Go

Special thanks [wtolson](https://github.com/wtolson) for writing the original [go-taglib](https://github.com/wtolson/go-taglib).
I've used his code as a start and extended it with:
- non-standard tag getting and setting
- functions for getting the album artist
- functions for getting and setting pictures

Also, wtolson's code used a global mutex for every operation, but I removed it, because (from my testing) I couldn't find any issues while running this concurrently.

# Documentation
Check the documentation at [pkg.go.dev/github.com/s1adem4n/go-taglib](https://pkg.go.dev/github.com/s1adem4n/go-taglib)


# Dependencies
To use this library you must have the static [taglib](https://taglib.org) libraries installed. Please check your distribution's package manager for further instructions.
On MacOS, you may install this with

    brew install taglib


# Installing
    
    go get github.com/s1adem4n/go-taglib


# Basic example
```go
import (
  taglib "github.com/s1adem4n/go-taglib"
  "fmt"
)

func main() {
  file, err := taglib.Read("test.mp3")
  if err != nil {
    fmt.Println(err)
    return
  }
  defer file.Close()

  // getting a tag
  fmt.Println(file.Title())

  picture, err := file.Picture()
  if err != nil {
    // no picture found
    fmt.Println(err)
    return
  }

  fmt.Println(picture.MimeType)

  // setting a normal tag
  file.SetTitle("Test")
  // setting a non-standard tag
  file.SetTag("banana", "text")

  // don't forget saving
  err = file.Save()
  if err != nil {
    fmt.Println(err)
    return
  }
}
```
