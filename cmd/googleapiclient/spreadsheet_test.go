package googleapiclient

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var (
	keyPath            = os.Getenv("KEYPATH")
	scopesSpreadsheets = os.Getenv("SPREADSHEET_SCOPE")
	projectId          = os.Getenv("PROJECT_ID")
	ssId               = os.Getenv("SPREADSHEET_TEST_ID")
)

func TestNewSheetService(t *testing.T) {

	t.Run("test client connection", func(t *testing.T) {
		expected := ssId
		s, err := NewSheetService(keyPath, projectId, scopesSpreadsheets)
		if err != nil {
			t.Fatal(err)
		}

		ss, err := s.spreadsheet.Spreadsheets.Get(ssId).Do()
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expected, ss.SpreadsheetId)
	})

	t.Run("test client connection with fake credentials", func(t *testing.T) {
		fakeCredPath := "credentials_example.json"
		_, err := NewSheetService(fakeCredPath, projectId, scopesSpreadsheets)
		assert.Error(t, err)
	})

}

func TestSpreadsheetAPI_GetSpreadsheetTitle(t *testing.T) {
	c, err := NewSheetService(keyPath, projectId, scopesSpreadsheets)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("get spreadsheet title", func(t *testing.T) {
		expected := "test-spreadsheet"
		actual, err := c.GetSpreadsheetTitle(ssId)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, expected, actual)
	})

	t.Run("get spreadsheet title with error", func(t *testing.T) {
		fakeSsId := "fakeSpreadsheetId"
		_, err := c.GetSpreadsheetTitle(fakeSsId)
		assert.Error(t, err)
	})
}

func TestSpreadsheetAPI_GetSpreadsheetSheetsList(t *testing.T) {
	c, err := NewSheetService(keyPath, projectId, scopesSpreadsheets)
	if err != nil {
		t.Fatal(err)
	}
	t.Run("get sheet list", func(t *testing.T) {
		sheetList, err := c.GetSpreadsheetSheetsList(ssId)
		if err != nil {
			t.Fatal(err)
		}

		assert.GreaterOrEqual(t, 2, len(sheetList))
	})
}

func TestSpreadsheetAPI_GetSpreadsheetValues(t *testing.T) {
	c, err := NewSheetService(keyPath, projectId, scopesSpreadsheets)
	if err != nil {
		t.Fatal(err)
	}

	type expected struct {
		cellsRange   []string
		expectedRows int
	}

	expectedList := []expected{
		{
			cellsRange:   []string{"'Sheet 2'!A1:D2", "A5:D5"},
			expectedRows: 3,
		},
		{
			cellsRange:   []string{"'Sheet 2'!A1:D2"},
			expectedRows: 2,
		},
	}

	for _, e := range expectedList {

		t.Run("get spreadsheet values", func(t *testing.T) {
			values, err := c.GetSpreadsheetValuesFormatted(ssId, e.cellsRange...)
			if err != nil {
				t.Fatal(err)
			}
			//for _, v := range values {
			//	log.Println(v)
			//}
			assert.Equal(t, e.expectedRows, len(values))
		})
	}

	t.Run("get spreadsheet values", func(t *testing.T) {
		cellsRange := []string{"'Sheet 2'!A1:D2", "A5:D5"}

		values, err := c.GetSpreadsheetValuesFormatted(ssId, cellsRange...)
		if err != nil {
			t.Fatal(err)
		}
		//for _, v := range values {
		//	log.Println(v)
		//}
		assert.Equal(t, 3, len(values))
	})

	t.Run("get spreadsheet values", func(t *testing.T) {
		cellsRange := []string{"'Sheet 2'!A1:D2", "A5:D5"}

		values, err := c.GetSpreadsheetValuesFormatted(ssId, cellsRange...)
		if err != nil {
			t.Fatal(err)
		}
		//for _, v := range values {
		//	log.Println(v)
		//}
		assert.Equal(t, 3, len(values))
	})
}

func TestSpreadsheetAPI_GetSpreadsheetValuesByColumn(t *testing.T) {
	c, err := NewSheetService(keyPath, projectId, scopesSpreadsheets)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("get spreadsheet values by column", func(t *testing.T) {

		values, err := c.GetSpreadsheetValuesByColumn(ssId)
		if err != nil {
			t.Fatal(err)
		}
		for _, v := range values {
			log.Println(v)
		}
		//assert.Equal(t, 3, len(values))
	})
}
