exiftool
======

Go library that simply shells out to `exiftool` and parses some fields.

Install
-----

`exiftool` is required to be in your `$PATH` env var.

```
go get github.com/bdotdub/exiftool
```

Example usage:
------

```go
package main

import (
  "os"
  "log"
  "fmt"

  "github.com/bdotdub/exiftool"
)

func main() {
	path := "sample1.jpg"

	exif, err := exiftool.DecodeFileAtPath(path)
	log.Println(exif.DateTimeOriginal)
	log.Println(exif.GPS.Latitude)
	log.Println(exif.GPS.Longitude)
}
```

Motivation
-------

I've tried a few native Go exif parsers and they have either all been
not robust enough, or have caused weird memory issues. Having used the
very mature `exiftool`, I thought why not just shell out to it ... sooo
here it is.

License
-------

MIT
