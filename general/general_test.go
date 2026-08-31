package general

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"fyne.io/fyne/v2/test"
)

func TestGetmemoryinfo(t *testing.T) {
	fmt.Println(Getmemoryinfo())
}

// TestMain körs automatiskt FÖRE alla tester och initierar Fynes testmiljö
func TestMain(m *testing.M) {
	// Skapa en dummy-applikation för Fyne så att fyne.CurrentApp() inte kraschar (nil dereference)
	test.NewApp()

	// Kör alla tester
	code := m.Run()

	// Avsluta
	os.Exit(code)
}

func TestAppendtotextfile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test_emaillog.txt")

	err := Appendtotextfile(tmpFile, "\r\nSome text1\r\n")
	if err != nil {
		t.Fatalf("Appendtotextfile misslyckades: %v", err)
	}
	Appendtotextfile(tmpFile, "\r\nSome text2\r\n")
	Appendtotextfile(tmpFile, "\r\nSome text3\r\n")
}

func TestGetAllCountries(t *testing.T) {
	s := GetAllCountries()
	if len(s) == 0 {
		t.Error("GetAllCountries returnerade en tom lista")
	}
}

func TestGetLastLines(t *testing.T) {
	// Skapa en temporär testfil istället för att lita på att pfsms.log finns
	tmpDir := t.TempDir()
	fn := filepath.Join(tmpDir, "pfsms_test.log")

	var logData string
	for i := 0; i < 30; i++ {
		logData += fmt.Sprintf("LINE %d\n", i)
	}
	err := os.WriteFile(fn, []byte(logData), 0644)
	if err != nil {
		t.Fatalf("Kunde inte skapa testfil: %v", err)
	}

	starttime := time.Now()
	m, err := ReadLastLineWithSeek(fn, 25)
	if err != nil {
		t.Fatalf("ReadLastLineWithSeek misslyckades: %v", err)
	}

	r := strings.Split(m, "\n")
	if len(r) == 0 {
		t.Error("Inga rader returnerades")
	}

	t.Logf("Tid för ReadLastLineWithSeek: %v", time.Since(starttime))
}

func TestFixphonenumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0736290839", "0046736290839"},
		{"+46736290839", "0046736290839"},
		{"+33736290839", "0033736290839"},
		{"+181736290839", "00181736290839"},
		{"0046736290839", "0046736290839"},
	}

	for i, tt := range tests {
		result := Fixphonenumber(tt.input)
		if result != tt.expected {
			t.Errorf("Test #%d misslyckades för input %q: förväntade %q, fick %q",
				i+1, tt.input, tt.expected, result)
		}
	}
}

func TestCheckGUI(t *testing.T) {
	_ = CheckGUI()
}
