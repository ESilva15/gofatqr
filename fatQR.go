package gofatqr

// TODO
// Check what we can take form SAF-T to put here, maybe work on a SAF-T
// implementation too, who knows amirite

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	dec "github.com/shopspring/decimal"
)

type ScanMode uint32

// Setting this up for the future
const (
	Strict ScanMode = 1 << iota
	NifValidation
)

type FiscalSpace struct {
	TaxCountryRegion        string       // 1 - fiscal space
	TaxableBase             *dec.Decimal // 2
	TaxableReduced          *dec.Decimal // 3
	VatTotalReducedTax      *dec.Decimal // 4
	TaxableIntermediateBase *dec.Decimal // 5
	VatTotalIntermediateTax *dec.Decimal // 6
	TaxableNormalBase       *dec.Decimal // 7
	VatTotalNormalBase      *dec.Decimal // 8
}

type FatQR struct {
	TaxRegistrationNumber     string       // A
	CustomerTaxID             string       // B
	Country                   string       // C customer country ISO 3166-1 alpha2
	InvoiceType               string       // D
	InvoiceStatus             string       // E
	InvoiceDate               time.Time    // F
	InvoiceNo                 string       // G
	ATCUD                     string       // H
	IFiscalSpace              FiscalSpace  // I1-I8
	JFiscalSpace              FiscalSpace  // I1-I8
	KFiscalSpace              FiscalSpace  // I1-I8
	NotTaxable                *dec.Decimal // L
	StampDuty                 *dec.Decimal // M "Imposto de Selo"
	TaxPayable                *dec.Decimal // N
	GrossTotal                *dec.Decimal // O
	WithholdingTaxAmount      *dec.Decimal // P
	HashQuartet               string       // Q
	SoftwareCertificateNumber int64        // R
	OtherInfo                 string       // S
}

type FieldCodec struct {
	Required bool
	Parse    func(f *FatQR, val string) error
	String   func(f *FatQR) string
	Empty    func(f *FatQR) bool
}

var (
	fieldOrder = []string{
		"A", "B", "C", "D", "E", "F", "G", "H",
		"I1", "I2", "I3", "I4", "I5", "I6", "I7", "I8",
		"J1", "J2", "J3", "J4", "J5", "J6", "J7", "J8",
		"K1", "K2", "K3", "K4", "K5", "K6", "K7", "K8",
		"L", "M", "N", "O", "P", "Q", "R", "S",
	}
	// lastKey = fieldOrder[len(fieldOrder)-1]
)

const (
	decCount  = 2
	separator = "*"
)

