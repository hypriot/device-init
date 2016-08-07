// Copyright © 2016 Govinda Fichtner <govinda.fichtner@googlemail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"compress/gzip"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var loggingPath = "/var/log/device-init"
var logFile = filepath.Join(loggingPath, "preloaded_images.log")

// preload-imagesCmd represents the preload-images command
var preload_imagesCmd = &cobra.Command{
	Use:   "preload-images",
	Short: "Preload Docker images",
	Long:  `Preload Docker images that are defined in device-init.yaml.`,
	Run: func(cmd *cobra.Command, args []string) {
		dockerPreloadImages()
	},
}

func dockerPreloadImages() {
	readDockerConfig()

	for _, imageFile := range myDockerConfig.Images {

		if _, err := os.Stat(imageFile); err == nil {

			if validFileType(imageFile) {
				importImage(imageFile)
			} else {
				continue
			}

		} else {
			fmt.Println("Image file does not exist:", imageFile)
		}

	}

}

func init() {
	dockerCmd.AddCommand(preload_imagesCmd)
}

func validFileType(path string) bool {
	gz, _ := filepath.Match("*.tar.gz", filepath.Base(path))
	tar, _ := filepath.Match("*.tar", filepath.Base(path))

	if gz || tar {
		return true
	} else {
		return false
	}
}

func importImage(path string) {
	var loadOptions docker.LoadImageOptions

	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)

	if !imageAlreadyImported(path) {

		plainFile, err := os.Open(path)
		if err != nil {
			fmt.Println("Could not open file:", err)
		}
		defer plainFile.Close()

		if isCompressed(path) {
			gzipFile, err := gzip.NewReader(plainFile)
			if err != nil {
				fmt.Println("Could not open zipped file:", err)
			}
			loadOptions.InputStream = gzipFile
		} else {
			loadOptions.InputStream = plainFile
		}

		err = client.LoadImage(loadOptions)
		if err != nil {
			fmt.Println("Could not import Docker image:", err)
		} else {
			logImportedImage(path)
		}
		fmt.Println("Imported image:", path)
	} else {
		fmt.Println("Already imported image:", path)
	}
}

func isCompressed(path string) bool {
	compressed := false
	if gz, _ := filepath.Match("*.tar.gz", filepath.Base(path)); gz {
		file, err := os.Open(path)
		if err != nil {
			fmt.Println("Could not open file:", err)
		}
		defer file.Close()

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			fmt.Println("Could not read file:", err)
		}

		filetype := http.DetectContentType(buff)
		if filetype == "application/x-gzip" {
			compressed = true
		}
	}
	return compressed
}

func logImportedImage(path string) {
	err := os.MkdirAll(loggingPath, 0755)
	if err != nil {
		fmt.Println("Could not create logging dir:", err)
	}

	if !imageAlreadyImported(path) {
		var f *os.File
		if _, err := os.Stat(logFile); err == nil {
			f, err = os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, 0666)
		} else {
			f, err = os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0666)
		}
		if err != nil {
			fmt.Println("Could not open logfile:", err)
		}
		defer f.Close()

		if _, err := f.WriteString(path); err != nil {
			fmt.Println("Could not write to logfile:", err)
		}
	}
}

func imageAlreadyImported(path string) bool {
	imported := false

	if _, err := os.Stat(logFile); err == nil {
		input, err := ioutil.ReadFile(logFile)
		if err != nil {
			fmt.Println("Could not open logfile:", err)
		}

		for _, line := range strings.Split(string(input), "\n") {
			if strings.Contains(line, path) {
				imported = true
			}
		}
	}

	return imported
}
