package owparser

//nolint:lll // dats a magic
const (
	userName = ".header-masthead"
	platform = "div.masthead-buttons.button-group.js-button-group > a.button.m-white-outline.m-sm.is-active"

	srPath = ".masthead-player-progression--mobile > .competitive-rank > .competitive-rank-role > .competitive-rank-section:last-child"

	endorsmentLvl           = "div.masthead-player > div > div.EndorsementIcon-tooltip > div.u-center"
	endorsmentShotcaller    = "div.masthead-player > div > div.EndorsementIcon-tooltip > div.endorsement-level > div > div > svg.EndorsementIcon-border.EndorsementIcon-border--shotcaller"
	endorsmentTeammate      = "div.masthead-player > div > div.EndorsementIcon-tooltip > div.endorsement-level > div > div > svg.EndorsementIcon-border.EndorsementIcon-border--teammate"
	endorsmentSportsmanship = "div.masthead-player > div > div.EndorsementIcon-tooltip > div.endorsement-level > div > div > svg.EndorsementIcon-border.EndorsementIcon-border--sportsmanship"

	baseComp = "#competitive"
	baseQP   = "#quickplay"

	statPath = "section > div > div[data-category-id=\"%s\"] > div > div > table > tbody > tr"
	namePath = "section:nth-child(2) > div > div > div > select > option[value=\"%s\"]"

	SelectorHeroes = "section:nth-child(2) > div > div.flex-container\\@md-min.m-bottom-items > div.flex-item\\@md-min.m-grow.u-align-right > div > select"
)
