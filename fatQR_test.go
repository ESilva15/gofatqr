package gofatqr

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	dec "github.com/shopspring/decimal"
)

func TestScan(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{
			name: "5.0. Example 0 - Some invoice I had laying around",
			input: "A:513847588*B:132772985*C:PT*D:FS*E:N*F:20250828*" +
				"G:FS 437251001/115666*H:0*I1:PT*I7:37.94*I8:8.73*" +
				"N:8.73*O:46.67*Q:mWUp*R:0177",
		},
		{
			name: "5.1. Example 1 - Invoice",
			input: "A:123456789*B:999999990*C:PT*D:FT*E:N*F:20191231*" +
				"G:FT AB2019/0035*H:CSDF7T5H-0035*I1:PT*I2:12000.00*I3:15000.00*" +
				"I4:900.00*I5:50000.00*I6:6500.00*I7:80000.00*I8:18400.00*J1:PT-AC*" +
				"J2:10000.00*J3:25000.56*J4:1000.02*J5:75000.00*J6:6750.00*" +
				"J7:100000.00*J8:18000.00*K1:PT-MA*K2:5000.00*K3:12500.00*K4:625.00*" +
				"K5:25000.00*K6:3000.00*K7:40000.00*K8:8800.00*L:100.00*M:25.00*" +
				"N:64000.02*O:513600.58*P:100.00*Q:kLp0*R:9999*" +
				"S:TB;PT00000000000000000000000;513500.58",
		},
		{
			name: "5.2. Example 2 - Simplified Invoice",
			input: "A:123456789*B:999999990*C:PT*D:FS*E:N*F:20190812*" +
				"G:FS CDVF/12345*H:CDF7T5HD-12345*I1:PT*I7:0.65*I8:0.15*N:0.15*O:0.80*" +
				"Q:YhGV*R:9999*S:NU;0.80",
		},
		{
			name: "5.3. Example 3 - Pro-form Invoice",
			input: "A:500000000*B:123456789*C:PT*D:PF*E:N*F:20190123*" +
				"G:PF G2019CB/145789*H:HB6FT7RV-145789*I1:PT*I2:12345.34*I3:12532.65*" +
				"I4:751.96*I5:52789.00*I6:6862.57*I7:32425.69*I8:7457.91*N:15072.44*" +
				"O:125165.12*Q:r/fY*R:9999",
		},
		{
			name: "5.4. Example 4 - Transport Document",
			input: "A:500000000*B:123456789*C:PT*D:GT*E:N*F:20190720*" +
				"G:GT G234CB/50987*H:GTVX4Y8B-50987*I1:0*N:0.00*O:0.00*Q:5uIg*R:9999",
		},
		{
			name: "5.5. Example 5 - Invoice with foreign tax",
			input: "A:123456789*B:4443332215*C:FR*D:FT*E:N*F:20190526*" +
				"G:ABC BNH/4561*H:DK5ZJ2HN-4561*I1:FR*I7:100.00*I8:20.00*N:20.00*" +
				"O:120.00*Q:YJRE*R:9999",
		},
		{
			name: "5.6. Example 6 - Tax Rectification Debit Note",
			input: "A:123456789*B:500000000*C:PT*D:ND*E:N*F:20190216*" +
				"G:M1F KLG/6145*H: RQD8L6DG-6145*I1:PT-MA*I6:26.50*N:26.50*O:26.50*" +
				"Q:h1rB*R:9999",
		},
		{
			name: "5.7. Example 7 - Margin",
			input: "A:500000000*B:123456789*C:PT*D:FT*E:N*F:20191124*" +
				"G:NF 19A/789145*H:JL9DS4TT-789145*I1:PT-AC*I7:50.00*I8:9.00*L:1000.00*" +
				"N:9.00*O:1059.00*Q:d8/K*R:9999",
		},
	}

	for _, tc := range testCases {
		fmt.Printf("Test: %s\n", tc.name)
		var qr FatQR
		err := qr.Scan(tc.input, 0)
		if err != nil {
			t.Errorf("Should not have occurred an error for: %s, %+v", tc.input, err)
		}

		stringified := qr.String()
		if diff := cmp.Diff(tc.input, stringified); diff != "" {
			t.Errorf("String mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestScanPart(t *testing.T) {
	var fq FatQR
	part := []string{"«", "123987654"}

	err := fq.scanPart(part)
	if err == nil {
		t.Errorf("Expected error, got %v", err)
	}
}

func TestScanParts(t *testing.T) {
	var fq FatQR
	input := "A*B:123654987"

	err := fq.Scan(input, 0)
	if err == nil {
		t.Errorf("Expected error about bad parts, got %v", err)
	}
}

func TestFatQRString(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for invalid key, but did not panic")
		}
	}()

	originalOrder := fieldOrder
	fieldOrder = []string{"IM_NOT_EVEN_THERE"}
	defer func() {
		fieldOrder = originalOrder
	}()

	var fq FatQR
	_ = fq.String()
}

