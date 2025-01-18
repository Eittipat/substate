package main

import (
	"github.com/0xsoniclabs/substate/db"
	"github.com/0xsoniclabs/substate/substate"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

type RLPtoProtobufCommand struct {
	src db.SubstateDB
	dst db.SubstateDB
	ctx *cli.Context
}

func (c *RLPtoProtobufCommand) Execute() error {
	segment, err := ParseBlockSegment(c.ctx.String(BlockSegmentFlag.Name))
	if err != nil {
		return err
	}

	taskPool := &db.SubstateTaskPool{
		Name:     "rlp-to-protobuf",
		TaskFunc: c.performSubstateUpgrade,

		First: segment.First,
		Last:  segment.Last,

		Workers:         c.ctx.Int(WorkersFlag.Name),
		SkipTransferTxs: c.ctx.Bool(SkipTransferTxsFlag.Name),
		SkipCallTxs:     c.ctx.Bool(SkipCallTxsFlag.Name),
		SkipCreateTxs:   c.ctx.Bool(SkipCreateTxsFlag.Name),

		Ctx: c.ctx,

		DB: c.src,
	}
	_, err = c.dst.SetSubstateEncoding("protobuf")
	if err != nil {
		return err
	}
	return taskPool.Execute()
}

func (c *RLPtoProtobufCommand) performSubstateUpgrade(block uint64, tx int, substate *substate.Substate, taskPool *db.SubstateTaskPool) error {
	err := c.dst.PutSubstate(substate)
	if err != nil {
		log.Fatalf("Failed to put substate: %v", err)
		return err
	}
	return nil
}

func action(ctx *cli.Context) error {
	// Open old DB
	src, err := db.NewCustomSubstateDB(ctx.String(SrcDbFlag.Name), 1024, 100, true)
	if err != nil {
		return err
	}
	defer src.Close()

	// Remove destination DB if exists
	err = RemoveIfExist(ctx.String(DstDbFlag.Name))
	if err != nil {
		return err
	}
	// Open new DB
	dst, err := db.NewCustomSubstateDB(ctx.String(DstDbFlag.Name), 1024, 100, false)
	if err != nil {
		return err
	}
	defer dst.Close()

	command := RLPtoProtobufCommand{
		src: src,
		dst: dst,
		ctx: ctx,
	}
	return command.Execute()
}

func main() {
	app := &cli.App{
		Name:   "rlp-to-protobuf",
		Usage:  "Convert RLP encoded substate to protobuf encoded substate",
		Action: action,
		Flags: []cli.Flag{
			&WorkersFlag,
			&SrcDbFlag,
			&DstDbFlag,
			&BlockSegmentFlag,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
