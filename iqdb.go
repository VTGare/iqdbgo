package iqdbgo

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/sirupsen/logrus"
)

var (
	baseURL = "https://iqdb.org/?url="
)

type Result struct {
	BestMatch       *Match
	PossibleMatches []*Match
}

type Match struct {
	Thumbnail  string
	URL        string
	Tags       string
	Similarity int
}

func Search(file string) (*Result, error) {
	var (
		col    = colly.NewCollector(colly.Async(true))
		err    error
		result = &Result{nil, make([]*Match, 0)}
	)

	col.OnRequest(func(r *colly.Request) {
		logrus.Infof("Making a request. URL: %s", r.URL)
	})

	col.OnError(func(r *colly.Response, e error) {
		logrus.Errorf("Request errored: %v", e)
		err = e
	})

	col.OnHTML("tbody", func(h *colly.HTMLElement) {
		var res *Match
		if t := h.ChildText("tr > th"); t == "Best match" {
			res = &Match{}
			result.BestMatch = res
		} else if t == "Additional match" || t == "Possible match" {
			res = &Match{}
			result.PossibleMatches = append(result.PossibleMatches, res)
		}

		if res != nil {
			res.URL = h.ChildAttr("tr > .image > a", "href")
			res.Thumbnail = "https://iqdb.org" + h.ChildAttr("tr > .image > a > img", "src")
			res.Tags = h.ChildAttr("tr > .image > a > img", "alt")

			for _, s := range h.ChildTexts("tr > td") {
				if strings.Contains(s, "similarity") {
					simStr := s[:strings.Index(s, "%")]
					sim, err := strconv.Atoi(simStr)
					if err != nil {
						logrus.Warnf("iqdb -> OnHTML(): %v", err)
					}

					res.Similarity = sim
				}
			}
		}
	})

	col.Visit(baseURL + url.QueryEscape(file))
	col.Wait()

	if err != nil {
		return nil, err
	}

	return result, nil
}
