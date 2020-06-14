package aws4_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/go-imsto/aws4"
)

func ExampleS3() {

	region := os.Getenv("AWS_S3_REGION")
	if region == "" {
		region = "ap-northeast-1"
	}
	text := "hello world"
	uri := "https://s3." + region + ".amazonaws.com/" + os.Getenv("AWS_S3_BUCKET") + "/test.txt"
	r, _ := http.NewRequest("PUT", uri, bytes.NewBufferString(text))
	sha256Sum := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
	r.Header.Set("x-amz-content-sha256", sha256Sum)
	r.Header.Set("content-type", "text/plain")
	r.Header.Set("content-length", fmt.Sprint(len(text)))
	resp, err := aws4.DefaultClient.Do(r)
	if err != nil {
		log.Fatal(err)
	}

	var buf []byte
	buf, _ = ioutil.ReadAll(resp.Body)
	log.Printf("put resp(%d) %v %s", resp.StatusCode, resp.Header, buf)
	resp.Body.Close()

	r, _ = http.NewRequest("GET", uri, nil)
	r.Header.Set("x-amz-content-sha256", aws4.EmptySum)
	resp, err = aws4.DefaultClient.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	buf, _ = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		log.Printf("get resp %s", buf)
	} else {
		log.Printf("get resp %s", buf[0:8])
	}
	resp.Body.Close()

	fmt.Println(resp.StatusCode)
	// Output:
	// 200
}

func Example_jSONBody() {
	data := strings.NewReader("{}")
	r, _ := http.NewRequest("POST", "https://dynamodb.us-east-1.amazonaws.com/", data)
	r.Header.Set("Content-Type", "application/x-amz-json-1.0")
	r.Header.Set("X-Amz-Target", "DynamoDB_20111205.ListTables")

	resp, err := aws4.DefaultClient.Do(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.StatusCode)
	// Output:
	// 200
}

func Example_formEncodedBody() {
	v := make(url.Values)
	v.Set("Action", "DescribeAutoScalingGroups")

	resp, err := aws4.PostForm("https://autoscaling.us-east-1.amazonaws.com/", v)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.StatusCode)
	// Output:
	// 200
}

func ExampleSignGlacier() {
	r, _ := http.NewRequest("GET", "https://glacier.us-east-1.amazonaws.com/-/vaults", nil)
	r.Header.Set("X-Amz-Glacier-Version", "2012-06-01")

	resp, err := aws4.DefaultClient.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	var buf []byte
	buf, _ = ioutil.ReadAll(resp.Body)
	log.Printf("resp %s", buf)

	fmt.Println(resp.StatusCode)
	// Output:
	// 200
}
