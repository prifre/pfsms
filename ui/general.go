package ui

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)
func NewBoldLabel(text string) *widget.Label {
	return &widget.Label{Text: text, TextStyle: fyne.TextStyle{Bold: true}}
}
func Appendtotextfile(fn string, m string) error {
	var err error
	var path string
	path, err = os.UserHomeDir()
	if err != nil {
		panic("path")
	}
	path = fmt.Sprintf("%s%c%s",path,os.PathSeparator,"pfsms")
	if _, err = os.Stat(path); err != nil {
		log.Println("#1 Adding folder data: " + path)
		if os.IsNotExist(err) {
			err = os.Mkdir(path, 0755)
			if err != nil {
				panic(err.Error())
			}
			// file does not exist
		} else {
			panic(err.Error())
			// other error
		}
	}
	fn = fmt.Sprintf("%s%c%s",path ,os.PathSeparator, fn)
	// m=strings.Replace(m,"\r","<CR>",-1)
	// m=strings.Replace(m,"\n","",-1)
	f, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(m); err != nil {
		log.Println(err)
	}
	return err
}
func Readtextfile(fn string) (string,error) {
	var err error
	var path string
	path, err = os.UserHomeDir()
	if err != nil {
		panic("path")
	}
	path = fmt.Sprintf("%s%c%s",path,os.PathSeparator,"pfsms")
	if _, err = os.Stat(path); err != nil {
		log.Println("#1 Adding folder data: " + path)
		if os.IsNotExist(err) {
			err = os.Mkdir(path, 0755)
			if err != nil {
				panic(err.Error())
			}
			// file does not exist
		} else {
			panic(err.Error())
			// other error
		}
	}
	fn = fmt.Sprintf("%s%c%s",path ,os.PathSeparator, fn)
	var b0 []byte
	b0, err = os.ReadFile(fn) // SQL to make tables!
	if err != nil {
		fmt.Print(err)
	}
	return string(b0),err
}
func GetAllCountries() []string {
// Countries taken from github.com/IftekherSunny/go_country and from
// https://en.wikipedia.org/wiki/List_of_country_calling_codes
	var  thecountries = []string{
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
		"Zimbabwe} (+263)"}
	return thecountries
}
func Fixphonenumber(pn string,cc string) string {
	// pn phonenumber  cc coutrycode
	// Sweden (+46) converts to 0046
	var cci string ="00"
	for i:=0;i<len(cc);i++ {
		if strings.Index("0123456789",string(cc[i]))>0 {
			cci +=string(cc[i])
		}
	}
	if pn[0:2]==string("00") {
		return pn
	}
	if string(pn[0])=="0" {
		return  cci+pn[1:]
	}
	if string(pn[0])=="+" {
		return "00"+pn[1:]
	}
	return cci+pn
}
func Setuplog() {
	var wrt io.Writer
	var path string
	var err error
	path, err = os.UserHomeDir()
	if err!=nil {
		panic("path")
	}
	path = fmt.Sprintf("%s%c%s",path ,os.PathSeparator,"pfsms")
	if _, err = os.Stat(path); err != nil {
		log.Println("#1 Adding folder data: " + path)
		if os.IsNotExist(err) {
			err = os.Mkdir(path, 0755)
			if err!=nil {
				panic(err.Error())
			}
			// file does not exist
		} else {
			panic(err.Error())
			// other error
		}
	}
	f, err := os.OpenFile( fmt.Sprintf("%s%c%s",path ,os.PathSeparator,"smslog.txt"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	//	defer f.Close()
	wrt = io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)
}
func Getcustomersfilename() string {
	var path string
	var err error
	path, err = os.UserHomeDir()
	if err!=nil {
		panic("path")
	}
	path = fmt.Sprintf("%s%c%s%c%s",path ,os.PathSeparator,"pfsms",os.PathSeparator,"customers.txt")
	return path
}
func Getgroupsfilename() string {
	var path string
	var err error
	path, err = os.UserHomeDir()
	if err!=nil {
		panic("path")
	}
	path = fmt.Sprintf("%s%c%s%c%s",path ,os.PathSeparator,"pfsms",os.PathSeparator,"groups.txt")
	return path
}
func ReadLastLineWithSeek(fn string,cnt int) (string, error) {
	var err error
	var path string
	path, err = os.UserHomeDir()
	if err != nil {
		panic("path")
	}
	path = fmt.Sprintf("%s%c%s",path,os.PathSeparator,"pfsms")
	if _, err = os.Stat(path); err != nil {
		log.Println("#1 Adding folder data: " + path)
		if os.IsNotExist(err) {
			err = os.Mkdir(path, 0755)
			if err != nil {
				panic(err.Error())
			}
			// file does not exist
		} else {
			panic(err.Error())
			// other error
		}
	}
	fn = fmt.Sprintf("%s%c%s",path ,os.PathSeparator, fn)
    fileHandle, err := os.Open(fn)

    if err != nil {
        panic("Cannot open file")
    }
    defer fileHandle.Close()

    line := ""
    var cursor int64 = 0
    stat, _ := fileHandle.Stat()
    cursor = stat.Size()
	linecount:=0
    for { 
        cursor --
        fileHandle.Seek(cursor, io.SeekStart)

        char := make([]byte, 1)
        fileHandle.Read(char)

        if  (char[0] == 10 || char[0] == 13) { // stop if we find a line
			if len(line)>0 {
				if !(line[0]==10 || line[0]==13) {
					linecount ++
					if linecount>=cnt {
						break
					}
				}
			} else {
				linecount++
			}
        }
        line = fmt.Sprintf("%s%s", string(char), line) // there is more efficient way
        if cursor <= 0 { // stop if we are at the begining
            break
        }
    }

    return line,err
}
