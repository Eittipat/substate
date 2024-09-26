package db

import (
	"strings"
	"testing"
	"math/big"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
	pb "github.com/Fantom-foundation/Substate/protobuf"
	"github.com/Fantom-foundation/Substate/rlp"
	trlp "github.com/Fantom-foundation/Substate/types/rlp"
)

type encTest struct {
	bytes []byte
	blk uint64
	tx int
}

var (
	o uint64 = 1
	one *uint64 = &o
	bigOne *big.Int = new(big.Int).SetUint64(1)
	typeOne = pb.Substate_TxMessage_TXTYPE_LEGACY
)

var testPbSubstate = &pb.Substate{
	InputAlloc: &pb.Substate_Alloc{},
	OutputAlloc: &pb.Substate_Alloc{},
	BlockEnv: &pb.Substate_BlockEnv{
		Coinbase:   []byte{1},
		Difficulty: []byte{1},
		GasLimit:   one,
		Number:     one,
		Timestamp:  one,
		BaseFee: wrapperspb.Bytes([]byte{1}),
	},
	TxMessage: &pb.Substate_TxMessage{
		Nonce: one,
		GasPrice: []byte{1},
		Gas: one,
		From: []byte{1},
		To: nil,
		Value: []byte{1},
		Input: &pb.Substate_TxMessage_InitCodeHash{
			InitCodeHash: []byte{1},
		},
		TxType: &typeOne,
		AccessList: []*pb.Substate_TxMessage_AccessListEntry{},
		GasFeeCap: wrapperspb.Bytes([]byte{1}),
		GasTipCap: wrapperspb.Bytes([]byte{1}),
		BlobGasFeeCap: wrapperspb.Bytes([]byte{1}),
		BlobHashes: [][]byte{},
	},
	Result: &pb.Substate_Result{
		Status: one,
		Bloom: []byte{1},
		Logs: []*pb.Substate_Result_Log{},
		GasUsed: one,
	},
}

var (
	simplePb, _ = proto.Marshal(testPbSubstate)
	testPb = encTest {
		bytes: simplePb,
		blk: testSubstate.Block,
		tx: testSubstate.Transaction,
	}

	simpleRlp, _ = trlp.EncodeToBytes(rlp.NewRLP(testSubstate))
	testRlp = encTest{
		bytes: simpleRlp,
		blk: testSubstate.Block,
		tx: testSubstate.Transaction,
	}

	supportedEncoding = map[string]encTest{
		"rlp": testRlp,
		"protobuf": testPb,
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
	_, err = db.decodeToSubstate(testRlp.bytes, testRlp.blk, testRlp.tx)
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

	_, err = db.decodeToSubstate(testRlp.bytes, testRlp.blk, testRlp.tx)
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

func TestSubstateEncoding_EncodePb(t *testing.T) {
	ss, err := proto.Marshal(testPbSubstate)
	if err != nil {
		t.Fatal(err)
	}

	pbSubstate := &pb.Substate{}
	if err := proto.Unmarshal(ss, pbSubstate); err != nil {
		t.Fatal(err)
	}
}


func TestSubstateEncoding_TestDb(t *testing.T) {
	for encoding, et := range supportedEncoding {
		path := t.TempDir() + "test-db-" + encoding
		db, err := newSubstateDB(path, nil, nil, nil)
		if err != nil {
			t.Errorf("cannot open db; %v", err)
		}

		db, err = db.SetSubstateEncoding(encoding)
		if err != nil {
			t.Error(err)
		}

		ss, err := db.decodeToSubstate(et.bytes, et.blk, et.tx)
		if err != nil {
			t.Error(err)
		}

		err = addCustomSubstate(db, et.blk, ss)
		if err != nil {
			t.Error(err)
		}
		
		testSubstateDB_GetSubstate(db, t)
	}
}

func TestSubstateEncoding_TestIterator(t *testing.T) {
	for encoding, et := range supportedEncoding {
		path := t.TempDir() + "test-db-" + encoding
		db, err := newSubstateDB(path, nil, nil, nil)
		if err != nil {
			t.Errorf("cannot open db; %v", err)
		}

		fmt.Println(encoding)
		_, err = db.SetSubstateEncoding(encoding)
		if err != nil {
			t.Error(err)
		}

		ss, err := db.decodeToSubstate(et.bytes, et.blk, et.tx)
		if err != nil {
			t.Error(err)
		}

		err = addCustomSubstate(db, et.blk, ss)
		if err != nil {
			t.Error(err)
		}

		testSubstatorIterator_Value(db, t)
	}
}
