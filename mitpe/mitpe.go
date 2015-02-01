package mitpe

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/andres-erbsen/mit"
)

type PESection struct {
	Name  string // complete course numnber, including "PE" prefix and section suffix
	Title string // free-form title of class
	Days  string // first letter of weekday names
	Time  string // time of the section, in hopefully in "Kitchen" format
}

func GetMePE(c *http.Client, mitid, course, day, timeSpec string) error {
	err := mit.TouchstoneLogin(c, "https://sisapp.mit.edu/mitpe/student")
	if err != nil {
		return err
	}
	log.Printf("Logged in.")
poll:
	for errors := 0; errors < 10; {
		log.Printf("Listing courses and sections...")
		t0 := time.Now()
		var sections map[PESection]string
		sections, err = ListCoursesAndSections(c)
		if err != nil {
			errors++
			continue poll
		}
		for s, v := range sections {
			log.Printf("%v: %s", s, v)
		}
		section, err := choose(sections, course, day, timeSpec)
		if err != nil {
			continue poll
		}
		log.Printf("Registering for %v", section)
		sectionID := sections[section]
		if err := Register(c, sectionID, mitid); err != nil {
			errors++
			continue poll
		}
		log.Printf("Everything seems to have gone fine.")
		log.Printf("Critical section time: %s\n", time.Since(t0))
		return nil
	}
	return err
}

func ListCoursesAndSections(c *http.Client) (map[PESection]string, error) {
	resp, err := c.Get("https://edu-apps.mit.edu/mitpe/student/registration/sectionList")
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	m := make(map[PESection]string, 100)
	doc.Find(".sectionContainer").Find("a[href]").Each(func(i int, s *goquery.Selection) {
		sectionLink, _ := s.Attr("href")
		if !strings.Contains(sectionLink, "sectionId") {
			return
		}
		sectionID := regexp.MustCompile("[0-9a-fA-F]{32}").FindString(sectionLink)
		title := sanespace(s.Closest("td").Next().Next().Find("p").First().Text())
		day := sanespace(s.Closest("td").Next().Next().Next().Text())
		time := sanespace(s.Closest("td").Next().Next().Next().Next().Text())
		m[PESection{s.Text(), title, day, time}] = sectionID
	})
	return m, nil
}

func Register(c *http.Client, section string, mitid string) error {
	resp, err := c.PostForm("https://edu-apps.mit.edu/mitpe/student/registration/create",
		url.Values{"sectionId": {section}, "mitId": {mitid}, "wf": {"/registration/quick"}})
	if err != nil {
		return err
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}
	if !strings.Contains(doc.Text(), "Click to Cancel") {
		return fmt.Errorf("Something went wrong:\n%s\n%s", fmt.Sprint(resp.Header), doc.Text())
	}
	return nil
}

func choose(sections map[PESection]string, name, days, time string) (PESection, error) {
	for s := range sections {
		if strings.Contains(s.Name, name) && strings.Contains(s.Days, days) && strings.Contains(s.Time, time) {
			return s, nil
		}
	}
	return PESection{}, fmt.Errorf("no suitable section found")
}

func sanespace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
