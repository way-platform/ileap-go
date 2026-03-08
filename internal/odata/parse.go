package odata

import (
	"strings"

	ileapv1 "github.com/way-platform/ileap-go/proto/gen/wayplatform/connect/ileap/v1"
)

// ParseFilter parses a best-effort OData subset into standalone filter pairs.
// Unsupported or invalid clauses are ignored.
func ParseFilter(raw string) []*ileapv1.Filter {
	data := strings.TrimSpace(raw)
	if data == "" {
		return nil
	}
	filters := make([]*ileapv1.Filter, 0)
	for _, clause := range splitTopLevelAndClauses(data) {
		clause = strings.TrimSpace(clause)
		if clause == "" {
			continue
		}
		filter, ok := parseClause(clause)
		if !ok {
			continue
		}
		filters = append(filters, filter)
	}
	return filters
}

func splitTopLevelAndClauses(filter string) []string {
	clauses := make([]string, 0, 4)
	start := 0
	depth := 0
	inString := false
	for i := 0; i < len(filter); i++ {
		ch := filter[i]
		if ch == '\'' {
			if inString && i+1 < len(filter) && filter[i+1] == '\'' {
				i++
				continue
			}
			inString = !inString
			continue
		}
		if inString {
			continue
		}
		switch ch {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		default:
			if depth == 0 && matchesTopLevelAndAt(filter, i) {
				clauses = append(clauses, filter[start:i])
				i += 2
				start = i + 1
			}
		}
	}
	clauses = append(clauses, filter[start:])
	return clauses
}

func matchesTopLevelAndAt(data string, i int) bool {
	if i < 0 || i+2 >= len(data) {
		return false
	}
	if !strings.EqualFold(data[i:i+3], "and") {
		return false
	}
	before := byte(' ')
	after := byte(' ')
	if i > 0 {
		before = data[i-1]
	}
	if i+3 < len(data) {
		after = data[i+3]
	}
	return isSeparator(before) && isSeparator(after)
}

func isSeparator(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r' || b == '(' || b == ')'
}

func parseClause(clause string) (*ileapv1.Filter, bool) {
	tokens, ok := lexClause(clause)
	if !ok {
		return nil, false
	}
	parser := clauseParser{tokens: tokens}
	return parser.parse()
}

type tokenKind int

const (
	tokenInvalid tokenKind = iota
	tokenEOF
	tokenIdent
	tokenString
	tokenSlash
	tokenDot
	tokenLParen
	tokenRParen
	tokenColon
)

type token struct {
	kind tokenKind
	text string
}

func lexClause(input string) ([]token, bool) {
	tokens := make([]token, 0, len(input)/2)
	for i := 0; i < len(input); {
		ch := input[i]
		switch {
		case ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r':
			i++
		case ch == '/':
			tokens = append(tokens, token{kind: tokenSlash, text: "/"})
			i++
		case ch == '.':
			tokens = append(tokens, token{kind: tokenDot, text: "."})
			i++
		case ch == '(':
			tokens = append(tokens, token{kind: tokenLParen, text: "("})
			i++
		case ch == ')':
			tokens = append(tokens, token{kind: tokenRParen, text: ")"})
			i++
		case ch == ':':
			tokens = append(tokens, token{kind: tokenColon, text: ":"})
			i++
		case ch == '\'':
			value, next, ok := readStringLiteral(input, i)
			if !ok {
				return nil, false
			}
			tokens = append(tokens, token{kind: tokenString, text: value})
			i = next
		case isIdentStart(ch):
			start := i
			i++
			for i < len(input) && isIdentPart(input[i]) {
				i++
			}
			tokens = append(tokens, token{
				kind: tokenIdent,
				text: input[start:i],
			})
		default:
			return nil, false
		}
	}
	tokens = append(tokens, token{kind: tokenEOF})
	return tokens, true
}

func readStringLiteral(input string, start int) (string, int, bool) {
	if start >= len(input) || input[start] != '\'' {
		return "", start, false
	}
	var b strings.Builder
	for i := start + 1; i < len(input); i++ {
		if input[i] != '\'' {
			b.WriteByte(input[i])
			continue
		}
		if i+1 < len(input) && input[i+1] == '\'' {
			b.WriteByte('\'')
			i++
			continue
		}
		return b.String(), i + 1, true
	}
	return "", start, false
}

func isIdentStart(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9') ||
		ch == '_'
}

func isIdentPart(ch byte) bool {
	return isIdentStart(ch) || ch == '-'
}

type clauseParser struct {
	tokens []token
	pos    int
}

func (p *clauseParser) parse() (*ileapv1.Filter, bool) {
	openParens := 0
	for p.matchKind(tokenLParen) {
		openParens++
	}
	mark := p.pos
	filter, ok := p.parseAnyEq()
	if ok && p.consumeClosing(openParens) && p.peek().kind == tokenEOF {
		return filter, true
	}
	p.pos = mark
	filter, ok = p.parseSimple()
	if ok && p.consumeClosing(openParens) && p.peek().kind == tokenEOF {
		return filter, true
	}
	return nil, false
}

