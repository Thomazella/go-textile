package wallet

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/crypto"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	"github.com/textileio/textile-go/util"
	"github.com/textileio/textile-go/wallet/thread"
	mh "gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
	libp2pc "gx/ipfs/QmaPbCnUMBohSGo3KnxEa2bHqyJVVeEEcwtqJAYxerieBo/go-libp2p-crypto"
)

// AddThread adds a thread with a given name and secret key
func (w *Wallet) AddThread(name string, secret libp2pc.PrivKey, join bool) (*thread.Thread, error) {
	skb, err := secret.Bytes()
	if err != nil {
		return nil, err
	}
	pkb, err := secret.GetPublic().Bytes()
	if err != nil {
		return nil, err
	}
	pk := libp2pc.ConfigEncodeKey(pkb)
	threadModel := &repo.Thread{
		Id:      pk,
		Name:    name,
		PrivKey: skb,
	}
	if err := w.datastore.Threads().Add(threadModel); err != nil {
		return nil, err
	}

	// load as active thread
	thrd, err := w.loadThread(threadModel)
	if err != nil {
		return nil, err
	}

	// we join here if we're the creator
	if join {
		if _, err := thrd.JoinInitial(); err != nil {
			return nil, err
		}
	}

	// notify listeners
	w.sendUpdate(Update{Id: thrd.Id, Name: thrd.Name, Type: ThreadAdded})

	log.Debugf("added a new thread %s with name %s", thrd.Id, name)

	return thrd, nil
}

// RemoveThread removes a thread
func (w *Wallet) RemoveThread(id string) (mh.Multihash, error) {
	if !w.IsOnline() {
		return nil, ErrOffline
	}

	i, thrd := w.GetThread(id) // gets the loaded thread
	if thrd == nil {
		return nil, errors.New("thread not found")
	}

	// notify peers
	addr, err := thrd.Leave()
	if err != nil {
		return nil, err
	}

	// remove model from db
	if err := w.datastore.Threads().Delete(thrd.Id); err != nil {
		return nil, err
	}

	// clean up
	copy(w.threads[*i:], w.threads[*i+1:])
	w.threads[len(w.threads)-1] = nil
	w.threads = w.threads[:len(w.threads)-1]

	// notify listeners
	w.sendUpdate(Update{Id: thrd.Id, Name: thrd.Name, Type: ThreadRemoved})

	log.Infof("removed thread %s with name %s", thrd.Id, thrd.Name)

	return addr, nil
}

// AcceptThreadInvite attemps to download an encrypted thread key from an internal invite,
// add the thread, and notify the inviter of the join
func (w *Wallet) AcceptThreadInvite(blockId string) (mh.Multihash, error) {
	if !w.IsOnline() {
		return nil, ErrOffline
	}

	// download
	messageb, err := util.GetDataAtPath(w.ipfs, fmt.Sprintf("%s", blockId))
	if err != nil {
		return nil, err
	}
	env := new(pb.Envelope)
	if err := proto.Unmarshal(messageb, env); err != nil {
		return nil, err
	}

	// verify author sig
	authorPk, err := libp2pc.UnmarshalPublicKey(env.Pk)
	if err != nil {
		return nil, err
	}
	if err := w.VerifyEnvelope(env); err != nil {
		return nil, err
	}

	// unpack invite
	signed := new(pb.SignedThreadBlock)
	if err := ptypes.UnmarshalAny(env.Message.Payload, signed); err != nil {
		return nil, err
	}
	invite := new(pb.ThreadInvite)
	if err := proto.Unmarshal(signed.Block, invite); err != nil {
		return nil, err
	}

	// verify invitee
	if invite.InviteeId != w.ipfs.Identity.Pretty() {
		return nil, errors.New("invalid invitee")
	}

	// decrypt thread key with private key
	key, err := w.GetPrivKey()
	if err != nil {
		return nil, err
	}
	skb, err := crypto.Decrypt(key, invite.SkCipher)
	if err != nil {
		return nil, err
	}
	sk, err := libp2pc.UnmarshalPrivateKey(skb)
	if err != nil {
		return nil, err
	}
	pkb, err := sk.GetPublic().Bytes()
	if err != nil {
		return nil, err
	}

	// ensure we dont already have this thread
	id := libp2pc.ConfigEncodeKey(pkb)
	if _, thrd := w.GetThread(id); thrd != nil {
		// thread exists, aborting
		return nil, nil
	}

	// verify thread sig
	if err := crypto.Verify(sk.GetPublic(), signed.Block, signed.ThreadSig); err != nil {
		return nil, err
	}

	// add it
	thrd, err := w.AddThread(invite.SuggestedName, sk, false)
	if err != nil {
		return nil, err
	}

	// follow parents, update head
	if err := thrd.HandleInviteMessage(invite); err != nil {
		return nil, err
	}

	// join the thread
	addr, err := thrd.Join(authorPk, blockId)
	if err != nil {
		return nil, err
	}

	// invite devices
	if err := w.InviteDevices(thrd); err != nil {
		return nil, err
	}

	return addr, nil
}

