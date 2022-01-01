package main

import (
	"bytes"
	"flag"
	"io"
	"math/rand"
	"os"
	"sort"
)

var (
	rotkOffset int
)

func init() {
	flag.IntVar(&rotkOffset, "k", 13, "rotK (if minus => all, if 0~25 => only k)")
}

type RotKTranslater struct {
	k int
	w *bytes.Buffer
}

func NewRotKTranslater(k int) *RotKTranslater {
	return &RotKTranslater{
		k: k,
		w: &bytes.Buffer{},
	}
}

func rot13(r rune, _k int) rune {
	k := rune(_k)
	switch {
	case 'a' <= r && r <= 'z':
		return 'a' + ((r-'a')+k)%26
	case 'A' <= r && r <= 'Z':
		return 'A' + ((r-'A')+k)%26
	default:
		return r
	}
}

func (t *RotKTranslater) Write(p []byte) (n int, err error) {
	return t.w.Write(p)
}

func (t *RotKTranslater) Read(p []byte) (n int, err error) {
	n, err = t.w.Read(p)
	if err != nil {
		return
	}
	runes := []rune(string(p))
	for i, v := range runes {
		runes[i] = rot13(v, t.k)
	}
	copy(p, []byte(string(runes)))
	return
}

func main() {
	flag.Parse()

	vals := []int{rotkOffset}
	if rotkOffset < 0 {
		vals = rand.Perm(26)
		sort.Ints(vals)
	}

	var r io.Reader
	r = os.Stdin
	for _, v := range vals {
		w := &bytes.Buffer{}
		rot := NewRotKTranslater(v)
		io.Copy(io.MultiWriter(w, rot), r)
		io.Copy(os.Stdout, rot)
		r = w
	}
}
