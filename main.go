package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Jeffail/gabs/v2"
	apisv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

func trimEmpty(key string, jsonParsed *gabs.Container, parent *gabs.Container) {
	if jsonParsed.Data() == nil {
		parent.Delete(key)
	} else if v, ok := jsonParsed.Data().(string); ok && v == "" {
		parent.Delete(key)
	} else if v, ok := jsonParsed.Data().(float64); ok && v == 0 {
		parent.Delete(key)
	} else if v, ok := jsonParsed.Data().(int64); ok && v == 0 {
		parent.Delete(key)
	} else if v, ok := jsonParsed.Data().(map[string]interface{}); ok {
		for k, child := range jsonParsed.ChildrenMap() {
			trimEmpty(k, child, jsonParsed)
		}

		if len(v) == 0 {
			parent.Delete(key)
		}
	} else if v, ok := jsonParsed.Data().([]interface{}); ok {
		for k, child := range jsonParsed.ChildrenMap() {
			trimEmpty(k, child, jsonParsed)
		}

		if len(v) == 0 {
			parent.Delete(key)
		}
	}
}

func defaultsForResource(resource []byte) *gabs.Container {
	var myscheme *runtime.Scheme
	var decoder runtime.Decoder

	myscheme = runtime.NewScheme()
	apisv1.AddToScheme(myscheme)
	decoder = scheme.Codecs.UniversalDeserializer()

	defaultB := make([]byte, len(resource))
	copy(defaultB, resource)

	defaultJson, err := gabs.ParseJSON(defaultB)
	if err != nil {
		panic(err)
	}

	defaultFlat, _ := defaultJson.Flatten()
	for key, _ := range defaultFlat {
		if strings.HasPrefix(key, "spec.") {
			defaultJson.SetP(nil, key)
		}

	}

	rawDefault, _ := json.Marshal(defaultJson.Data())
	candidate, _, _ := decoder.Decode(rawDefault, nil, nil)
	myscheme.Default(candidate)
	default_, _ := json.Marshal(candidate)
	defaultJson2, _ := gabs.ParseJSON(default_)

	return defaultJson2
}

func main() {

	input, err := io.ReadAll(os.Stdin)

	if err != nil {
		panic(err)
	}

	inputJson, err := gabs.ParseJSON(input)
	if err != nil {
		panic(err)
	}
	defaultJson := defaultsForResource(input)

	flat, _ := inputJson.Flatten()
	for key, _ := range flat {
		left := inputJson.Path(key).Data()
		right := defaultJson.Path(key).Data()

		if left == right {
			inputJson.DeleteP(key)
		}
	}

	inputJson.DeleteP("status")

	trimEmpty("", inputJson, nil)

	fmt.Println(inputJson.String())

	prettyJSON, _ := json.MarshalIndent(inputJson.Data(), "", "  ")

	fmt.Print(string(prettyJSON))

}
