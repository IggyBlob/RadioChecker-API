package endpoint

import (
	"net/http"
	"fmt"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "RadioChecker API")
}
