package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var (
	style           string
	a               float64
	b               float64
	c               float64
	d               float64
	min             int
	max             int
	concurrentLimit int
	array1          = []float64{75, 60, 45, 30, 15, 0}
	array2          = [][]float64{
		{-0.0015702102444, 111320.7020616939, 1704480524535203, -10338987376042340, 26112667856603880, -35149669176653700, 26595700718403920, -10725012454188240, 1800819912950474, 82.5},
		{0.0008277824516172526, 111320.7020463578, 647795574.6671607, -4082003173.641316, 10774905663.51142, -15171875531.51559, 12053065338.62167, -5124939663.577472, 913311935.9512032, 67.5},
		{0.00337398766765, 111320.7020202162, 4481351.045890365, -23393751.19931662, 79682215.47186455, -115964993.2797253, 97236711.15602145, -43661946.33752821, 8477230.501135234, 52.5},
		{0.00220636496208, 111320.7020209128, 51751.86112841131, 3796837.749470245, 992013.7397791013, -1221952.21711287, 1340652.697009075, -620943.6990984312, 144416.9293806241, 37.5},
		{-0.0003441963504368392, 111320.7020576856, 278.2353980772752, 2485758.690035394, 6070.750963243378, 54821.18345352118, 9540.606633304236, -2710.55326746645, 1405.483844121726, 22.5},
		{-0.0003218135878613132, 111320.7020701615, 0.00369383431289, 823725.6402795718, 0.46104986909093, 2351.343141331292, 1.58060784298199, 8.77738589078284, 0.37238884252424, 7.45},
	}
)

type LatLngPoint struct {
	Lat, Lng float64
}

type PointF struct {
	X, Y float64
}

func init() {
	// 默认值
	style = "normal"
	a = 116.199599
	b = 40.033261
	c = 116.537074
	d = 39.830986
	min = 3
	max = 19
	concurrentLimit = 50

	// 用户交互式输入
	fmt.Println("地图风格")
	fmt.Println("常规 normal")
	fmt.Println("清新蓝 light")
	fmt.Println("黑夜 dark")
	fmt.Println("红色警戒 redalert")
	fmt.Println("精简 googlelite")
	fmt.Println("自然绿 grassgreen")
	fmt.Println("午夜蓝 midnight")
	fmt.Println("浪漫粉 pink")
	fmt.Println("青春绿 darkgreen")
	fmt.Println("清新蓝绿 bluish")
	fmt.Println("高端灰 grayscale")
	fmt.Println("强边界 hardedge")

	fmt.Printf("地图风格 (默认 %s): ", style)
	styleInput := ""
	fmt.Scanln(&styleInput)
	if styleInput != "" {
		style = styleInput
	}

	fmt.Printf("百度地图左上角经度 (默认 %.10f): ", a)
	aInput := ""
	fmt.Scanln(&aInput)
	if aInput != "" {
		a, _ = strconv.ParseFloat(aInput, 64)
	}

	fmt.Printf("百度地图左上角纬度 (默认 %.10f): ", b)
	bInput := ""
	fmt.Scanln(&bInput)
	if bInput != "" {
		b, _ = strconv.ParseFloat(bInput, 64)
	}

	fmt.Printf("百度地图右下角经度 (默认 %.10f): ", c)
	cInput := ""
	fmt.Scanln(&cInput)
	if cInput != "" {
		c, _ = strconv.ParseFloat(cInput, 64)
	}

	fmt.Printf("百度地图右下角纬度 (默认 %.10f): ", d)
	dInput := ""
	fmt.Scanln(&dInput)
	if dInput != "" {
		d, _ = strconv.ParseFloat(dInput, 64)
	}

	fmt.Printf("最小层级 (默认 %d): ", min)
	minInput := ""
	fmt.Scanln(&minInput)
	if minInput != "" {
		min, _ = strconv.Atoi(minInput)
	}

	fmt.Printf("最大层级 (默认 %d): ", max)
	maxInput := ""
	fmt.Scanln(&maxInput)
	if maxInput != "" {
		max, _ = strconv.Atoi(maxInput)
	}

	fmt.Printf("最大并发请求数量 (默认 %d): ", concurrentLimit)
	concurrentInput := ""
	fmt.Scanln(&concurrentInput)
	if concurrentInput != "" {
		concurrentLimit, _ = strconv.Atoi(concurrentInput)
	}
}

