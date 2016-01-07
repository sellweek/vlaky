package vlak

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	nameRegexp  = regexp.MustCompile(`(.+) (\d+) ?(.*)?`)
	routeRegexp = regexp.MustCompile(`\((\d\d\.\d\d\. \d\d:\d\d)\) (.*) -> (.*) \((\d\d\.\d\d\. \d\d:\d\d)\)`)
	delayRegexp = regexp.MustCompile(`\d+`)
)

type TrainInfo struct {
	Category, Name string
	Number         int
	Current        Delay
	From, To       Location
}

type Location struct {
	Station string
	Time    time.Time
}

type Delay struct {
	Location
	Actually time.Time
	Delay    int
}

func Parse() (locations []TrainInfo, err error) {
	doc, err := goquery.NewDocument("http://tis.zsr.sk/elis/pohybvlaku?jazyk_stranky=sk")
	if err != nil {
		return err
	}
	delayHeaders := doc.Find(".accordionHeader")
	delayTables := doc.Find(".trainDelayTable")

	locations = make([]TrainInfo, len(delayHeaders.Nodes), len(delayHeaders.Nodes))
	delayHeaders.Each(func(i int, element *goquery.Selection) {
		parseHeader(element, &locations[i])
		parseTable(delayTables.Get(i), &locations[i])
	})
}

func parseHeader(element *goquery.Selection, info *TrainInfo) {
	element.Find("span").Each(func(i int, element *goquery.Selection) {
		switch i {
		case 0:
			info.Category, info.Number, info.Name = parseTrainDenomination(element.Text())
		case 2:
			info.From, info.To = parseTrainRoute(element.Text())
		}
	})
}

func parseTable(element *html.Node, info *TrainInfo) {
	rows := flattenTable(element)
	for i, row := range rows {
		cells := row.Find("td")
		switch i {
		case 0:
			if len(cells.Nodes) == 1 {
				info.Current.Delay = 0
			} else {
				match := delayRegexp.FindString(cells.Text())
				info.Current.Delay, _ = strconv.Atoi(match)
			}
		case 1:
			info.Current.Station = cells.Get(1).FirstChild.Data
		case 2:
			info.Current.Actually = parseTime(cells.Get(1).FirstChild.Data)
		case 3:
			info.Current.Time = parseTime(cells.Get(1).FirstChild.Data)
		}
	}
}

func flattenTable(table *html.Node) (selections []*goquery.Selection) {
	doc := goquery.NewDocumentFromNode(table)
	selections = make([]*goquery.Selection, 0)
	doc.Find("tr").Each(func(i int, row *goquery.Selection) {
		row.RemoveFiltered("tr")
		if row.Text() != "" {
			selections = append(selections, row)
		}
	})
	return
}

func parseTrainDenomination(d string) (category string, number int, name string) {
	result := nameRegexp.FindStringSubmatch(d)
	category = result[1]
	number, _ = strconv.Atoi(result[2])
	if len(result) > 3 {
		name = strings.Trim(result[3], " ")
	}
	return
}

func parseTrainRoute(r string) (from, to Location) {
	result := routeRegexp.FindStringSubmatch(r)
	from.Time = parseTime(result[1])
	from.Station = result[2]
	to.Time = parseTime(result[4])
	to.Station = result[3]
	return
}

func parseTime(str string) time.Time {
	t, err := time.ParseInLocation("02.01. 15:04", str, time.UTC)
	if err != nil {
		fmt.Println(err)
		return time.Time{}
	}
	return t.AddDate(time.Now().Year(), 0, 0)
}
