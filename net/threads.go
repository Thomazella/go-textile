package net

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/segmentio/ksuid"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/ipfs"
	"github.com/textileio/textile-go/net/service"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/thread"
	"gx/ipfs/QmVzK524a2VWLqyvtBeiHKsUAWYgeAk4DBeZoY7vpNPNRx/go-block-format"
	"gx/ipfs/QmYVNvtQkeZ6AKSwDrjQTs432QtL6umrrK41EBq3cu7iSP/go-cid"
	"gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
	libp2pc "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	"gx/ipfs/QmebqVUQQqQFhg74FtQFszUJo22Vpr3e8qBAkvvV4ho9HH/go-ipfs/core"
)

// ThreadService is a libp2p service for orchestrating a collection of files
// with annotations amongst a group of peers
type ThreadsService struct {
	service          *service.Service
	getThread        func(id string) (*int, *thread.Thread)
	sendNotification func(note *repo.Notification) error
}

// NewThreadsService returns a new threads service
func NewThreadsService(
	node *core.IpfsNode,
	datastore repo.Datastore,
	getThread func(id string) (*int, *thread.Thread),
	sendNotification func(note *repo.Notification) error,
) *ThreadsService {
	handler := &ThreadsService{
		getThread:        getThread,
		sendNotification: sendNotification,
	}
	handler.service = service.NewService(handler, node, datastore)
	return handler
}

// Protocol returns the handler protocol
func (h *ThreadsService) Protocol() protocol.ID {
	return protocol.ID("/textile/threads/1.0.0")
}

// Node returns the underlying ipfs Node
func (h *ThreadsService) Node() *core.IpfsNode {
	return h.service.Node
}

// Datastore returns the underlying datastore
func (h *ThreadsService) Datastore() repo.Datastore {
	return h.service.Datastore
}

// Ping pings another peer
func (h *ThreadsService) Ping(pid peer.ID) (service.PeerStatus, error) {
	return h.service.Ping(pid)
}

// VerifyEnvelope calls service verify
func (h *ThreadsService) VerifyEnvelope(env *pb.Envelope) error {
	return h.service.VerifyEnvelope(env)
}

// NewBlock returns a thread-signed block in an envelope
func (h *ThreadsService) NewBlock(sk libp2pc.PrivKey, mtype pb.Message_Type, msg proto.Message) (*pb.Envelope, error) {
	ser, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	threadSig, err := sk.Sign(ser)
	if err != nil {
		return nil, err
	}
	signed := &pb.SignedThreadBlock{
		Block:     ser,
		ThreadSig: threadSig,
	}
	return h.service.NewEnvelope(mtype, signed, nil, false)
}

// SendMessage sends a message to a peer
func (h *ThreadsService) SendMessage(env *pb.Envelope, pid peer.ID) error {
	ctx, cancel := context.WithTimeout(context.Background(), service.DefaultTimeout)
	defer cancel()
	return h.service.SendMessage(ctx, pid, env)
}

// Handle is called by the underlying service handler method
func (h *ThreadsService) Handle(mtype pb.Message_Type) func(peer.ID, *pb.Envelope) (*pb.Envelope, error) {
	switch mtype {
	case pb.Message_THREAD_INVITE:
		return h.handleInvite
	case pb.Message_THREAD_JOIN:
		return h.handleJoin
	case pb.Message_THREAD_LEAVE:
		return h.handleLeave
	case pb.Message_THREAD_DATA:
		return h.handleData
	case pb.Message_THREAD_ANNOTATION:
		return h.handleAnnotation
	case pb.Message_THREAD_IGNORE:
		return h.handleIgnore
	case pb.Message_THREAD_MERGE:
		return h.handleMerge

	// TODO: Move these to pin service
	case pb.Message_BLOCK:
		return h.handleBlock
	case pb.Message_STORE:
		return h.handleStore
	default:
		return nil
	}
}