func main() {
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrentLimit)

	for z := min; z <= max; z++ {
		x1, y1 := getID(a, b, z)
		x2, y2 := getID(c, d, z)

		// 确保 x2 > x1 和 y2 > y1
		if x1 > x2 {
			x1, x2 = x2, x1
		}
		if y1 > y2 {
			y1, y2 = y2, y1
		}

		fmt.Printf("\n缩放级别 %d: 瓦片范围: x 从 %d 到 %d, y 从 %d 到 %d\n", z, x1, x2, y1, y2)
		totalTiles := (x2 - x1 + 1) * (y2 - y1 + 1)
		fmt.Printf("缩放级别 %d: 瓦片总数: %d\n", z, totalTiles)

		for x := x1; x <= x2; x++ {
			for y := y1; y <= y2; y++ {
				wg.Add(1)
				sem <- struct{}{}
				go func(x, y, z int) {
					defer wg.Done()
					defer func() { <-sem }()
					downloadTile(x, y, z)
				}(x, y, z)
			}
		}
	}

	wg.Wait()
	fmt.Println("下载完成")
}

func LatLng2Mercator(p LatLngPoint) PointF {
	var arr []float64
	if p.Lat > 74 {
		p.Lat = 74
	} else if p.Lat < -74 {
		p.Lat = -74
	}

	for i, latLimit := range array1 {
		if p.Lat >= latLimit {
			arr = array2[i]
			break
		}
	}
	if arr == nil {
		for i := len(array1) - 1; i >= 0; i-- {
			if p.Lat <= -array1[i] {
				arr = array2[i]
				break
			}
		}
	}

	res := Convertor(p.Lng, p.Lat, arr)
	return PointF{res[0], res[1]}
}

func Convertor(x, y float64, param []float64) []float64 {
	T := param[0] + param[1]*math.Abs(x)
	cC := math.Abs(y) / param[9]
	cF := param[2] + param[3]*cC + param[4]*cC*cC + param[5]*cC*cC*cC + param[6]*cC*cC*cC*cC + param[7]*cC*cC*cC*cC*cC + param[8]*cC*cC*cC*cC*cC*cC
	T *= sign(x)
	cF *= sign(y)
	return []float64{T, cF}
}

func sign(value float64) float64 {
	if value < 0 {
		return -1
	}
	return 1
}

func getTileIndices(x, y float64, zoom int) (int, int) {
	scale := math.Pow(2, float64(18-zoom))
	pixelX := x / scale
	pixelY := y / scale
	return int(pixelX) / 256, int(pixelY) / 256
}

func getID(lng, lat float64, z int) (int, int) {
	mercator := LatLng2Mercator(LatLngPoint{Lat: lat, Lng: lng})
	return getTileIndices(mercator.X, mercator.Y, z)
}

func downloadTile(x, y, z int) {
	url := fmt.Sprintf("https://api.map.baidu.com/customimage/tile?&x=%d&y=%d&z=%d&customid=%s", x, y, z, style)
	filePath := fmt.Sprintf("map/%d/%d/%d.png", z, x, y)

	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		fmt.Printf("创建目录失败: %v\n", err)
		return
	}

	for {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("下载失败: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("下载失败, 状态码: %d\n", resp.StatusCode)
			time.Sleep(1 * time.Second)
			continue
		}

		file, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("创建文件失败: %v\n", err)
			return
		}
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			fmt.Printf("保存文件失败: %v\n", err)
			time.Sleep(1 * time.Second)
			continue
		}

		fmt.Printf("下载成功: %s\n", filePath)
		break
	}
}
