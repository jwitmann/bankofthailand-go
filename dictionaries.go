package bankofthailand

var AuctionNameTranslation = map[string]string{
	"พันธบัตรรัฐบาล":                 "Government Bond",
	"บัตรเงินคลัง":                   "Treasury Bill",
	"ตั๋วเงินคลัง":                   "Treasury Bill",
	"พันธบัตรรัฐบาลพันธบัตรเงินคลัง": "Government Bond / Treasury Bill",
}

var BusinessTypeTranslation = map[string]string{
	"ธนาคารพาณิชย์ไทย": "Thai Commercial Bank",
	"ธนาคารออมสิน":     "Government Savings Bank",
	"ธนาคารพาณิชย์ที่เป็นบริษัทลูกของธนาคารต่างประเทศ":   "Foreign Bank Subsidiary",
	"สาขาธนาคารต่างประเทศ":                               "Foreign Bank Branch",
	"สถาบันการเงินในต่างประเทศ":                          "Foreign Financial Institution",
	"สำนักงานผู้แทนธนาคารต่างประเทศ":                     "Foreign Bank Representative Office",
	"นิติบุคคลในต่างประเทศ":                              "Foreign Legal Entity",
	"ผู้ประกอบธุรกิจการให้สินเชื่อที่มิใช่สถาบันการเงิน": "Non-Bank Lender",
	"บรรษัทตลาดรองสินเชื่อที่อยู่อาศัย":                  "Secondary Mortgage Corporation",
	"บรรษัทประกันสินเชื่ออุตสาหกรรมขนาดย่อม":             "Small Business Credit Guarantee Corporation",
}
