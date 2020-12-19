package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handleSearch(searcher))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Listening on port %s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

type Searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}
		results := searcher.Search(query[0])
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}

func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)
	s.SuffixArray = suffixarray.New([]byte(strings.ToLower(string(dat))))
	return nil
}

func (s *Searcher) Search(query string) []string {
	query = strings.ToLower(query)
	idxs := s.SuffixArray.Lookup([]byte(query), -1)
	results := []string{}
	for _, idx := range idxs {
		results = append(results, s.processResult(idx, query))
	}

	if len(results) == 0 {
		ch := make(chan []Result, 2*(len(query)-3))
		fQuery := query
		for len(query) >= 4 {
			query = query[0 : len(query)-1]
			go s.partialResults(query, ch)
		}

		count := 1
		length := len(fQuery)
		tQuery := fQuery

		for len(tQuery) >= 4 {
			tQuery = fQuery[count:length]
			count++
			go s.partialResults(tQuery, ch)
		}

		added := make(map[int]bool)

		for i := 0; i < cap(ch); i++ {
			result := <-ch
			for _, res := range result {
				idx := res.idx
				if _, ok := added[idx]; !ok {
					results = append(results, res.res)
					added[idx] = true
				}
			}
		}
	}
	return results
}

type Result struct {
	res string
	idx int
}

func (s *Searcher) partialResults(query string, ch chan []Result) {
	idxs := s.SuffixArray.Lookup([]byte(query), -1)
	var results []Result
	for _, idx := range idxs {
		results = append(results, Result{s.processResult(idx, query), idx})
	}
	ch <- results
}

func (s *Searcher) processResult(from int, query string) string {
	const twoNewLines string = "\n\r\n"
	const oneNewLine string = "\n"
	const prevLine string = "\r\n"
	prevIdx := from - 250
	left := prevIdx
	prevFirstIdx := -1

	for i := from; i >= 3 && (i >= left || prevIdx > left-200); i-- { //todo check end conditions for i=0,1,2
		if s.CompleteWorks[i-3:i] == twoNewLines {
			prevIdx = i
		}

		if prevFirstIdx == -1 && s.CompleteWorks[i-2:i] == prevLine {
			prevFirstIdx = i
		}

	}

	nextIdx := from + 250
	right := nextIdx
	nextFirstIdx := -1

	for i := from; (i <= right || (nextIdx == right || nextIdx < right-50)) && i < len(s.CompleteWorks)-3; i++ { //todo check end conditions for i=len-1,len-2,len-3
		if s.CompleteWorks[i+1:i+4] == twoNewLines {
			nextIdx = i
		}
		if nextFirstIdx == -1 && s.CompleteWorks[i+1:i+2] == oneNewLine {
			nextFirstIdx = i
		}
	}

	//return s.CompleteWorks[prevIdx:from] + "<b>" + s.CompleteWorks[from:from+len(query)] + "</b>" + s.CompleteWorks[from+len(query):nextIdx] //todo check this how it works without including nextindex
	return "<span class='hide'>" + s.CompleteWorks[prevIdx:prevFirstIdx] + "</span><b>" + s.CompleteWorks[prevFirstIdx:from] + "<mark>" + s.CompleteWorks[from:from+len(query)] + "</mark>" + s.CompleteWorks[from+len(query):nextFirstIdx] + "</b><span class='hide'>" + s.CompleteWorks[nextFirstIdx:nextIdx] + "</span>" //todo check this how it works without including nextindex
}


