// +build !embeds

package htmldir

import "net/http"

var FS http.FileSystem = http.Dir("./ui/html")
