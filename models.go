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

type ThaiFAHoliday struct {
	HolidayWeekDay         string `json:"HolidayWeekDay"`
	HolidayWeekDayThai     string `json:"HolidayWeekDayThai"`
	Date                   string `json:"Date"`
	DateThai               string `json:"DateThai"`
	HolidayDescription     string `json:"HolidayDescription"`
	HolidayDescriptionThai string `json:"HolidayDescriptionThai"`
}

type ThaiFAResponse struct {
	Result struct {
		API       string          `json:"api"`
		Timestamp string          `json:"timestamp"`
		Data      []ThaiFAHoliday `json:"data"`
	} `json:"result"`
}
