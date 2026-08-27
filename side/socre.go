package side

import (
	. "github.com/sterfuny/receiptexam"

	"fmt"
)

func Socre(r int, w int) string {
	return fmt.Sprintf("%s%s\n",
		Right.Render(fmt.Sprintf("✓ %d", r)),
		Wrong.Render(fmt.Sprintf(" ✗ %d", w)),
	)
}

func Tip(parts ...string) string {
	tip := ""
	for i, s := range parts {
		if i > 0 {
			tip += " • "
		}
		tip += s
	}
	return tip
}
