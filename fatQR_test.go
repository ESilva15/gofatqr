package gofatqr

import (
	"testing"
)

func TestScan(t *testing.T) {
	testStr := "A:513847588*B:239752988*C:PT*D:FS*E:N*F:20250828*G:FS 437251001/115666*H:0*I1:PT*I7:37.94*I8:8.73*N:8.73*O:46.67*Q:mWUp*R:0177"

	var qr FatQR
	err := qr.Scan(testStr)
	if err != nil {
		t.Errorf("Should not have occurred an error for: %s, %+v", testStr, err)
	}
}
