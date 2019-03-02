package core

import (
	"fmt"
	"strings"

	"github.com/textileio/go-textile/pb"
)

func (t *Textile) Flags(target string) (*pb.FlagList, error) {
	flags := make([]*pb.Flag, 0)

	query := fmt.Sprintf("type=%d and target='%s'", pb.Block_FLAG, target)
	for _, block := range t.Blocks("", -1, query).Items {
		info, err := t.flag(block, feedItemOpts{annotations: true})
		if err != nil {
			continue
		}
		flags = append(flags, info)
	}

	return &pb.FlagList{Items: flags}, nil
}

func (t *Textile) Flag(blockId string) (*pb.Flag, error) {
	block, err := t.Block(blockId)
	if err != nil {
		return nil, err
	}

	return t.flag(block, feedItemOpts{annotations: true})
}

func (t *Textile) flag(block *pb.Block, opts feedItemOpts) (*pb.Flag, error) {
	if block.Type != pb.Block_FLAG {
		return nil, ErrBlockWrongType
	}

	targetId := strings.TrimPrefix(block.Target, "flag-")
	target, err := t.feedItem(t.datastore.Blocks().Get(targetId), feedItemOpts{})
	if err != nil {
		return nil, err
	}

	return &pb.Flag{
		Block:  block.Id,
		Date:   block.Date,
		User:   t.User(block.Author),
		Target: target,
	}, nil
}
