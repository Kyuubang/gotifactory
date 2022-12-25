package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"github.com/Jeffail/gabs/v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	version string
	pathbin string
	channel string
	commit  string
}

type Manifesto struct {
	Package string `json:"package"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
	URL     string `json:"url"`
	Sha256  string `json:"sha256"`
	Channel string `json:"channel"`
}

var config = new(Config)
var server string

func init() {
	flag.StringVar(&config.version, "version", "dev", "specify version binary")
	flag.StringVar(&config.pathbin, "pathbin", "", "binary location")
	flag.StringVar(&config.channel, "channel", "latest", "type of build [nightly, latest, stable]")
	flag.StringVar(&config.commit, "commit", "", "commit hash")
	flag.StringVar(&server, "server", "http://localhost", "specify gotifactory server url")
}

func main() {
	flag.Parse()

	var pkgIndex int
	var pkgExist bool
	var out string

	_, err := os.Stat("repo")
	if os.IsNotExist(err) {
		err = os.Mkdir("repo", 0775)
		if err != nil {
			log.Fatal(err)
		}
	}

	_, packageName := filepath.Split(config.pathbin)

	manifest := Manifesto{
		Package: packageName,
		Version: config.version,
		Commit:  config.commit,
		Sha256:  getHash(config.pathbin),
		Channel: config.channel,
		URL:     fmt.Sprintf(server+"repo/%s/%s", packageName, packageName),
	}

	_, err = os.Stat("repo/manifest.json")
	if os.IsNotExist(err) {
		newManifest := gabs.New()
		newManifest.Array("gotifactory")
		newManifest.ArrayAppend(manifest, "gotifactory")
		out = newManifest.StringIndent(" ", "  ")
	} else if err == nil {
		jsonParsed, err := gabs.ParseJSONFile("repo/manifest.json")
		if err != nil {
			panic(err)
		}
		//forLoop:
		for index, child := range jsonParsed.S("gotifactory").Children() {

			if child.Search("package").Data() == packageName &&
				child.Search("channel").Data() == config.channel {
				pkgIndex, pkgExist = index, true
				break
			} else {
				continue
			}
		}

		if pkgExist {
			log.Println("package found")
			err = jsonParsed.ArrayRemove(pkgIndex, "gotifactory")
			if err != nil {
				log.Fatal(err)
			}
			jsonParsed.ArrayAppend(manifest, "gotifactory")
		} else {
			log.Println("package not found")
			jsonParsed.ArrayAppend(manifest, "gotifactory")
		}
		out = jsonParsed.StringIndent(" ", "  ")
	} else {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("repo/manifest.json", []byte(out), 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("upload to repo")
	upRepo(config, packageName)
	log.Println("repo updated with version", config.version)
}

func getHash(filepath string) string {
	file, err := os.Open(filepath) // Open the file for reading
	if err != nil {
		panic(err)
	}
	defer file.Close() // Be sure to close your file!

	hash := sha256.New() // Use the Hash in crypto/sha256

	if _, err := io.Copy(hash, file); err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)) // Get encoded hash sum
}

func upRepo(c *Config, pkgName string) {
	var destDirPath = "repo/" + pkgName

	_, err := os.Stat(c.pathbin)
	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(destDirPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(destDirPath, 0775)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	source, err := os.Open(c.pathbin)
	if err != nil {
		log.Fatal(err)
	}
	defer source.Close()

	destination, err := os.Create(destDirPath + "/" + pkgName)
	if err != nil {
		log.Fatal(err)
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		log.Fatal(err)
	}
}
