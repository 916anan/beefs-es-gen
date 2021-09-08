package main

import (
	"context"
	"fmt"
	"github.com/bingoohuang/gg/pkg/chinaid"
	"github.com/bingoohuang/gg/pkg/ctl"
	"github.com/bingoohuang/gg/pkg/flagparse"
	"github.com/bingoohuang/gg/pkg/randx"
	"github.com/bingoohuang/gg/pkg/rest"
	"github.com/bingoohuang/jj"
	"github.com/prometheus/common/log"
	"golang.org/x/sync/semaphore"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Config struct {
	Filer     string `flag:"f" val:"http://localhost:9335" usage:"Filer 地址"`
	Bucket    string `flag:"b" val:"" usage:"桶"`
	Port      string `flag:"p" val:"9337" usage:"端口"`
	Version   bool   `flag:"v" usage:"Print version info and exit"`
	Thread    int64  `flag:"t" val:"8" usage:"线程数"`
	Num       int    `flag:"n" val:"1" usage:"数量"`
	BasicAuth string `flag:"u" val:"scott:tiger" usage:"数量"`
}

func main() {
	c := &Config{}
	flagparse.Parse(c)
	ctl.Config{PrintVersion: c.Version}.ProcessInit()

	s := semaphore.NewWeighted(c.Thread)
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(c.Num)
	for i := 0; i < c.Num; i++ {
		go func() {
			s.Acquire(context.Background(), 1)
			uploadRandomFile(c)
			s.Release(1)
			waitGroup.Done()
		}()
	}
	waitGroup.Wait()
	log.Infof("上传完成 %+v", c)
}

func uploadRandomFile(c *Config) {
	fileName, data := randomImg()
	var addr string
	if c.Bucket == "" {
		addr = fmt.Sprintf("%s/%s", c.Filer, fileName)
	} else {
		addr = fmt.Sprintf("%s/%s/%s/%s", c.Filer, "buckets", c.Bucket, fileName)
	}

	upload, err := rest.Rest{
		Headers:   randomHeaders(),
		BasicAuth: c.BasicAuth,
		Addr:      addr,
		Timeout:   10 * time.Second,
	}.Upload(fileName, data)
	if err != nil || upload.Status > 300 {
		log.Errorln("上传文件失败", err, upload.Status)
	}
}

func randomImg() (string, []byte) {
	config := randx.ImgConfig{
		Width:      1024,
		Height:     1024,
		RandomText: "Tom",
		FastMode:   true,
		PixelSize:  256,
	}
	data, _ := config.Gen(".jpg")

	return jj.Regex(`[A-Za-z0-9]{10}`).(string) + ".jpg", data

}

func randomHeaders() map[string]string {
	colors := strings.Split("red,orange,green,blue,yellow,purple", ",")
	headers := make(map[string]string)
	headers["beefs-color"] = colors[randx.IntBetween(0, len(colors)-1)]
	headers["beefs-mobile"] = chinaid.Mobile()
	headers["beefs-chinaid"] = chinaid.ChinaID()
	headers["beefs-weight"] = strconv.Itoa(randx.IntBetween(1, 99))
	headers["beefs-doctime"] = jj.RandomTime("yyyy-MM-dd").(string)
	return headers
}
