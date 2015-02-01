package evals

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func sanespace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func Search(c *http.Client, termId, departmentId, subjectCode, instructorName string) (map[string]string, error) {
	resp, err := c.Get(fmt.Sprintf("https://edu-apps.mit.edu/ose-rpt/subjectEvaluationSearch.htm?termId=%s&departmentId=%s&subjectCode=%s&instructorName=%s&search=Search", termId, departmentId, subjectCode, instructorName))
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		hrefP, err := url.Parse(href)
		if err != nil {
			return
		}
		q, err := url.ParseQuery(hrefP.RawQuery)
		if err != nil {
			return
		}

		subjectIds, ok := q["subjectId"]
		if !ok {
			return
		}
		if len(subjectIds) < 1 {
			return
		}
		description := sanespace(s.Closest("p").Text())
		m[description] = "https://edu-apps.mit.edu/ose-rpt/" + href
	})

	return m, nil
}

func Report(c *http.Client, reportLink string) (map[string]string, error) {
	resp, err := c.Get(reportLink)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	doc.Find(".avg").Each(func(i int, s *goquery.Selection) {
		description := sanespace(s.Siblings().First().Text())
		// TODO: recognize fields from other sections than "SUBJECT"
		if description != "" {
			m[description] = s.Text()
		}
	})

	return m, nil
}
