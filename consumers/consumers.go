// Copyright © 2021 Jacob Hansen. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

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
type DataPrinterConsumer struct{}

// ErrorPrinterConsumer is a consumer that prints an error.
type ErrorPrinterConsumer struct{}

func (d DataPrinterConsumer) Consume(i interface{}) {
	fmt.Printf("worker value received: %v\n", i)
}

func (e ErrorPrinterConsumer) Consume(err error) {
	if err != nil {
		fmt.Printf("worker error: %s", err.Error())
	}
}
