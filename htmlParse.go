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
//	txt := false
//	par := nil
	if dbg {PrintAst(top, "top")}
	cnod := top
	istate := 2
	nlev := 0
	for ni:=1; ni< 50; ni++ {
		tt, data := l.Next()
		if dbg {fmt.Printf("%d -- token type: %s data [%d]: %s\n", ni, tt.String(), len(data), data)}

		switch tt {
		case html.ErrorToken:
		// error or EOF set in l.Err()
			return

		case html.StartTagToken:
			if dbg {fmt.Printf("dbg -- token: %s\n", data[1:])}
			switch istate {
			case 0:
				istate = 1
			// nested token
			case 2,3:
				nlev++
				istate = 1
			default:
				log.Fatalf("error -- invalid state %d for StartTagToken!\n", istate)
			}
//	PrintAst(par)
			n := new(node)
			(*n).typ = string(data[1:])
			(*n).att = make (map[string]string)
			(*n).ni = ni
			(*n).pnode = cnod
			cnod.children = append(cnod.children, n)
//fmt.Printf("parent typ: %s children: %d\n", (*par).typ, len((*par).children))
//PrintAst(top)
			cnod = n

		case html.EndTagToken:
			if dbg {fmt.Printf("dbg -- end token: %s\n", data)}
			switch istate {
			// attributes
			case 2, 3:
				if nlev>0 {
					nlev--
				}
				istate = 0
			default:
				log.Fatalf("error -- invalid state %d for StartTagCloseToken!\n", istate)
			}
//			txt = false
			cnod = cnod.pnode
			if cnod.pnode == nil {last = true}

		case html.StartTagCloseToken:
			if dbg {fmt.Printf("dbg -- close token token: %s\n", data)}
			switch istate {
			// attributes
			case 1, 2:
				istate = 3
			default:
				log.Fatalf("error -- invalid state %d for StartTagCloseToken!\n", istate)
			}

		case html.StartTagVoidToken:
			if dbg {fmt.Printf("dbg -- void token: %s\n", data)}

		case html.AttributeToken:
			if dbg {fmt.Printf("dbg -- att token: %s\n", data)}
			switch istate {
			case 1, 2:
				istate = 2
				if dbg {fmt.Printf("  key: %s val: %s\n", l.AttrKey(), l.AttrVal())}
				if cnod.att == nil {cnod.att = make (map[string]string)}
				cnod.att[string(l.AttrKey())] = string(l.AttrVal())

			default:
				log.Fatalf("error -- invalid state %d for StartTagCloseToken!\n", istate)
			}

		case html.TextToken:
			if dbg {fmt.Printf("dbg -- text[%d]: %s\n", len(data), data)}
			switch istate {
			case 3:
				cnod.txt = string(data)
			default:
				log.Fatalf("error -- invalid state %d for TextToken!\n", istate)
			}

		case html.CommentToken:
			if dbg {fmt.Printf("dbg -- comment: %s\n", data)}

		case html.DoctypeToken:
			if dbg {fmt.Printf("dbg -- doc token: %s\n", data)}

		default:
			log.Fatalf("unknown token: %s\n", data)

		}

		if last {
			if dbg {fmt.Println("*** last ***")}
			break
		}
	} //for loop end

	PrintAst(top, "ast")

}

func PrintAst(n *node, msg string) {
	fmt.Printf("******** %s Node Start: %s **********\n", msg, (*n).typ)
	prnode(n)
	chnum := len((*n).children)
//    fmt.Printf("children [%d]\n", chnum)
    for i:= 0; i< chnum; i++ {
		ch := (*n).children[i]
		chmsg := fmt.Sprintf("child %d",i)
		PrintAst(ch, chmsg)
	}
	fmt.Printf("********* %s Node End: %s ***********\n", msg, (*n).typ)
}

func prnode (n *node) {
	fmt.Printf("type: %s\n", n.typ)
	if len((*n).txt) > 1 {
		fmt.Printf("  text: %s\n", (*n).txt)
	} else {
		fmt.Printf("  no text\n")
	}
	attr:= (*n).att
	if len(attr) > 0 {
		fmt.Println("   attributes:")
		for key, val := range (*n).att {
    		fmt.Printf("      %s:%s\n", key, val)
		}
	}
	par := (*n).pnode
	if par == nil {
		fmt.Println("  no parent")
	} else {
		fmt.Printf("  parent typ: %s\n", par.typ)
	}

	chnum := len((*n).children)
	if chnum > 0 {
		fmt.Printf("  children [%d]\n", chnum)
	} else {
		fmt.Printf("  no children\n")
	}
	for i:= 0; i< chnum; i++ {
		ch := (*n).children[i]
		fmt.Printf("    child[%d] -- typ: %s\n", i,(*ch).typ)
//		fmt.Printf("child[%d] -- txt: %s\n", i,(*ch).txt)
	}
}
