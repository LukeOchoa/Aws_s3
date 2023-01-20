package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	uuid "github.com/satori/go.uuid"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"encoding/json"
)

var s3session *s3.S3

const (
	BUCKET_NAME = "borkcraftbucket"
	REGION      = "us-west-1"
)

func init() {
	s3session = s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(REGION),
	})))
}

func messageJSON() []byte {
	var message = map[string]string{
		"message": "There is Nothing to show!!!",
	}
	r, err := json.Marshal(message)
	if err != nil {
		fmt.Println("A messageJSON() caused a panic.")
		panic(err)
	}
	return r
}

func makeMessageJSON(message map[int]string) []byte {
	r, err := json.Marshal(message)
	if err != nil {
		fmt.Println("A makeMessageJSON() caused a panic.")
		panic(err)
	}
	return r
}

func corsHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		switch r.Method {
		case "OPTIONS":
			fmt.Println("OPTIONS")
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		case "GET":
			fmt.Println("GET")
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			h.ServeHTTP(w, r)
		default:
			fmt.Println("Default")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			h.ServeHTTP(w, r)
		}
		//if r.Method == "OPTIONS" {
		//	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		//	w.Header().Set("Access-Control-Allow-Credentials", "true")
		//	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		//}
		// else {
		//	w.Header().Set("Access-Control-Allow-Credentials", "true")
		//	h.ServeHTTP(w, r)
		//}
	}
}

func getObject(filename string) io.ReadCloser {
	fmt.Println("Downloading: ", filename)

	resp, err := s3session.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(filename),
	})

	if err != nil {
		panic(err)
	}

	return resp.Body

}

func listObjects() (resp *s3.ListObjectsV2Output) {
	resp, err := s3session.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(BUCKET_NAME),
	})

	if err != nil {
		panic(err)
	}

	return resp
}

func listObjectsKeys() {
	resp, err := s3session.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(BUCKET_NAME),
	})

	if err != nil {
		panic(err)
	}
	for _, object := range resp.Contents {
		fmt.Println(*object.Key)
	}
}

func listObjectsKeysAsMap() map[int]string {
	resp, err := s3session.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(BUCKET_NAME),
	})

	if err != nil {
		panic(err)
	}

	names := make(map[int]string)
	for index, object := range resp.Contents {
		names[index] = *object.Key
	}

	return names
}

func allPictureNames(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Write(makeMessageJSON(listObjectsKeysAsMap()))
}

//func awsPictures(w http.ResponseWriter, r *http.Request) {
//}


// This one is something else
func SpecificPicture(picture string) []byte {

	file, err := ioutil.ReadFile(picture)
	if err != nil {
		panic(err)
	}

	return file
}

// Real, The Route specified
func specificPicture(w http.ResponseWriter, r *http.Request) {

	var nj = r.URL.Query()
	fmt.Println("filename: ", nj["name"][0])
	picture := getObject(nj["name"][0])
	file, err := ioutil.ReadAll(picture)
	if err != nil {
		panic(err)
	}
	w.Write(file)

}
func getNetherPortalImage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SOMEOMEOEMOEMEJK")


	var nj = r.URL.Query()
	fmt.Println("filename: ", nj["name"][0])
	picture := getObject(nj["name"][0])
	file, err := ioutil.ReadAll(picture)
	if err != nil {
		panic(err)
	}
	w.Write(file)

}
func pictures(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside Pictures...")

	//pics := []string{"SpiderCowboy.png", "RyGyDuggy.png"}
	pics := []string{"SpiderCowboy.png"}
	for x := range pics {
		w.Write(SpecificPicture(pics[x]))
		//w.Write([]byte("Write me twice..."))
		fmt.Println("passage: ", x)
	}

}

func picturename(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("the URL from the REQUEST: ", r.URL.Query())
	for k, v := range r.URL.Query() {
		fmt.Println(k, v)
	}

	pictureNames := make(map[int]string)
	var cnt = 0
	for _, object := range listObjects().Contents {
		fmt.Println(*object.Key)
		pictureNames[cnt] = "http://localhost:1234/specificpicture?name=" + *object.Key
		cnt = cnt + 1
		//getObject(*object.Key)
	}

	w.Write(makeMessageJSON(pictureNames))
}

func uploadPicture(w http.ResponseWriter, r *http.Request) {
	sess := session.Must(session.NewSession(&aws.Config{
				Region: aws.String(REGION),
			}))

	uploader := s3manager.NewUploader(sess)

	//body, err := ioutil.ReadAll(r.Body)
	//panik(err)

	// get a key name from the url params array of json and byte arrays
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(BUCKET_NAME),
		Key: aws.String(r.URL.Query()["name"][0]),
		Body: r.Body,
	})
	panik(err)

	fmt.Printf("\nFile uploaded to: %s\n", aws.StringValue(&result.Location))
}

func saveImage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("we got an image request...?")
	sess := session.Must(session.NewSession(&aws.Config{
				Region: aws.String(REGION),
			}))

	uploader := s3manager.NewUploader(sess)

	name := r.URL.Query()["name"][0] + uuid.NewV4().String()
	// get a key name from the url params array of json and byte arrays
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(BUCKET_NAME),
		Key: aws.String(name),
		Body: r.Body,
	})
	panik(err)

	fmt.Printf("\nFile uploaded to: %s\n", aws.StringValue(&result.Location))

	w.WriteHeader(http.StatusAccepted)

	var nameMap = map[string]string {
		"name": name,
	}
	bytes, err := json.Marshal(nameMap)
	panik(err)
	w.Write(bytes)
}

func deleteImage(writer http.ResponseWriter, request *http.Request) {
	imageName := request.URL.Query()["name"][0]

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(REGION),
	}))
	svc := s3.New(sess)

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key: aws.String(imageName),
	}
	result, err := svc.DeleteObject(input)
	if err != nil {
		writer.WriteHeader(http.StatusForbidden)
		panic(err)
	}
	writer.WriteHeader(http.StatusAccepted)
	fmt.Println("The result was successful?... -> ", result)
}


func panik(err error) {
	if err != nil {
		panic(err)
	}
}

func doNothing(w http.ResponseWriter, r *http.Request) {}
func main() {
	fmt.Println("Server running @ localhost:1234...")

	http.HandleFunc("/favicon.ico", doNothing)
	http.HandleFunc("/ping", ping)
	http.Handle("/pictures", corsHandler(http.HandlerFunc(pictures)))
	http.Handle("/picturename", corsHandler(http.HandlerFunc(picturename)))
	http.Handle("/specificpicture", corsHandler(http.HandlerFunc(specificPicture)))
	http.Handle("/allpicturenames", corsHandler(http.HandlerFunc(allPictureNames)))
	http.HandleFunc("/uploadpicture", uploadPicture)
	http.HandleFunc("/getnetherportalimage", getNetherPortalImage)
	http.HandleFunc("/saveimage", saveImage)
	http.HandleFunc("/deleteimage", deleteImage)
	http.ListenAndServe(":1234", nil)
}

func ping(writer http.ResponseWriter, request *http.Request) {

	writer.WriteHeader(http.StatusOK)
	message, err := json.Marshal(map[string]string{"ping": "ping"})
	panik(err)
	writer.Write(message)
}