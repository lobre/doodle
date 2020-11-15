// +build !embeds

package staticdir

import "net/http"

var FS http.FileSystem = http.Dir("./ui/static")
