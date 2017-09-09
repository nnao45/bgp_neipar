package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"./getter"
	"./peer"
	"github.com/codeskyblue/go-sh"
)

const (
	NEIPARDIR = "/usr/local/bgp_neipar/"

	NOWCONNECT  = "/usr/local/bgp_neipar/.nowconn.txt"
	LASTCONNECT = "/usr/local/bgp_neipar/.lastconn.txt"

	NOWLIST = "/usr/local/bgp_neipar/raw.txt"

	DIFFDIR  = "/usr/local/bgp_neipar/diff/"
	NOWDIFF  = "/usr/local/bgp_neipar/diff/diff.txt"
	LASTDIFF = "/usr/local/bgp_neipar/diff/lastdiff.txt"

	RESULTDIR = "/usr/local/bgp_neipar/result/"
	RESULTCSV = "/usr/local/bgp_neipar/result/result.csv"
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func cat(filename string) string {
	buff, err := ioutil.ReadFile(filename)
	fatal(err)
	return string(buff)
}

func addog(text string, filename string) {
	var writer *bufio.Writer
	text_data := []byte(text)

	write_file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	writer = bufio.NewWriter(write_file)
	writer.Write(text_data)
	writer.Flush()
	fatal(err)
	defer write_file.Close()
}

func deleteLine(filename string) string {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File %s could not read: %v\n", filename, err)
		os.Exit(1)
	}

	defer f.Close()
	var b bytes.Buffer
	var x int
	x = 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := scanner.Text()
		if x == 0 {
			if !strings.Contains(s, "Neighbor") {
				s = ""
				b.WriteString(s)
				continue
			} else {
				s = ""
				x++
				continue
			}
		}
		s = s + "\n"
		b.WriteString(s)
		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}
	return b.String()
}

type Neighbor struct {
	Peer   string
	AS     int
	LastUP time.Duration
	Pfx    int
	Desc   string
}

type showNei []Neighbor

func initialZero(num string) bool {
	if numm, err := regexp.MatchString("^0", num); numm {
		return true
		if err != nil {
			return false
		}
	}
	return false
}

func timeconv(t string) time.Duration {
	if strings.Contains(t, "y") {
		t_ary := strings.Split(t, "y")
		t_years_int, err := strconv.Atoi(t_ary[0])
		fatal(err)
		t_ary[1] = strings.Trim(t_ary[1], "w")
		if initialZero(t_ary[1]) && t_ary[1] != "0" && t_ary[1] != "00" {
			t_ary[1] = strings.Trim(t_ary[1], "0")
		}
		t_weeks_int, err := strconv.Atoi(t_ary[1])
		fatal(err)
		t_weeks_int = t_weeks_int * 7 * 24
		t_years_int = t_years_int * 365 * 24
		t_int := (t_years_int + t_weeks_int)
		t = strconv.Itoa(t_int) + "h"
	} else if strings.Contains(t, "h") {
		t_ary := strings.Split(t, "d")
		t_days_int, err := strconv.Atoi(t_ary[0])
		fatal(err)
		t_days_int = t_days_int * 24
		t_ary[1] = strings.Trim(t_ary[1], "h")
		if initialZero(t_ary[1]) && t_ary[1] != "0" && t_ary[1] != "00" {
			t_ary[1] = strings.Trim(t_ary[1], "0")
		}
		t_hour_int, err := strconv.Atoi(t_ary[1])
		fatal(err)
		t_int := t_days_int + t_hour_int
		t = strconv.Itoa(t_int) + "h"
	} else if strings.Contains(t, "w") {
		t_ary := strings.Split(t, "w")
		t_week_int, err := strconv.Atoi(t_ary[0])
		fatal(err)
		t_ary[1] = strings.Trim(t_ary[1], "d")
		if initialZero(t_ary[1]) && t_ary[1] != "0" && t_ary[1] != "00" {
			t_ary[1] = strings.Trim(t_ary[1], "0")
		}
		t_days_int, err := strconv.Atoi(t_ary[1])
		fatal(err)
		t_week_int = t_week_int * 7
		t_int := (t_week_int + t_days_int) * 24
		t = strconv.Itoa(t_int) + "h"

	} else {
		t = strings.Replace(t, ":", "h", 1)
		t = strings.Replace(t, ":", "m", 1)
		t = t + "s"
	}
	d, err := time.ParseDuration(t)
	fatal(err)
	return d
}

