package db

import (
	"strings"
	"testing"

	"github.com/Fantom-foundation/Substate/rlp"
	trlp "github.com/Fantom-foundation/Substate/types/rlp"
)

var (
	testRlp, _ = trlp.EncodeToBytes(rlp.NewRLP(testSubstate))
	testBlk    = testSubstate.Block
	testTx     = testSubstate.Transaction

	supportedEncoding = map[string][]byte{
		"rlp": testRlp,
	}
)

func TestSubstateEncoding_NilEncodingDefaultsToRlp(t *testing.T) {
	path := t.TempDir() + "test-db"
	db, err := newSubstateDB(path, nil, nil, nil)
	if err != nil {
		t.Errorf("cannot open db; %v", err)
	}

	if got := db.GetSubstateEncoding(); got != "" {
		t.Fatalf("substate encoding should be nil, got: %s", got)
	}

	// purposely never set encoding
	_, err = db.decodeToSubstate(testRlp, testBlk, testTx)
	if err != nil {
		t.Fatal(err)
	}

	if got := db.GetSubstateEncoding(); got != "rlp" {
		t.Fatalf("db should default to rlp, got: %s", got)
	}
}

func TestSubstateEncoding_DefaultEncodingDefaultsToRlp(t *testing.T) {
	path := t.TempDir() + "test-db"
	db, err := newSubstateDB(path, nil, nil, nil)
	if err != nil {
		t.Errorf("cannot open db; %v", err)
	}

	_, err = db.SetSubstateEncoding("default")
	if err != nil {
		t.Fatal("default is supportet, but error")
	}

	_, err = db.decodeToSubstate(testRlp, testBlk, testTx)
	if err != nil {
		t.Fatal(err)
	}

	if got := db.GetSubstateEncoding(); got != "rlp" {
		t.Fatalf("db should default to rlp, got: %s", got)
	}
}

func TestSubstateEncoding_UnsupportedEncodingThrowsError(t *testing.T) {
	path := t.TempDir() + "test-db"
	db, err := newSubstateDB(path, nil, nil, nil)
	if err != nil {
		t.Errorf("cannot open db; %v", err)
	}

	_, err = db.SetSubstateEncoding("EncodingNotSupported")
	if err == nil || !strings.Contains(err.Error(), "Encoding not supported") {
		t.Error("Encoding not supported, but no error")
	}
}

func TestSubstateEncoding_TestDb(t *testing.T) {
	path := t.TempDir() + "test-db"
	db, err := newSubstateDB(path, nil, nil, nil)
	if err != nil {
		t.Errorf("cannot open db; %v", err)
	}

	for encoding, bytes := range supportedEncoding {
		_, err = db.SetSubstateEncoding(encoding)
		if err != nil {
			t.Error(err)
		}

		ss, err := db.decodeToSubstate(bytes, testBlk, testTx)
		if err != nil {
			t.Error(err)
		}

		err = addCustomSubstate(db, testBlk, ss)
		if err != nil {
			t.Error(err)
		}

		testSubstateDB_GetSubstate(db, t)
	}
}

func TestSubstateEncoding_TestIterator(t *testing.T) {
	path := t.TempDir() + "test-db"
	db, err := newSubstateDB(path, nil, nil, nil)
	if err != nil {
		t.Errorf("cannot open db; %v", err)
	}

	for encoding, bytes := range supportedEncoding {
		_, err = db.SetSubstateEncoding(encoding)
		if err != nil {
			t.Error(err)
		}

		ss, err := db.decodeToSubstate(bytes, testBlk, testTx)
		if err != nil {
			t.Error(err)
		}

		err = addCustomSubstate(db, testBlk, ss)
		if err != nil {
			t.Error(err)
		}

		testSubstatorIterator_Value(db, t)
	}
}