func (h *ThreadsService) handleInvite(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	signed, err := unpackThreadMessage(env)
	if err != nil {
		return nil, err
	}
	invite := new(pb.ThreadInvite)
	if err := proto.Unmarshal(signed.Block, invite); err != nil {
		return nil, err
	}

	// load thread
	threadId := libp2pc.ConfigEncodeKey(invite.Header.ThreadPk)
	if _, thrd := h.getThread(threadId); thrd != nil {
		// thread exists, aborting
		return nil, nil
	}

	// check if it'h meant for us (should be, but safety first)
	if invite.InviteeId != h.Node().Identity.Pretty() {
		return nil, errors.New("invalid invite block")
	}

	// unknown thread and invite meant for us
	// unpack new thread secret that should be encrypted with our key
	skb, err := crypto.Decrypt(h.Node().PrivateKey, invite.SkCipher)
	if err != nil {
		return nil, err
	}
	sk, err := libp2pc.UnmarshalPrivateKey(skb)
	if err != nil {
		return nil, err
	}

	// verify thread sig
	if err := crypto.Verify(sk.GetPublic(), signed.Block, signed.ThreadSig); err != nil {
		return nil, err
	}

	// add to local ipfs for later use when joining
	envb, err := proto.Marshal(env)
	if err != nil {
		return nil, err
	}
	ci, err := ipfs.PinData(h.Node(), bytes.NewReader(envb))
	if err != nil {
		return nil, err
	}
	id := ci.Hash().B58String()

	// send notification
	notification, err := newThreadNotification(sk, invite.Header, repo.ReceivedInviteNotification)
	if err != nil {
		return nil, err
	}
	notification.Subject = invite.SuggestedName
	notification.SubjectId = libp2pc.ConfigEncodeKey(invite.Header.ThreadPk)
	notification.BlockId = id // invite block
	notification.Body = "invited you to join"
	if err := h.sendNotification(notification); err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *ThreadsService) handleJoin(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	signed, err := unpackThreadMessage(env)
	if err != nil {
		return nil, err
	}
	join := new(pb.ThreadJoin)
	if err := proto.Unmarshal(signed.Block, join); err != nil {
		return nil, err
	}

	// load thread
	threadId := libp2pc.ConfigEncodeKey(join.Header.ThreadPk)
	_, thrd := h.getThread(threadId)
	if thrd == nil {
		return nil, errors.New("invalid join block")
	}

	// verify thread sig
	if err := thrd.Verify(signed); err != nil {
		return nil, err
	}

	// handle
	addr, _, err := thrd.HandleJoinBlock(&pid, env, signed, join, false)
	if err != nil {
		return nil, err
	}

	// send notification
	notification, err := newThreadNotification(thrd.PrivKey, join.Header, repo.PeerJoinedNotification)
	if err != nil {
		return nil, err
	}
	notification.Subject = thrd.Name
	notification.SubjectId = thrd.Id
	notification.BlockId = addr.B58String()
	notification.Body = "joined"
	if err := h.sendNotification(notification); err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *ThreadsService) handleLeave(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	signed, err := unpackThreadMessage(env)
	if err != nil {
		return nil, err
	}
	leave := new(pb.ThreadLeave)
	if err := proto.Unmarshal(signed.Block, leave); err != nil {
		return nil, err
	}

	// load thread
	threadId := libp2pc.ConfigEncodeKey(leave.Header.ThreadPk)
	_, thrd := h.getThread(threadId)
	if thrd == nil {
		return nil, errors.New("invalid leave block")
	}

	// verify thread sig
	if err := thrd.Verify(signed); err != nil {
		return nil, err
	}

	// handle
	addr, err := thrd.HandleLeaveBlock(&pid, env, signed, leave, false)
	if err != nil {
		return nil, err
	}

	// send notification
	notification, err := newThreadNotification(thrd.PrivKey, leave.Header, repo.PeerLeftNotification)
	if err != nil {
		return nil, err
	}
	notification.Subject = thrd.Name
	notification.SubjectId = thrd.Id
	notification.BlockId = addr.B58String()
	notification.Body = "left"
	if err := h.sendNotification(notification); err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *ThreadsService) handleData(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	signed, err := unpackThreadMessage(env)
	if err != nil {
		return nil, err
	}
	data := new(pb.ThreadData)
	if err := proto.Unmarshal(signed.Block, data); err != nil {
		return nil, err
	}

	// load thread
	threadId := libp2pc.ConfigEncodeKey(data.Header.ThreadPk)
	_, thrd := h.getThread(threadId)
	if thrd == nil {
		return nil, errors.New("invalid data block")
	}

	// verify thread sig
	if err := thrd.Verify(signed); err != nil {
		return nil, err
	}

	// handle
	addr, err := thrd.HandleDataBlock(&pid, env, signed, data, false)
	if err != nil {
		return nil, err
	}

	// send notification
	// check for old username format
	if data.Header.AuthorUnCipher == nil {
		data.Header.AuthorUnCipher = data.UsernameCipher
	}
	var notification *repo.Notification
	switch data.Type {
	case pb.ThreadData_PHOTO:
		notification, err = newThreadNotification(thrd.PrivKey, data.Header, repo.PhotoAddedNotification)
		if err != nil {
			return nil, err
		}
		notification.DataId = data.DataId
		notification.Body = "added a photo"
	case pb.ThreadData_TEXT:
		notification, err = newThreadNotification(thrd.PrivKey, data.Header, repo.TextAddedNotification)
		if err != nil {
			return nil, err
		}
		body, err := thrd.Decrypt(data.CaptionCipher)
		if err != nil {
			return nil, err
		}
		notification.Body = string(body)
	}
	notification.BlockId = addr.B58String()
	notification.Subject = thrd.Name
	notification.SubjectId = thrd.Id
	if err := h.sendNotification(notification); err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *ThreadsService) handleAnnotation(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	signed, err := unpackThreadMessage(env)
	if err != nil {
		return nil, err
	}
	annotation := new(pb.ThreadAnnotation)
	if err := proto.Unmarshal(signed.Block, annotation); err != nil {
		return nil, err
	}

	// load thread
	threadId := libp2pc.ConfigEncodeKey(annotation.Header.ThreadPk)
	_, thrd := h.getThread(threadId)
	if thrd == nil {
		return nil, errors.New("invalid annotation block")
	}

	// verify thread sig
	if err := thrd.Verify(signed); err != nil {
		return nil, err
	}

	// handle
	addr, err := thrd.HandleAnnotationBlock(&pid, env, signed, annotation, false)
	if err != nil {
		return nil, err
	}

	// find dataId block locally
	dataBlock := h.Datastore().Blocks().Get(annotation.DataId)
	if dataBlock == nil {
		return nil, nil
	}
	var target string
	authorId, err := ipfs.IdFromEncodedPublicKey(dataBlock.AuthorPk)
	if err != nil {
		return nil, err
	}
	// NOTE: not restricted to photo annotations here, just currently only thing possible
	if authorId.Pretty() == h.Node().Identity.Pretty() {
		target = "your photo"
	} else {
		target = "a photo"
	}

	// send notification
	var notification *repo.Notification
	switch annotation.Type {
	case pb.ThreadAnnotation_COMMENT:
		notification, err = newThreadNotification(thrd.PrivKey, annotation.Header, repo.CommentAddedNotification)
		if err != nil {
			return nil, err
		}
		body, err := thrd.Decrypt(annotation.CaptionCipher)
		if err != nil {
			return nil, err
		}
		notification.Body = fmt.Sprintf("commented on %s: \"%s\"", target, string(body))
	case pb.ThreadAnnotation_LIKE:
		notification, err = newThreadNotification(thrd.PrivKey, annotation.Header, repo.LikeAddedNotification)
		if err != nil {
			return nil, err
		}
		notification.Body = "liked " + target
	}
	notification.BlockId = addr.B58String()
	notification.DataId = dataBlock.DataId
	notification.Subject = thrd.Name
	notification.SubjectId = thrd.Id
	if err := h.sendNotification(notification); err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *ThreadsService) handleIgnore(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	signed, err := unpackThreadMessage(env)
	if err != nil {
		return nil, err
	}
	ignore := new(pb.ThreadIgnore)
	if err := proto.Unmarshal(signed.Block, ignore); err != nil {
		return nil, err
	}

	// load thread
	threadId := libp2pc.ConfigEncodeKey(ignore.Header.ThreadPk)
	_, thrd := h.getThread(threadId)
	if thrd == nil {
		return nil, errors.New("invalid ignore block")
	}

	// verify thread sig
	if err := thrd.Verify(signed); err != nil {
		return nil, err
	}

	// handle
	if _, err := thrd.HandleIgnoreBlock(&pid, env, signed, ignore, false); err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *ThreadsService) handleMerge(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	signed, err := unpackThreadMessage(env)
	if err != nil {
		return nil, err
	}
	merge := new(pb.ThreadMerge)
	if err := proto.Unmarshal(signed.Block, merge); err != nil {
		return nil, err
	}

	// load thread
	threadId := libp2pc.ConfigEncodeKey(merge.ThreadPk)
	_, thrd := h.getThread(threadId)
	if thrd == nil {
		return nil, errors.New("invalid merge block")
	}

	// verify thread sig
	if err := thrd.Verify(signed); err != nil {
		return nil, err
	}

	// handle
	if _, err := thrd.HandleMergeBlock(&pid, env.Message, signed, merge, false); err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *ThreadsService) handleBlock(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	pbblock := new(pb.Block)
	err := ptypes.UnmarshalAny(env.Message.Payload, pbblock)
	if err != nil {
		return nil, err
	}
	id, err := cid.Decode(pbblock.Cid)
	if err != nil {
		return nil, err
	}
	block, err := blocks.NewBlockWithCid(pbblock.RawData, id)
	if err != nil {
		return nil, err
	}
	if err := h.Node().Blocks.AddBlock(block); err != nil {
		return nil, err
	}
	log.Debugf("received IPFS_BLOCK message from %s", pid.Pretty())
	return nil, nil
}

