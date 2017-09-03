package getter

import (
	//	"bufio"

	"fmt"

	"github.com/golang/glog"
	expect "github.com/google/goexpect"
	"github.com/google/goterm/term"
	//	"github.com/joho/godotenv"
	"bufio"
	"io/ioutil"
	"os"

	"github.com/ziutek/telnet"
	"golang.org/x/crypto/ssh/terminal"
	//	"regexp"
	"strings"
	"syscall"
	"time"
)

const (
	NOWCONNECT  = "/usr/local/bgp_neipar/.nowconn.txt"
	LASTCONNECT = "/usr/local/bgp_neipar/.lastconn.txt"
	network     = "tcp"
	timeout     = 10 * time.Second
	command     = "show bgp ipv4 unicast summary"
	outputfile  = "/usr/local/bgp_neipar/raw.txt"
)

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func dog(text string, filename string) {
	text_data := []byte(text)
	err := ioutil.WriteFile(filename, text_data, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
}

func cat(filename string) string {
	buff, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	return string(buff)
}

/*
func IsIP(ip string) bool {
	if ipm, err := regexp.MatchString("^(([1-9]?[0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([1-9]?[0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$", ip); ipm {
		return true
		if err != nil {
			return false
		}
	}
	return false
}
*/
func Showgetter(outfile string) {
	//flag.Parse()

	var address string
	fmt.Println(term.Bluef("Getting %s", command))
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("router address: ")
	address, _ = reader.ReadString('\n')
	address = strings.Trim(address, "\n")
	//		if !IsIP(address) {
	//			fmt.Printf("%s is bad address format\n", address)
	//			continue
	//		}
	//		break
	//	}
	address = address + `:23`
	fmt.Print("username: ")
	username, _ := reader.ReadString('\n')
	username = strings.Trim(username, "\n")
	fmt.Print("password: ")
	PassWord, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(PassWord)
	password = strings.Trim(password, "\n")
	fmt.Print("\n")
	exp, _, err := telnetSpawn(address, timeout, expect.Verbose(true))
	if err != nil {
		glog.Exitf("telnetSpawn(%q,%v) failed: %v", address, timeout, err)
	}

	defer func() {
		if err := exp.Close(); err != nil {
			glog.Infof("exp.Close failed: %v", err)
		}
	}()

	res, err := exp.ExpectBatch([]expect.Batcher{
		/*
			&expect.BExp{R: `\n\.`},
			&expect.BSnd{S: command + "\r\n"},
			&expect.BExp{R: `\n\.`},
		*/
		//For example, IOS-XR series.
		&expect.BExp{R: `Username:`},
		&expect.BSnd{S: username + "\r\n"},
		&expect.BExp{R: `Password:`},
		&expect.BSnd{S: password + "\r\n"},
		&expect.BExp{R: `#`},
		&expect.BSnd{S: "terminal length 0" + "\r\n"},
		&expect.BExp{R: `#`},
		&expect.BSnd{S: command + "\r\n"},
		&expect.BExp{R: `#`},
		&expect.BSnd{S: "exit" + "\r\n"},
	}, timeout)
	if err != nil {
		glog.Exitf("exp.ExpectBatch failed: %v , res: %v", err, res)
	}
	//        fmt.Println(term.Greenf("Res: %s", res[len(res)-1].Output))
	if exists(LASTCONNECT) {
		os.Rename(NOWCONNECT,LASTCONNECT)
	}
	address = strings.Trim(address,":23")
	dog(address, NOWCONNECT)
	dog(res[len(res)-1].Output, outfile)

}

func telnetSpawn(addr string, timeout time.Duration, opts ...expect.Option) (expect.Expecter, <-chan error, error) {
	conn, err := telnet.Dial(network, addr)
	if err != nil {
		return nil, nil, err
	}

	resCh := make(chan error)

	return expect.SpawnGeneric(&expect.GenOptions{
		In:  conn,
		Out: conn,
		Wait: func() error {
			return <-resCh
		},
		Close: func() error {
			close(resCh)
			return conn.Close()
		},
		Check: func() bool { return true },
	}, timeout, opts...)
}
