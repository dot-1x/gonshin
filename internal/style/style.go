// Package style provides a small Tailwind class merger used by the view
// components. It mirrors the behaviour of the `tailwind-merge` npm package
// (used by the original React components via `cn`) for the utility groups
// this app relies on: later classes win within a conflicting group, while
// unrelated groups are preserved.
package style

import (
	"strings"
)

// CN joins class strings, resolving Tailwind utility conflicts so that the
// last class in a conflicting group wins.
func CN(classes ...string) string {
	type token struct {
		raw    string
		prefix string
		group  string
	}

	kept := make([]token, 0, 16)

	for _, class := range classes {
		for _, raw := range strings.Fields(class) {
			prefix, base := splitVariant(raw)
			group := classify(base)

			remove := conflictSet(group)
			filtered := kept[:0]
			for _, k := range kept {
				if k.prefix == prefix && (k.group == group || remove[k.group]) {
					continue
				}
				filtered = append(filtered, k)
			}
			kept = append(filtered, token{raw: raw, prefix: prefix, group: group})
		}
	}

	parts := make([]string, len(kept))
	for i, k := range kept {
		parts[i] = k.raw
	}
	return strings.Join(parts, " ")
}

// splitVariant separates a variant prefix from the utility. The split point is
// the last top-level ':' (one that is not nested in square brackets or quotes).
func splitVariant(class string) (prefix, base string) {
	depth := 0
	var quote rune
	last := -1
	for i, r := range class {
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
			}
		case r == '\'' || r == '"':
			quote = r
		case r == '[':
			depth++
		case r == ']':
			if depth > 0 {
				depth--
			}
		case r == ':' && depth == 0:
			last = i
		}
	}
	if last == -1 {
		return "", class
	}
	return class[:last], class[last+1:]
}

// conflictSet returns the set of groups that group overrides, plus itself.
func conflictSet(group string) map[string]bool {
	set := map[string]bool{group: true}
	for _, g := range overrides[group] {
		set[g] = true
	}
	return set
}

var overrides = map[string][]string{
	"p":   {"px", "py", "pt", "pr", "pb", "pl", "ps", "pe"},
	"px":  {"pl", "pr", "ps", "pe"},
	"py":  {"pt", "pb"},
	"gap": {"gap-x", "gap-y"},
	"rounded": {
		"rounded-t", "rounded-r", "rounded-b", "rounded-l",
		"rounded-tl", "rounded-tr", "rounded-br", "rounded-bl",
		"rounded-s", "rounded-e", "rounded-ss", "rounded-se", "rounded-ee", "rounded-es",
	},
	"rounded-t": {"rounded-tl", "rounded-tr"},
	"rounded-r": {"rounded-tr", "rounded-br"},
	"rounded-b": {"rounded-br", "rounded-bl"},
	"rounded-l": {"rounded-tl", "rounded-bl"},
	"border-w": {
		"border-w-x", "border-w-y", "border-w-t", "border-w-r", "border-w-b", "border-w-l",
		"border-w-s", "border-w-e",
	},
	"border-w-x": {"border-w-l", "border-w-r"},
	"border-w-y": {"border-w-t", "border-w-b"},
}

var (
	paddingSides = map[byte]string{'x': "px", 'y': "py", 't': "pt", 'r': "pr", 'b': "pb", 'l': "pl", 's': "ps", 'e': "pe"}

	fontSizes = map[string]bool{
		"xs": true, "sm": true, "base": true, "lg": true, "xl": true,
		"2xl": true, "3xl": true, "4xl": true, "5xl": true, "6xl": true,
		"7xl": true, "8xl": true, "9xl": true,
	}

	fontWeights = map[string]bool{
		"thin": true, "extralight": true, "light": true, "normal": true,
		"medium": true, "semibold": true, "bold": true, "extrabold": true, "black": true,
	}

	textNonColor = map[string]bool{
		"balance": true, "pretty": true, "nowrap": true, "wrap": true,
		"left": true, "center": true, "right": true, "justify": true,
		"start": true, "end": true, "clip": true, "ellipsis": true, "truncate": true,
	}

	bgNonColor = []string{
		"bg-gradient", "bg-clip-", "bg-origin-", "bg-blend-", "bg-repeat",
		"bg-no-repeat", "bg-none", "bg-fixed", "bg-local", "bg-scroll",
		"bg-cover", "bg-contain", "bg-bottom", "bg-top", "bg-left", "bg-right", "bg-center",
	}

	roundedSizes = map[string]bool{
		"": true, "none": true, "sm": true, "md": true, "lg": true, "xl": true,
		"2xl": true, "3xl": true, "4xl": true, "full": true,
	}
)

