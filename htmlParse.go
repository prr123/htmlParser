// https://github.com/tdewolff/parse/tree/master/html

package main

import (
	"fmt"
	"log"
	"os"
	"bytes"

	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/html"
    cliutil "github.com/prr123/utility/utilLib"
)

type node struct {
	pnode *node
	typ string
	children []*node
	att map[string]string
	txt string
	ni int
}

func main() {

    numarg := len(os.Args)
    flags:=[]string{"dbg", "in", "out"}

    useStr := " /in=infile /out=outfile [/dbg]"
    helpStr := "markdown to html conversion program"

    if numarg > len(flags) +1 {
        fmt.Println("too many arguments in cl!")
        fmt.Println("usage: %s %s\n", os.Args[0], useStr)
        os.Exit(-1)
    }

    if numarg == 1 || (numarg > 1 && os.Args[1] == "help") {
        fmt.Printf("help: %s\n", helpStr)
        fmt.Printf("usage is: %s %s\n", os.Args[0], useStr)
        os.Exit(1)
    }

    flagMap, err := cliutil.ParseFlags(os.Args, flags)
    if err != nil {log.Fatalf("util.ParseFlags: %v\n", err)}

    dbg:= false
    _, ok := flagMap["dbg"]
    if ok {dbg = true}


    inFil := ""
    inval, ok := flagMap["in"]
    if !ok {
        log.Fatalf("error -- no in flag provided!\n")
    } else {
        if inval.(string) == "none" {log.Fatalf("error -- no input file name provided!\n")}
        inFil = inval.(string)
    }

    outFil := ""
    outval, ok := flagMap["out"]
    if !ok {
		outFil = inFil
//      log.Fatalf("error -- no out flag provided!\n")
    } else {
        if outval.(string) == "none" {outFil = inFil}
//{log.Fatalf("error -- no output file name provided!\n")}
        outFil = outval.(string)
    }

    inFilnam := "html/" + inFil + ".html"
    outFilnam := "dump/" + outFil + ".txt"

    if dbg {
        fmt.Printf("input:  %s\n", inFilnam)
        fmt.Printf("output: %s\n", outFilnam)
    }

    source, err := os.ReadFile(inFilnam)
    if err != nil {log.Fatalf("error -- open file: %v\n",err)}

	r := bytes.NewReader(source)
	l := html.NewLexer(parse.NewInput(r))

	tt, data := l.Next()
	if tt !=  html.StartTagToken {log.Fatalf("invalid start token!\n")}
	top := new(node)
	top.typ = string(data[1:])
	top.att = make (map[string]string)
	top.ni = 0
	last := false
	txt := false
	par := top
	PrintAst(top)

	for ni:=1; ni< 50; ni++ {
		tt, data := l.Next()
		fmt.Printf("%d -- token type: %s data [%d]: %s\n", ni, tt.String(), len(data), data)

		switch tt {
		case html.ErrorToken:
		// error or EOF set in l.Err()

		return
		case html.StartTagToken:
			fmt.Printf("dbg -- token: %s\n", data[1:])
//	PrintAst(par)
			n := new(node)
			(*n).typ = string(data[1:])
			(*n).att = make (map[string]string)
			(*n).ni = ni
			(*n).pnode = par
			(*par).children = append((*par).children, n)
//fmt.Printf("parent typ: %s children: %d\n", (*par).typ, len((*par).children))
			par = n
//PrintAst(top)

		case html.EndTagToken:
			fmt.Printf("dbg -- end token: %s\n", data)
			txt = false
			par = par.pnode
			if par == nil {last = true}

		case html.StartTagCloseToken:
			txt = true
			fmt.Printf("dbg -- close token token: %s\n", data)

		case html.StartTagVoidToken:
			fmt.Printf("dbg -- void token: %s\n", data)

		case html.AttributeToken:
			fmt.Printf("dbg -- att token: %s\n", data)
			fmt.Printf("  ttAttr: %s dataAttr: %s\n", tt, data)
			par.att[string(tt)] = string(data)

		case html.TextToken:
			fmt.Printf("dbg -- text[%d]: %s\n", len(data), data)
			if txt {
				par.txt = string(data)
			}
		case html.CommentToken:
			fmt.Printf("dbg -- comment: %s\n", data)

		case html.DoctypeToken:
			fmt.Printf("dbg -- doc token: %s\n", data)

		default:
			fmt.Printf("dbg -- unknown token: %s\n", data)

		}

		if last {
			fmt.Println("*** last ***")
			break
		}
	} //for loop end

	fmt.Println ("*** ast ***")
	PrintAst(top)

}

func PrintAst(n *node) {

	prnode(n)
	chnum := len((*n).children)
//    fmt.Printf("children [%d]\n", chnum)
    for i:= 0; i< chnum; i++ {
		ch := (*n).children[i]
		PrintAst(ch)
	}
}

func prnode (n *node) {
	fmt.Printf("type: %s\n", n.typ)
	if len((*n).txt) > 1 {
		fmt.Printf("text: %s\n", (*n).txt)
	} else {
		fmt.Printf("no text\n")
	}
	par := (*n).pnode
	if par == nil {
		fmt.Println("no parent")
	} else {
		fmt.Printf("parent: %s\n", par.typ)
	}
	chnum := len((*n).children)
	if chnum > 0 {
		fmt.Printf("children [%d]\n", chnum)
	} else {
		fmt.Printf("no children\n")
	}
	for i:= 0; i< chnum; i++ {
		ch := (*n).children[i]
		fmt.Printf("child[%d] -- typ: %s\n", i,(*ch).typ)
//		fmt.Printf("child[%d] -- txt: %s\n", i,(*ch).txt)
	}
}
