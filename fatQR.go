package gofatqr

// TODO
// Check what we can take form SAF-T to put here, maybe work on a SAF-T
// implementation too, who knows amirite

import (
	"fmt"
	"math/big"
	"strings"
	"time"
)

type FiscalSpace struct {
	TaxCountryRegion        string   // 1 - fiscal space
	TaxableBase             *big.Rat // 2
	TaxabledReduced         *big.Rat // 3
	VatTotalReducedTax      *big.Rat // 4
	TaxableIntermediateBase *big.Rat // 5
	VatTotalIntermediateTax *big.Rat // 6
	TaxableNormalBase       *big.Rat // 7
	VatTotalNormalBase      *big.Rat // 8
}

type FatQR struct {
	TaxRegistrationNumber     string      // A
	CustomerTaxID             string      // B
	Country                   string      // C - customer country
	InvoiceType               string      // D
	InvoiceStatus             bool        // E
	InvoiceDate               time.Time   // F
	InvoiceNo                 string      // G
	ATCUD                     string      // H
	IFiscalSpace              FiscalSpace // I1 - I8
	JFiscalSpace              FiscalSpace // I1 - I8
	KFiscalSpace              FiscalSpace // I1 - I8
	NotTaxable                *big.Rat    // L
	StampDuty                 *big.Rat    // M - "Imposto de Selo"
	TaxPayable                *big.Rat    // N
	GrossTotal                *big.Rat    // O
	WithholdingTaxAmount      *big.Rat    // P
	HashQuartet               string      // Q
	SoftwareCertificateNumber int         // R
	OtherInfo                 string      // S
}

type part struct {
	key   string
	value string
}

func (fq *FatQR) scanParts(parts []string) error {
	fmt.Println(parts)

	for k := range parts {
		keyValuePair := strings.Split(parts[k], ":")
		if len(keyValuePair) != 2 {
			return fmt.Errorf("failed to parse key value pair: %s", parts[k])
		}
	}

	return nil
}

func (fq *FatQR) Scan(s string) error {
	parts := strings.Split(s, "*")

	fq.scanParts(parts)

	return nil
}

func (fq *FatQR) String() string {
	return ""
}
