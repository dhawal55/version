# version

Package 'version' returns a HTTP request multiplexer with a "/version"  pattern that return the application's version.

## Usage

Implement the Versioner interface and pass it to the New method. Add the returned handler to your request multiplexer.

```
package main

import (
    "log"
    "net/http"

    "github.com/dhawal55/version"
)

type Version struct{}

func (v *Version) GetVersion() string {
    return "1.0"
}

func main() {
    mux := http.NewServeMux()
    //Add application handlers
    //mux.Handle("/users", userHandler)

    v := &Version{}
    versionMux := version.New(v)
    mux.Handle("/", versionMux)

    log.Fatal(http.ListenAndServe(":8080", mux))
}
```
