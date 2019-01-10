package entity

type ContainerSummary struct {
	simno      string
	carno      string
	carid      string
	commmode   string
	unitcode   string
	cartype    string
	saveflag   string
	calcflag   string
	changeflag string
	changetime string
	regtime    string
	devtype    string
	useacc     string
	groupname  string
	checkflag  string
	boxtype    string
	boxsize    string
	tableName  string // default: 'Tblcarbaseinfo'
}
