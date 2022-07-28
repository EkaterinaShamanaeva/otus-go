package hw10programoptimization

import (
	"archive/zip"
	"testing"
)

func BenchmarkStates(b *testing.B) {
	r, err := zip.OpenReader("testdata/users.dat.zip")
	if err != nil {
		b.Errorf("open zip error: %v", err)
	}

	for i := 0; i < b.N; i++ {
		data, errF := r.File[0].Open()
		if errF != nil {
			b.Errorf("zip content open error: %v", errF)
		}
		_, err = GetDomainStat(data, "biz")
		if err != nil {
			b.Errorf("get users error: %v", err)
		}
	}

	err = r.Close()
	if err != nil {
		b.Errorf("close zip error: %v", err)
	}
}
