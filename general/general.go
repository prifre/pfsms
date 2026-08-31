package general

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// Förkompilerade regex på paketnivå för bättre prestanda
var (
	reDigits = regexp.MustCompile(`\d+`)
	rePhone  = regexp.MustCompile(`[^\d+]`)
	logfile  *os.File
)

func Getmemoryinfo() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	var r string
	r += fmt.Sprintf("Memory Usage = %v MB", (m.Alloc / 1024 / 1024))
	r += fmt.Sprintf("\r\nApplication Memory = %v MB", (m.Sys / 1024 / 1024))
	return r
}

func Setupfiles() {
	var err error
	home := GetHomeDir()

	logPath := filepath.Join(home, "pfsms.log")
	prefs := fyne.CurrentApp().Preferences()
	prefs.SetString("customersfile", filepath.Join(home, "customers.txt"))
	prefs.SetString("groupsfile", filepath.Join(home, "groups.txt"))
	prefs.SetString("historyfile", filepath.Join(home, "history.txt"))
	prefs.SetString("pfsmsdb", filepath.Join(home, "pfsms.db"))
	prefs.SetString("pfsmslog", logPath)

	// O_WRONLY används för att undvika fillåsningsproblem på Windows vid loggning
	logfile, err = os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Varning: Kunde inte öppna loggfil %s: %v\n", logPath, err)
		return
	}

	wrt := io.MultiWriter(os.Stdout, logfile)
	log.SetOutput(wrt)
}

func GetHomeDir() string {
	var path string
	var err error

	if fyne.CurrentApp() != nil && fyne.CurrentApp().Preferences().Bool("debug") {
		exePath, err := os.Executable()
		if err == nil {
			path = filepath.Dir(exePath)
		} else {
			path, _ = os.Getwd()
		}
	} else {
		path, err = os.UserHomeDir()
		if err != nil {
			log.Printf("#2 GetHomeDir misslyckades med UserHomeDir: %v", err)
			path, _ = os.Getwd()
		}
	}

	dataDir := filepath.Join(path, "pfsmsdata")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("GetHomeDir kunde inte skapa datakatalog %s: %v", dataDir, err)
	}

	return dataDir
}

func CloseLog() {
	if logfile != nil {
		logfile.Close()
		log.SetOutput(os.Stdout)
	}
}

func NewBoldLabel(text string) *widget.Label {
	return &widget.Label{Text: text, TextStyle: fyne.TextStyle{Bold: true}}
}

func Appendtotextfile(fn string, m string) error {
	f, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("kunde inte öppna fil %s: %w", fn, err)
	}
	defer f.Close()

	_, err = f.WriteString(m)
	return err
}

func Readtextfile(fn string) (string, error) {
	b0, err := os.ReadFile(fn)
	if err != nil {
		return "", err
	}
	return string(b0), nil
}

func CheckGUI() bool {
	return os.Getenv("APP_MODE") == "GUI"
}

func Fixphonenumber(pn string) string {
	var cc string
	if !CheckGUI() {
		cc = "Sweden(+46)"
	} else {
		cc = fyne.CurrentApp().Preferences().StringWithFallback("mobilecountry", "Sweden(+46)")
	}

	ccDigits := strings.Join(reDigits.FindAllString(cc, -1), "")
	cleanPn := rePhone.ReplaceAllString(pn, "")

	if len(cleanPn) < 5 {
		return ""
	}

	if strings.HasPrefix(cleanPn, "00") {
		return cleanPn
	}

	if strings.HasPrefix(cleanPn, "+") {
		return "00" + cleanPn[1:]
	}

	if strings.HasPrefix(cleanPn, "0") {
		return "00" + ccDigits + cleanPn[1:]
	}

	return "00" + ccDigits + cleanPn
}

func Showdebugmsg(s string) string {
	r2 := s
	r2 = strings.ReplaceAll(r2, "\r", "\\r")
	r2 = strings.ReplaceAll(r2, "\n", "\\n")
	r2 = strings.ReplaceAll(r2, "\x00", "\\0")
	r2 = strings.ReplaceAll(r2, "\t", "\\t")
	r2 = strings.ReplaceAll(r2, "\x1a", "\\z")
	return r2
}