func TestScanA(t *testing.T) {
	testStr := "A:500000000*B:123456789*C:PT*D:FT*E:N*F:20191124*" +
		"G:NF 19A/789145*H: JL9DS4TT-789145*I1:PT-AC*I7:50.00*I8:9.00*" +
		"L:1000.00*N:9.00*O:1059.00*Q:d8/K*R:9999"

	var qr FatQR
	err := qr.Scan(testStr, 0)
	if err != nil {
		t.Errorf("Should not have occurred an error for: %s, %+v", testStr, err)
	}

	stringified := qr.String()
	if diff := cmp.Diff(testStr, stringified); diff != "" {
		fmt.Printf("t: %v\n", stringified)
		t.Errorf("String mismatch (-want +got):\n%s", diff)
	}
}

func TestScanPartA(t *testing.T) {
	f := &FatQR{}

	input := "123456789"
	err := fatQRFieldMap["A"].Parse(f, input)
	if err != nil {
		t.Errorf("Expected no err for field A with input `%s`, got %+v", input, err)
	}

	// input = "23456789"
	// err = fatQRFieldMap["A"].Parse(f, input)
	// if err == nil {
	// 	t.Errorf("Expected err for field A with input `%s`, got %+v", input, err)
	// }
}

func TestScanPartB(t *testing.T) {
	f := &FatQR{}

	input := "123456789"
	err := fatQRFieldMap["B"].Parse(f, input)
	if err != nil {
		t.Errorf("Expected no err for field B with input `%s`, got %+v", input, err)
	}

	// input = "23456789"
	// err = fatQRFieldMap["B"].Parse(f, input)
	// if err == nil {
	// 	t.Errorf("Expected err for field B with input `%s`, got %+v", input, err)
	// }
}

func TestScanPartC(t *testing.T) {
	f := &FatQR{}

	input := "PT"
	err := fatQRFieldMap["C"].Parse(f, input)
	if err != nil {
		t.Errorf("Expected no err for field C with input `%s`, got %+v", input, err)
	}
}

func TestScanPartD(t *testing.T) {
	f := &FatQR{}

	input := "something"
	err := fatQRFieldMap["D"].Parse(f, input)
	if err != nil {
		t.Errorf("Expected no err for field D with input `%s`, got %+v", input, err)
	}
}

func TestScanPartE(t *testing.T) {
	f := &FatQR{}

	input := "anothersomething"
	err := fatQRFieldMap["E"].Parse(f, input)
	if err != nil {
		t.Errorf("Expected no err for field E with input `%s`, got %+v", input, err)
	}
}

func TestScanPartF(t *testing.T) {
	f := &FatQR{}

	loc, err := time.LoadLocation("Europe/Lisbon")
	if err != nil {
		t.Errorf("Failed to prepare timezone var: %+v", err)
	}
	expected := time.Date(2025, 12, 31, 0, 0, 0, 0, loc)

	input := "20251231"
	err = fatQRFieldMap["F"].Parse(f, input)
	if err != nil {
		t.Errorf("Expected no err for field F with input `%s`, got %+v", input, err)
	}

	if !f.InvoiceDate.Equal(expected) {
		t.Errorf("Result date doesnt match expected date: %v\nGot: %+v, Want: %+v",
			f.InvoiceDate, expected, err)
	}
}

func TestScanPartG(t *testing.T) {
	f := &FatQR{}

	input := "FS 81239/99999"
	err := fatQRFieldMap["G"].Parse(f, input)
	if err != nil {
		t.Errorf("Expected no err for field G with input `%s`, got %+v", input, err)
	}

	if f.InvoiceNo != input {
		t.Errorf("Expected `%s`, got %+v", input, f.InvoiceNo)
	}
}

func TestScanPartH(t *testing.T) {
	f := &FatQR{}

	input := "CSDF7T5H-0035"
	err := fatQRFieldMap["H"].Parse(f, input)
	if err != nil {
		t.Errorf("Expected no err for field H with input `%s`, got %+v", input, err)
	}

	if f.ATCUD != input {
		t.Errorf("Expected `%s`, got %+v", input, f.ATCUD)
	}
}

func TestScanPartI1(t *testing.T) {
	f := &FatQR{}

	input := "PT"
	err := fatQRFieldMap["I1"].Parse(f, input)
	if err != nil {
		t.Errorf("Expected no err for field I1 with input `%s`, got %+v", input, err)
	}

	if f.IFiscalSpace.TaxCountryRegion != input {
		t.Errorf("Expected `%s`, got %+v", input, f.IFiscalSpace.TaxCountryRegion)
	}
}

