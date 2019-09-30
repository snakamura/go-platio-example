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

	for _, record := range records {
		fmt.Printf("Id: %s, Name: %s, Age: %d\n", record.Id, record.Name(), record.Age())
	}

	type RecordIdWithError struct {
		recordId platio.RecordId
		error    error
	}
	recordIdWithErrors := make(chan RecordIdWithError)

	for _, record := range records {
		record := record
		go func() {
			err = api.UpdateRecord(record.Id, &platio.Values{
				Age: &platio.NumberValue{float64(record.Age() + 1)},
			})
			recordIdWithErrors <- RecordIdWithError{record.Id, err}
		}()
	}

	for range records {
		recordIdWithError := <-recordIdWithErrors
		if recordIdWithError.error != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", recordIdWithError.recordId, recordIdWithError.error)
		}
	}
}
