package isdomain

// TLDs is a set of TLDs, according to ICANN in 2014.
var TLDs = map[string]bool{
	"AC":                       true,
	"ACADEMY":                  true,
	"ACCOUNTANTS":              true,
	"ACTIVE":                   true,
	"ACTOR":                    true,
	"AD":                       true,
	"AE":                       true,
	"AERO":                     true,
	"AF":                       true,
	"AG":                       true,
	"AGENCY":                   true,
	"AI":                       true,
	"AIRFORCE":                 true,
	"AL":                       true,
	"ALLFINANZ":                true,
	"AM":                       true,
	"AN":                       true,
	"AO":                       true,
	"AQ":                       true,
	"AR":                       true,
	"ARCHI":                    true,
	"ARMY":                     true,
	"ARPA":                     true,
	"AS":                       true,
	"ASIA":                     true,
	"ASSOCIATES":               true,
	"AT":                       true,
	"ATTORNEY":                 true,
	"AU":                       true,
	"AUCTION":                  true,
	"AUDIO":                    true,
	"AUTOS":                    true,
	"AW":                       true,
	"AX":                       true,
	"AXA":                      true,
	"AZ":                       true,
	"BA":                       true,
	"BAR":                      true,
	"BARGAINS":                 true,
	"BAYERN":                   true,
	"BB":                       true,
	"BD":                       true,
	"BE":                       true,
	"BEER":                     true,
	"BERLIN":                   true,
	"BEST":                     true,
	"BF":                       true,
	"BG":                       true,
	"BH":                       true,
	"BI":                       true,
	"BID":                      true,
	"BIKE":                     true,
	"BIO":                      true,
	"BIZ":                      true,
	"BJ":                       true,
	"BLACK":                    true,
	"BLACKFRIDAY":              true,
	"BLUE":                     true,
	"BM":                       true,
	"BMW":                      true,
	"BN":                       true,
	"BNPPARIBAS":               true,
	"BO":                       true,
	"BOO":                      true,
	"BOUTIQUE":                 true,
	"BR":                       true,
	"BRUSSELS":                 true,
	"BS":                       true,
	"BT":                       true,
	"BUDAPEST":                 true,
	"BUILD":                    true,
	"BUILDERS":                 true,
	"BUSINESS":                 true,
	"BUZZ":                     true,
	"BV":                       true,
	"BW":                       true,
	"BY":                       true,
	"BZ":                       true,
	"BZH":                      true,
	"CA":                       true,
	"CAB":                      true,
	"CAL":                      true,
	"CAMERA":                   true,
	"CAMP":                     true,
	"CANCERRESEARCH":           true,
	"CAPETOWN":                 true,
	"CAPITAL":                  true,
	"CARAVAN":                  true,
	"CARDS":                    true,
	"CARE":                     true,
	"CAREER":                   true,
	"CAREERS":                  true,
	"CASA":                     true,
	"CASH":                     true,
	"CAT":                      true,
	"CATERING":                 true,
	"CC":                       true,
	"CD":                       true,
	"CENTER":                   true,
	"CEO":                      true,
	"CERN":                     true,
	"CF":                       true,
	"CG":                       true,
	"CH":                       true,
	"CHANNEL":                  true,
	"CHEAP":                    true,
	"CHRISTMAS":                true,
	"CHROME":                   true,
	"CHURCH":                   true,
	"CI":                       true,
	"CITIC":                    true,
	"CITY":                     true,
	"CK":                       true,
	"CL":                       true,
	"CLAIMS":                   true,
	"CLEANING":                 true,
	"CLICK":                    true,
	"CLINIC":                   true,
	"CLOTHING":                 true,
	"CLUB":                     true,
	"CM":                       true,
	"CN":                       true,
	"CO":                       true,
	"CODES":                    true,
	"COFFEE":                   true,
	"COLLEGE":                  true,
	"COLOGNE":                  true,
	"COM":                      true,
	"COMMUNITY":                true,
	"COMPANY":                  true,
	"COMPUTER":                 true,
	"CONDOS":                   true,
	"CONSTRUCTION":             true,
	"CONSULTING":               true,
	"CONTRACTORS":              true,
	"COOKING":                  true,
	"COOL":                     true,
	"COOP":                     true,
	"COUNTRY":                  true,
	"CR":                       true,
	"CREDIT":                   true,
	"CREDITCARD":               true,
	"CRUISES":                  true,
	"CU":                       true,
	"CUISINELLA":               true,
	"CV":                       true,
	"CW":                       true,
	"CX":                       true,
	"CY":                       true,
	"CYMRU":                    true,
	"CZ":                       true,
	"DAD":                      true,
	"DANCE":                    true,
	"DATING":                   true,
	"DAY":                      true,
	"DE":                       true,
	"DEALS":                    true,
	"DEGREE":                   true,
	"DEMOCRAT":                 true,
	"DENTAL":                   true,
	"DENTIST":                  true,
	"DESI":                     true,
	"DIAMONDS":                 true,
	"DIET":                     true,
	"DIGITAL":                  true,
	"DIRECT":                   true,
	"DIRECTORY":                true,
	"DISCOUNT":                 true,
	"DJ":                       true,
	"DK":                       true,
	"DM":                       true,
	"DNP":                      true,
	"DO":                       true,
	"DOMAINS":                  true,
	"DURBAN":                   true,
	"DVAG":                     true,
	"DZ":                       true,
	"EAT":                      true,
	"EC":                       true,
	"EDU":                      true,
	"EDUCATION":                true,
	"EE":                       true,
	"EG":                       true,
	"EMAIL":                    true,
	"ENGINEER":                 true,
	"ENGINEERING":              true,
	"ENTERPRISES":              true,
	"EQUIPMENT":                true,
	"ER":                       true,
	"ES":                       true,
	"ESQ":                      true,
	"ESTATE":                   true,
	"ET":                       true,
	"EU":                       true,
	"EUS":                      true,
	"EVENTS":                   true,
	"EXCHANGE":                 true,
	"EXPERT":                   true,
	"EXPOSED":                  true,
	"FAIL":                     true,
	"FARM":                     true,
	"FEEDBACK":                 true,
	"FI":                       true,
	"FINANCE":                  true,
	"FINANCIAL":                true,
	"FISH":                     true,
	"FISHING":                  true,
	"FITNESS":                  true,
	"FJ":                       true,
	"FK":                       true,
	"FLIGHTS":                  true,
	"FLORIST":                  true,
	"FLY":                      true,
	"FM":                       true,
	"FO":                       true,
	"FOO":                      true,
	"FORSALE":                  true,
	"FOUNDATION":               true,
	"FR":                       true,
	"FRL":                      true,
	"FROGANS":                  true,
	"FUND":                     true,
	"FURNITURE":                true,
	"FUTBOL":                   true,
	"GA":                       true,
	"GAL":                      true,
	"GALLERY":                  true,
	"GB":                       true,
	"GBIZ":                     true,
	"GD":                       true,
	"GE":                       true,
	"GENT":                     true,
	"GF":                       true,
	"GG":                       true,
	"GH":                       true,
	"GI":                       true,
	"GIFT":                     true,
	"GIFTS":                    true,
	"GIVES":                    true,
	"GL":                       true,
	"GLASS":                    true,
	"GLE":                      true,
	"GLOBAL":                   true,
	"GLOBO":                    true,
	"GM":                       true,
	"GMAIL":                    true,
	"GMO":                      true,
	"GMX":                      true,
	"GN":                       true,
	"GOOGLE":                   true,
	"GOP":                      true,
	"GOV":                      true,
	"GP":                       true,
	"GQ":                       true,
	"GR":                       true,
	"GRAPHICS":                 true,
	"GRATIS":                   true,
	"GREEN":                    true,
	"GRIPE":                    true,
	"GS":                       true,
	"GT":                       true,
	"GU":                       true,
	"GUIDE":                    true,
	"GUITARS":                  true,
	"GURU":                     true,
	"GW":                       true,
	"GY":                       true,
	"HAMBURG":                  true,
	"HAUS":                     true,
	"HEALTHCARE":               true,
	"HELP":                     true,
	"HERE":                     true,
	"HIPHOP":                   true,
	"HIV":                      true,
	"HK":                       true,
	"HM":                       true,
	"HN":                       true,
	"HOLDINGS":                 true,
	"HOLIDAY":                  true,
	"HOMES":                    true,
	"HORSE":                    true,
	"HOST":                     true,
	"HOSTING":                  true,
	"HOUSE":                    true,
	"HOW":                      true,
	"HR":                       true,
	"HT":                       true,
	"HU":                       true,
	"IBM":                      true,
	"ID":                       true,
	"IE":                       true,
	"IL":                       true,
	"IM":                       true,
	"IMMO":                     true,
	"IMMOBILIEN":               true,
	"IN":                       true,
	"INDUSTRIES":               true,
	"INFO":                     true,
	"ING":                      true,
	"INK":                      true,
	"INSTITUTE":                true,
	"INSURE":                   true,
	"INT":                      true,
	"INTERNATIONAL":            true,
	"INVESTMENTS":              true,
	"IO":                       true,
	"IQ":                       true,
	"IR":                       true,
	"IS":                       true,
	"IT":                       true,
	"JE":                       true,
	"JETZT":                    true,
	"JM":                       true,
	"JO":                       true,
	"JOBS":                     true,
	"JOBURG":                   true,
	"JP":                       true,
	"JUEGOS":                   true,
	"KAUFEN":                   true,
	"KE":                       true,
	"KG":                       true,
	"KH":                       true,
	"KI":                       true,
	"KIM":                      true,
	"KITCHEN":                  true,
	"KIWI":                     true,
	"KM":                       true,
	"KN":                       true,
	"KOELN":                    true,
	"KP":                       true,
	"KR":                       true,
	"KRD":                      true,
	"KRED":                     true,
	"KW":                       true,
	"KY":                       true,
	"KZ":                       true,
	"LA":                       true,
	"LACAIXA":                  true,
	"LAND":                     true,
	"LAWYER":                   true,
	"LB":                       true,
	"LC":                       true,
	"LEASE":                    true,
	"LGBT":                     true,
	"LI":                       true,
	"LIFE":                     true,
	"LIGHTING":                 true,
	"LIMITED":                  true,
	"LIMO":                     true,
	"LINK":                     true,
	"LK":                       true,
	"LOANS":                    true,
	"LONDON":                   true,
	"LOTTO":                    true,
	"LR":                       true,
	"LS":                       true,
	"LT":                       true,
	"LTDA":                     true,
	"LU":                       true,
	"LUXE":                     true,
	"LUXURY":                   true,
	"LV":                       true,
	"LY":                       true,
	"MA":                       true,
	"MAISON":                   true,
	"MANAGEMENT":               true,
	"MANGO":                    true,
	"MARKET":                   true,
	"MARKETING":                true,
	"MC":                       true,
	"MD":                       true,
	"ME":                       true,
	"MEDIA":                    true,
	"MEET":                     true,
	"MELBOURNE":                true,
	"MEME":                     true,
	"MENU":                     true,
	"MG":                       true,
	"MH":                       true,
	"MIAMI":                    true,
	"MIL":                      true,
	"MINI":                     true,
	"MK":                       true,
	"ML":                       true,
	"MM":                       true,
	"MN":                       true,
	"MO":                       true,
	"MOBI":                     true,
	"MODA":                     true,
	"MOE":                      true,
	"MONASH":                   true,
	"MORTGAGE":                 true,
	"MOSCOW":                   true,
	"MOTORCYCLES":              true,
	"MOV":                      true,
	"MP":                       true,
	"MQ":                       true,
	"MR":                       true,
	"MS":                       true,
	"MT":                       true,
	"MU":                       true,
	"MUSEUM":                   true,
	"MV":                       true,
	"MW":                       true,
	"MX":                       true,
	"MY":                       true,
	"MZ":                       true,
	"NA":                       true,
	"NAGOYA":                   true,
	"NAME":                     true,
	"NAVY":                     true,
	"NC":                       true,
	"NE":                       true,
	"NET":                      true,
	"NETWORK":                  true,
	"NEUSTAR":                  true,
	"NEW":                      true,
	"NEXUS":                    true,
	"NF":                       true,
	"NG":                       true,
	"NGO":                      true,
	"NHK":                      true,
	"NI":                       true,
	"NINJA":                    true,
	"NL":                       true,
	"NO":                       true,
	"NP":                       true,
	"NR":                       true,
	"NRA":                      true,
	"NRW":                      true,
	"NU":                       true,
	"NYC":                      true,
	"NZ":                       true,
	"OKINAWA":                  true,
	"OM":                       true,
	"ONG":                      true,
	"ONL":                      true,
	"OOO":                      true,
	"ORG":                      true,
	"ORGANIC":                  true,
	"OTSUKA":                   true,
	"OVH":                      true,
	"PA":                       true,
	"PARIS":                    true,
	"PARTNERS":                 true,
	"PARTS":                    true,
	"PE":                       true,
	"PF":                       true,
	"PG":                       true,
	"PH":                       true,
	"PHARMACY":                 true,
	"PHOTO":                    true,
	"PHOTOGRAPHY":              true,
	"PHOTOS":                   true,
	"PHYSIO":                   true,
	"PICS":                     true,
	"PICTURES":                 true,
	"PINK":                     true,
	"PIZZA":                    true,
	"PK":                       true,
	"PL":                       true,
	"PLACE":                    true,
	"PLUMBING":                 true,
	"PM":                       true,
	"PN":                       true,
	"POHL":                     true,
	"POST":                     true,
	"PR":                       true,
	"PRAXI":                    true,
	"PRESS":                    true,
	"PRO":                      true,
	"PROD":                     true,
	"PRODUCTIONS":              true,
	"PROF":                     true,
	"PROPERTIES":               true,
	"PROPERTY":                 true,
	"PS":                       true,
	"PT":                       true,
	"PUB":                      true,
	"PW":                       true,
	"PY":                       true,
	"QA":                       true,
	"QPON":                     true,
	"QUEBEC":                   true,
	"RE":                       true,
	"REALTOR":                  true,
	"RECIPES":                  true,
	"RED":                      true,
	"REHAB":                    true,
	"REISE":                    true,
	"REISEN":                   true,
	"REN":                      true,
	"RENTALS":                  true,
	"REPAIR":                   true,
	"REPORT":                   true,
	"REPUBLICAN":               true,
	"REST":                     true,
	"RESTAURANT":               true,
	"REVIEWS":                  true,
	"RICH":                     true,
	"RIO":                      true,
	"RO":                       true,
	"ROCKS":                    true,
	"RODEO":                    true,
	"RS":                       true,
	"RSVP":                     true,
	"RU":                       true,
	"RUHR":                     true,
	"RW":                       true,
	"RYUKYU":                   true,
	"SA":                       true,
	"SAARLAND":                 true,
	"SARL":                     true,
	"SB":                       true,
	"SC":                       true,
	"SCA":                      true,
	"SCB":                      true,
	"SCHMIDT":                  true,
	"SCHULE":                   true,
	"SCOT":                     true,
	"SD":                       true,
	"SE":                       true,
	"SERVICES":                 true,
	"SEXY":                     true,
	"SG":                       true,
	"SH":                       true,
	"SHIKSHA":                  true,
	"SHOES":                    true,
	"SI":                       true,
	"SINGLES":                  true,
	"SJ":                       true,
	"SK":                       true,
	"SL":                       true,
	"SM":                       true,
	"SN":                       true,
	"SO":                       true,
	"SOCIAL":                   true,
	"SOFTWARE":                 true,
	"SOHU":                     true,
	"SOLAR":                    true,
	"SOLUTIONS":                true,
	"SOY":                      true,
	"SPACE":                    true,
	"SPIEGEL":                  true,
	"SR":                       true,
	"ST":                       true,
	"SU":                       true,
	"SUPPLIES":                 true,
	"SUPPLY":                   true,
	"SUPPORT":                  true,
	"SURF":                     true,
	"SURGERY":                  true,
	"SUZUKI":                   true,
	"SV":                       true,
	"SX":                       true,
	"SY":                       true,
	"SYSTEMS":                  true,
	"SZ":                       true,
	"TATAR":                    true,
	"TATTOO":                   true,
	"TAX":                      true,
	"TC":                       true,
	"TD":                       true,
	"TECHNOLOGY":               true,
	"TEL":                      true,
	"TF":                       true,
	"TG":                       true,
	"TH":                       true,
	"TIENDA":                   true,
	"TIPS":                     true,
	"TIROL":                    true,
	"TJ":                       true,
	"TK":                       true,
	"TL":                       true,
	"TM":                       true,
	"TN":                       true,
	"TO":                       true,
	"TODAY":                    true,
	"TOKYO":                    true,
	"TOOLS":                    true,
	"TOP":                      true,
	"TOWN":                     true,
	"TOYS":                     true,
	"TP":                       true,
	"TR":                       true,
	"TRADE":                    true,
	"TRAINING":                 true,
	"TRAVEL":                   true,
	"TT":                       true,
	"TUI":                      true,
	"TV":                       true,
	"TW":                       true,
	"TZ":                       true,
	"UA":                       true,
	"UG":                       true,
	"UK":                       true,
	"UNIVERSITY":               true,
	"UNO":                      true,
	"UOL":                      true,
	"US":                       true,
	"UY":                       true,
	"UZ":                       true,
	"VA":                       true,
	"VACATIONS":                true,
	"VC":                       true,
	"VE":                       true,
	"VEGAS":                    true,
	"VENTURES":                 true,
	"VERSICHERUNG":             true,
	"VET":                      true,
	"VG":                       true,
	"VI":                       true,
	"VIAJES":                   true,
	"VILLAS":                   true,
	"VISION":                   true,
	"VLAANDEREN":               true,
	"VN":                       true,
	"VODKA":                    true,
	"VOTE":                     true,
	"VOTING":                   true,
	"VOTO":                     true,
	"VOYAGE":                   true,
	"VU":                       true,
	"WALES":                    true,
	"WANG":                     true,
	"WATCH":                    true,
	"WEBCAM":                   true,
	"WEBSITE":                  true,
	"WED":                      true,
	"WF":                       true,
	"WHOSWHO":                  true,
	"WIEN":                     true,
	"WIKI":                     true,
	"WILLIAMHILL":              true,
	"WME":                      true,
	"WORK":                     true,
	"WORKS":                    true,
	"WORLD":                    true,
	"WS":                       true,
	"WTC":                      true,
	"WTF":                      true,
	"XN--1QQW23A":              true,
	"XN--3BST00M":              true,
	"XN--3DS443G":              true,
	"XN--3E0B707E":             true,
	"XN--45BRJ9C":              true,
	"XN--4GBRIM":               true,
	"XN--55QW42G":              true,
	"XN--55QX5D":               true,
	"XN--6FRZ82G":              true,
	"XN--6QQ986B3XL":           true,
	"XN--80ADXHKS":             true,
	"XN--80AO21A":              true,
	"XN--80ASEHDB":             true,
	"XN--80ASWG":               true,
	"XN--90A3AC":               true,
	"XN--C1AVG":                true,
	"XN--CG4BKI":               true,
	"XN--CLCHC0EA0B2G2A9GCD":   true,
	"XN--CZR694B":              true,
	"XN--CZRU2D":               true,
	"XN--D1ACJ3B":              true,
	"XN--FIQ228C5HS":           true,
	"XN--FIQ64B":               true,
	"XN--FIQS8S":               true,
	"XN--FIQZ9S":               true,
	"XN--FPCRJ9C3D":            true,
	"XN--FZC2C9E2C":            true,
	"XN--GECRJ9C":              true,
	"XN--H2BRJ9C":              true,
	"XN--I1B6B1A6A2E":          true,
	"XN--IO0A7I":               true,
	"XN--J1AMH":                true,
	"XN--J6W193G":              true,
	"XN--KPRW13D":              true,
	"XN--KPRY57D":              true,
	"XN--KPUT3I":               true,
	"XN--L1ACC":                true,
	"XN--LGBBAT1AD8J":          true,
	"XN--MGB9AWBF":             true,
	"XN--MGBA3A4F16A":          true,
	"XN--MGBAAM7A8H":           true,
	"XN--MGBAB2BD":             true,
	"XN--MGBAYH7GPA":           true,
	"XN--MGBBH1A71E":           true,
	"XN--MGBC0A9AZCG":          true,
	"XN--MGBERP4A5D4AR":        true,
	"XN--MGBX4CD0AB":           true,
	"XN--NGBC5AZD":             true,
	"XN--NQV7F":                true,
	"XN--NQV7FS00EMA":          true,
	"XN--O3CW4H":               true,
	"XN--OGBPF8FL":             true,
	"XN--P1ACF":                true,
	"XN--P1AI":                 true,
	"XN--PGBS0DH":              true,
	"XN--Q9JYB4C":              true,
	"XN--RHQV96G":              true,
	"XN--S9BRJ9C":              true,
	"XN--SES554G":              true,
	"XN--UNUP4Y":               true,
	"XN--VERMGENSBERATER-CTB":  true,
	"XN--VERMGENSBERATUNG-PWB": true,
	"XN--VHQUV":                true,
	"XN--WGBH1C":               true,
	"XN--WGBL6A":               true,
	"XN--XHQ521B":              true,
	"XN--XKC2AL3HYE2A":         true,
	"XN--XKC2DL3A5EE0H":        true,
	"XN--YFRO4I67O":            true,
	"XN--YGBI2AMMX":            true,
	"XN--ZFR164B":              true,
	"XXX":                      true,
	"XYZ":                      true,
	"YACHTS":                   true,
	"YANDEX":                   true,
	"YE":                       true,
	"YOKOHAMA":                 true,
	"YOUTUBE":                  true,
	"YT":                       true,
	"ZA":                       true,
	"ZIP":                      true,
	"ZM":                       true,
	"ZONE":                     true,
	"ZW":                       true,
}

// ExtendedTLDs is a set of additional "TLDs", allowing decentralized name
// systems, like TOR and Namecoin.
var ExtendedTLDs = map[string]bool{
	"BIT":   true,
	"ONION": true,
}
