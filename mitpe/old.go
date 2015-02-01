package mitpe

import (
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

/*
func OldGetMePE(c *http.Client, coursePattern, sectionPattern, mitid string) error {
	err := mit.TouchstoneLogin(c, "https://sisapp.mit.edu/mitpe/student")
	if err != nil {
		return err
	}
	log.Printf("Logged in. ")
poll:
	for errors := 0; errors < 10; {
		// TODO: rewrite this to use the sections page, requiring one fewer http request
		log.Printf("Listing courses...")
		var courses map[string]string
		courses, err = ListCourses(c)
		if err != nil {
			errors++
			continue poll
		}
		courseName := choose(coursePattern, courses)
		if courseName == "" {
			continue poll
		}
		course := courses[courseName]

		log.Printf("Listing sections...")
		sections, err := ListSections(c, course)
		if err != nil {
			continue poll
		}
		for s, v := range sections {
			log.Printf("%s: %s", s, v)
		}
		sectionName := choose("MW 1:00", sections)
		if sectionName == "" {
			continue poll
		}
		section := sections[sectionName]
		log.Printf("Registering for %s %s", courseName, sectionName)
		if err := Register(c, section, mitid); err != nil {
			continue poll
		}
		log.Printf("Everything seems to have gone fine.")
		return nil
	}
	return err
}
*/

func ListCourses(c *http.Client) (map[string]string, error) {
	resp, err := c.Get("https://edu-apps.mit.edu/mitpe/student/registration/quick?wf=%2fregistration%2fhome")
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string, 50)
	doc.Find("#course").Find("optgroup").Find("option[value]").Each(func(i int, s *goquery.Selection) {
		m[s.Text()], _ = s.Attr("value")
	})
	return m, nil
}

func ListSections(c *http.Client, section string) (map[string]string, error) {
	resp, err := c.Get("https://edu-apps.mit.edu/mitpe/student/registration/quick?wf=%2Fregistration%2Fquick&course=" + section)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string, 50)
	doc.Find(".sectionContainer").Find("button[onclick]").Each(func(i int, s *goquery.Selection) {
		onclick, _ := s.Attr("onclick")
		sectionID := regexp.MustCompile("[0-9a-fA-F]{32}").FindString(onclick)
		m[sanespace(s.Closest("td.data").Next().Text())] = sectionID
	})
	return m, nil
}
