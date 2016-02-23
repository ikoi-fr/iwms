package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	bbox := r.Form.Get("BBOX")
	width := r.Form.Get("WIDTH")
	height := r.Form.Get("HEIGHT")

	if bbox == "" || width == "" || height == "" {
		log.Println("wrong call", r.URL)
		w.WriteHeader(500)
		return
	}

	path := "http://wms.francethd.fr/geoserver/Observatoire/wms?LAYERS=Observatoire%3Ainfra_d_c_f_mercator&" +
		"TRANSPARENT=true&TILED=true&SERVICE=WMS&VERSION=1.1.1&REQUEST=GetMap&STYLES=&FORMAT=image%2Fpng" +
		"&SRS=EPSG%3A900913&BBOX=" + bbox + "&WIDTH=" + width + "&HEIGHT=" + height

	log.Println("path", path)

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(500)
		return
	}

	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Accept", "image/webp,image/*,*/*;q=0.8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.109 Safari/537.36")
	req.Header.Add("Referer", "http://observatoire.francethd.fr/")
	req.Header.Add("Accept-Encoding", "gzip, deflate, sdch")
	req.Header.Add("Accept-Language", "en-US,en;q=0.8,fr;q=0.6")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	for k, _ := range resp.Header {
		w.Header().Add(k, resp.Header.Get(k))
	}

	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Println(err)
	}

	resp.Body.Close()

}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+port, nil)
}
