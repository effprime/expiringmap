# Expiringmap

Expiringmap is a clone of the Python library [Expiringdict](https://pypi.org/project/expiringdict/)

# Usage

```go
package main

import (
    "fmt"
    "github.com/effprime/expiringmap"
    "time"
)

func main() {
    // Create a new ExpiringMap with a lifespan of 5 minutes and a maximum length of 100 items
    emap := expiringmap.NewExpiringMap[int](expiringmap.Settings{
        Age:       5 * time.Minute,
        MaxLength: 100,
    })

    // Set a value
    emap.Set("key1", 123)

    // Retrieve a value
    if val, ok := emap.Get("key1"); ok {
        fmt.Println("Value:", val)
    }
}
```
