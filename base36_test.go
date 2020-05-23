package base36

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"testing"
)

var testEncoders = []func([]byte) string{
	EncodeToStringLc,
	EncodeToStringUc,
}

func testBasicTrip(t *testing.T, src, enc, result string) {

	if enc != result {
		t.Fatalf("Encoding failed: expected '%s', got '%s'", result, enc)
	}

	trip, err := DecodeString(enc)
	if err != nil {
		t.Fatalf("Unexpected error when decoding '%s': %s", enc, err)
	}

	if !bytes.Equal(trip, []byte(src)) {
		t.Fatalf("Encode-Decode roundtrip failed: expected '%s', got '%s'", src, trip)
	}
}

func TestBasicTripLc(t *testing.T) {
	testBasicTrip(t,
		"Decentralize everything!!!",
		EncodeToStringLc([]byte("Decentralize everything!!!")),
		"m552ng4dabi4neu1oo8l4i5mndwmpc3mkukwtxy9",
	)
}

func TestBasicTripUc(t *testing.T) {
	testBasicTrip(t,
		"Decentralize everything!!!",
		EncodeToStringUc([]byte("Decentralize everything!!!")),
		"M552NG4DABI4NEU1OO8L4I5MNDWMPC3MKUKWTXY9",
	)
}

func TestPermute(t *testing.T) {

	buf := make([]byte, 137+16) // sufficiently large prime number of bytes + another 16 to test leading 0s
	rand.Read(buf[16:])

	for _, encoder := range testEncoders {

		// test roundtrip from the full zero-prefixed buffer down to a single byte
		for i := 0; i < len(buf); i++ {

			// use a copy to verify we are not overwriting the supplied buffer
			newBuf := make([]byte, len(buf)-i)
			copy(newBuf, buf[i:])

			enc := encoder(newBuf)

			out, err := DecodeString(enc)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(newBuf, buf[i:]) {
				t.Fatal("the provided buffer was modified", buf[i:], out)
			}

			if !bytes.Equal(buf[i:], out) {
				t.Fatal("input wasn't the same as output", buf[i:], out)
			}
		}
	}
}

var benchmarkBuf [36]byte // typical CID size
var benchmarkCodecs []string

func init() {
	rand.Read(benchmarkBuf[:])
}

func BenchmarkRoundTrip(b *testing.B) {
	b.ResetTimer()

	for i, encoder := range testEncoders {

		b.Run(fmt.Sprintf("encoder #%d", i+1), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				enc := encoder(benchmarkBuf[:])

				out, err := DecodeString(enc)
				if err != nil {
					b.Fatal(err)
				}

				if !bytes.Equal(benchmarkBuf[:], out) {
					b.Fatal("input wasnt the same as output", benchmarkBuf, out)
				}
			}
		})
	}
}

func BenchmarkEncode(b *testing.B) {
	b.ResetTimer()

	for i, encoder := range testEncoders {
		b.Run(fmt.Sprintf("encoder #%d", i+1), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				encoder(benchmarkBuf[:])
			}
		})
	}
}

func BenchmarkDecode(b *testing.B) {
	b.ResetTimer()

	enc := testEncoders[0](benchmarkBuf[:])
	b.Run("Decoding", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			DecodeString(enc)
		}
	})
}
