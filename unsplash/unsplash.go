package unsplash

import (
	"encoding/json"
	"fmt"
	_ "github.com/dustin/go-humanize"
	"io"
	"net/http"
	"os"
	"strings"
)

const BaseUrl = "https://api.unsplash.com"
const APIKEY = "e16e0a00f80aa7f1d491201d5db32bfdfd801d9be57b05a2b959436432e55d71"

// Typings
type PhotoURLS struct {
	Raw string `json:"raw"`
	Full string `json:"full"`
	Regular string `json:"regular"`
	Small string `json:"small"`
	Thumb string `json:"thumb"`
}

type Photo struct {
	ID string `json:"id"`
	Width int `json:"width"`
	Height int `json:"height"`
	Description string `json:"description"`
	Color string `json:"color"`
	Urls PhotoURLS `json:"urls"`
}

type Collection struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	TotalPhotos int `json:"total_photos"`
	Links struct {
		Self string `json:"self"`
		Html string `json:"html"`
		Photos string `json:"photos"`
	} `json:"links"`
}


func GetCollectionPhotos(collection Collection, page int, perPage int) ([]Photo, error) {
	var photos []Photo

	if page < 1 {
		page = 1
	}

	if perPage < 1 {
		perPage = 1
	}

	url := fmt.Sprintf("%s?page=%d&per_page=%d&client_id=%s", collection.Links.Photos, page, perPage, APIKEY)

	r, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&photos)

	if err != nil {
		return nil, err
	}

	return photos, nil
}

func GetCollection(collectionID string) (*Collection, error) {
	var data *Collection
	url := fmt.Sprintf("%s/collections/%s?client_id=%s", BaseUrl, collectionID, APIKEY)

	res, err := http.Get(url)

	if err != nil {
		res.Body.Close()
		return nil, err
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&data)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func DownloadPhoto (photo Photo, subFolder string) error {
	filename := fmt.Sprintf("%s.jpg", photo.ID)
	return download(photo.Urls.Raw, fmt.Sprintf("photos/%s/%s", subFolder, filename))
}


func download(url string, filename string) error {
	var r *http.Response

	pathComponents := strings.Split(filename, "/")
	filepath := strings.Join(pathComponents[:len(pathComponents) - 1], "/")

	e := os.MkdirAll(filepath, os.ModePerm)

	if e != nil {
		panic(e.Error())
	}

	out, err := os.Create(filename + ".tmp")

	if err != nil {
		_ = out.Close()
		return err
	}

	r, err = http.Get(url)

	if err != nil {
		_ = r.Body.Close()
		return err
	}

	defer r.Body.Close()

	if _, err = io.Copy(out, r.Body); err != nil {
		_ = out.Close()
		return err
	}

	_ = out.Close()

	if err = os.Rename(filename + ".tmp", filename); err != nil {
		return err
	}

	return nil
}