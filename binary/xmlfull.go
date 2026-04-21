package binary

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// XMLStringFull is like XMLString but outputs full hex for large byte content
// instead of truncating to "<!-- N bytes -->".
func (n *Node) XMLStringFull() string {
	content := n.contentStringFull()
	if len(content) == 0 {
		return fmt.Sprintf("<%[1]s%[2]s/>", n.Tag, n.attributeString())
	}
	newline := "\n"
	if len(content) == 1 || !IndentXML {
		newline = ""
	}
	return fmt.Sprintf("<%[1]s%[2]s>%[4]s%[3]s%[4]s</%[1]s>", n.Tag, n.attributeString(), strings.Join(content, newline), newline)
}

func (n *Node) contentStringFull() []string {
	split := make([]string, 0)
	switch content := n.Content.(type) {
	case []Node:
		for _, item := range content {
			split = append(split, strings.Split(item.XMLStringFull(), "\n")...)
		}
	case []byte:
		if strContent := printable(content); len(strContent) > 0 {
			if IndentXML {
				split = append(split, strings.Split(string(content), "\n")...)
			} else {
				split = append(split, strings.ReplaceAll(string(content), "\n", "\\n"))
			}
		} else if len(content) > MaxBytesToPrintAsHex {
			hexData := hex.EncodeToString(content)
			split = append(split, hexData)
		} else if !IndentXML {
			split = append(split, hex.EncodeToString(content))
		} else {
			hexData := hex.EncodeToString(content)
			split = append(split, hexData)
		}
	case nil:
		// don't append anything
	default:
		strContent := fmt.Sprintf("%s", content)
		if IndentXML {
			split = append(split, strings.Split(strContent, "\n")...)
		} else {
			split = append(split, strings.ReplaceAll(strContent, "\n", "\\n"))
		}
	}
	if len(split) > 1 && IndentXML {
		for i, line := range split {
			split[i] = "  " + line
		}
	}
	return split
}
