package gofatqr_test

import (
	"github.com/ESilva15/gofatqr"
	"testing"
)

func TestFatQRAPI(t *testing.T) {
	input := "A:123456789*B:999999990*C:PT*D:FT*E:N*F:20191231*" +
		"G:FT AB2019/0035*H:CSDF7T5H-0035*I1:PT*I2:12000.00*I3:15000.00*" +
		"I4:900.00*I5:50000.00*I6:6500.00*I7:80000.00*I8:18400.00*J1:PT-AC*" +
		"J2:10000.00*J3:25000.56*J4:1000.02*J5:75000.00*J6:6750.00*" +
		"J7:100000.00*J8:18000.00*K1:PT-MA*K2:5000.00*K3:12500.00*K4:625.00*" +
		"K5:25000.00*K6:3000.00*K7:40000.00*K8:8800.00*L:100.00*M:25.00*" +
		"N:64000.02*O:513600.58*P:100.00*Q:kLp0*R:9999*" +
		"S:TB;PT00000000000000000000000;513500.58"

	var qr gofatqr.FatQR
	err := qr.Scan(input, 0)
	if err != nil {
		t.Errorf("Should not have occurred an error for: %s, %+v", input, err)
	}
}