func makeTmp(f string) showNei {
	NeighborLine := make([]Neighbor, 0)
	scanner := bufio.NewScanner(strings.NewReader(f))
	for scanner.Scan() {
		s := scanner.Text()
		if strings.Contains(s, " ") {
			s = strings.Replace(s, "      ", " ", -1)
			s = strings.Replace(s, "     ", " ", -1)
			s = strings.Replace(s, "    ", " ", -1)
			s = strings.Replace(s, "   ", " ", -1)
			s = strings.Replace(s, "  ", " ", -1)
			s_ary := strings.Split(s, " ")

			s_diff := s_ary[0] + " " + s_ary[2] + " " + s_ary[9] + " " + peer.Peer(s_ary[0]) + "\n"
			addog(s_diff, NOWDIFF)
			s_csv := s_ary[0] + "," + s_ary[2] + "," + s_ary[8] + "," + s_ary[9] + "\n"
			addog(s_csv, RESULTCSV)

			as_i, err := strconv.Atoi(s_ary[2])
			fatal(err)
			if s_ary[9] == "Active" {
				s_ary[9] = "-1"
			}
			if s_ary[9] == "Idle" {
				s_ary[9] = "-2"
			}
			pfx_i, err := strconv.Atoi(s_ary[9])
			fatal(err)
			NeighborLine = append(NeighborLine, Neighbor{s_ary[0], as_i, timeconv(s_ary[8]), pfx_i, peer.Peer(s_ary[0])})
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return NeighborLine
}

func (s showNei) Len() int {
	return len(s)
}

func (s showNei) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type ByPfx struct {
	showNei
}

func (b ByPfx) Less(i, j int) bool {
	return b.showNei[i].Pfx < b.showNei[j].Pfx
}

type ByAS struct {
	showNei
}

func (b ByAS) Less(i, j int) bool {
	return b.showNei[i].AS < b.showNei[j].AS
}

type ByLastUP struct {
	showNei
}

func (b ByLastUP) Less(i, j int) bool {
	return b.showNei[i].LastUP < b.showNei[j].LastUP
}

type ByDesc struct {
	showNei
}

func (b ByDesc) Less(i, j int) bool {
	return b.showNei[i].Desc < b.showNei[j].Desc
}

func shortDur(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}

func printNei(s showNei) {
	const format = "%v\t%v\t%v\t%v\t%v\t\n"
	tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(tw, format, "Peer", "AS", "LastUP", "Pfx/Stat", "Description")
	fmt.Fprintf(tw, format, "---------------", "------", "----------", "------", "------------")
	for _, t := range s {
		if t.Pfx == -1 {
			fmt.Fprintf(tw, format, t.Peer, t.AS, shortDur(t.LastUP), "Active", t.Desc)
		} else if t.Pfx == -2 {
			fmt.Fprintf(tw, format, t.Peer, t.AS, shortDur(t.LastUP), "Idle", t.Desc)
		} else {
			fmt.Fprintf(tw, format, t.Peer, t.AS, shortDur(t.LastUP), t.Pfx, t.Desc)
		}
	}
	tw.Flush() // calculate column widths and print table
}

func showAll(flag int) {
	if exists(RESULTCSV) {
		t := time.Now().Format("2006-01-02_030405")
		LASTCSV := RESULTDIR + "/result_" + t + ".csv"
		os.Rename(RESULTCSV, LASTCSV)
	}
	if exists(NOWDIFF) {
		os.Rename(NOWDIFF, LASTDIFF)
	}
	/*
		if !exists(NOWLIST){
			addog("test",NOWLIST)
		}
	*/

	getter.Showgetter()

	var s showNei = makeTmp(deleteLine(NOWLIST))

	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	if flag == 0 {
		sort.Sort(sort.Reverse(ByPfx{s}))
		fmt.Println("\n################## Sort by Pfx or Status ##################\n")
		printNei(s)
	} else if flag == 1 {
		sort.Sort(ByAS{s})
		fmt.Println("\n#################### Sort by AS Number ####################\n")
		printNei(s)
	} else if flag == 2 {
		sort.Sort(ByLastUP{s})
		fmt.Println("\n################### Sort by Last UP/Down ##################\n")
		printNei(s)
	} else if flag == 3 {
		sort.Sort(ByDesc{s})
		fmt.Println("\n################### Sort by Description ###################\n")
		printNei(s)

	}
	if exists(LASTCONNECT) {
		if cat(NOWCONNECT) == cat(LASTCONNECT) {
			fmt.Println("\n################ diff Now and Last show cmd ###############\n")
			sh.Command("colordiff", "-u", LASTDIFF, NOWDIFF).Run()
			fmt.Println("\n")
		}
	}

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	//fmt.Print(out)
	cmd := sh.Command("less", "-R")
	cmd.Stdin = strings.NewReader(out)
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	fatal(err)

}

var aflag bool
var pflag bool
var uflag bool
var dflag bool

func init() {
	flag.BoolVar(&pflag, "p", false, "Sort Recieved-routes Number.")
	flag.BoolVar(&aflag, "a", false, "Sort AS Number.")
	flag.BoolVar(&uflag, "u", false, "Sort LastUP Neighbor.")
	flag.BoolVar(&dflag, "d", false, "Sort Description.")
}

func main() {
	if !exists(NEIPARDIR) {
		os.MkdirAll(NEIPARDIR, 0777)
	}
	if !exists(DIFFDIR) {
		os.MkdirAll(DIFFDIR, 0777)
	}
	if !exists(RESULTDIR) {
		os.MkdirAll(RESULTDIR, 0777)
	}
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Parse()
	if pflag {
		showAll(0)
	} else if aflag {
		showAll(1)
	} else if uflag {
		showAll(2)
	} else if dflag {
		showAll(3)
	}
}
