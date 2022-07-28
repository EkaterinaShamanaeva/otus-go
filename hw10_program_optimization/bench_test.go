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

	defer func(r *zip.ReadCloser) {
		errClose := r.Close()
		if errClose != nil {
			b.Errorf("close zip error: %v", errClose)
		}
	}(r)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		data, errF := r.File[0].Open()
		if errF != nil {
			b.Errorf("zip content open error: %v", errF)
		}

		b.StartTimer()
		_, err = GetDomainStat(data, "biz")
		b.StopTimer()
		if err != nil {
			b.Errorf("get users error: %v", err)
		}
	}
}