// classify returns the conflict group for a base utility (variant stripped).
func classify(base string) string {
	if base == "" {
		return "unique:" + base
	}

	// Padding: p, px, py, pt, pr, pb, pl, ps, pe (+ arbitrary values).
	if g := classifyPadding(base); g != "" {
		return g
	}

	// Gap.
	if g := classifyGap(base); g != "" {
		return g
	}

	// Text: size vs color vs alignment/wrap.
	if strings.HasPrefix(base, "text-") {
		return classifyText(base)
	}

	// Background colour.
	if strings.HasPrefix(base, "bg-") {
		for _, p := range bgNonColor {
			if strings.HasPrefix(base, p) {
				return "unique:" + base
			}
		}
		return "bg-color"
	}

	// Border width / colour.
	if g := classifyBorder(base); g != "" {
		return g
	}

	// Rounded corners.
	if strings.HasPrefix(base, "rounded") {
		return classifyRounded(base)
	}

	// Font weight vs family.
	if strings.HasPrefix(base, "font-") {
		rest := strings.TrimPrefix(base, "font-")
		if fontWeights[rest] {
			return "font-weight"
		}
		if rest == "mono" || rest == "sans" || rest == "serif" {
			return "unique:" + base
		}
		return "unique:" + base
	}

	// Box sizing: size / width / height are kept separate (tailwind-merge does
	// not let `size-*` override `w-*`/`h-*`).
	for _, p := range []string{"size-", "w-", "h-"} {
		if strings.HasPrefix(base, p) {
			return p[:len(p)-1]
		}
	}

	// max/min constraints.
	if strings.HasPrefix(base, "max-w-") {
		return "max-w"
	}
	if strings.HasPrefix(base, "max-h-") {
		return "max-h"
	}
	if strings.HasPrefix(base, "min-w-") {
		return "min-w"
	}
	if strings.HasPrefix(base, "min-h-") {
		return "min-h"
	}

	// Flex shorthand / direction / wrap. `flex` alone is a display utility.
	switch {
	case base == "flex":
		return "unique:" + base
	case strings.HasPrefix(base, "flex-"):
		rest := strings.TrimPrefix(base, "flex-")
		switch rest {
		case "row", "row-reverse", "col", "col-reverse":
			return "flex-dir"
		case "wrap", "wrap-reverse", "nowrap":
			return "flex-wrap"
		default:
			return "flex"
		}
	}

	// Display utilities that could be swapped.
	switch base {
	case "block", "inline-block", "inline", "grid", "inline-grid", "contents", "hidden", "table", "flow-root":
		return "display"
	}

	// Overflow axes.
	if strings.HasPrefix(base, "overflow") {
		return base
	}

	// z-index (strict, so animation utilities like `zoom-in-95` aren't caught).
	if base == "z-auto" || isNumericValue(base, "z-") {
		return "z"
	}

	// shadow / opacity / transition / duration / ease / animate / leading /
	// tracking / whitespace / items / justify / self / place / col / row.
	if g := classifyPrefixGroup(base, []string{
		"shadow", "opacity", "duration", "ease", "animate", "leading",
		"tracking", "whitespace", "items", "justify", "self", "place",
		"col-", "row-", "order", "basis", "grow", "shrink", "object",
	}); g != "" {
		return g
	}
	if base == "transition" || strings.HasPrefix(base, "transition-") {
		return "transition"
	}

	return "unique:" + base
}