func (h *ThreadsService) handleStore(pid peer.ID, env *pb.Envelope) (*pb.Envelope, error) {
	errRes := func(error string) *pb.Message {
		payload := &any.Any{Value: []byte(error)}
		message := &pb.Message{
			Type:    pb.Message_ERROR,
			Payload: payload,
		}
		return message
	}
	cList := new(pb.CidList)
	err := ptypes.UnmarshalAny(env.Message.Payload, cList)
	if err != nil {
		return newEnvelope(h.Node().PrivateKey, errRes("could not unmarshall message"))
	}
	var need []string
	for _, id := range cList.Cids {
		decoded, err := cid.Decode(id)
		if err != nil {
			continue
		}
		has, err := h.Node().Blockstore.Has(decoded)
		if err != nil || !has {
			need = append(need, decoded.String())
		}
	}
	log.Debugf("requesting %d blocks from %s", len(need), pid.Pretty())
	resp := new(pb.CidList)
	resp.Cids = need
	payload, err := ptypes.MarshalAny(resp)
	if err != nil {
		return newEnvelope(h.Node().PrivateKey, errRes("error marshalling response"))
	}
	message := &pb.Message{
		Type:    pb.Message_STORE,
		Payload: payload,
	}
	return newEnvelope(h.Node().PrivateKey, message)
}