var fatQRFieldMap = map[string]FieldCodec{
	"A": {
		Required: true,
		Parse: func(f *FatQR, val string) error {
			// if !isValidNIF(val) {
			// 	return fmt.Errorf("invalid NIF `%s`", val)
			// }

			f.TaxRegistrationNumber = val
			return nil
		},
		String: func(f *FatQR) string {
			return "A:" + f.TaxRegistrationNumber
		},
		Empty: func(f *FatQR) bool {
			return f.TaxRegistrationNumber == ""
		},
	},
	"B": {
		Required: true,
		Parse: func(f *FatQR, val string) error {
			// if !isValidNIF(val) {
			// 	return fmt.Errorf("invalid NIF `%s`", val)
			// }

			f.CustomerTaxID = val
			return nil
		},
		String: func(f *FatQR) string {
			return "B:" + f.CustomerTaxID
		},
		Empty: func(f *FatQR) bool {
			return f.CustomerTaxID == ""
		},
	},
	"C": {
		Required: true,
		Parse: func(f *FatQR, val string) error {
			// We can use biter777/countries package around here for verification
			// but this comes from a code, so eeehhh
			f.Country = val
			return nil
		},
		String: func(f *FatQR) string {
			return "C:" + f.Country
		},
		Empty: func(f *FatQR) bool {
			return f.Country == ""
		},
	},
	"D": {
		Required: true,
		Parse: func(f *FatQR, val string) error {
			f.InvoiceType = val
			return nil
		},
		String: func(f *FatQR) string {
			return "D:" + f.InvoiceType
		},
		Empty: func(f *FatQR) bool {
			return f.InvoiceType == ""
		},
	},
	"E": {
		Required: true,
		Parse: func(f *FatQR, val string) error {
			f.InvoiceStatus = val
			return nil
		},
		String: func(f *FatQR) string {
			return "E:" + f.InvoiceStatus
		},
		Empty: func(f *FatQR) bool {
			return f.InvoiceStatus == ""
		},
	},
	"F": {
		Required: true,
		Parse: func(f *FatQR, val string) error {
			tz, err := time.LoadLocation("Europe/Lisbon")
			if err != nil {
				return err
			}

			f.InvoiceDate, err = time.ParseInLocation("20060102", val, tz)
			if err != nil {
				return err
			}

			return nil
		},
		String: func(f *FatQR) string {
			return "F:" + f.InvoiceDate.Format("20060102")
		},
		Empty: func(f *FatQR) bool {
			return f.InvoiceStatus == ""
		},
	},
	"G": {
		Required: true,
		Parse: func(f *FatQR, val string) error {
			f.InvoiceNo = val
			return nil
		},
		String: func(f *FatQR) string {
			return "G:" + f.InvoiceNo
		},
		Empty: func(f *FatQR) bool {
			return f.InvoiceNo == ""
		},
	},
	"H": {
		Required: true,
		Parse: func(f *FatQR, val string) error {
			// TODO
			// Add some ATCUD validation
			f.ATCUD = val
			return nil
		},
		String: func(f *FatQR) string {
			return "H:" + f.ATCUD
		},
		Empty: func(f *FatQR) bool {
			return f.ATCUD == ""
		},
	},
	"I1": {
		Required: true,
		Parse: func(f *FatQR, val string) error {
			f.IFiscalSpace.TaxCountryRegion = val
			return nil
		},
		String: func(f *FatQR) string {
			return "I1:" + f.IFiscalSpace.TaxCountryRegion
		},
		Empty: func(f *FatQR) bool {
			return f.IFiscalSpace.TaxCountryRegion == ""
		},
	},
	"I2": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.IFiscalSpace.TaxableBase, val)
		},
		String: func(f *FatQR) string {
			return "I2:" + stringDecimal(f.IFiscalSpace.TaxableBase)
		},
		Empty: func(f *FatQR) bool {
			return f.IFiscalSpace.TaxableBase == nil
		},
	},
	"I3": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.IFiscalSpace.TaxableReduced, val)
		},
		String: func(f *FatQR) string {
			return "I3:" + stringDecimal(f.IFiscalSpace.TaxableReduced)
		},
		Empty: func(f *FatQR) bool {
			return f.IFiscalSpace.TaxableReduced == nil
		},
	},
	"I4": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.IFiscalSpace.VatTotalReducedTax, val)
		},
		String: func(f *FatQR) string {
			return "I4:" + stringDecimal(f.IFiscalSpace.VatTotalReducedTax)
		},
		Empty: func(f *FatQR) bool {
			return f.IFiscalSpace.VatTotalReducedTax == nil
		},
	},
	"I5": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.IFiscalSpace.TaxableIntermediateBase, val)
		},
		String: func(f *FatQR) string {
			return "I5:" + stringDecimal(f.IFiscalSpace.TaxableIntermediateBase)
		},
		Empty: func(f *FatQR) bool {
			return f.IFiscalSpace.TaxableIntermediateBase == nil
		},
	},
	"I6": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.IFiscalSpace.VatTotalIntermediateTax, val)
		},
		String: func(f *FatQR) string {
			return "I6:" + stringDecimal(f.IFiscalSpace.VatTotalIntermediateTax)
		},
		Empty: func(f *FatQR) bool {
			return f.IFiscalSpace.VatTotalIntermediateTax == nil
		},
	},
	"I7": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.IFiscalSpace.TaxableNormalBase, val)
		},
		String: func(f *FatQR) string {
			return "I7:" + stringDecimal(f.IFiscalSpace.TaxableNormalBase)
		},
		Empty: func(f *FatQR) bool {
			return f.IFiscalSpace.TaxableNormalBase == nil
		},
	},
	"I8": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.IFiscalSpace.VatTotalNormalBase, val)
		},
		String: func(f *FatQR) string {
			return "I8:" + stringDecimal(f.IFiscalSpace.VatTotalNormalBase)
		},
		Empty: func(f *FatQR) bool {
			return f.IFiscalSpace.VatTotalNormalBase == nil
		},
	},
	"J1": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			f.JFiscalSpace.TaxCountryRegion = val
			return nil
		},
		String: func(f *FatQR) string {
			return "J1:" + f.JFiscalSpace.TaxCountryRegion
		},
		Empty: func(f *FatQR) bool {
			return f.JFiscalSpace.TaxCountryRegion == ""
		},
	},
	"J2": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.JFiscalSpace.TaxableBase, val)
		},
		String: func(f *FatQR) string {
			return "J2:" + stringDecimal(f.JFiscalSpace.TaxableBase)
		},
		Empty: func(f *FatQR) bool {
			return f.JFiscalSpace.TaxableBase == nil
		},
	},
	"J3": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.JFiscalSpace.TaxableReduced, val)
		},
		String: func(f *FatQR) string {
			return "J3:" + stringDecimal(f.JFiscalSpace.TaxableReduced)
		},
		Empty: func(f *FatQR) bool {
			return f.JFiscalSpace.TaxableReduced == nil
		},
	},
	"J4": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.JFiscalSpace.VatTotalReducedTax, val)
		},
		String: func(f *FatQR) string {
			return "J4:" + stringDecimal(f.JFiscalSpace.VatTotalReducedTax)
		},
		Empty: func(f *FatQR) bool {
			return f.JFiscalSpace.VatTotalReducedTax == nil
		},
	},
	"J5": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.JFiscalSpace.TaxableIntermediateBase, val)
		},
		String: func(f *FatQR) string {
			return "J5:" + stringDecimal(f.JFiscalSpace.TaxableIntermediateBase)
		},
		Empty: func(f *FatQR) bool {
			return f.JFiscalSpace.TaxableIntermediateBase == nil
		},
	},
	"J6": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.JFiscalSpace.VatTotalIntermediateTax, val)
		},
		String: func(f *FatQR) string {
			return "J6:" + stringDecimal(f.JFiscalSpace.VatTotalIntermediateTax)
		},
		Empty: func(f *FatQR) bool {
			return f.JFiscalSpace.VatTotalIntermediateTax == nil
		},
	},
	"J7": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.JFiscalSpace.TaxableNormalBase, val)
		},
		String: func(f *FatQR) string {
			return "J7:" + stringDecimal(f.JFiscalSpace.TaxableNormalBase)
		},
		Empty: func(f *FatQR) bool {
			return f.JFiscalSpace.TaxableNormalBase == nil
		},
	},
	"J8": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.JFiscalSpace.VatTotalNormalBase, val)
		},
		String: func(f *FatQR) string {
			return "J8:" + stringDecimal(f.JFiscalSpace.VatTotalNormalBase)
		},
		Empty: func(f *FatQR) bool {
			return f.JFiscalSpace.VatTotalNormalBase == nil
		},
	},
	"K1": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			f.KFiscalSpace.TaxCountryRegion = val
			return nil
		},
		String: func(f *FatQR) string {
			return "K1:" + f.KFiscalSpace.TaxCountryRegion
		},
		Empty: func(f *FatQR) bool {
			return f.KFiscalSpace.TaxCountryRegion == ""
		},
	},
	"K2": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.KFiscalSpace.TaxableBase, val)
		},
		String: func(f *FatQR) string {
			return "K2:" + stringDecimal(f.KFiscalSpace.TaxableBase)
		},
		Empty: func(f *FatQR) bool {
			return f.KFiscalSpace.TaxableNormalBase == nil
		},
	},
	"K3": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.KFiscalSpace.TaxableReduced, val)
		},
		String: func(f *FatQR) string {
			return "K3:" + stringDecimal(f.KFiscalSpace.TaxableReduced)
		},
		Empty: func(f *FatQR) bool {
			return f.KFiscalSpace.TaxableReduced == nil
		},
	},
	"K4": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.KFiscalSpace.VatTotalReducedTax, val)
		},
		String: func(f *FatQR) string {
			return "K4:" + stringDecimal(f.KFiscalSpace.VatTotalReducedTax)
		},
		Empty: func(f *FatQR) bool {
			return f.KFiscalSpace.VatTotalReducedTax == nil
		},
	},
	"K5": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.KFiscalSpace.TaxableIntermediateBase, val)
		},
		String: func(f *FatQR) string {
			return "K5:" + stringDecimal(f.KFiscalSpace.TaxableIntermediateBase)
		},
		Empty: func(f *FatQR) bool {
			return f.KFiscalSpace.TaxableIntermediateBase == nil
		},
	},
	"K6": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.KFiscalSpace.VatTotalIntermediateTax, val)
		},
		String: func(f *FatQR) string {
			return "K6:" + stringDecimal(f.KFiscalSpace.VatTotalIntermediateTax)
		},
		Empty: func(f *FatQR) bool {
			return f.KFiscalSpace.VatTotalIntermediateTax == nil
		},
	},
	"K7": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.KFiscalSpace.TaxableNormalBase, val)
		},
		String: func(f *FatQR) string {
			return "K7:" + stringDecimal(f.KFiscalSpace.TaxableNormalBase)
		},
		Empty: func(f *FatQR) bool {
			return f.KFiscalSpace.TaxableNormalBase == nil
		},
	},
	"K8": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.KFiscalSpace.VatTotalNormalBase, val)
		},
		String: func(f *FatQR) string {
			return "K8:" + stringDecimal(f.KFiscalSpace.VatTotalNormalBase)
		},
		Empty: func(f *FatQR) bool {
			return f.KFiscalSpace.VatTotalNormalBase == nil
		},
	},
	"L": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.NotTaxable, val)
		},
		String: func(f *FatQR) string {
			return "L:" + stringDecimal(f.NotTaxable)
		},
		Empty: func(f *FatQR) bool {
			return f.NotTaxable == nil
		},
	},
	"M": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.StampDuty, val)
		},
		String: func(f *FatQR) string {
			return "M:" + stringDecimal(f.StampDuty)
		},
		Empty: func(f *FatQR) bool {
			return f.StampDuty == nil
		},
	},
	"N": {
		Required: true,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.TaxPayable, val)
		},
		String: func(f *FatQR) string {
			return "N:" + stringDecimal(f.TaxPayable)
		},
		Empty: func(f *FatQR) bool {
			return f.TaxPayable == nil
		},
	},
	"O": {
		Required: true,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.GrossTotal, val)
		},
		String: func(f *FatQR) string {
			return "O:" + stringDecimal(f.GrossTotal)
		},
		Empty: func(f *FatQR) bool {
			return f.GrossTotal == nil
		},
	},
	"P": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			return parseDecimal(&f.WithholdingTaxAmount, val)
		},
		String: func(f *FatQR) string {
			return "P:" + stringDecimal(f.WithholdingTaxAmount)
		},
		Empty: func(f *FatQR) bool {
			return f.WithholdingTaxAmount == nil
		},
	},
	"Q": {
		Required: true,
		Parse: func(f *FatQR, val string) error {
			f.HashQuartet = val
			return nil
		},
		String: func(f *FatQR) string {
			return "Q:" + f.HashQuartet
		},
		Empty: func(f *FatQR) bool {
			return f.HashQuartet == ""
		},
	},
	"R": {
		Required: true,
		Parse: func(f *FatQR, val string) error {
			var err error
			f.SoftwareCertificateNumber, err = strconv.ParseInt(val, 10, 16)
			if err != nil {
				return err
			}

			return nil
		},
		String: func(f *FatQR) string {
			return "R:" + fmt.Sprintf("%04d", f.SoftwareCertificateNumber)
		},
		Empty: func(f *FatQR) bool {
			return f.SoftwareCertificateNumber < 1
		},
	},
	"S": {
		Required: false,
		Parse: func(f *FatQR, val string) error {
			f.OtherInfo = val
			return nil
		},
		String: func(f *FatQR) string {
			return "S:" + f.OtherInfo
		},
		Empty: func(f *FatQR) bool {
			return f.OtherInfo == ""
		},
	},
}

