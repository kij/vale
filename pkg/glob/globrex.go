package glob

import (
	"runtime"

	"github.com/dlclark/regexp2"
)

var sep = "\\/"
var defaultOpts = rexOpts{
	extended: false,
	globstar: false,
	strict:   false,
	fp:       false,
}

func init() {
	if runtime.GOOS == "windows" {
		sep = "\\\\+"
	}
}

type rexOpt func(r *rex, opts *rexOpts)

type rexOpts struct {
	extended bool
	globstar bool
	strict   bool
	fp       bool
}

type pathLike struct {
	regex    string
	segments []*regexp2.Regexp
}

type rex struct {
	fp      bool
	glob    string
	inGroup bool
	inRange bool
	path    pathLike
	regex   string
	segment string
}

func newRex(glob string, opts ...rexOpt) (*rex, error) {
	var curr byte
	var next byte
	var stack []string

	base := defaultOpts

	r := rex{glob: glob}
	for i := range glob {
		curr = glob[i]
		next = glob[i+1]

		s := string(curr)
		c := `\\` + s

		inStack := len(stack) > 0
		switch curr {
		case '\\', '$', '^', '.', '=':
			r.add(c, "", false, false)
		case '/':
			r.add(c, "", true, false)
			if next == '/' && !base.strict {
				r.regex += "?"
			}
		case '(':
			if inStack {
				r.add(s, "", false, false)
			} else {
				r.add(c, "", false, false)
			}
		case ')':
			if inStack {
				r.add(s, "", false, false)
				x := stack[len(stack)-1]
				if x == "@" {
					r.add("{1}", "", false, false)
				} else if x == "!" {
					r.add(`([^\/]*)`, "", false, false)
				} else {
					r.add(x, "", false, false)
				}
				stack = stack[:len(stack)-1]
			} else {
				r.add(c, "", false, false)
			}
		case '|':
			if inStack {
				r.add(s, "", false, false)
			} else {
				r.add(c, "", false, false)
			}
		case '+':
			if next == '(' && base.extended {
				stack = append(stack, s)
			} else {
				r.add(c, "", false, false)
			}
		case '@':
			if base.extended && next == '(' {
				stack = append(stack, s)
			}
		case '!':
			if base.extended {
				if r.inRange {
					r.add("^", "", false, false)
					continue
				}
				if next == '(' {
					stack = append(stack, s)
					r.add("(?!", "", false, false)
					i++
					continue
				}
			}
			r.add(c, "", false, false)
		case '?':
			if base.extended {
				if next == '(' {
					stack = append(stack, s)
				} else {
					r.add(".", "", false, false)
				}
			} else {
				r.add(c, "", false, false)
			}
		case '[':
			if base.extended {
				r.inRange = true
				r.add(s, "", false, false)
			} else {
				r.add(c, "", false, false)
			}
		case ']':
			if base.extended {
				r.inRange = false
				r.add(s, "", false, false)
			} else {
				r.add(c, "", false, false)
			}
		case '{':
			if base.extended {
				r.inGroup = true
				r.add("(", "", false, false)
			} else {
				r.add(c, "", false, false)
			}
		case '}':
			if base.extended {
				r.inGroup = false
				r.add(")", "", false, false)
			} else {
				r.add(c, "", false, false)
			}
		case '*':
			if next == '(' && base.extended {
				stack = append(stack, s)
				continue
			}

		}
	}

	return &r, nil
}

func (r *rex) add(s, only string, split, last bool) {
	if only != "path" {
		r.regex += s
	}
	if r.fp && only != "regex" {
		if s == "\\/" {
			r.path.regex += sep
		} else {
			r.path.regex += s
		}
		if split {
			if last {
				r.segment += s
			}
			if r.segment != "" {
				r.path.segments = append(
					r.path.segments,
					regexp2.MustCompile(r.segment, regexp2.ECMAScript),
				)
			}
			r.segment = ""
		} else {
			r.segment += s
		}
	}
}
