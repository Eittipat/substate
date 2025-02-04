package db

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/syndtr/goleveldb/leveldb"

	"github.com/0xsoniclabs/substate/substate"
	"github.com/0xsoniclabs/substate/types"
)

var testSubstate = &substate.Substate{
	InputSubstate:  substate.NewWorldState(),
	OutputSubstate: substate.NewWorldState(),
	Env: &substate.Env{
		Coinbase:   types.Address{1},
		Difficulty: new(big.Int).SetUint64(1),
		GasLimit:   1,
		Number:     1,
		Timestamp:  1,
		BaseFee:    new(big.Int).SetUint64(1),
	},
	Message:     substate.NewMessage(1, true, new(big.Int).SetUint64(1), 1, types.Address{1}, new(types.Address), new(big.Int).SetUint64(1), []byte{1}, nil, types.AccessList{}, new(big.Int).SetUint64(1), new(big.Int).SetUint64(1), new(big.Int).SetUint64(1), make([]types.Hash, 0)),
	Result:      substate.NewResult(1, types.Bloom{}, []*types.Log{}, types.Address{}, 1),
	Block:       37_534_834,
	Transaction: 1,
}

func TestSubstateDB_PutSubstate(t *testing.T) {
	dbPath := t.TempDir() + "test-db"
	db, err := createDbAndPutSubstate(dbPath)
	if err != nil {
		t.Fatal(err)
	}

	s := new(leveldb.DBStats)
	err = db.backend.Stats(s)
	if err != nil {
		t.Fatalf("cannot get db stats; %v", err)
	}

	// 54 is the base write when creating levelDB
	if s.IOWrite <= 54 {
		t.Fatal("db file should have something inside")
	}
}

func TestSubstateDB_HasSubstate(t *testing.T) {
	dbPath := t.TempDir() + "test-db"
	db, err := createDbAndPutSubstate(dbPath)
	if err != nil {
		t.Fatal(err)
	}

	has, err := db.HasSubstate(37_534_834, 1)
	if err != nil {
		t.Fatalf("has substate returned error; %v", err)
	}

	if !has {
		t.Fatal("substate is not within db")
	}
}

func TestSubstateDB_GetSubstate(t *testing.T) {
	dbPath := t.TempDir() + "test-db"
	db, err := createDbAndPutSubstate(dbPath)
	if err != nil {
		t.Fatal(err)
	}

	testSubstateDB_GetSubstate(db, t)
}

func testSubstateDB_GetSubstate(db *substateDB, t *testing.T) {
	ss, err := db.GetSubstate(37_534_834, 1)
	if err != nil {
		t.Fatalf("get substate returned error; %v", err)
	}

	if ss == nil {
		t.Fatal("substate is nil")
	}

	if err = ss.Equal(testSubstate); err != nil {
		t.Fatalf("substates are different; %v", err)
	}

}

func TestSubstateDB_DeleteSubstate(t *testing.T) {
	dbPath := t.TempDir() + "test-db"
	db, err := createDbAndPutSubstate(dbPath)
	if err != nil {
		t.Fatal(err)
	}

	err = db.DeleteSubstate(37_534_834, 1)
	if err != nil {
		t.Fatalf("delete substate returned error; %v", err)
	}

	ss, err := db.GetSubstate(37_534_834, 1)
	if err == nil {
		t.Fatal("get substate must fail")
	}

	if got, want := err, leveldb.ErrNotFound; !errors.Is(got, want) {
		t.Fatalf("unexpected err, got: %v, want: %v", got, want)
	}

	if ss != nil {
		t.Fatal("substate was not deleted")
	}
}

func TestSubstateDB_getLastBlock(t *testing.T) {
	dbPath := t.TempDir() + "test-db"
	db, err := createDbAndPutSubstate(dbPath)
	if err != nil {
		t.Fatal(err)
	}

	// add one more substate
	if err = addSubstate(db, testSubstate.Block+1); err != nil {
		t.Fatal(err)
	}

	block, err := db.getLastBlock()
	if err != nil {
		t.Fatal(err)
	}

	if block != 37534835 {
		t.Fatalf("incorrect block number\ngot: %v\nwant: %v", block, testSubstate.Block+1)
	}

}

func TestSubstateDB_GetFirstSubstate(t *testing.T) {
	// save data for comparison
	want := *testSubstate
	want.Block = 1

	dbPath := t.TempDir() + "test-db"
	db, err := createDbAndPutSubstate(dbPath)
	if err != nil {
		t.Fatal(err)
	}

	// add one more substate
	if err = addSubstate(db, 2); err != nil {
		t.Fatal(err)
	}

	got := db.GetFirstSubstate()

	if err = got.Equal(&want); err != nil {
		t.Fatalf("substates are different\nerr: %v\ngot: %s\nwant: %s", err, got, &want)
	}

}

func TestSubstateDB_GetLastSubstate(t *testing.T) {
	// save data for comparison
	want := *testSubstate
	want.Block = 2

	dbPath := t.TempDir() + "test-db"
	db, err := createDbAndPutSubstate(dbPath)
	if err != nil {
		t.Fatal(err)
	}

	// add one more substate
	if err = addSubstate(db, 2); err != nil {
		t.Fatal(err)
	}

	got, err := db.GetLastSubstate()
	if err != nil {
		t.Fatal(err)
	}

	if err = got.Equal(&want); err != nil {
		t.Fatalf("substates are different\nerr: %v\ngot: %s\nwant: %s", err, got, &want)
	}

}

func createDbAndPutSubstate(dbPath string) (*substateDB, error) {
	db, err := newSubstateDB(dbPath, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot open db; %v", err)
	}

	if err = addSubstate(db, testSubstate.Block); err != nil {
		return nil, err
	}

	return db, nil
}

func addSubstate(db *substateDB, blk uint64) error {
	return addCustomSubstate(db, blk, testSubstate)
}

func addCustomSubstate(db *substateDB, blk uint64, ss *substate.Substate) error {
	h1 := types.Hash{}
	h1.SetBytes(nil)

	h2 := types.Hash{}
	h2.SetBytes(nil)

	s := *ss

	s.InputSubstate[types.Address{1}] = substate.NewAccount(1, new(big.Int).SetUint64(1), h1[:])
	s.OutputSubstate[types.Address{2}] = substate.NewAccount(2, new(big.Int).SetUint64(2), h2[:])
	s.Block = blk

	return db.PutSubstate(&s)
}
