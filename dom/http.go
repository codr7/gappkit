package dom

import (
	"fmt"
	"io"
	"strings"
)

func PostJSON(url string, body string, ok, fail string) string {
	var js strings.Builder
	fmt.Fprintf(&js, "fetch('%v', {\n", url)
	io.WriteString(&js, "method: 'POST',\n")
	io.WriteString(&js, "headers: {'Content-Type': 'application/json'},\n")
	io.WriteString(&js, "body: JSON.stringify(\n")
	io.WriteString(&js, body)
	io.WriteString(&js, "\n)}).then((response) => {\n")
	io.WriteString(&js, "if (response.ok) {\n")
	fmt.Fprintf(&js, "%v(response);\n", ok)
	io.WriteString(&js, "} else {\n")
	fmt.Fprintf(&js, "%v(response);\n", fail)
	io.WriteString(&js, "}\n")
	io.WriteString(&js, "});\n")
	return js.String()
}
