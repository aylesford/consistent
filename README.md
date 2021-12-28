# consistent
This package of `consistent` is the most concise and efficient algorithm of consistent hash based on golang.

# Example
Quick start:
```go
package main

import (
	"fmt"
	"github.com/aylesford/consistent"
)

func main() {
	c := consistent.New()

	c.Add("aa", "bb", "cc")

	fmt.Println(
		c.Get("AAAA"), 
		c.Get("BBBB"), 
		c.Get("CCCC"), 
		c.Get("DDDD"),
		c.Get("EEEE"),
		c.Get("FFFF"),
	)
}

```
