package controller

import (
	"cassandra"
	"encoding/json"
	"fmt"
	"generateshortcode"
	"net/http"

	"redis"

	"github.com/gorilla/mux"
)

//LENGTHSHORTCODE độ dài shortcode
const LENGTHSHORTCODE = 6

//DOMAIN của server
const DOMAIN = "http://k8sfresh.misa.com.vn/url/"

//Redirect người dùng đến link gốc
func Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortcode := vars["shortcode"]
	defer r.Body.Close()
	longurl, err := redis.Get(shortcode)
	if err != nil {
		fmt.Println("Redirect: Redis error", err)
	} else {
		if longurl != "" {
			fmt.Println("Redirect: Cache", longurl)
			http.Redirect(w, r, longurl, 301)
			return
		}
	}

	longurl, err = cassandra.Get(shortcode)
	if err != nil {
		fmt.Println("Redirect: Cassandra error", err)
		w.WriteHeader(503)
		fmt.Fprintln(w, "Server error")
	} else {
		if longurl != "" {
			fmt.Println("Redirect: Cassandra", longurl)
			http.Redirect(w, r, longurl, 301)
			if err = redis.Set(shortcode, longurl); err != nil {
				fmt.Println("Redirect: Redis error", err)
			}
		} else {
			w.WriteHeader(404)
			fmt.Fprintln(w, "Page not found")
			fmt.Println("Page not found")
		}
	}
	return
}

//Dữ liệu respone trả về
type res struct {
	Code int    `json:"code"`
	Data string `json:"data"` // Link rút gọn
}

//Dữ liệu gửi lên server
type data struct {
	LongURL    string `json:"longurl"`
	CustomCode string `json:"customcode"`
}

//Shorten Rút gọn link
func Shorten(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	encoder := json.NewEncoder(w)
	var body data

	w.Header().Set("Content-Type", "application/json")
	if err := decoder.Decode(&body); err != nil {
		encoder.Encode(res{Code: http.StatusBadRequest})
		return
	}
	defer r.Body.Close()

	var custom bool
	if body.CustomCode != "" {
		custom = true
	} else {
		custom = false
	}

	longurl := body.LongURL
	var shortcode string
	if custom {
		shortcode = body.CustomCode
		url, err := cassandra.Get(shortcode)
		if err != nil {
			encoder.Encode(res{Code: http.StatusServiceUnavailable})
			return
		}
		if url != "" {
			encoder.Encode(res{Code: http.StatusConflict})
			return
		}

	} else {
		shortcode = generateshortcode.GenerateShortCode(longurl, LENGTHSHORTCODE)
		url, err := cassandra.Get(shortcode)
		for {
			if err != nil {
				encoder.Encode(res{Code: http.StatusServiceUnavailable})
				return
			}

			if url == "" {
				break
			} else {
				if url == body.LongURL {
					encoder.Encode(res{Code: http.StatusOK, Data: shortcode})
					return
				}
				longurl = "m" + longurl
				shortcode = generateshortcode.GenerateShortCode(longurl, LENGTHSHORTCODE)
				url, err = cassandra.Get(shortcode)
			}
		}
	}

	if cassandra.Insert(shortcode, body.LongURL) != nil {
		encoder.Encode(res{Code: http.StatusServiceUnavailable})
		return
	}
	encoder.Encode(res{Code: http.StatusCreated, Data: DOMAIN + shortcode})
	return
}
