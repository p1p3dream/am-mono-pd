package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"abodemine/lib/consts"
	"abodemine/lib/errors"
	"abodemine/lib/gconf"
	"abodemine/lib/stringutil"
	"abodemine/lib/val"
)

var txtCmd = &cobra.Command{
	Use:          "txt",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := gconf.LoadZerolog("", false); err != nil {
			return errors.Forward(err, "23196b80-4883-4065-af5a-2737b87d0194")
		}

		if err := statToPartnerDataFile(&statToPartnerDataFileInput{
			PackageName: viper.GetString("txt.package"),
			StatFile:    viper.GetString("txt.stat"),
			StructName:  viper.GetString("txt.struct"),
			TableName:   viper.GetString("txt.table"),
		}); err != nil {
			return errors.Forward(err, "84c96f23-9d0d-46ef-99a6-5b16d6f9dc19")
		}

		return nil
	},
}

func init() {
	txtCmd.PersistentFlags().String("package", "", "The go package name.")
	if err := viper.BindPFlag("txt.package", txtCmd.PersistentFlags().Lookup("package")); err != nil {
		panic(err)
	}

	txtCmd.PersistentFlags().String("stat", "", "Input csvstat json file.")
	if err := viper.BindPFlag("txt.stat", txtCmd.PersistentFlags().Lookup("stat")); err != nil {
		panic(err)
	}

	txtCmd.PersistentFlags().String("struct", "", "The struct.")
	if err := viper.BindPFlag("txt.struct", txtCmd.PersistentFlags().Lookup("struct")); err != nil {
		panic(err)
	}

	txtCmd.PersistentFlags().String("table", "", "The table name.")
	if err := viper.BindPFlag("txt.table", txtCmd.PersistentFlags().Lookup("table")); err != nil {
		panic(err)
	}

	mainCmd.AddCommand(txtCmd)
}

type statToPartnerDataFileInput struct {
	PackageName string
	StatFile    string
	StructName  string
	TableName   string
}

