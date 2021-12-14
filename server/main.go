package main

import (
  "fmt"
  "net/http"
  "log"
  "path"
  "io"
  "os"
  "math/rand"

  "github.com/streadway/amqp"
)
const imagePath = "C:\\Games\\"

func exitOnError(err error, msg string) {
  if err != nil {
    log.Panicf("%s: %s", err,msg)
  }
}


func imageHandler(w http.ResponseWriter, r *http.Request)  {

 if r.Method != "POST" {
   w.WriteHeader(http.StatusNotFound)
   fmt.Fprint(w, "Not found")
   return
 }

  r.ParseMultipartForm(32 << 20)

  file, meta, err := r.FormFile("image")
  if err != nil  {
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprint(w,"Error loading image")
  }

  if meta.Size > 5000000  {
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprint(w,"File is too big")
    return
  }

  ext := path.Ext(meta.Filename)
  acceptedFiles := [4]string{".jpg",".png",".jpeg",".gif"}
  isValidExt := false
  for _, val := range acceptedFiles {
    if val == ext {
      isValidExt = true
      break
    }
  }

  if !isValidExt {
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprint(w,"Invalid image format")
    return
  }

  imageName := generateName(15) + path.Ext(meta.Filename)

  newFile, err := os.Create(imagePath + imageName)
  if err != nil  {
    log.Println(err)
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Fprint(w,"Something went wrong")
    return
  }
  defer newFile.Close()

  io.Copy(newFile, file)

  conn, err :=  amqp.Dial("amqp://guest:guest@localhost:5672/")
  exitOnError(err,"Failed to connect to rabbitmq server")
  defer conn.Close()
  ch, err := conn.Channel()
  exitOnError(err, "Failed to create a channle")
  defer ch.Close()

  q, err := ch.QueueDeclare(
    "image-queue",
    true,
    false,
    false,
    false,
    nil,
  )
  exitOnError(err, "Failed to declare queue")

  err =  ch.Publish(
     "",
     q.Name,
     false,
     false,
     amqp.Publishing{
                       DeliveryMode: amqp.Persistent,
                       ContentType:  "text/plain",
                       Body:         []byte(imageName),
               })
  exitOnError(err, "Failed to public message to queue")

  w.WriteHeader(http.StatusOK)
  w.Write([]byte("Image uploaded successfully"))
}


func main() {

 http.HandleFunc("/image", imageHandler)
 log.Fatal(http.ListenAndServe(":8080", nil))

}

func generateName(n int) string {
  chars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
  var str string
  for i := 0; i < n; i++ {
     p := rand.Intn(len(chars))
     str += string(chars[p])
  }
  return str
}
