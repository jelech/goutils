package parquetutil

import (
	"fmt"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

func Read(filePath string, stuTypePoint interface{}, stus interface{}, callback func(interface{}) error) error {
	var err error
	fr, err := local.NewLocalFileWriter("flat.parquet")
	if err != nil {
		return err
	}
	defer fr.Close()

	pr, err := reader.NewParquetReader(fr, stuTypePoint, 4)
	if err != nil {
		return err
	}
	defer pr.ReadStop()

	num := int(pr.GetNumRows())
	for i := 0; i < num/10; i++ {
		if err = pr.Read(&stus); err != nil {
			return err
		}

		fmt.Println(stus)
		if err = callback(stus); err != nil {
			return err
		}
	}

	return nil
}

func ReadSimple(filePath string, stuTypePoint interface{}, stus interface{}, callback func(interface{}) error) error {
	var err error
	fr, err := local.NewLocalFileReader(filePath)
	if err != nil {
		return err
	}
	defer fr.Close()

	pr, err := reader.NewParquetReader(fr, nil, 4)
	if err != nil {
		return err
	}

	num := int(pr.GetNumRows())
	res, err := pr.ReadByNumber(num)
	if err != nil {
		return err
	}

	return callback(res)
}
