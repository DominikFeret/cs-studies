package main

import (
	"testing"
)

func TestCryptoanalysisStd(t *testing.T) {
	engDict := getEnglishDict()
	encStr := []rune("envnbkkyakazgn")

	decStr, key := cryptoanalysis(engDict, encStr, decryptCharStd)

	if string(decStr) != "inertiacatered" {
		t.Errorf("Expected 'inertiacatered', got '%s'", string(decStr))
	}

	if string(key) != "warwick" {
		t.Errorf("Expected 'warwick', got '%s'", string(key))
	}
}

func TestCryptoanalysisXor(t *testing.T) {
	engDict := getEnglishDict()
	encStr := []rune("envnbkkyakazgn")

	decStr, key := cryptoanalysis(engDict, encStr, decryptCharXor)

	if string(decStr) != "inertiacatered" {
		t.Errorf("Expected 'inertiacatered', got '%s'", string(decStr))
	}

	if string(key) != "warwick" {
		t.Errorf("Expected 'warwick', got '%s'", string(key))
	}

}
