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
	record, err := api.GetLatestRecord()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	if record == nil {
		fmt.Println("No records found.")
		return
	}

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

	fmt.Printf("Id: %s, Name: %s, Age: %d\n", record.Id, name, age)

	err = api.UpdateRecord(record.Id, &platio.Values{
		Age: &platio.NumberValue{float64(age + 1)},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
}
