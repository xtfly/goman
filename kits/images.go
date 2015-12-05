package kits

import (
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
)

// 水印
func cmd_watermark(file string, to string) {
	// 打开图片并解码
	file_origin, _ := os.Open(file)
	origin, _ := jpeg.Decode(file_origin)
	defer file_origin.Close()

	// 打开水印图并解码
	file_watermark, _ := os.Open("watermark.png")
	watermark, _ := png.Decode(file_watermark)
	defer file_watermark.Close()

	//原始图界限
	origin_size := origin.Bounds()

	//创建新图层
	canvas := image.NewNRGBA(origin_size)
	// 贴原始图
	draw.Draw(canvas, origin_size, origin, image.ZP, draw.Src)
	// 贴水印图
	draw.Draw(canvas, watermark.Bounds().Add(image.Pt(30, 30)), watermark, image.ZP, draw.Over)

	//生成新图片
	create_image, _ := os.Create(to)
	jpeg.Encode(create_image, canvas, &jpeg.Options{Quality: 95})
	defer create_image.Close()
}