func ReadLastLineWithSeek(fn string, maxLines int) (string, error) {
	file, err := os.Open(fn)
	if err != nil {
		return "", err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "", err
	}

	fileSize := stat.Size()
	var cursor int64 = fileSize
	var result []byte
	lineCount := 0

	const bufferSize = 4096
	buf := make([]byte, bufferSize)

	for cursor > 0 && lineCount < maxLines {
		readSize := int64(bufferSize)
		if cursor < readSize {
			readSize = cursor
		}

		cursor -= readSize
		_, err = file.Seek(cursor, io.SeekStart)
		if err != nil {
			return "", err
		}

		_, err = io.ReadFull(file, buf[:readSize])
		if err != nil {
			return "", err
		}

		for i := int(readSize) - 1; i >= 0; i-- {
			if buf[i] == '\n' {
				lineCount++
				if lineCount > maxLines {
					result = append(buf[i+1:readSize], result...)
					break
				}
			}
		}

		if lineCount <= maxLines {
			result = append(buf[:readSize], result...)
		}
	}

	return strings.TrimSpace(string(result)), nil
}

func GetAllCountries() []string {
	return []string{
		"Afghanistan (+93)",
		"Albania (+355)",
		"Algeria (+213)",
		"American Samoa (+1684)",
		"Andorra (+376)",
		"Angola (+244)",
		"Anguilla (+1264)",
		"Antarctica (+672)",
		"Antigua and Barbuda (+1268)",
		"Argentina (+54)",
		"Armenia (+374)",
		"Aruba (+297)",
		"Australia (+61)",
		"Austria (+43)",
		"Azerbaijan (+994)",
		"Bahrain (+973)",
		"The Bahamas (+1242)",
		"Bangladesh (+880)",
		"Barbados (+1 246)",
		"Belarus (+375)",
		"Belgium (+32)",
		"Belize (+501)",
		"Benin (+229)",
		"Bermuda (+1441)",
		"Bhutan (+975)",
		"Bolivia (+591)",
		"Bonaire (+599)",
		"Bosnia and Herzegovina (+387)",
		"Botswana (+267)",
		"Bouvet (+47)",
		"Brazil (+55)",
		"British Indian Ocean Territory (+246)",
		"British Virgin Islands (+1284)",
		"Brunei (+673)",
		"Bulgaria (+359)",
		"Burkina Faso (+226)",
		"Myanmar (+95)",
		"Burundi (+257)",
		"Cambodia (+855)",
		"Cameroon (+237)",
		"Canada (+1)",
		"Cape Verde (+238)",
		"Cayman Islands (+1345)",
		"Central African Republic (+236)",
		"Chad (+235)",
		"Chile (+56)",
		"China (+86)",
		"Christmas Island (+61)",
		"Cocos-Keeling Islands (+672)",
		"Colombia (+57)",
		"Comoros (+269)",
		"Congo (+242)",
		"Congo, Dem. Rep. of (Zaire) (+243)",
		"Cook Islands (+682)",
		"Costa Rica (+506)",
		"Cote d'Ivoire (+225)",
		"Croatia (+385)",
		"Curacao (+599)",
		"Cuba (+53)",
		"Cyprus (+357)",
		"Czech Republic (+420)",
		"Denmark (+45)",
		"Djibouti (+253)",
		"Dominica (+1767)",
		"Dominican Republic (+1809)",
		"East Timor (+670)",
		"Ecuador (+593)",
		"Egypt (+20)",
		"El Salvador (+503)",
		"Equatorial Guinea (+240)",
		"Eritrea (+291)",
		"Estonia (+372)",
		"Ethiopia (+251)",
		"Falkland Islands (+500)",
		"Fiji (+679)",
		"Finland (+358)",
		"France (+33)",
		"French Guiana (+594)",
		"French Polynesia (+689)",
		"French Southern and Antarctic Lands (+262)",
		"Gabon (+241)",
		"The Gambia (+220)",
		"Georgia (+995)",
		"Germany (+49)",
		"Ghana (+233)",
		"Greece (+30)",
		"Greenland (+299)",
		"Grenada (+1473)",
		"Guadeloupe (+590)",
		"Guam (+1671)",
		"Guatemala (+502)",
		"Guernsey (+44)",
		"Guinea (+224)",
		"Guinea-Bissau (+245)",
		"Guyana (+592)",
		"Haiti (+509)",
		"Heard Island and McDonald Islands (+0)",
		"Holy See (Vatican City) (+39)",
		"Honduras (+504)",
		"Hong Kong SAR China (+852)",
		"Hungary (+36)",
		"Iceland (+354)",
		"India (+91)",
		"Indonesia (+62)",
		"Iran (+98)",
		"Iraq (+964)",
		"Ireland (+353)",
		"Isle of Man (+44)",
		"Israel (+972)",
		"Italy (+39)",
		"Jamaica (+1876)",
		"Japan (+81)",
		"Jordan (+962)",
		"Kazakhstan (+7)",
		"Kenya (+254)",
		"Kiribati (+686)",
		"Kuwait (+965)",
		"Kyrgyzstan (+996)",
		"Laos (+856)",
		"Latvia (+371)",
		"Lebanon (+961)",
		"Lesotho (+266)",
		"Liberia (+231)",
		"Libya (+218)",
		"Liechtenstein (+423)",
		"Lithuania (+370)",
		"Luxembourg (+352)",
		"Macau SAR China (+853)",
		"Macedonia (+389)",
		"Madagascar (+261)",
		"Malawi (+265)",
		"Malaysia (+60)",
		"Maldives (+960)",
		"Mali (+223)",
		"Malta (+356)",
		"Marshall Islands (+692)",
		"Martinique (+596)",
		"Mauritania (+222)",
		"Mauritius (+230)",
		"Mayotte (+262)",
		"Mexico (+52)",
		"Micronesia, Federated States Of (+691)",
		"Midway Island (+1808)",
		"Moldova (+373)",
		"Monaco (+377)",
		"Mongolia (+976)",
		"Montenegro (+382)",
		"Montserrat (+1664)",
		"Morocco (+212)",
		"Mozambique (+258)",
		"Namibia (+264)",
		"Nauru (+674)",
		"Nepal (+977)",
		"Netherlands (+31)",
		"Netherlands Antilles (+599)",
		"New Caledonia (+687)",
		"New Zealand (+64)",
		"Nicaragua (+505)",
		"Niger (+227)",
		"Nigeria (+234)",
		"Niue (+683)",
		"Norfolk Island (+672)",
		"North Korea (+850)",
		"Northern Mariana Islands (+1670)",
		"Norway (+47)",
		"Oman (+968)",
		"Pakistan (+92)",
		"Palau (+680)",
		"Panama (+507)",
		"Papua New Guinea (+675)",
		"Paraguay (+595)",
		"Peru (+51)",
		"Philippines (+63)",
		"Pitcairn Islands (+870)",
		"Poland (+48)",
		"Portugal (+351)",
		"Puerto Rico (+1787)",
		"Qatar (+974)",
		"Reunion (+262)",
		"Romania (+40)",
		"Russia (+7)",
		"Rwanda (+250)",
		"Saint Barthelemy (+590)",
		"Saint Helena (+290)",
		"Saint Kitts and Nevis (+1869)",
		"Saint Lucia (+1758)",
		"Saint Martin (+1)",
		"Saint Pierre and Miquelon (+508)",
		"Saint tome and principle (+239)",
		"Saint Vincent and the Grenadines (+1784)",
		"Samoa (+684)",
		"San Marino (+378)",
		"Saudi Arabia (+966)",
		"Senegal (+221)",
		"Serbia (+381)",
		"Seychelles (+248)",
		"Sierra Leone (+232)",
		"Singapore (+65)",
		"Sint Maarten (+721)",
		"Slovakia (+421)",
		"Slovenia (+386)",
		"Solomon Islands (+677)",
		"South Africa (+27)",
		"South Georgia and the South Sandwich Islands (+500)",
		"South Korea (+82)",
		"South Sudan (+211)",
		"Spain (+34)",
		"Sri Lanka (+94)",
		"Sudan (+249)",
		"Suriname (+597)",
		"Svalbard (+47)",
		"Swaziland (+268)",
		"Sweden (+46)",
		"Switzerland (+41)",
		"Syria (+963)",
		"Taiwan (+886)",
		"Tajikistan (+992)",
		"Tanzania (+255)",
		"Thailand (+66)",
		"Togo (+228)",
		"Tokelau (+690)",
		"Tonga (+676)",
		"Trinidad and Tobago (+1868)",
		"Tunisia (+216)",
		"Turkey (+90)",
		"Turkmenistan (+7370)",
		"Turks and Caicos Islands (+1649)",
		"Tuvalu (+688)",
		"Uganda (+256)",
		"Ukraine (+380)",
		"United Arab Emirates (+971)",
		"United Kingdom (+44)",
		"United States Minor Outlying Islands (+1)",
		"United States (+1)",
		"Uruguay (+598)",
		"Uzbekistan (+998)",
		"Vanuatu (+678)",
		"Venezuela (+58)",
		"Vietnam (+84)",
		"Virgin Islands (+1340)",
		"Wallis and Futuna (+681)",
		"Western Sahara (+212)",
		"Yemen (+967)",
		"Zambia (+260)",
		"Zimbabwe (+263)",
	}
}