// AcceptExternalThreadInvite attemps to download an encrypted thread key from an external invite,
// add the thread, and notify the inviter of the join
func (w *Wallet) AcceptExternalThreadInvite(blockId string, key []byte) (mh.Multihash, error) {
	if !w.IsOnline() {
		return nil, ErrOffline
	}

	// download
	envb, err := util.GetDataAtPath(w.ipfs, fmt.Sprintf("%s", blockId))
	if err != nil {
		return nil, err
	}
	env := new(pb.Envelope)
	if err := proto.Unmarshal(envb, env); err != nil {
		return nil, err
	}

	// verify author sig
	authorPk, err := libp2pc.UnmarshalPublicKey(env.Pk)
	if err != nil {
		return nil, err
	}
	if err := w.VerifyEnvelope(env); err != nil {
		return nil, err
	}

	// unpack invite
	signed := new(pb.SignedThreadBlock)
	if err := ptypes.UnmarshalAny(env.Message.Payload, signed); err != nil {
		return nil, err
	}
	invite := new(pb.ThreadExternalInvite)
	if err := proto.Unmarshal(signed.Block, invite); err != nil {
		return nil, err
	}

	// decrypt thread key
	skb, err := crypto.DecryptAES(invite.SkCipher, key)
	if err != nil {
		return nil, err
	}
	sk, err := libp2pc.UnmarshalPrivateKey(skb)
	if err != nil {
		return nil, err
	}
	pkb, err := sk.GetPublic().Bytes()
	if err != nil {
		return nil, err
	}

	// ensure we dont already have this thread
	id := libp2pc.ConfigEncodeKey(pkb)
	if _, thrd := w.GetThread(id); thrd != nil {
		// thread exists, aborting
		return nil, nil
	}

	// verify thread sig
	if err := crypto.Verify(sk.GetPublic(), signed.Block, signed.ThreadSig); err != nil {
		return nil, err
	}

	// add it
	thrd, err := w.AddThread(invite.SuggestedName, sk, false)
	if err != nil {
		return nil, err
	}

	// follow parents, update head
	if err := thrd.HandleExternalInviteMessage(invite); err != nil {
		return nil, err
	}

	// join the thread
	addr, err := thrd.Join(authorPk, blockId)
	if err != nil {
		return nil, err
	}

	// invite devices
	if err := w.InviteDevices(thrd); err != nil {
		return nil, err
	}

	return addr, nil
}

// Threads lists loaded threads
func (w *Wallet) Threads() []*thread.Thread {
	return w.threads
}

// GetThread get a thread by id from loaded threads
func (w *Wallet) GetThread(id string) (*int, *thread.Thread) {
	for i, thrd := range w.threads {
		if thrd.Id == id {
			return &i, thrd
		}
	}
	return nil, nil
}

// ThreadInfo gets thread info
func (w *Wallet) ThreadInfo(id string) (*thread.Info, error) {
	_, thrd := w.GetThread(id)
	if thrd == nil {
		return nil, errors.New(fmt.Sprintf("cound not find thread: %s", id))
	}
	return thrd.Info()
}
