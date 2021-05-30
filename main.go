package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func main()  {
	var outPath string
	var dataPath string

	flag.StringVar(&outPath, "outPath", "", "输入目录，默认为空")
	flag.StringVar(&dataPath, "dataPath", "", "数据目录，默认为空")

	flag.Parse()
	if (len(outPath) == 0 || len(dataPath) == 0) {
		logger("参数都要填噢 郭郭!!")
		return
	}

	var files []string
	files, err := GetAllFile(dataPath, files)
	if err != nil {
		logger("读取文件出错")
		return
	}

	for _, file := range files {
		length := strings.Index(file, "_")
		outPutDirSeg := file[0:length]
		desPath := outPath + "/" + outPutDirSeg
		_, err := os.Stat(desPath)

		if os.IsNotExist(err) {
			ret := getGenArray(outPutDirSeg, files)
			file1 := ret[0]
			file2 := ret[1]

			cmd := "shovill --trim --outdir " + outPath + "/" + outPutDirSeg + " --force --R1 " + dataPath  + "/" + file1 + " --R2 " + dataPath + "/" + file2

			logger(cmd)
			Command(cmd)
		}

	}



}

func Command(cmd string) error {
	c := exec.Command("bash", "-c", cmd)
	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	go func() {
		reader := bufio.NewReader(stdout)
		for {
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				break
			}
			fmt.Print(readString)
		}
	}()
	return c.Run()
}

func getGenArray(seg string, files []string) ([]string) {
	var ret []string
	for _, file := range files {
		length := strings.Index(file, "_")
		outPutDirSeg := file[0:length]
		if seg == outPutDirSeg {
			ret = append(ret, file)
		}
	}

	sort.Sort(sort.StringSlice(ret))

	return ret
}

func GetAllFile(pathname string, s []string) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := pathname + "/" + fi.Name()
			s, err = GetAllFile(fullDir, s)
			if err != nil {
				fmt.Println("read dir fail:", err)
				return s, err
			}
		} else {
			//仅返回文件名
			//fullName := pathname + "/" + fi.Name()
			fullName := fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}

func logger(logthis string) {
	filePath := "logfile.log"
	fmt.Println(logthis)
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		_, err := os.Create(filePath)
		if err != nil {
			fmt.Println(err)
		}

	}
	file, err := os.OpenFile(filePath, os.O_WRONLY | os.O_TRUNC, 0666)
	// file, err := os.OpenFile(filePath, os.O_WRONLY | os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("open file err=%v\n", err)
		return
	}

	defer file.Close()

	now := time.Now()
	time := now.Format("2006/01/02 15:04:05")
	str := fmt.Sprintf("%v %v ", time, logthis)
	writer := bufio.NewWriter(file)
	writer.WriteString(str)

	writer.Flush()
}
