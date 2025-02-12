package app

type SvEn struct {
	Sv string `json:"sv"`
	En string `json:"en"`
}

type SuperGroup struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	PrettyName string `json:"prettyName"`
	Type       string `json:"type"`
	Email      string `json:"email"`
}

type Group struct {
	Id              string     `json:"id"`
	BecomesActive   int        `json:"becomesActive"`
	BecomesInactive int        `json:"becomesInactive"`
	Description     SvEn       `json:"description"`
	Email           string     `json:"email"`
	Function        SvEn       `json:"function"`
	Name            string     `json:"name"`
	PrettyName      string     `json:"prettyName"`
	SuperGroup      SuperGroup `json:"superGroup"`
	AvatarURL       string     `json:"avatarURL"`
	Active          bool       `json:"active"`
}

type Relationship struct {
	SuperGroup SuperGroup `json:"superGroup"`
	Group      Group      `json:"group"`
}

type Authority struct {
	Id        string `json:"id"`
	Authority string `json:"authority"`
}

type User struct {
	Cid     string `json:"cid"`
	Nick    string `json:"nickname"`
	Picture string `json:"picture"`
}
