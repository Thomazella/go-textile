package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"syscall/js"

	// "time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/util"
)

func getBlocks(inputs []js.Value) {
	opts := map[string]string{
		"thread": inputs[0].String(),
		"offset": "",
		"limit":  strconv.Itoa(20),
	}
	callback := inputs[len(inputs)-1:][0]
	go func() {
		res, err := executeJSONCmd(GET, "blocks", params{opts: opts}, nil)
		if err != nil {
			callback.Invoke(err, js.Null())
			return
		}
		callback.Invoke(js.Null(), res)
		c <- true
		return
	}()
}

func getFiles(inputs []js.Value) {
	opts := map[string]string{
		"thread": inputs[0].String(),
		"offset": inputs[1].String(),
		"limit":  strconv.Itoa(20),
	}
	var list pb.FilesList
	callback := inputs[len(inputs)-1:][0]
	go func() {
		res, err := executeJSONPbCmd(GET, "files", params{opts: opts}, &list)
		if err != nil {
			callback.Invoke(err, js.Null())
			return
		}
		callback.Invoke(js.Null(), res)
		c <- true
		return
	}()
}

// ClientOptions control how to access the running Textile daemon's API
type ClientOptions struct {
	APIAddr    string `long:"api" description:"API address to use" default:"http://127.0.0.1:40600"`
	APIVersion string `long:"api-version" description:"API version to use" default:"v0"`
}

var apiAddr = "http://127.0.0.1:40600"
var apiVersion = "v0"

type method string

const (
	// GET represents the GET HTTP method
	GET method = "GET"
	// POST represents the POST HTTP method
	POST method = "POST"
	// PUT represents the PUT HTTP method
	PUT method = "PUT"
	// DEL represents the DELETE HTTP method
	DEL method = "DELETE"
	// PATCH represents the PATCH HTTP method
	PATCH method = "PATCH"
)

type params struct {
	args []string
	opts map[string]string
	// payload io.Reader
	ctype string
}

func executeStringCmd(meth method, pth string, pars params) (string, error) {
	res, _, err := request(meth, pth, pars)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := util.UnmarshalString(res.Body)
	if err != nil {
		return "", err
	}
	return body, nil
}

func executeJSONCmd(meth method, pth string, pars params, target interface{}) (string, error) {
	res, _, err := request(meth, pth, pars)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, err := util.UnmarshalString(res.Body)
		if err != nil {
			return "", err
		}
		return "", errors.New(body)
	}

	if target == nil {
		target = new(interface{})
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(data, target); err != nil {
		return "", err
	}
	jsn, err := json.MarshalIndent(target, "", "    ")
	if err != nil {
		return "", err
	}
	return string(jsn), nil
}

func executeJSONPbCmd(meth method, pth string, pars params, target proto.Message) (map[string]interface{}, error) {
	res, _, err := request(meth, pth, pars)
	if err != nil {
		return map[string]interface{}{}, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, err := util.UnmarshalString(res.Body)
		if err != nil {
			return map[string]interface{}{}, err
		}
		return map[string]interface{}{}, errors.New(body)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return map[string]interface{}{}, err
	}
	if err := pbUnmarshaler.Unmarshal(bytes.NewReader(data), target); err != nil {
		return map[string]interface{}{}, err
	}
	jsn, err := pbMarshaler.MarshalToString(target)
	if err != nil {
		return map[string]interface{}{}, err
	}
	var out map[string]interface{}
	err = json.Unmarshal([]byte(jsn), &out)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return out, nil
}

func request(meth method, pth string, pars params) (*http.Response, func(), error) {
	apiURL := fmt.Sprintf("%s/api/%s/%s", apiAddr, apiVersion, pth)
	req, err := http.NewRequest(string(meth), apiURL, nil) //, pars.payload)
	if err != nil {
		return nil, nil, err
	}

	if len(pars.args) > 0 {
		var args []string
		for _, arg := range pars.args {
			args = append(args, url.PathEscape(arg))
		}
		req.Header.Set("X-Textile-Args", strings.Join(args, ","))
	}

	if len(pars.opts) > 0 {
		var items []string
		for k, v := range pars.opts {
			items = append(items, k+"="+url.PathEscape(v))
		}
		req.Header.Set("X-Textile-Opts", strings.Join(items, ","))
	}

	if pars.ctype != "" {
		req.Header.Set("Content-Type", pars.ctype)
	}

	tr := &http.Transport{}
	res, err := tr.RoundTrip(req)
	cancel := func() {
		tr.CancelRequest(req)
	}
	return res, cancel, err
}

var pbMarshaler = jsonpb.Marshaler{
	EnumsAsInts: false,
	Indent:      "    ",
}
var pbUnmarshaler = jsonpb.Unmarshaler{
	AllowUnknownFields: true,
}

var c chan bool

func init() {
	c = make(chan bool)
}

func main() {
	js.Global().Set("getFiles", js.NewCallback(getFiles))
	fmt.Println("Hello from WASM/Go!")
	// fmt.Println("account/address")
	// res, err := executeStringCmd(GET, "account/address", params{})
	// if err != nil {
	// 	return
	// }
	// output(res)
	<-c
}
