package googleapiclient

import (
	"context"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Spreadsheets struct {
	spreadsheet *sheets.Service
}

// NewSheetService create a new service. It needs a path to the service account token and the required spreadsheets scopes
func NewSheetService(keypath, projectId string, scopes ...string) (*Spreadsheets, error) {

	c, err := getCredentials(keypath, projectId, scopes...)
	if err != nil {
		return nil, err
	}

	s, err := sheets.NewService(context.Background(), option.WithHTTPClient(c))
	if err != nil {
		return nil, err
	}

	return &Spreadsheets{spreadsheet: s}, nil
}

func (s *Spreadsheets) GetSpreadsheetTitle(spreadsheetId string) (string, error) {

	ss, err := s.spreadsheet.Spreadsheets.Get(spreadsheetId).Do()
	if err != nil {
		return "", err
	}

	return ss.Properties.Title, nil
}

func (s *Spreadsheets) GetSpreadsheetSheetsList(spreadsheetId string) ([]string, error) {

	ss, err := s.spreadsheet.Spreadsheets.Get(spreadsheetId).Do()
	if err != nil {
		return nil, err
	}

	var sheetList []string
	for _, sh := range ss.Sheets {
		sheetList = append(sheetList, sh.Properties.Title)
	}

	return sheetList, nil
}

// GetSpreadsheetValuesFormatted gets a range of values from a specific sheet. Use this format "Sheet!A1:H9"
func (s *Spreadsheets) GetSpreadsheetValuesFormatted(spreadsheetId string, valuesRange ...string) ([][]interface{}, error) {
	return s.getSpreadsheetValues(spreadsheetId, "FORMATTED_VALUE", valuesRange...)
}

// GetSpreadsheetValuesUnFormatted gets a range of values from a specific sheet. Use this format "Sheet!A1:H9"
func (s *Spreadsheets) GetSpreadsheetValuesUnFormatted(spreadsheetId string, valuesRange ...string) ([][]interface{}, error) {
	return s.getSpreadsheetValues(spreadsheetId, "UNFORMATTED_VALUE", valuesRange...)
}
func (s *Spreadsheets) getSpreadsheetValues(spreadsheetId, valueRenderOptions string, valuesRange ...string) ([][]interface{}, error) {

	resp, err := s.spreadsheet.Spreadsheets.Values.BatchGet(spreadsheetId).Ranges(valuesRange...).ValueRenderOption(valueRenderOptions).Do()
	if err != nil {
		return nil, err
	}

	var respValues [][]interface{}
	for _, v := range resp.ValueRanges {
		for _, vv := range v.Values {
			respValues = append(respValues, vv)
		}
	}

	return respValues, nil
}

func (s *Spreadsheets) GetSpreadsheetValuesByColumn(spreadsheetId string, cells ...string) ([][]interface{}, error) {

	datafilter := &sheets.DataFilter{
		A1Range: "Title",
	}
	filters := &sheets.BatchGetValuesByDataFilterRequest{
		DataFilters: []*sheets.DataFilter{datafilter},
	}
	resp, err := s.spreadsheet.Spreadsheets.Values.BatchGetByDataFilter(spreadsheetId, filters).Do()
	if err != nil {
		return nil, err
	}

	var respValues [][]interface{}
	for _, v := range resp.ValueRanges {
		for _, vv := range v.ValueRange.Values {
			respValues = append(respValues, vv)
		}
	}

	return respValues, nil
}
