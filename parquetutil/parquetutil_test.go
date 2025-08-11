package parquetutil

import (
	"fmt"
	"testing"
)

type RLData struct {
	Idx             string    `parquet:"name=idx, type=BYTE_ARRAY, convertedtype=UTF8"`
	Date            int64     `parquet:"name=date, type=INT64, convertedtype=DATE"`
	Leadtime        int64     `parquet:"name=leadtime, type=INT64, convertedtype=INT_32"`
	Actual          int64     `parquet:"name=actual, type=INT64, convertedtype=INT_32"`
	PredY           float64   `parquet:"name=pred_y, type=DOUBLE"`
	PredictedDemand []float64 `parquet:"name=predicted_demand, type=DOUBLE, repetitiontype=REPEATED"`
	InitialStock    float64   `parquet:"name=initial_stock, type=DOUBLE"`
}

func TestRead(t *testing.T) {
	typeData := RLData{}
	arrData := make([]RLData, 10)
	err := ReadSimple("", &typeData, arrData, func(i interface{}) error {
		x, ok := i.([]RLData)
		if !ok {
			return fmt.Errorf("type assertion failed")
		}
		fmt.Println(x)
		return nil
	})
	if err != nil {
		t.Errorf("Read error: %v", err)
		return
	}
}
