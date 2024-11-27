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

	for {
		tt, data := l.Next()
		fmt.Printf("token type: %s data [%d]: %s\n", tt.String(), len(data), data)

		switch tt {
		case html.ErrorToken:
		// error or EOF set in l.Err()

		return
		case html.StartTagToken:
		// ...
			for {
				ttAttr, dataAttr := l.Next()
				if ttAttr != html.AttributeToken {break}
				fmt.Printf("  ttAttr: %s dataAttr: %s\n", ttAttr, dataAttr)
			}

		case html.EndTagToken:

		case html.StartTagCloseToken:

		case html.StartTagVoidToken:

		case html.AttributeToken:

		case html.TextToken:


		case html.CommentToken:


		case html.DoctypeToken:

		default:

		}

	} //for loop end


}
