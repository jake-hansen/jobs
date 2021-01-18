package consumers

import "fmt"

// DataConsumer performs some action upon data.
type DataConsumer interface {
	Consume(data interface{})
}

// ErrorConsumer performs some action upon an error.
type ErrorConsumer interface {
	Consume(err error)
}

// DataPrinterConsumer is a consumer that prints data.
type DataPrinterConsumer struct {}

// ErrorPrinterConsumer is a consumer that prints an error.
type ErrorPrinterConsumer struct {}

func (d DataPrinterConsumer) Consume(i interface{}) {
	fmt.Printf("value %d\n", i)
}

func (e ErrorPrinterConsumer) Consume(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}
