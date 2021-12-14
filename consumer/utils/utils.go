package utils

import (
  "log"
  "os"
  "image"
  "image/jpeg"
  _"image/png"
  _"image/gif"

  "github.com/nfnt/resize"
)

const imagePath = "C:\\Games\\"

var imageMap = map[string]struct{width uint}{"large":{width:200}, "medium":{width:150}, "small":{width:60}}

func FailOnError(err error, msg string) {
   if err != nil {
     log.Panicf("%s: %s", err, msg)
   }
}

func Resize(fileName, option string) {

  file, err := os.Open(imagePath + fileName)
  FailOnError(err, "Failed to open file " + fileName)
  image, _, err := image.Decode(file)
  file.Close()
  FailOnError(err, "Failed to decode image")
  size := imageMap[option]
  newImage := resize.Resize(size.width, 0, image, resize.Lanczos3)
  dst, err := os.Create(imagePath + option + "_" + fileName)
  FailOnError(err, "could create file")
  defer dst.Close()
  err = jpeg.Encode(dst, newImage, nil)
}
