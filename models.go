package bankofthailand

type Holiday struct {
	Date                   string `json:"Date"`
	DateThai               string `json:"DateThai"`
	HolidayWeekDay         string `json:"HolidayWeekDay"`
	HolidayWeekDayThai     string `json:"HolidayWeekDayThai"`
	HolidayDescription     string `json:"HolidayDescription"`
	HolidayDescriptionThai string `json:"HolidayDescriptionThai"`
	HolidayType            string `json:"HolidayType,omitempty"`
}

type HolidaysResult struct {
	API       string    `json:"api"`
	Timestamp string    `json:"timestamp"`
	Data      []Holiday `json:"data"`
}

type HolidaysResponse struct {
	Result HolidaysResult `json:"result"`
}
