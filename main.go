package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"
	"unicode"
)

func main() {

	dataFileName := flag.String("d", "", "json or yaml data file")
	templateFolder := flag.String("t", "", "go template folder")
	outputDirectory := flag.String("o", ".", "output directory")
	flag.Parse()

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])
		fmt.Printf("    servgen -d data-file (Json or Yaml) -t template-file [-o output-directory] \n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *dataFileName == "" || *templateFolder == "" {
		flag.Usage()
	}

	dataFile, err := os.Open(*dataFileName)
	if err != nil {
		fmt.Println(err)
	}

	defer dataFile.Close()

	byteValue, err := ioutil.ReadAll(dataFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	byteValue, _ = ToJSON(byteValue)

	var data interface{}
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		fmt.Println(err)
		return
	}

	cleanRoot := filepath.Clean(*templateFolder)
	pfx := len(cleanRoot) + 1
	root := template.New("").Funcs(fm)

	err = filepath.Walk(cleanRoot, func(path string, info os.FileInfo, e1 error) error {
		if !info.IsDir() {
			if e1 != nil {
				return e1
			}

			b, e2 := ioutil.ReadFile(path)
			if e2 != nil {
				return e2
			}

			name := path[pfx:]
			t := root.New(name).Funcs(fm)
			_, e2 = t.Parse(string(b))
			if e2 != nil {
				return e2
			}

			outputFileName := strings.Replace(name, ".gotmpl", ".go", -1)
			outputFileName = strings.Replace(outputFileName, ".modtmpl", ".mod", -1)

			// process filename with template
			outputFileName = processFilename(outputFileName, data)

			// generate processed file
			generateFile(t, *outputDirectory, outputFileName, data)
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

}

func processFilename(outputFileName string, data interface{}) string {

	t := template.New("").Funcs(fm)
	_, err := t.Parse(string(outputFileName))
	if err != nil {
		fmt.Println(err)
		return outputFileName
	}
	buf := new(bytes.Buffer)
	// process filename with template
	err = t.Execute(buf, data)
	if err != nil {
		fmt.Println(err)
	}
	return buf.String()
}

func generateFile(template *template.Template, outputDirectory string, outputFileName string, data interface{}) {
	absOutputFileName := path.Join(outputDirectory, outputFileName)
	_ = os.MkdirAll(path.Dir(absOutputFileName), os.ModePerm)
	outputFile, err := os.Create(absOutputFileName)
	fmt.Println("Generating file : " + absOutputFileName)
	defer outputFile.Close()
	if err != nil {
		fmt.Println(err)
	}
	err = template.Execute(outputFile, data)
	if err != nil {
		fmt.Println(err)
	}
}

func ToJSON(data []byte) ([]byte, error) {
	if hasJSONPrefix(data) {
		return data, nil
	}
	return yaml.YAMLToJSON(data)
}

func hasJSONPrefix(buf []byte) bool {
	trim := bytes.TrimLeftFunc(buf, unicode.IsSpace)
	return bytes.HasPrefix(trim, []byte("{"))
}

// functions run from the template context
var fm = map[string]interface{}{
	"Title":   strings.Title,
	"ToLower": strings.ToLower,
	"ToUpper": strings.ToUpper,
	"Now":     func() string { return time.Now().Round(time.Second).Format("20060102150405") },
}
