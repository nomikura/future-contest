package future

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"cloud.google.com/go/storage"
	"google.golang.org/appengine"
	"google.golang.org/appengine/file"
	"google.golang.org/appengine/log"
)

type demo struct {
	client     *storage.Client
	bucketName string
	bucket     *storage.BucketHandle

	w   io.Writer
	ctx context.Context
	// cleanUp is a list of filenames that need cleaning up at the end of the demo.
	cleanUp []string
	// failed indicates that one or more of the demo steps failed.
	failed bool
}

func (atcoder *ContestModule) FileIO(operate string) {
	var r *http.Request = atcoder.Context.Request
	var w http.ResponseWriter = atcoder.Context.Writer
	ctx := appengine.NewContext(r)

	// デフォルトのバケットを指定する(App Engineのコンテストから取得できる)
	bucket, err := file.DefaultBucketName(ctx)
	if err != nil {
		log.Errorf(ctx, "Faild to get default GCS bucket name: %v", err)
	}

	// clientをつくる
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Errorf(ctx, "Faild to create client: %v", err)
	}
	defer client.Close()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	buf := &bytes.Buffer{}
	d := &demo{
		w:          buf,
		ctx:        ctx,
		client:     client,
		bucket:     client.Bucket(bucket),
		bucketName: bucket,
	}

	fileName := "demo-testfile-go"

	if operate == "write" {
		// コンテストデータをバイナリにエンコード
		var binaryData []byte
		Encode(atcoder.AllContest, &binaryData)
		// ファイルに書き込む
		d.createFile(fileName, binaryData)
		log.Infof(ctx, "Write to file")
	} else if operate == "read" {
		// ファイルをバイナリで読み込む
		var binaryData []byte
		d.readFile(fileName, &binaryData)
		// バイナリをコンテストデータにデコード
		Decode(binaryData, &atcoder.AllContest)
		log.Infof(ctx, "Read file")
	} else {
		log.Infof(ctx, "operation is not read and write")
	}

	if d.failed {
		w.WriteHeader(http.StatusInternalServerError)
		buf.WriteTo(w)
	} else {
		w.WriteHeader(http.StatusOK)
		buf.WriteTo(w)
	}
}

//[START write]
func (d *demo) createFile(fileName string, byteDataToWrite []byte) {
	wc := d.bucket.Object(fileName).NewWriter(d.ctx)
	wc.ContentType = "text/plain"
	wc.Metadata = map[string]string{
		"x-goog-meta-foo": "foo",
		"x-goog-meta-bar": "bar",
	}
	d.cleanUp = append(d.cleanUp, fileName)

	// 書き込む
	if _, err := wc.Write(byteDataToWrite); err != nil {
		d.errorf("createFile: unable to write data to bucket %q, file %q: %v", d.bucketName, fileName, err)
		return
	}

	// ファイル閉じてるのかな？これがないと書き込めない
	if err := wc.Close(); err != nil {
		d.errorf("createFile: unable to close bucket %q, file %q: %v", d.bucketName, fileName, err)
		return
	}
}

//[END write]

//[START read]
func (d *demo) readFile(fileName string, data *[]byte) {
	// ファイルを開く
	rc, err := d.bucket.Object(fileName).NewReader(d.ctx)
	if err != nil {
		d.errorf("readFile: unable to open file from bucket %q, file %q: %v", d.bucketName, fileName, err)
		return
	}
	defer rc.Close()

	// データを読み込む
	slurp, err := ioutil.ReadAll(rc)
	if err != nil {
		d.errorf("readFile: unable to read data from bucket %q, file %q: %v", d.bucketName, fileName, err)
		return
	}

	*data = slurp
}

//[END read]

// dataをbyteArrayにエンコードする
func Encode(data interface{}, byteArray *[]byte) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		// log.Print("encode: ", err)
	}

	*byteArray = buffer.Bytes()
}

// byteArrayをdataにデコードする(dataはポインタ型)
func Decode(byteArray []byte, data interface{}) {
	buffer := bytes.NewBuffer(byteArray)
	dec := gob.NewDecoder(buffer)
	err := dec.Decode(data)
	if err != nil {
		// log.Print("decode: ", err)
	}
}

func (d *demo) errorf(format string, args ...interface{}) {
	d.failed = true
	fmt.Fprintln(d.w, fmt.Sprintf(format, args...))
	log.Errorf(d.ctx, format, args...)
}