func classifyPadding(base string) string {
	if !strings.HasPrefix(base, "p") {
		return ""
	}
	rest := base[1:]
	if rest == "" {
		return ""
	}
	if rest[0] == '-' {
		// All sides: p-<value>.
		return "p"
	}
	// Side-specific: px-, py-, pt-, ...
	if g, ok := paddingSides[rest[0]]; ok && len(rest) > 1 && rest[1] == '-' {
		return g
	}
	return ""
}

func classifyGap(base string) string {
	if !strings.HasPrefix(base, "gap-") {
		return ""
	}
	switch {
	case strings.HasPrefix(base, "gap-x-"):
		return "gap-x"
	case strings.HasPrefix(base, "gap-y-"):
		return "gap-y"
	default:
		return "gap"
	}
}

func classifyText(base string) string {
	rest := strings.TrimPrefix(base, "text-")
	if textNonColor[rest] {
		return "unique:" + base
	}
	if fontSizes[rest] {
		return "text-size"
	}
	if strings.HasPrefix(rest, "[") && isLength(rest) {
		return "text-size"
	}
	return "text-color"
}

// isLength reports whether an arbitrary value like "[0.65rem]" is a length
// (font-size) rather than a colour.
func isLength(value string) bool {
	v := strings.TrimSuffix(strings.TrimPrefix(value, "["), "]")
	if v == "" {
		return false
	}
	for _, suffix := range []string{"rem", "px", "em", "%", "vw", "vh", "ch", "ex"} {
		if strings.HasSuffix(v, suffix) {
			return true
		}
	}
	// Bare number (e.g. text-[12]).
	for _, r := range v {
		if (r < '0' || r > '9') && r != '.' {
			return false
		}
	}
	return true
}

func classifyBorder(base string) string {
	if base == "border" {
		return "border-w"
	}
	if !strings.HasPrefix(base, "border-") {
		return ""
	}
	rest := strings.TrimPrefix(base, "border-")

	// Optional side: border-t, border-b-2, border-x, ...
	sideKey := ""
	if len(rest) >= 1 {
		switch rest[0] {
		case 't', 'r', 'b', 'l', 's', 'e', 'x', 'y':
			if len(rest) == 1 {
				// Bare side utility is the 1px width.
				return "border-w-" + string(rest[0])
			}
			if rest[1] == '-' {
				sideKey = string(rest[0])
				rest = rest[2:]
			}
		}
	}

	// Width: numeric or arbitrary length.
	if isNumeric(rest) || (strings.HasPrefix(rest, "[") && isLength(rest)) {
		group := "border-w"
		if sideKey != "" {
			if sideKey == "x" || sideKey == "y" {
				group = "border-w-" + sideKey
			} else {
				group = "border-w-" + sideKey
			}
		}
		return group
	}

	// Otherwise a colour.
	return "border-color"
}

// isNumericValue reports whether base is `prefix` followed by a number or an
// arbitrary value, e.g. `z-50` or `z-[100]`.
func isNumericValue(base, prefix string) bool {
	rest, ok := strings.CutPrefix(base, prefix)
	if !ok || rest == "" {
		return false
	}
	return isNumeric(rest) || strings.HasPrefix(rest, "[")
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func classifyRounded(base string) string {
	rest := strings.TrimPrefix(base, "rounded")
	if rest == "" {
		return "rounded"
	}
	if rest[0] != '-' {
		return "unique:" + base
	}
	rest = rest[1:]
	// Side-specific variants.
	for _, side := range []string{"tl", "tr", "br", "bl", "ss", "se", "ee", "es", "t", "r", "b", "l", "s", "e"} {
		if rest == side || strings.HasPrefix(rest, side+"-") {
			// Fall back to the base `rounded` group for sizes; tailwind-merge
			// treats a plain `rounded-<size>` as the all-corners group.
			return "rounded-" + side
		}
	}
	if roundedSizes[rest] || strings.HasPrefix(rest, "[") {
		return "rounded"
	}
	return "unique:" + base
}

func classifyPrefixGroup(base string, prefixes []string) string {
	for _, p := range prefixes {
		if strings.HasPrefix(base, p) {
			return p
		}
	}
	return ""
}
