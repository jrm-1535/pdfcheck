
package main

import (
    "fmt"
    "flag"
    "os"
    "github.com/jrm-1535/pdf"
)

const (
    VERSION     = "0.1"
    VERBOSE     = false
    FIXPARSE    = false
    FIXSTREAM   = false
    SUMMARY     = false
    DUMP        = false

    HELP        = `pdfcheck [-h] [-v] [-vp] [-fp] [-vs] [-fs] [-o=name] filepath

    Checks if a file is a valid pdf document, allowing to fix some minor syntactic error
    in the pdf file format or in various stream formats

    Options:
        -h          print this help message and exits
        -v          print pdfCheck version and exits

        -s          display a summary (info & catalog)
        -d          dump all defined indirect objects
 
        -vp         verbose during parsing
        -fp         try fixing errors during parsing, instead of stopping

        -vs         verbose checking the embedded streams
        -fs         try fixing stream errors (e.g. invalid jpeg streams)

        -o  name    output the hopefully fixed data to a new file
                    this option is meaningful if -fp or -fs is specified
                    if nothing was fixed, the files will be similar if not identical

    filepath is the path to the original file to process

`
)

type processArgs struct {
    output  string
    summary bool
    dump    bool
}

func getArgs( ) (* pdf.ParseArgs, *pdf.StreamArgs, *processArgs) {
    pArgs := new( pdf.ParseArgs )

    var version bool
    flag.BoolVar( &version, "v", false, "print pdfCheck version and exits" )
    flag.BoolVar( &pArgs.Verbose, "vp", VERBOSE, "verbose during parsing" )
    flag.BoolVar( &pArgs.Fix, "fp", FIXPARSE, "try fixing errors during parsing, instead of stopping" )

    sArgs := new( pdf.StreamArgs )
    flag.BoolVar( &sArgs.Verbose, "vs", VERBOSE, "verbose checking the embedded streams" )
    flag.BoolVar( &sArgs.Fix, "fs", FIXSTREAM, "try fixing stream errors (e.g. invalid jpeg streams)" )

    processArgs := new( processArgs )
    flag.StringVar( &processArgs.output, "o", "", "output the hopefully fixed data to the file`name`" )
    flag.BoolVar( &processArgs.summary, "s", SUMMARY, "display a summary (info & catalog)" )
    flag.BoolVar( &processArgs.dump, "d", DUMP, "dump all defined indirect objects" )

    flag.Usage = func() {
        fmt.Fprintf( flag.CommandLine.Output(), HELP )
    }
    flag.Parse()
    if version {
        fmt.Printf( "pdfCheck version %s\n", VERSION )
        os.Exit(0)
    }
    arguments := flag.Args()
    if len( arguments ) < 1 {
        fmt.Printf( "Missing the name of the file to process\n" )
        os.Exit(2)
    }
    if len( arguments ) > 1 {
        fmt.Printf( "Too many files specified (only 1 file at a time)\n" )
        os.Exit(2)
    }
    if pArgs.Fix || sArgs.Fix {
        if processArgs.output == "" {
            fmt.Printf( "Warning: although fixing the original file is requested, NO output file is requested\n" )
            fmt.Printf( "         proceeding anyway\n" )
        }
    } else if processArgs.output != "" {
        fmt.Printf( "Warning: although an output file is requested, fixing the original file is NOT requested\n" )
        fmt.Printf( "         proceeding anyway\n" )
    }
    pArgs.Path = arguments[0]
    return pArgs, sArgs, processArgs
}

func main() {
    pArgs, cArgs, process := getArgs()

    fmt.Printf( "Parsing pdf file: %s\n", pArgs.Path )
    pdfData, err := pdf.Parse( pArgs )
    if err != nil {
        fmt.Printf( "%v", err )
    } else {
        err = pdfData.Check( cArgs )
        if err != nil {
            os.Exit(2)
        }

        if process.summary {
//            fmt.Printf("PDF header: \"%s\"\n", pdfData.Header )
            fmt.Printf("PDF version: <%s>\n", pdfData.Version )
            fmt.Printf("number of indirect objects: %d (%d)\n", len( pdfData.Objects ), pdfData.Size )
            pdfData.PrintFileIds()
//            pdfData.PrintFileTrailer()
            pdfData.PrintInfo()
//            pdfData.PrintEncryption()
            pdfData.PrintCatalog()
        }
        if process.dump {
            pdfData.PrintAllObjs()
        }
        if process.output != "" {
            fmt.Printf( "Writing to file: %s\n", process.output )
            pdfData.Serialize( process.output )
        }
    }
// use diff -ua <original-file> <new-file> to check differences

}
