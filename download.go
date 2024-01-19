package main

import (
	"archive/zip"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/fmartingr/go-cbz"
)

func isExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func createDir(fpath string) error {
	return os.MkdirAll(fpath, 0o755)
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

// thanks fam https://gosamples.dev/remove-non-alphanumeric/
func clearString(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}

func titleCleanup(cname string) string {
	cname = strings.ToLower(cname)
	cname = strings.ReplaceAll(cname, " ", "-")
	cname = strings.TrimSpace(cname)
	return cname
}

func addFileToZip(zipWriter *zip.Writer, filePath, zpath string) error {
	// Open the image file
	imgFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer imgFile.Close()

	// Create a new zip file entry
	zipFile, err := zipWriter.Create(zpath)
	if err != nil {
		return err
	}

	// Copy the image file content to the zip entry
	_, err = io.Copy(zipFile, imgFile)
	if err != nil {
		return err
	}

	return nil
}

func makeCbz(file, name string, imgs []string) error {
	file = file + ".cbz"

	out, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer out.Close()

	zr := zip.NewWriter(out)
	defer zr.Close()

	for _, img := range imgs {
		ip := filepath.Base(img)
		err := addFileToZip(zr, img, ip)
		if err != nil {
			return err
		}
	}
	return nil
}

// filename - series name - array of relative image paths
func makeCBZ(file string, name string, imgs []string) error {
	comic, err := cbz.New()
	if err != nil {
		return err
	}

	comic.ComicInfo().Series = name

	for _, img := range imgs {
		err := comic.AppendPage(img)
		if err != nil {
			return err
		}
	}

	return comic.Save(file)
}

var badReq = errors.New("Http request status not OK")

func getImg(url, fpath string) error {
	out, err := os.OpenFile(fpath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer out.Close()

	cl := getClient(time.Second * 10)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", opts.UserAgent)

	resp, err := cl.Do(req)
	if resp.StatusCode != http.StatusOK {
		return badReq
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func Download(urls []string, dir string) ([]string, error) {
	var fails []string
	var imgs []string

	for _, url := range urls {
		fpath := filepath.Base(url)
		fpath = filepath.Join(dir, fpath)

		if isExists(fpath) {
			continue
		}

		if err := getImg(url, fpath); err != nil {
			fails = append(fails, url)
		} else {
			imgs = append(imgs, fpath)
		}
	}

	// gotta be a better way lol but idk yet, we wait a sec and retry
	if len(fails) > 1 {
		time.Sleep(time.Second * 5)
		for _, furl := range fails {
			fpath := filepath.Base(furl)
			fpath = filepath.Join(dir, fpath)

			if err := getImg(furl, fpath); err != nil {
				return nil, err
			} else {
				imgs = append(imgs, fpath)
			}
		}
	}

	return imgs, nil
}
