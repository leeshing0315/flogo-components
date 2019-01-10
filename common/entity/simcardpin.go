package entity

type SimcardPin struct {
	pinno     string
	simno     string
	regtype   string
	regtime   string
	tableName string // default: 'TblPin2SimNo'
}
