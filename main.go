package main

import (
	"flag"
	"fmt"
	"go-platio-example/platio"
	"os"
)

func main() {
	authorization := flag.String("a", "", "Authorization header")
	flag.Parse()

	const collectionUrl = "https://api.plat.io/v1/pwdhds3gsg5chpc6p4oes3af2ki/collections/t1c7d21c"

	api := platio.NewAPI(collectionUrl, *authorization)
	records, err := api.GetLatestRecords(10)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	if len(records) == 0 {
		fmt.Println("No records found.")
		return
	}

	type RecordWithAge struct {
		platio.Record
		age int
	}
	recordWithAges := make([]RecordWithAge, 0, len(records))

	for _, record := range records {
		name := func() string {
			if record.Values.Name != nil {
				return record.Values.Name.Value
			} else {
				return ""
			}
		}()
		age := func() int {
			if record.Values.Age != nil {
				return int(record.Values.Age.Value)
			} else {
				return 0
			}
		}()

		recordWithAges = append(recordWithAges, RecordWithAge{record, age})

		fmt.Printf("Id: %s, Name: %s, Age: %d\n", record.Id, name, age)
	}

	type RecordIdWithError struct {
		recordId platio.RecordId
		error    error
	}
	recordIdWithErrors := make(chan RecordIdWithError)

	for _, recordWithAge := range recordWithAges {
		recordWithAge := recordWithAge
		go func() {
			err = api.UpdateRecord(recordWithAge.Id, &platio.Values{
				Age: &platio.NumberValue{float64(recordWithAge.age + 1)},
			})
			recordIdWithErrors <- RecordIdWithError{recordWithAge.Id, err}
		}()
	}

	for n := 0; n < len(recordWithAges); n++ {
		recordIdWithError := <-recordIdWithErrors
		if recordIdWithError.error != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", recordIdWithError.recordId, recordIdWithError.error)
		}
	}
}
