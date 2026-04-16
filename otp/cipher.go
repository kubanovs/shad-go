//go:build !solution

package otp

import (
	"errors"
	"fmt"
	"io"
)

type myStreamCipherReader struct {
	target io.Reader
	prng   io.Reader
}

func (obj myStreamCipherReader) Read(p []byte) (n int, err error) {
	nTarget, errTarget := obj.target.Read(p)

	if errTarget != nil && errTarget != io.EOF {
		return nTarget, errTarget
	}

	pPrng := make([]byte, nTarget)
	nPrng, errPrng := obj.prng.Read(pPrng)

	if errPrng != nil {
		return nPrng, errPrng
	}

	for i := 0; i < nTarget; i++ {
		p[i] = p[i] ^ pPrng[i]
	}

	return nTarget, errTarget
}

type myStreamCipherWriter struct {
	target io.Writer
	prng   io.Reader
}

func (obj myStreamCipherWriter) Write(p []byte) (n int, err error) {
	// Создаем буфер для ключевых байтов того же размера
	pPrng := make([]byte, len(p))

	// Читаем ключевые байты из PRNG
	nPrng, errPrng := obj.prng.Read(pPrng)
	if errPrng != nil {
		return 0, errPrng
	}

	// Проверяем, что прочитали достаточно байтов
	if nPrng < len(p) {
		return 0, fmt.Errorf("prng returned only %d bytes, need %d", nPrng, len(p))
	}

	// Шифруем данные: XOR с ключевыми байтами
	encrypted := make([]byte, len(p))
	for i := range p {
		encrypted[i] = p[i] ^ pPrng[i]
	}

	errors.As()
	// Записываем зашифрованные данные в target
	return obj.target.Write(encrypted)
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return myStreamCipherReader{r, prng}
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return myStreamCipherWriter{w, prng}
}