func (p *clauseParser) consumeClosing(n int) bool {
	for i := 0; i < n; i++ {
		if !p.matchKind(tokenRParen) {
			return false
		}
	}
	return true
}

func (p *clauseParser) parseSimple() (*ileapv1.Filter, bool) {
	path, ok := p.parsePath(false)
	if !ok {
		return nil, false
	}
	operator, ok := p.matchOperator()
	if !ok {
		return nil, false
	}
	valueToken, ok := p.matchString()
	if !ok {
		return nil, false
	}
	filter := new(ileapv1.Filter)
	filter.SetFieldPath(strings.Join(path, "."))
	filter.SetOperator(operator)
	filter.SetValue(valueToken.text)
	return filter, true
}

func (p *clauseParser) parseAnyEq() (*ileapv1.Filter, bool) {
	collectionPath, ok := p.parsePath(true)
	if !ok || !p.matchKind(tokenSlash) || !p.matchKeyword("any") {
		return nil, false
	}
	if !p.matchKind(tokenLParen) {
		return nil, false
	}
	alias, ok := p.matchIdent()
	if !ok || !p.matchKind(tokenColon) || !p.matchKind(tokenLParen) {
		return nil, false
	}
	innerPath, ok := p.parsePath(false)
	if !ok || len(innerPath) == 0 || !strings.EqualFold(innerPath[0], alias.text) {
		return nil, false
	}
	operator, ok := p.matchOperator()
	if !ok {
		return nil, false
	}
	valueToken, ok := p.matchString()
	if !ok || !p.matchKind(tokenRParen) || !p.matchKind(tokenRParen) {
		return nil, false
	}
	fieldPath := append([]string{}, collectionPath...)
	fieldPath = append(fieldPath, innerPath[1:]...)
	filter := new(ileapv1.Filter)
	filter.SetFieldPath(strings.Join(fieldPath, "."))
	filter.SetOperator(operator)
	filter.SetValue(valueToken.text)
	return filter, true
}

func (p *clauseParser) parsePath(stopBeforeAny bool) ([]string, bool) {
	first, ok := p.matchIdent()
	if !ok {
		return nil, false
	}
	segments := []string{first.text}
	for {
		next := p.peek()
		if next.kind != tokenSlash && next.kind != tokenDot {
			break
		}
		if stopBeforeAny && next.kind == tokenSlash {
			after := p.peekN(1)
			if after.kind == tokenIdent && strings.EqualFold(after.text, "any") {
				break
			}
		}
		p.pos++
		part, ok := p.matchIdent()
		if !ok {
			return nil, false
		}
		segments = append(segments, part.text)
	}
	return segments, true
}

func (p *clauseParser) matchKeyword(keyword string) bool {
	token := p.peek()
	if token.kind != tokenIdent || !strings.EqualFold(token.text, keyword) {
		return false
	}
	p.pos++
	return true
}

func (p *clauseParser) matchOperator() (ileapv1.Filter_Operator, bool) {
	tok := p.peek()
	if tok.kind != tokenIdent {
		return ileapv1.Filter_OPERATOR_UNSPECIFIED, false
	}
	operator, ok := operatorFromToken(tok.text)
	if !ok {
		return ileapv1.Filter_OPERATOR_UNSPECIFIED, false
	}
	p.pos++
	return operator, true
}

func operatorFromToken(raw string) (ileapv1.Filter_Operator, bool) {
	switch {
	case strings.EqualFold(raw, "eq"):
		return ileapv1.Filter_EQ, true
	case strings.EqualFold(raw, "ne"):
		return ileapv1.Filter_NE, true
	case strings.EqualFold(raw, "lt"):
		return ileapv1.Filter_LT, true
	case strings.EqualFold(raw, "le"):
		return ileapv1.Filter_LE, true
	case strings.EqualFold(raw, "gt"):
		return ileapv1.Filter_GT, true
	case strings.EqualFold(raw, "ge"):
		return ileapv1.Filter_GE, true
	default:
		return ileapv1.Filter_OPERATOR_UNSPECIFIED, false
	}
}

func (p *clauseParser) matchKind(kind tokenKind) bool {
	if p.peek().kind != kind {
		return false
	}
	p.pos++
	return true
}

func (p *clauseParser) matchIdent() (token, bool) {
	tok := p.peek()
	if tok.kind != tokenIdent {
		return token{}, false
	}
	p.pos++
	return tok, true
}

func (p *clauseParser) matchString() (token, bool) {
	tok := p.peek()
	if tok.kind != tokenString {
		return token{}, false
	}
	p.pos++
	return tok, true
}

func (p *clauseParser) peek() token {
	return p.peekN(0)
}

func (p *clauseParser) peekN(offset int) token {
	index := p.pos + offset
	if index >= len(p.tokens) {
		return token{kind: tokenEOF}
	}
	return p.tokens[index]
}
