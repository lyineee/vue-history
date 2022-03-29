package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type ReleaseData struct {
	TagName string `json:"tag_name"`
}

var (
	etag        = ""
	delayTime   = 2
	extractPath = "/usr/share/nginx/html"
	dlFilename  = "dist.tar.gz"
	ghproxy     = "https://ghproxy.com/"
)

func main() {
	mainLoop()
}

func mainLoop() {
	for {
		time.Sleep(time.Duration(delayTime) * time.Second)
		retry := 5
		req, err := http.NewRequest("GET", "https://api.github.com/repos/lyineee/vue-history/releases/latest", nil)
		if err != nil {
			log.Println("new request error: ", err)
			continue
		}
		req.Header.Add("If-None-Match", fmt.Sprintf(`"%s"`, etag))
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode == 304 {
			//log.Println("resources Not Modified")
			continue
		}
		if resp.StatusCode == 403 {
			log.Println("API rate limit exceeded", " ratelimit reset: ", resp.Header.Get("x-ratelimit-reset"))
			if resetTime := resp.Header.Get("x-ratelimit-reset"); resetTime != "" {
				unixTime, err := strconv.Atoi(resetTime)
				if err != nil {
					log.Println("err parse reset time to int, reset time string: ", resetTime)
				} else {
					resetTime := time.Unix(int64(unixTime), 0)
					sleepTime := resetTime.Sub(time.Now())
					log.Printf("sleep until reset time %s, sleep time: %s", resetTime, sleepTime)
					time.Sleep(sleepTime)
				}

			}
			continue
		}

		etagHeader := resp.Header["Etag"]
		if len(etagHeader) != 1 {
			log.Println("no etag header")
			continue
		}
		newEtag := etagHeader[0][3 : len(etagHeader[0])-1]
		log.Println("etag: ", newEtag)
		// if newEtag == etag {
		// 	continue
		// }
		for re := retry; re > 0; re-- {
			time.Sleep(time.Duration(delayTime) * time.Second)
			req, err := http.NewRequest("GET", "https://api.github.com/repos/lyineee/vue-history/releases/latest", nil)
			if err != nil {
				log.Println("new content request error: ", err)
				continue
			}
			req.Header.Add("Accept", "application/vnd.github")
			resp, err := client.Do(req)
			if err != nil {
				log.Println("get content request content: ", err)
			}
			release := ReleaseData{}
			err = json.NewDecoder(resp.Body).Decode(&release)
			if err != nil {
				log.Println("json recode error: ", err)
			}
			log.Println("release: ", release.TagName)
			err = downloadAndExtract(extractPath, fmt.Sprintf("%shttps://github.com/lyineee/vue-history/releases/download/%s/%s", ghproxy, release.TagName, dlFilename))
			if err != nil {
				log.Println("download and extract file err: ", err)
				continue
			}
			etag = newEtag
			break
		}
	}
}

func downloadAndExtract(dst, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gzipReader.Close()
	tarReader := tar.NewReader(gzipReader)

	for {
		hdr, err := tarReader.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case hdr == nil:
			continue
		}

		// 处理下保存路径，将要保存的目录加上 header 中的 Name
		// 这个变量保存的有可能是目录，有可能是文件，所以就叫 FileDir 了……
		baseLength := len(filepath.Ext(filepath.Ext(dlFilename)))
		dstFileDir := filepath.Join(dst, hdr.Name[baseLength+1:])
		// 根据 header 的 Typeflag 字段，判断文件的类型
		switch hdr.Typeflag {
		case tar.TypeDir: // 如果是目录时候，创建目录
			// 判断下目录是否存在，不存在就创建
			if b := ExistDir(dstFileDir); !b {
				// 使用 MkdirAll 不使用 Mkdir ，就类似 Linux 终端下的 mkdir -p，
				// 可以递归创建每一级目录
				if err := os.MkdirAll(dstFileDir, 0775); err != nil {
					return err
				}
				log.Printf("make directory: %s", dstFileDir)
			}
		case tar.TypeReg: // 如果是文件就写入到磁盘
			// 创建一个可以读写的文件，权限就使用 header 中记录的权限
			// 因为操作系统的 FileMode 是 int32 类型的，hdr 中的是 int64，所以转换下
			file, err := os.OpenFile(dstFileDir, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			n, err := io.Copy(file, tarReader)
			if err != nil {
				return err
			}
			// 将解压结果输出显示
			log.Printf("成功解压： %s , 共处理了 %d 个字符\n", dstFileDir, n)

			// 不要忘记关闭打开的文件，因为它是在 for 循环中，不能使用 defer
			// 如果想使用 defer 就放在一个单独的函数中
			file.Close()
		}
	}
}

// 判断目录是否存在
func ExistDir(dirname string) bool {
	fi, err := os.Stat(dirname)
	return (err == nil || os.IsExist(err)) && fi.IsDir()
}
