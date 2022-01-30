package owparser

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kr/pretty"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

func NewPlayerByLink(playerURL *url.URL) *Player {
	return &Player{
		url:    *playerURL,
		Heroes: map[string][]Stat{},
	}
}

func (p *Player) Gather(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.url.String(), http.NoBody)
	if err != nil {
		return errors.Wrap(err, "cant create request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "cant do request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return errors.Wrap(err, "cant read profile")
	}

	p.Name = doc.Find(userName).Text()
	p.Platform = doc.Find(platform).Text()

	var rawString string

	doc.Find(srPath).Each(func(i int, selection *goquery.Selection) {
		rawString, e := selection.Find(".competitive-rank-tier.competitive-rank-tier-tooltip").Attr("data-ow-tooltip-text")
		if e {
			switch rawString {
			case "Tank Skill Rating":
				p.Rank.Tank = stringToInt(selection.Text())
			case "Damage Skill Rating":
				p.Rank.DD = stringToInt(selection.Text())
			case "Support Skill Rating":
				p.Rank.Heal = stringToInt(selection.Text())
			}
		}
	})

	rawString = doc.Find(endorsmentLvl).Text()
	if rawString != "" {
		i, err := strconv.Atoi(rawString)
		if err != nil {
			log.Println(err)
		}

		p.Endorsment.Level = i
	}

	rawString, exists := doc.Find(endorsmentShotcaller).Attr("data-value")
	if exists {
		rawEndorsment, err := strconv.ParseFloat(rawString, 64)
		if err == nil {
			p.Endorsment.Shotcaller = rawEndorsment
		}
	}

	rawString, exists = doc.Find(endorsmentTeammate).Attr("data-value")
	if exists {
		rawEndorsment, err := strconv.ParseFloat(rawString, 64)
		if err == nil {
			p.Endorsment.Teammate = rawEndorsment
		}
	}

	rawString, exists = doc.Find(endorsmentSportsmanship).Attr("data-value")
	if exists {
		rawEndorsment, err := strconv.ParseFloat(rawString, 64)
		if err == nil {
			p.Endorsment.Sportsmanship = rawEndorsment
		}
	}

	p.parseStats(doc)

	return nil
}

func (p *Player) parseStats(s *goquery.Document) {
	defer pretty.Println(p)

	var (
		sel   *goquery.Selection
		str   string
		value float64
	)

	switcher := []bool{true, false}
	heroes := make(map[string]string)

	s.Find("section:nth-child(2)").Find("option").Each(func(i int, s *goquery.Selection) {
		code, e := s.Attr("value")
		if e {
			heroes[code] = s.Text()
		}
	})

	for _, isComp := range switcher {
		if isComp {
			sel = s.Find(baseComp)
		} else {
			sel = s.Find(baseQP)
		}

		for code := range heroes {
			// fmt.Println(name)

			s.Find(baseQP).Find(fmt.Sprintf("div[data-category-id=\"%s\"]", code)).Find("table.DataTable").Each(func(i int, s *goquery.Selection) {
				// fmt.Printf("\t%s\n", s.Find(".stat-title").Text())
				s.Find("tbody").Find("tr").Each(func(i int, s *goquery.Selection) {
					var exStat Stat
					stat, e := s.Attr("data-stat-id")
					if e {
						exStat.ID = stat
						exStat.Name = s.Find("td:nth-child(1)").Text()
						exStat.Value.QP = 6.6

						p.Heroes[heroes[code]] = append(p.Heroes[heroes[code]], exStat)
					}
					// fmt.Printf("\t\t%s\t%s\n", s.Find("td:nth-child(1)").Text(), s.Find("td:nth-child(2)").Text())
				})
			})
		}

		for _, heroCode := range heroes {

			heroName, e := sel.Find(fmt.Sprintf(namePath, heroCode)).Attr("option-id")

			if !e {
				continue
			}

			sel.Find(fmt.Sprintf(statPath, heroCode)).Each(func(i int, s *goquery.Selection) {
				var stat Stat

				stat.Name = s.Find("td:nth-child(1)").Text()
				str = s.Find("td:nth-child(2)").Text()

				switch {
				case strings.Contains(str, "%"):
					value = stringToFloat64(strings.Trim(str, "%"))
				case strings.Contains(str, ":"):
					value = timeToSec(str)
				default:
					value = stringToFloat64(str)
				}

				if isComp {
					stat.Value.Competitive = value
				} else {
					stat.Value.QP = value
				}

				p.Heroes[heroName] = append(p.Heroes[heroName], stat)
			})
		}
	}
}

//nolint:gomnd,varnamelen // MAGIC
func timeToSec(s string) (time float64) {
	switch len(s) {
	case 8:
		time = float64(
			(((int(s[0])-48)*10+int(s[1])-48)*60+((int(s[3])-48)*10+(int(s[4])-48)))*60 + (int(s[6])-48)*10 + int(s[7]) - 48,
		)
	case 5:
		time = float64(
			((int(s[0])-48)*10+(int(s[1])-48))*60 + (int(s[3])-48)*10 + int(s[4]) - 48,
		)
	case 2:
		time = float64(
			(int(s[0])-48)*10 + int(s[1]) - 48,
		)
	default:
		time = 0
	}

	return
}

func stringToFloat64(s string) (u float64) {
	// no reason to check this err
	u, _ = strconv.ParseFloat(s, 64)

	return
}

func stringToInt(s string) (i int) {
	// no reason to check this err
	i, _ = strconv.Atoi(s)

	return
}