// unpackThreadMessage returns an envelope's signed thread block
func unpackThreadMessage(env *pb.Envelope) (*pb.SignedThreadBlock, error) {
	signed := new(pb.SignedThreadBlock)
	if err := ptypes.UnmarshalAny(env.Message.Payload, signed); err != nil {
		return nil, err
	}
	return signed, nil
}

// newThreadNotification returns new thread notification
func newThreadNotification(
	threadKey libp2pc.PrivKey,
	header *pb.ThreadBlockHeader,
	ntype repo.NotificationType) (*repo.Notification, error) {
	date, err := ptypes.Timestamp(header.Date)
	if err != nil {
		return nil, err
	}
	authorPk, err := libp2pc.UnmarshalPublicKey(header.AuthorPk)
	if err != nil {
		return nil, err
	}
	authorId, err := peer.IDFromPublicKey(authorPk)
	if err != nil {
		return nil, err
	}
	var authorUn string
	if header.AuthorUnCipher != nil {
		authorUnb, err := crypto.Decrypt(threadKey, header.AuthorUnCipher)
		if err != nil {
			return nil, err
		}
		authorUn = string(authorUnb)
	}
	return &repo.Notification{
		Id:            ksuid.New().String(),
		Date:          date,
		ActorId:       authorId.Pretty(),
		ActorUsername: authorUn,
		Type:          ntype,
	}, nil
}
