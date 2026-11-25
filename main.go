package main

import (
	"database/sql"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3" //import for side effects
	"github.com/xuri/excelize/v2"
)

func main() { //taken from class demo on GitHub
	myDatabase := OpenDataBase("./Demo.db")
	defer myDatabase.Close()
	create_tables(myDatabase)
	//add_sample_data(myDatabase)
	InsertClassTableData(myDatabase)
}
func OpenDataBase(dbfile string) *sql.DB { //also taken from class demo
	database, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		log.Fatal(err)
	}
	return database
}

func create_tables(database *sql.DB) {
	//createStatement1 := "CREATE TABLE IF NOT EXISTS students(    " + //carried over to test formatting
	//	"banner_id INTEGER PRIMARY KEY," +
	//	"first_name TEXT NOT NULL," +
	//	"last_name TEXT NOT NULL," +
	//	"gpa REAL DEFAULT 0," +
	//	"credits INTEGER DEFAULT 0);"
	createStatement1 := "CREATE TABLE IF NOT EXISTS salaries(	" +
		"area INTEGER NOT NULL," +
		"area_title TEXT NOT NULL," +
		"area_type INTEGER NOT NULL," +
		"prim_state TEXT NOT NULL," +
		"naics INTEGER NOT NULL," +
		"naics_title TEXT NOT NULL," +
		"i_group TEXT NOT NULL," +
		"own_code INTEGER NOT NULL," +
		"occ_code TEXT NOT NULL," +
		"occ_title TEXT NOT NULL," +
		"o_group TEXT NOT NULL," +
		"tot_emp VARCHAR," +
		"emp_prse VARCHAR," +
		"jobs_1000 VARCHAR," +
		"loc_quotient VARCHAR," +
		"pct_total VARCHAR," +
		"pct_rpt VARCHAR," +
		"h_mean VARCHAR," +
		"a_mean VARCHAR," +
		"mean_prse VARCHAR," +
		"h_pct10 VARCHAR," +
		"h_pct25 VARCHAR," +
		"h_median VARCHAR," +
		"h_pct75 VARCHAR," +
		"h_pct90 VARCHAR," +
		"a_pct10 VARCHAR," +
		"a_pct25 VARCHAR," +
		"a_median VARCHAR," +
		"a_pct75 VARCHAR," +
		"a_pct90 VARCHAR," +
		"annual BOOLEAN DEFAULT NULL," +
		"hourly BOOLEAN DEFAULT NULL," +
		"FOREIGN KEY (area_title) REFERENCES states (state_name) ON DELETE CASCADE ON UPDATE NO ACTION);"
	database.Exec(createStatement1) //creates salary data table

	createStatement2 := "CREATE TABLE IF NOT EXISTS rentdata( " +
		"name TEXT NOT NULL," +
		"period VARCHAR," +
		"value INTEGER NOT NULL," +
		"preceeding_period INTEGER NOT NULL," +
		"year_ago_from_period INTEGER NOT NULL);"
	database.Exec(createStatement2) //creates rent data table

	createStatement3 := "CREATE TABLE IF NOT EXISTS homeprices( " +
		"hpi_type TEXT NOT NULL," +
		"hpi_flavor TEXT NOT NULL," +
		"frequency TEXT NOT NULL," +
		"level TEXT NOT NULL," +
		"place_name TEXT NOT NULL," +
		"place_id TEXT NOT NULL," +
		"yr INTEGER NOT NULL," +
		"period INTEGER NOT NULL," +
		"index_nsa REAL," +
		"index_sa REAL," +
		"rstderr VARCHAR," +
		"note VARCHAR);"
	database.Exec(createStatement3) //creates home price table

	createStatement4 := "CREATE TABLE IF NOT EXISTS states( " +
		"state_name TEXT PRIMARY KEY," +
		"state_code TEXT NOT NULL," +
		"nationwide BOOLEAN NOT NULL);"
	database.Exec(createStatement4) //creates state data table
}
func InsertClassTableData(database *sql.DB) {
	//insert state data
	file, err := excelize.OpenFile("statedata.xlsx")
	if err != nil {
		log.Fatalf("Failed at loading xlsx", err)
	}
	defer file.Close()
	rows, err := file.GetRows("Sheet1") //had to look up this command
	if err != nil {
		log.Fatalf("Failed at loading statedata sheet", err)
	}
	//add to table
	for i, row := range rows {
		if i != 0 {
			//fmt.Println(row)
			stateName := row[0]
			stateCode := row[1]
			nationwide, err := strconv.ParseBool(row[2])
			if err != nil {
				log.Fatalf("Failed at parsing nationwide row", err)
			}
			insertStatement := "INSERT INTO states (state_name, state_code, nationwide) VALUES (?, ?, ?);"
			preppedStatement, err := database.Prepare(insertStatement)
			if err != nil {
				log.Fatalf("Failed at prepping", err)
			}
			_, err = preppedStatement.Exec(stateName, stateCode, nationwide)
		}
	}

	//insert salary data (this was a nightmare)
	file, err = excelize.OpenFile("BLSsalaryDara2025.xlsx")
	if err != nil {
		log.Fatalf("Failed at loading xlsx", err)
	}
	defer file.Close()
	rows, err = file.GetRows("state_M2025_dl") //had to look up this command
	if err != nil {
		log.Fatalf("Failed at loading salarydata sheet", err)
	}
	//add to table
	for i, row := range rows {
		if i != 0 {
			for len(row) < 32 { //had to add this to lengthen the row if it got no value for annual or hourly (not sure if this is good solution, but otherwise instantly crashed)
				row = append(row, "")
			}
			/*fmt.Println(expectedLength)
			fmt.Println(realLength)
			for _, cell := range row {
				fmt.Println(cell)
			}*/
			//fmt.Println(row)
			area, err := strconv.Atoi(row[0])
			if err != nil {
				log.Fatalf("Failed at parsing area", err)
			}
			areaTitle := row[1]
			areaType, err := strconv.Atoi(row[2])
			if err != nil {
				log.Fatalf("Failed at parsing areaType", err)
			}
			primState := row[3]
			naics, err := strconv.Atoi(row[4])
			if err != nil {
				log.Fatalf("Failed at parsing naics", err)
			}
			naicsTitle := row[5]
			iGroup := row[6]
			ownCode, err := strconv.Atoi(row[7])
			if err != nil {
				log.Fatalf("Failed at parsing ownCode", err)
			}
			occCode := row[8]
			occTitle := row[9]
			oGroup := row[10]
			totEmp := row[11]
			emp := row[12]
			jobs := row[13]
			loc := row[14]
			pctTotal := row[15]
			pctRpt := row[16]
			hMean := row[17]
			aMean := row[18]
			meanPrse := row[19]
			hPct10 := row[20]
			hPct25 := row[21]
			hMedian := row[22]
			hPct75 := row[23]
			hPct90 := row[24]
			aPct10 := row[25]
			aPct25 := row[26]
			aMedian := row[27]
			aPct75 := row[28]
			aPct90 := row[29]

			//ran into a lot of issues with these columns being mostly empty
			var annual, hourly bool
			if row[30] != "" {
				annual, err = strconv.ParseBool(row[30])
				if err != nil {
					log.Fatalf("Failed at parsing annual row", err)
				}
			}
			if row[31] != "" {
				hourly, err = strconv.ParseBool(row[31])
				if err != nil {
					log.Fatalf("Failed at parsing hourly row", err)
				}
			}

			insertStatement := "INSERT INTO salaries (area, area_title, area_type, prim_state, naics, naics_title, i_group, own_code, occ_code, occ_title, o_group, tot_emp, emp_prse, jobs_1000, loc_quotient, pct_total, pct_rpt, h_mean, a_mean, mean_prse, h_pct10, h_pct25, h_median, h_pct75, h_pct90, a_pct10, a_pct25, a_median, a_pct75, a_pct90, annual, hourly) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);"
			preppedStatement, err := database.Prepare(insertStatement)
			if err != nil {
				log.Fatalf("Failed at prepping", err)
			}
			_, err = preppedStatement.Exec(area, areaTitle, areaType, primState, naics, naicsTitle, iGroup, ownCode, occCode, occTitle, oGroup, totEmp, emp, jobs, loc, pctTotal, pctRpt, hMean, aMean, meanPrse, hPct10, hPct25, hMedian, hPct75, hPct90, aPct10, aPct25, aMedian, aPct75, aPct90, annual, hourly)
		}
	}

	//insert rent data
	file, err = excelize.OpenFile("RentPer10SquareFt2026.xlsx")
	if err != nil {
		log.Fatalf("Failed at loading xlsx", err)
	}
	defer file.Close()
	rows, err = file.GetRows("RentPerSquareFootAug2026")
	if err != nil {
		log.Fatalf("Failed at loading statedata sheet", err)
	}
	//add to table
	for i, row := range rows {
		if i > 1 { // this table has 2 header rows???
			name := row[0]
			period := row[1]
			value, err := strconv.Atoi(row[2])
			if err != nil {
				log.Fatalf("Failed at parsing value", err)
			}
			preceedingPeriod, err := strconv.Atoi(row[3])
			if err != nil {
				log.Fatalf("Failed at parsing preceedingPeriod", err)
			}
			yearAgoFromPeriod, err := strconv.Atoi(row[4])
			if err != nil {
				log.Fatalf("Failed at parsing yearAgoFromPeriod", err)
			}

			insertStatement := "INSERT INTO rentdata (name, period, value, preceeding_period, year_ago_from_period) VALUES (?,?,?,?,?);"
			preppedStatement, err := database.Prepare(insertStatement)
			if err != nil {
				log.Fatalf("Failed at prepping", err)
			}
			_, err = preppedStatement.Exec(name, period, value, preceedingPeriod, yearAgoFromPeriod)
		}
	}

	//insert homeprices data
	file, err = excelize.OpenFile("homepriceindex.xlsx")
	if err != nil {
		log.Fatalf("Failed at loading xlsx", err)
	}
	defer file.Close()
	rows, err = file.GetRows("hpi_master")
	if err != nil {
		log.Fatalf("Failed at loading statedata sheet", err)
	}
	//add to table
	for i, row := range rows {
		if i != 0 && row[6] == "2026" {
			for len(row) < 12 { //took a guess that the index issue would happen again with the blank columns at the end
				row = append(row, "")
			}
			hpiType := row[0]
			hpiFlavor := row[1]
			frequency := row[2]
			level := row[3]
			placeName := row[4]
			placeId := row[5]
			yr, err := strconv.Atoi(row[6])
			if err != nil {
				log.Fatalf("Failed at parsing yr", err)
			}
			period, err := strconv.Atoi(row[7])
			if err != nil {
				log.Fatalf("Failed at parsing period", err)
			}
			indexNsa, err := strconv.ParseFloat(row[8], 64)
			if err != nil {
				log.Fatalf("Failed at parsing indexNsa", err)
			}

			var indexSa *float64 //* so that it defaults to null instead of filling everything with 0's
			//fmt.Println(indexSa) prints nil as it should, troubleshooting (program filling indexSa with random numbers seemingly)
			//Turns out I had indexSa set as not nullable. I spent an hour figuring that out.
			if row[9] != "" {
				val, err := strconv.ParseFloat(row[9], 64)
				if err != nil {
					log.Fatalf("Failed at parsing indexSa", err)
				}
				indexSa = &val
			}
			rstderr := row[10]
			note := row[11]

			insertStatement := "INSERT INTO homeprices (hpi_type, hpi_flavor, frequency, level, place_name, place_id, yr, period, index_nsa, index_sa, rstderr, note) VALUES (?,?,?,?,?,?,?,?,?,?,?,?);"
			preppedStatement, err := database.Prepare(insertStatement)
			if err != nil {
				log.Fatalf("Failed at prepping", err)
			}
			_, err = preppedStatement.Exec(hpiType, hpiFlavor, frequency, level, placeName, placeId, yr, period, indexNsa, indexSa, rstderr, note)
		}
	}
}