func statToPartnerDataFile(in *statToPartnerDataFileInput) error {
	b, err := os.ReadFile(in.StatFile)
	if err != nil {
		return &errors.Object{
			Id:     "8dd62b05-b834-4584-910a-76c819418d62",
			Code:   errors.Code_FAILED_PRECONDITION,
			Detail: "Failed to read stat file.",
		}
	}

	var statFileBody []*ColumnStat

	if err := json.Unmarshal(b, &statFileBody); err != nil {
		return &errors.Object{
			Id:     "2f4347c2-dee7-436e-b6f1-6f68e2d01f9f",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to unmarshal stat file.",
			Cause:  err.Error(),
		}
	}

	inferDateFormat := func(s string) (string, error) {
		if s == "" {
			return "", &errors.Object{
				Id:     "f0a2b1c4-5d3e-4b8c-9f7d-6e1a2f3b4c5d",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Empty input.",
			}
		}

		if _, err := time.Parse(consts.RFC3339Date, s); err == nil {
			return "consts.RFC3339Date", nil
		}

		if _, err := time.Parse(consts.USSlashDate, s); err == nil {
			return "consts.USSlashDate", nil
		}

		// Ensure we can read int dates with leading and trailing zeros.
		if intDate, err := decimal.NewFromString(s); err == nil {
			if _, err := time.Parse(consts.IntegerDate, intDate.String()); err == nil {
				return "consts.IntegerDate", nil
			}
		}

		return "", &errors.Object{
			Id:     "2c97d0ef-e279-464b-b08a-84f88b924747",
			Code:   errors.Code_INVALID_ARGUMENT,
			Detail: "Failed to infer date format.",
			Meta: map[string]any{
				"value": s,
			},
		}
	}

	inferNumberType := func(colName string, s string, maxPrecision int) (string, error) {
		v, err := decimal.NewFromString(s)
		if err != nil {
			return "", &errors.Object{
				Id:     "11b1ac94-acfb-4f70-bb53-870ee9ea8697",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Failed to parse number.",
				Cause:  err.Error(),
				Meta: map[string]any{
					"column_name": colName,
					"value":       s,
				},
			}
		}

		// Check maxprecision to ensure we're not comparing only
		// against min/max values that happens to be integers.
		if maxPrecision == 0 && v.IsInteger() {
			int64value := v.IntPart()

			if int64value > math.MaxInt32 {
				return "int64", nil
			} else {
				return "int", nil
			}
		}

		// Probably non-monetary values.
		if maxPrecision > 4 {
			return "float64", nil
		}

		// Ensure best precision.
		return "decimal.Decimal", nil
	}

	dataFileTemplateIn := &DataFileTemplateInput{
		PackageName: in.PackageName,
		StructName:  in.StructName,
		TableName:   in.TableName,
	}

	var maxPostgresNameLength int

	for _, stat := range statFileBody {
		goName := new(strings.Builder)

		for i, c := range stat.ColumnName {
			if i > 0 && c == '_' {
				goName.WriteRune(c)
				continue
			}

			// Only write alphanumeric characters.
			if (c < 'A' || c > 'Z') && (c < 'a' || c > 'z') && (c < '0' || c > '9') {
				continue
			}

			goName.WriteRune(c)
		}

		field := &DataFileTemplateField{
			ColumnName:   stat.ColumnName,
			GoName:       goName.String(),
			PostgresName: stringutil.ASCIIToSnakeCase(goName.String()),
		}

		switch stat.Type {
		case "decimal":
			field.GoType = "decimal.Decimal"
			field.PostgresType = "decimal"

			if !stat.Nulls {
				field.GoValidationType = GoValidationMustDecimal
			} else {
				field.GoPtr = true
				field.GoValidationType = GoValidationDecimalPtrFromStringIfNonZero
			}
		case "float64":
			field.GoType = "float64"
			field.PostgresType = "double precision"

			if !stat.Nulls {
				field.GoValidationType = GoValidationMustFloat64
			} else {
				field.GoPtr = true
				field.GoValidationType = GoValidationFloat64PtrFromStringIfNonZero
			}
		case "int":
			field.GoType = "int"
			field.PostgresType = "integer"

			if !stat.Nulls {
				field.GoValidationType = GoValidationMustInt
			} else {
				field.GoPtr = true
				field.GoValidationType = GoValidationIntPtrFromStringIfNonZero
			}
		case "int64":
			field.GoType = "int64"
			field.PostgresType = "bigint"

			if !stat.Nulls {
				field.GoValidationType = GoValidationMustInt64
			} else {
				field.GoPtr = true
				field.GoValidationType = GoValidationInt64PtrFromStringIfNonZero
			}
		case "Boolean":
			field.GoType = "bool"
			field.PostgresType = "boolean"

			if !stat.Nulls {
				field.GoValidationType = GoValidationMustBool
			} else {
				field.GoPtr = true
				field.GoValidationType = GoValidationBoolPtrFromStringIfNonZero
			}
		case "date", "Date":
			field.GoType = "time.Time"
			field.PostgresType = "date"

			// Min and max should always be either both present or absent,
			// but we check them anyway just to be safe.
			if v, err := inferDateFormat(val.PtrDeref(stat.Min.StringValue)); err == nil {
				field.GoTimeFormat = v
			} else if v, err := inferDateFormat(val.PtrDeref(stat.Max.StringValue)); err == nil {
				field.GoTimeFormat = v
			} else {
				return &errors.Object{
					Id:     "4092c815-9d43-446c-b226-893580cf207f",
					Code:   errors.Code_INVALID_ARGUMENT,
					Detail: "Failed to infer date format.",
					Cause:  err.Error(),
				}
			}

			if !stat.Nulls {
				field.GoValidationType = GoValidationMustTime
			} else {
				field.GoPtr = true
				field.GoValidationType = GoValidationTimePtrFromStringIfNonZero
			}
		case "Number":
			minType, err := inferNumberType(stat.ColumnName, val.PtrDeref(stat.Min.StringValue), stat.MaxPrecision)
			if err != nil {
				return &errors.Object{
					Id:     "a3f9e521-3901-4f1d-a3e5-8c1fc3cc564e",
					Code:   errors.Code_INVALID_ARGUMENT,
					Detail: "Failed to infer number type.",
					Cause:  err.Error(),
				}
			}

			maxType, err := inferNumberType(stat.ColumnName, val.PtrDeref(stat.Max.StringValue), stat.MaxPrecision)
			if err != nil {
				return &errors.Object{
					Id:     "c99d2a3a-b076-4934-b691-915ea75764bc",
					Code:   errors.Code_INVALID_ARGUMENT,
					Detail: "Failed to infer number type.",
					Cause:  err.Error(),
				}
			}

			if stat.Nulls {
				field.GoPtr = true
			}

			switch {
			case minType == "decimal.Decimal" || maxType == "decimal.Decimal":
				field.GoType = "decimal.Decimal"
				field.PostgresType = "decimal"

				if stat.Nulls {
					field.GoValidationType = GoValidationDecimalPtrFromStringIfNonZero
				} else {
					field.GoValidationType = GoValidationMustDecimal
				}
			case minType == "float64" || maxType == "float64":
				field.GoType = "float64"
				field.PostgresType = "double precision"

				if stat.Nulls {
					field.GoValidationType = GoValidationFloat64PtrFromStringIfNonZero
				} else {
					field.GoValidationType = GoValidationMustFloat64
				}
			case minType == "int64" || maxType == "int64":
				field.GoType = "int64"
				field.PostgresType = "bigint"

				if stat.Nulls {
					field.GoValidationType = GoValidationInt64PtrFromStringIfNonZero
				} else {
					field.GoValidationType = GoValidationMustInt64
				}
			default:
				field.GoType = "int"
				field.PostgresType = "integer"

				if stat.Nulls {
					field.GoValidationType = GoValidationIntPtrFromStringIfNonZero
				} else {
					field.GoValidationType = GoValidationMustInt
				}
			}
		case "string", "Text":
			field.GoType = "string"
			field.PostgresType = "text"

			if !stat.Nulls {
				field.GoValidationType = GoValidationMustString
			} else {
				field.GoPtr = true
				field.GoValidationType = GoValidationStringPtrIfNonZero
			}
		default:
			return &errors.Object{
				Id:     "f058761f-0b8b-4329-bd4b-417c8843dd16",
				Code:   errors.Code_INVALID_ARGUMENT,
				Detail: "Unsupported type.",
				Meta: map[string]any{
					"column_name": stat.ColumnName,
					"type":        stat.Type,
				},
			}
		}

		if l := len(field.PostgresName); l > maxPostgresNameLength {
			maxPostgresNameLength = l
		}

		dataFileTemplateIn.Fields = append(dataFileTemplateIn.Fields, field)
	}

	if err := os.MkdirAll(".tmp", os.ModePerm); err != nil {
		return &errors.Object{
			Id:     "e039212b-bf7b-4dcc-ac5e-7154c74089fe",
			Code:   errors.Code_FAILED_PRECONDITION,
			Detail: "Failed to create tmp directory.",
			Cause:  err.Error(),
		}
	}

	goFile, err := os.Create(filepath.Join(".tmp", fmt.Sprintf("df_%s.go", stringutil.ASCIIToSnakeCase(dataFileTemplateIn.StructName))))
	if err != nil {
		return &errors.Object{
			Id:     "df0a45e7-4741-4d34-9ccd-386d7f2d3d32",
			Code:   errors.Code_FAILED_PRECONDITION,
			Detail: "Failed to create go file.",
			Cause:  err.Error(),
		}
	}

	defer goFile.Close()

	t := template.Must(template.New("datafile").Funcs(sprig.FuncMap()).Parse(dataFileTemplate))

	if err := t.Execute(goFile, dataFileTemplateIn); err != nil {
		return &errors.Object{
			Id:     "3239fe51-2e21-4811-916b-b790f6279e4a",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to execute template.",
			Cause:  err.Error(),
		}
	}

	// Calculate paddings.
	for _, field := range dataFileTemplateIn.Fields {
		field.PostgresNamePadding = strings.Repeat(" ", maxPostgresNameLength-len(field.PostgresName)+1)
	}

	sqlFile, err := os.Create(filepath.Join(".tmp", dataFileTemplateIn.TableName+".sql"))
	if err != nil {
		return &errors.Object{
			Id:     "616b2a40-219e-4408-aa76-42aa621074a1",
			Code:   errors.Code_FAILED_PRECONDITION,
			Detail: "Failed to create go file.",
			Cause:  err.Error(),
		}
	}

	defer sqlFile.Close()

	t = template.Must(template.New("sqlfile").Funcs(sprig.FuncMap()).Parse(sqlFileTemplate))

	if err := t.Execute(sqlFile, dataFileTemplateIn); err != nil {
		return &errors.Object{
			Id:     "425eb3ad-916b-4f5b-8d21-bccb2cca0f8a",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to execute template.",
			Cause:  err.Error(),
		}
	}

	return nil
}