func parseDecimal(f **dec.Decimal, val string) error {
	if *f == nil {
		*f = new(dec.Decimal)
	}
	return (*f).Scan(val)
}

func stringDecimal(d *dec.Decimal) string {
	return d.Truncate(decCount).StringFixed(decCount)
}

func (fq *FatQR) scanPart(part []string) error {
	val, ok := fatQRFieldMap[part[0]]
	if !ok {
		return fmt.Errorf("invalid key `%s`", part[0])
	}

	err := val.Parse(fq, part[1])
	if err != nil {
		return err
	}

	return nil
}

func (fq *FatQR) scanParts(parts []string) error {
	for k := range parts {
		keyValuePair := strings.Split(parts[k], ":")
		if len(keyValuePair) != 2 {
			return fmt.Errorf("failed to parse key value pair: %s", parts[k])
		}

		err := fq.scanPart(keyValuePair)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO
// Use flags like STRICT | VALIDATE to do conditional stuff
func (fq *FatQR) Scan(s string, mode ScanMode) error {
	parts := strings.Split(s, "*")

	err := fq.scanParts(parts)
	if err != nil {
		return err
	}

	return nil
}

func (fq *FatQR) String() string {
	parts := make([]string, 0, len(fieldOrder))

	for _, key := range fieldOrder {
		codec, ok := fatQRFieldMap[key]
		if !ok {
			panic("invalid key in order")
		}

		// TODO
		// Using the flags I have to add, use it to enforce or not obligatory fields
		// For now we are only going to stringify the data we have
		if !codec.Empty(fq) {
			parts = append(parts, codec.String(fq))
		}
	}

	return strings.Join(parts, separator)
}
