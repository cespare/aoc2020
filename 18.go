package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func init() {
	addSolutions(18, problem18)
}

func problem18(ctx *problemContext) {
	var exprToks [][]token
	scanner := ctx.scanner()
	for scanner.scan() {
		exprToks = append(exprToks, tokenizeExpr(scanner.text()))
	}
	ctx.reportLoad()

	var sum int64
	prec1 := map[tokenType]int{tokAdd: 0, tokMul: 0}
	for _, toks := range exprToks {
		expr := newExprParser(toks, prec1).parse()
		sum += expr.eval()
	}
	ctx.reportPart1(sum)

	sum = 0
	prec2 := map[tokenType]int{tokAdd: 1, tokMul: 0}
	for _, toks := range exprToks {
		expr := newExprParser(toks, prec2).parse()
		sum += expr.eval()
	}
	ctx.reportPart2(sum)
}

type exprNode interface {
	eval() int64
}

type binaryExprNode struct {
	op    tokenType
	left  exprNode
	right exprNode
}

func (n *binaryExprNode) eval() int64 {
	switch n.op {
	case tokAdd:
		return n.left.eval() + n.right.eval()
	case tokMul:
		return n.left.eval() * n.right.eval()
	default:
		panic("unexpected op")
	}
}

type intExprNode int64

func (n intExprNode) eval() int64 { return int64(n) }

type exprParser struct {
	toks []token
	i    int
	prec map[tokenType]int
}

func newExprParser(toks []token, prec map[tokenType]int) *exprParser {
	return &exprParser{
		toks: toks,
		prec: prec,
	}
}

func (p *exprParser) next() token {
	if p.i == len(p.toks) {
		return token{typ: tokEOF}
	}
	tok := p.toks[p.i]
	p.i++
	return tok
}

func (p *exprParser) backUp() {
	if p.i == 0 {
		panic("backUp before a token was consumed")
	}
	p.i--
}

func (p *exprParser) parse() exprNode {
	expr := p.parseBinary(0)
	if p.i < len(p.toks) {
		panic("trailing junk")
	}
	return expr
}

func (p *exprParser) parseBinary(precedence int) exprNode {
	left := p.parse1()
	for {
		tok := p.next()
		switch tok.typ {
		case tokAdd, tokMul:
		case tokEOF:
			return left
		case tokRightParen:
			p.backUp()
			return left
		default:
			panic(fmt.Sprintf("unexected token %q", tok.v))
		}
		prec, ok := p.prec[tok.typ]
		if !ok {
			panic("no precedence given for op")
		}
		if prec < precedence {
			// We found a lower-precedent op, so we're done with the
			// current expression.
			p.backUp()
			return left
		}
		left = &binaryExprNode{
			op:    tok.typ,
			left:  left,
			right: p.parseBinary(prec + 1),
		}
	}
}

func (p *exprParser) parse1() exprNode {
	tok := p.next()
	if tok.typ == tokLeftParen {
		expr := p.parseBinary(0)
		if p.next().typ != tokRightParen {
			panic("unmatched parens")
		}
		return expr
	}
	if tok.typ != tokInt {
		panic(fmt.Sprintf("unexected token %q", tok.v))
	}
	return intExprNode(parseInt(tok.v, 10, 64))
}

type tokenType int

const (
	tokEOF tokenType = iota
	tokLeftParen
	tokRightParen
	tokAdd
	tokMul
	tokInt
)

type token struct {
	typ tokenType
	v   string
}

func tokenizeExpr(s string) []token {
	var toks []token
	for len(s) > 0 {
		r, w := utf8.DecodeRuneInString(s)
		switch r {
		case '(', ')', '*', '+':
			s = s[w:]
			tok := token{v: string(r)}
			switch r {
			case '(':
				tok.typ = tokLeftParen
			case ')':
				tok.typ = tokRightParen
			case '+':
				tok.typ = tokAdd
			case '*':
				tok.typ = tokMul
			}
			toks = append(toks, tok)
			continue
		}
		if unicode.IsSpace(r) {
			s = s[w:]
			continue
		}
		if r >= '0' && r <= '9' {
			i := w
			for i < len(s) {
				r, w := utf8.DecodeRuneInString(s[i:])
				if r < '0' || r > '9' {
					break
				}
				i += w
			}
			tok := token{typ: tokInt, v: s[:i]}
			toks = append(toks, tok)
			s = s[i:]
			continue
		}
		panic("bad input")
	}
	return toks
}