func TestScanPartI2(t *testing.T) {
	f := &FatQR{}

	input := "1565.20"
	expect := dec.NewFromFloat(1565.20)
	err := fatQRFieldMap["I2"].Parse(f, input)
	if err != nil {
		t.Errorf("Expected no err for field I2 with input `%s`, got %+v", input, err)
	}

	if f.IFiscalSpace.TaxableBase.Truncate(2).String() != expect.Truncate(2).String() {
		t.Errorf("Expected `%s`, got %+v", expect, f.IFiscalSpace.TaxableBase)
	}

	input = "10.5A"
	err = fatQRFieldMap["I2"].Parse(f, input)
	if err == nil {
		t.Errorf("Expected err for field I2 with input `%s`, got %+v", input, err)
	}
}

func TestMarshalling(t *testing.T) {
	input := "A:123456789*B:999999990*C:PT*D:FT*E:N*F:20191231*" +
		"G:FT AB2019/0035*H:CSDF7T5H-0035*I1:PT*I2:12000.00*I3:15000.00*" +
		"I4:900.00*I5:50000.00*I6:6500.00*I7:80000.00*I8:18400.00*J1:PT-AC*" +
		"J2:10000.00*J3:25000.56*J4:1000.02*J5:75000.00*J6:6750.00*" +
		"J7:100000.00*J8:18000.00*K1:PT-MA*K2:5000.00*K3:12500.00*K4:625.00*" +
		"K5:25000.00*K6:3000.00*K7:40000.00*K8:8800.00*L:100.00*M:25.00*" +
		"N:64000.02*O:513600.58*P:100.00*Q:kLp0*R:9999*" +
		"S:TB;PT00000000000000000000000;513500.58"

	expected := "{\"TaxRegistrationNumber\":\"123456789\",\"CustomerTaxID\":\"999999990\"," +
		"\"Country\":\"PT\",\"InvoiceType\":\"FT\",\"InvoiceStatus\":\"N\"," +
		"\"InvoiceDate\":\"2019-12-31T00:00:00Z\",\"InvoiceNo\":\"FT AB2019/0035\"," +
		"\"ATCUD\":\"CSDF7T5H-0035\",\"IFiscalSpace\":{\"TaxCountryRegion\":\"PT\"," +
		"\"TaxableBase\":\"12000\",\"TaxableReduced\":\"15000\",\"VatTotalReducedTax\":\"900\"," +
		"\"TaxableIntermediateBase\":\"50000\",\"VatTotalIntermediateTax\":\"6500\"," +
		"\"TaxableNormalBase\":\"80000\",\"VatTotalNormalBase\":\"18400\"}," +
		"\"JFiscalSpace\":{\"TaxCountryRegion\":\"PT-AC\",\"TaxableBase\":\"10000\"," +
		"\"TaxableReduced\":\"25000.56\",\"VatTotalReducedTax\":\"1000.02\"," +
		"\"TaxableIntermediateBase\":\"75000\",\"VatTotalIntermediateTax\":\"6750\"," +
		"\"TaxableNormalBase\":\"100000\",\"VatTotalNormalBase\":\"18000\"}," +
		"\"KFiscalSpace\":{\"TaxCountryRegion\":\"PT-MA\",\"TaxableBase\":\"5000\"," +
		"\"TaxableReduced\":\"12500\",\"VatTotalReducedTax\":\"625\"," +
		"\"TaxableIntermediateBase\":\"25000\",\"VatTotalIntermediateTax\":\"3000\"," +
		"\"TaxableNormalBase\":\"40000\",\"VatTotalNormalBase\":\"8800\"},\"NotTaxable\":\"100\"," +
		"\"StampDuty\":\"25\",\"TaxPayable\":\"64000.02\",\"GrossTotal\":\"513600.58\"," +
		"\"WithholdingTaxAmount\":\"100\",\"HashQuartet\":\"kLp0\",\"SWCertNo\":9999," +
		"\"OtherInfo\":\"TB;PT00000000000000000000000;513500.58\"}"

	var f FatQR
	err := f.Scan(input, 0)
	if err != nil {
		t.Errorf("Scan should've been successful: %+v", err)
	}

	jsonData := f.ToJSON()
	if diff := cmp.Diff(string(jsonData), expected); diff != "" {
		t.Errorf("String mismatch (-want +got):\n%s", diff)
	}

	var newF FatQR
	err = newF.FromJSON(jsonData)
	if err != nil {
		t.Errorf("FromJSON should've been successful: %+v", err)
	}

	newJsonData := newF.ToJSON()
	if diff := cmp.Diff(string(newJsonData), expected); diff != "" {
		t.Errorf("String mismatch (-want +got):\n%s", diff)
	}
}

// TODO: Add more UTs
