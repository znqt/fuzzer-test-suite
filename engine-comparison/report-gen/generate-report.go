package main

import (
        "fmt"
//        "io"
	"path"
        "io/ioutil"
//        "strings"
        "strconv"
        "encoding/csv"
        "os"
)
func checkErr(e error){
	if (e != nil) {
		fmt.Println("POOP!")
		os.Exit(1)
	}
}

// Removes non-directory elements of any []os.FileInfo (a helper function, for convenience)
// Makes use of idiomatic Golang return: out_ls is declared and returned implicitly
func onlyDirectories(potential_files []os.FileInfo) (out_ls []os.FileInfo) {
	for _, fd := range potential_files {
		if fd.IsDir() {
			out_ls = append(out_ls, fd)
		}
	}
	return
}

// Extend "records" matrix to have rows until time "desired_time"
// Return: Extended version of record
func extendRecordsToTime(records [][]string, desired_time int, recordCols int) ([][]string)  {
	lenr := len(records)
	// records[1] stores cycle [1], as records[0] is column names
	for j := lenr; j < desired_time + 1; j++ {
		records = append(records, make([]string, recordCols))
		records[j][0] = strconv.Itoa(j)
		for k := 1; k < recordCols; k++ {
			records[j][k] = ""
		}
	}
	return records
}

// Enters all report subdirectories, from benchmark to fengine to trial;
// composes individual CSVs (only two columns) into larger CSVs
func composeAllNamed(desired_report_fname string) {
	master_path := "./reports"
	bmarks, err := ioutil.ReadDir(master_path)
	checkErr(err)
	for _, bmark := range bmarks {
		// all_fe_file, err := os.Create(path.Join(master_path, bmark.Name(), desired_report_fname))
		// checkErr(err)
		// defer all_fe_file.Close()
		// all_fe_writer := csv.NewWriter(all_fe_file)
		// meta_records := [][]string{{"time"}}
		potential_fengines, err := ioutil.ReadDir(path.Join(master_path, bmark.Name()))
		checkErr(err)
		// narrow potential_fengines to fengines so the indices of `range fengines` are useful
		fengines := onlyDirectories(potential_fengines)

		for _, fengine := range fengines {
			// Create fds
			this_fe_file, err := os.Create(path.Join(master_path, bmark.Name(), fengine.Name(), desired_report_fname))
			checkErr(err)
			defer this_fe_file.Close()
			this_fe_writer := csv.NewWriter(this_fe_file)

			// Create matrix, to eventually become a CSV
			records := [][]string{{"time"}}

			// Enter sub-directories
			potential_trials, err := ioutil.ReadDir(path.Join(master_path, bmark.Name(), fengine.Name()))
			checkErr(err)
			trials := onlyDirectories(potential_trials)

			num_record_columns := len(trials) + 1
			for j, trial := range trials {
				// Create fds
				this_file, err := os.Open(path.Join(master_path, bmark.Name(), fengine.Name(), trial.Name(), desired_report_fname))
				checkErr(err)
				defer this_file.Close()
				this_reader := csv.NewReader(this_file)

				// Read whole CSV to an array
				experiment_records, err := this_reader.ReadAll()
				checkErr(err)
				// Add the name of this new column to records[0]
				records[0] = append(records[0], fengine.Name() + trial.Name())

				for _, row := range experiment_records {
					// row[0] is time, on the x-axis; row[1] is value, on the y-axis
					time_now, err := strconv.Atoi(row[0])
					checkErr(err)
					//If this test went longer than all of the others, so far
					if(len(records) < time_now + 1) {
						records = extendRecordsToTime(records, time_now, num_record_columns)
					}
					records[time_now][j+1] = row[1]
				}
			}
			this_fe_writer.WriteAll(records)
			// Potentially put this fengine into a broader comparison CSV
		}
		// TODO: create comparison between fengines, having already composed trials
		// Do this by identifying the max (or potentially median) performing trial
		// For each fengine, and putting them all into a CSV which can be graphed
	}
}

func main(){
	composeAllNamed("coverage-graph.csv")
	composeAllNamed("corpus-size-graph.csv")
	composeAllNamed("corpus-elems-graph.csv")
	// createIFramesFor("setOfFrames.html")
	// <iframe width="960" height="500" src="benchmarkN/report.html" frameborder="0"></iframe>
}

