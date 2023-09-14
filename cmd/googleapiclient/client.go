package googleapiclient

type Client struct {
	Calendar     *Calendar
	Spreadsheets *Spreadsheets
	Events       *GEvent
}
