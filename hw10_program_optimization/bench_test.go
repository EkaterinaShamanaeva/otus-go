package hw10programoptimization

import (
	"archive/zip"
	"testing"
)

func BenchmarkStates(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r, _ := zip.OpenReader("testdata/users.dat.zip")
		defer r.Close()

		data, _ := r.File[0].Open()

		GetDomainStat(data, "biz")
	}
}
