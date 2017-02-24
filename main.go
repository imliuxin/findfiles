package main

import (
	"flag"
	"fmt"

	"os"
	"strings"
"bufio"
"io"
"strconv"
)

var (
	//gConf   config.Configer
	cfgFile = flag.String("c", "", "config file")
	// for log
	logPath string
	prefix  string

	// for sha1 scan
	outputChan     = make(chan string, 100)
	syncChan       = make(chan int, 100)
	quitChan       = make(chan int)
	walkPath       string
	filterDirName  string
	filterFileName string
	resultFile     string

	// concurrentFileNumber
	concurrentNumber int64
	concurrentChan   chan int
)
const middle = "========="

type Config struct {
	Mymap  map[string]string
	strcet string
}

func (c *Config) InitConfig(path string) {
	c.Mymap = make(map[string]string)

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		s := strings.TrimSpace(string(b))
		//fmt.Println(s)
		if strings.Index(s, "#") == 0 {
			continue
		}

		n1 := strings.Index(s, "[")
		n2 := strings.LastIndex(s, "]")
		if n1 > -1 && n2 > -1 && n2 > n1+1 {
			c.strcet = strings.TrimSpace(s[n1+1 : n2])
			continue
		}

		if len(c.strcet) == 0 {
			continue
		}
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}

		frist := strings.TrimSpace(s[:index])
		if len(frist) == 0 {
			continue
		}
		second := strings.TrimSpace(s[index+1:])

		pos := strings.Index(second, "\t#")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " #")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, "\t//")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " //")
		if pos > -1 {
			second = second[0:pos]
		}

		if len(second) == 0 {
			continue
		}

		key := c.strcet + middle + frist
		c.Mymap[key] = strings.TrimSpace(second)
	}
}

func (c Config) Read(node, key string) string {

	key = node + middle + key
println(key)
	v, found := c.Mymap[key]
	if !found {
		return ""
	}
	return v
}

func GetBetweenStr(str, start, end string) string {
    n := strings.Index(str, start)
    if n == -1 {
        n = 0
    }
    str = string([]byte(str)[n:])
    m := strings.Index(str, end)
    if m == -1 {
        m = len(str)
    }
    str = string([]byte(str)[:m])
    return str
}
func Substr(str string, start, length int) string {
    rs := []rune(str)
    rl := len(rs)
    end := 0

    if start < 0 {
        start = rl - 1 + start
    }
    end = start + length

    if start > end {
        start, end = end, start
    }

    if start < 0 {
        start = 0
    }
    if start > rl {
        start = rl
    }
    if end < 0 {
        end = 0
    }
    if end > rl {
        end = rl
    }

    return string(rs[start:end])
}
func main() {
	flag.Parse()
	if *cfgFile == "" {
		fmt.Printf("Usage: %s -c=conf/findfiles.conf", os.Args[0])
println("\n")
		os.Exit(1)
	}

	err := Init(*cfgFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	concurrentChan = make(chan int, concurrentNumber)

	defer func() {

	}()

println(walkPath)
println(filterDirName)
println(filterFileName)
println(resultFile)
	dirScanner, err := NewDirScanner(walkPath, filterDirName, filterFileName, resultFile)
	if err != nil {
		
		return
	}
	go dirScanner.fileStore()
println("get")
	// 获取文件列表
	dirScanner.ScanFileInfo()
	<-quitChan
	
	// 生成文件
}

func Init(fileName string) (err error) {
	myConfig := new(Config)
	myConfig.InitConfig("conf/dirscan.conf")
	fmt.Println(myConfig.Read("path", "concurrentNumber"))
	//fmt.Printf("%v", myConfig.Mymap)

println("44444")
/*
    f, err := os.Open("conf/dirscan.conf")
    if err != nil {
        panic(err)
    }
    defer f.Close()

    rd := bufio.NewReader(f)
    for {
        line, err := rd.ReadString('\n') //以'\n'为结束符读入一行
        //println(line)
if err != nil || io.EOF == err {
            break
        }
if strings.Contains(line, "concurrentNumber")==true{

}

}

	//gConf, err = config.NewConfiger(fileName)
	if err != nil {
		return
	}
*/
/*
	logPath, err = gConf.GetSetting("log_conf", "logdir")
	if err != nil {
		return err
	}
	prefix, err = gConf.GetSetting("log_conf", "prefix")
	if err != nil {
		return err
	}
	walkPath, err = gConf.GetSetting("path", "walkPath")
	if err != nil {
		return
	}
*/
err=nil
logPath=myConfig.Read("log_conf", "logdir")
prefix=myConfig.Read("log_conf", "prefix");//"dirscan"
walkPath=myConfig.Read("path", "walkPath");;//"/usr/src/dirscan_tool-master"
filterDirName=myConfig.Read("path", "filterDir");;//"src/";//"ratelimit.[v]?"
filterFileName=myConfig.Read("path", "filterFile");//"*.go"
b,error := strconv.Atoi(myConfig.Read("path", "concurrentNumber"))
println(error)
concurrentNumber=int64( b);//10
resultFile=myConfig.Read("path", "resultFile");//"/tmp/sha1.out"
println("55555")
println(logPath)
println(prefix)
println(walkPath)
println(filterDirName)
println(filterFileName)
println(concurrentNumber)
/*
	// currentFile
	if strings.Contains(walkPath, `~`) {
		home := os.Getenv("HOME")
		if len(home) > 0 {
			walkPath = strings.Replace(walkPath, `~`, home, -1)
		}
	}
	fmt.Println(walkPath)

	filterDirName, err = gConf.GetSetting("path", "filterDir")
	if err != nil {
		return
	}
	filterFileName, err = gConf.GetSetting("path", "filterFile")
	if err != nil {
		return
	}
	concurrentNumber, err = gConf.GetIntSetting("path", "concurrentNumber", 64)
	if err != nil {
		return
	}
	resultFile, err = gConf.GetSetting("path", "resultFile")
	if err != nil {
		return
	}
*/

	return
}
