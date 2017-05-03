# movabletype

Package movabletype provides parsing "The Movable Type Import / Export Format".

## Link

  * [MovableType.org – Documentation: The Movable Type Import / Export Format](https://movabletype.org/documentation/appendices/import-export-format.html)
  * (in Japanese) [記事のインポートフォーマット : Movable Type 6 ドキュメント](https://www.movabletype.jp/documentation/mt6/tools/import-export-format.html)

## Sample

``` go
package main

import (
	"fmt"
	"os"

	"github.com/catatsuy/movabletype"
)

func main() {
	entries, _ := movabletype.Parse(os.Stdin)

	for _, e := range entries {
		fmt.Println(e.Title)
	}
}
```
