package main

import (
	"hash"
	"os"
	"os/exec"
	"testing"

	"crypto/rand"

	"crypto/sha256"
	"crypto/sha512"

	"github.com/AudriusButkevicius/gohashcompare/crypto/blake2b"
	"github.com/AudriusButkevicius/gohashcompare/crypto/blake2s"
	"github.com/AudriusButkevicius/gohashcompare/crypto/skein"

	"github.com/AudriusButkevicius/gohashcompare/crypto/blake2bmodified"
	"github.com/AudriusButkevicius/gohashcompare/crypto/blake2smodified"

	"github.com/AudriusButkevicius/gohashcompare/crypto/blake2cgo"
)

var hashvalue = make([]byte, 64)
var blocksize int64 = 1 << 17

func must(hash hash.Hash, err error) hash.Hash {
	if err != nil {
		panic(err)
	}
	return hash
}

func wrap(hasher hash.Hash) func(*testing.B) {
	return func(b *testing.B) {
		buf := make([]byte, blocksize)
		rand.Read(buf)
		b.SetBytes(blocksize)
		b.StopTimer()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			rand.Read(buf)
			b.StartTimer()
			hasher.Write(buf)
			hashvalue = hasher.Sum(nil)
			hasher.Reset()
			b.StopTimer()
		}
	}
}

func main() {
	if len(os.Args) == 1 {
		cmd := exec.Command(os.Args[0], "-test.bench=.*")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		return
	}

	hashers := map[string]hash.Hash{
		"SHA256":               sha256.New(),
		"SHA512":               sha512.New(),
		"Blake2b 256":          blake2b.New256(),
		"Blake2b 512":          blake2b.New512(),
		"Blake2s 256":          blake2s.New256(),
		"Blake2b 256 Modified": blake2bmodified.New256(),
		"Blake2b 512 Modified": blake2bmodified.New512(),
		"Blake2s 256 Modified": blake2smodified.New256(),
		"Skein 256":            must(skein.New(skein.Skein256, 1<<8)),
		"Skein 512":            must(skein.New(skein.Skein512, 1<<8)),
		"Skein 1024":           must(skein.New(skein.Skein1024, 1<<8)),
		"Blake2b CGO 256":      New(&blake2bcgo.Config{Size: 32}),
		"Blake2b CGO 512":      blake2cgo.NewBlake2B(),
	}

	benchmarks := []testing.InternalBenchmark{}
	for name, hash := range hashers {
		benchmarks = append(benchmarks, testing.InternalBenchmark{
			Name: name,
			F:    wrap(hash),
		})
	}
	testing.Main(func(pat, str string) (bool, error) {
		return true, nil
	}, nil, benchmarks, nil)
}
