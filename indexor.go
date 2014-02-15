package main

import (
	"bufio"
	"fmt"
	"os"
	"io"
	"path"
	"path/filepath"
	"strings"

	"github.com/gonuts/commander"
	"github.com/gonuts/flag"
)

var g_cmd = &commander.Command{
	UsageLine: os.Args[0] + " does cool things",
}

func visit(path string, f os.FileInfo, err error) error {
	fmt.Printf("Visited: %s\n", path)
	return nil

}

// should return true if you want to descend further, false if stop
type FastWalkCallback func(string, []os.FileInfo) bool

func walk_one(pth string, cb FastWalkCallback) {
	dir, err := os.Open(pth)
	if err != nil {
		pe := err.(*os.PathError)
		fmt.Printf("Path error: %s (%s)\n", pe, pth)
		return
	}

	fis, err := dir.Readdir(-1)

	//fmt.Printf("%s", pth)
	r := cb(pth, fis)
	if !r {
		return
	}

	for _, fi := range fis {
		if fi.IsDir() {
			walk_one(pth+"/"+fi.Name(), cb)
		}
		//fmt.Println(fi.Name())
	}

}

func write_inline(wr io.Writer, fname string, lines int) {
	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}

	fmt.Printf("inline %s\n", fname)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l := scanner.Bytes()
		//fmt.Printf("scan %v\n", l)

		wr.Write([]byte{32})
		wr.Write(l)
		wr.Write([]byte{'\n'})
		if lines--; lines < 1 {
			wr.Write([]byte("# file too long") )
			return
		}

		//fmt.Println(l)




	}


}

func create_index_cmd(cmd *commander.Command, args []string) error {

	root, err := filepath.Abs(args[0])

	if err != nil {
		panic(err)

	}
	fmt.Printf("Indexing: [%s]\n", root)

	inline_pats := strings.Split(cmd.Flag.Lookup("inline").Value.Get().(string), ":")

	path.Match("","")
	fmt.Printf("Patterns: [%v]\n", inline_pats)

	of, err := os.Create("index.txt")
	if err != nil {
		panic(err)
	}
	wr := bufio.NewWriter(of)
	defer wr.Flush()
	counter := 0
	mycb := func(pth string, fis []os.FileInfo) bool {
		//if strings.HasPrefix(pth, ".") {
		//	return false
		//}
		//fmt.Printf("Got %d at %s\n", len(fis), pth)
		bname := filepath.Base(pth)
		if strings.HasPrefix(bname, ".") {
			return false
		}
		wr.WriteString(pth + "/\n")
		for _, fi := range fis {
			name := fi.Name()
			if !fi.IsDir() {
				fullname := pth + "/" + name
				_, err := wr.WriteString( fullname + "\n")
				if err != nil {
					panic(err)
				}
				counter += 1

				for _, pat := range inline_pats {
					m, _ := path.Match(pat, name)
					if m {
						fmt.Printf("Match %s %v\n", name, m)

						write_inline(wr, fullname, 50)
					}
				}

			}
			if counter%200 == 0 {
				fmt.Printf("%d\n", counter)
			}

		}
		return true

	}
	walk_one(root, mycb)

	return nil
}

func subcommands() {
	cmd_index := &commander.Command{
		Run:       create_index_cmd,
		UsageLine: "index <path>",
		Short:     "Creates index for specified path",
		Long:      "Will create index.txt at root path",
		Flag:      *flag.NewFlagSet("my-cmd-cmd1", flag.ExitOnError),
	}

	g_cmd.Subcommands = []*commander.Command{
		cmd_index,
	}

	cmd_index.Flag.String("inline", "", "List of file patterns to inline in the index")
}

func main() {

	subcommands()

	err := g_cmd.Dispatch(os.Args[1:])
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	//fmt.Printf("Hello, world.\n")
	//flag.Parse()

	//fo, err := os.Create("out.txt")
	//defer fo.Close()

}
