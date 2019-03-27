package db

import (
	"database/sql"
	"sync"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/textileio/go-textile/pb"
	"github.com/textileio/go-textile/repo"
)

var threadStore repo.ThreadStore

func init() {
	setupThreadDB()
}

func setupThreadDB() {
	conn, _ := sql.Open("sqlite3", ":memory:")
	initDatabaseTables(conn, "")
	threadStore = NewThreadStore(conn, new(sync.Mutex))
}

func TestThreadDB_Add(t *testing.T) {
	if err := threadStore.Add(&pb.Thread{
		Id:        "Qmabc123",
		Key:       ksuid.New().String(),
		Sk:        make([]byte, 8),
		Name:      "boom",
		Schema:    "Qm...",
		Initiator: "123",
		Type:      pb.Thread_OPEN,
		Members:   []string{"P1,P2"},
		Sharing:   pb.Thread_SHARED,
		State:     pb.Thread_LOADED,
	}); err != nil {
		t.Error(err)
	}
	stmt, err := threadStore.PrepareQuery("select id from threads where id=?")
	if err != nil {
		t.Error(err)
	}
	defer stmt.Close()
	var id string
	if err := stmt.QueryRow("Qmabc123").Scan(&id); err != nil {
		t.Error(err)
	}
	if id != "Qmabc123" {
		t.Errorf(`expected "Qmabc123" got %s`, id)
	}
}

func TestThreadDB_Get(t *testing.T) {
	setupThreadDB()
	if err := threadStore.Add(&pb.Thread{
		Id:        "Qmabc",
		Key:       ksuid.New().String(),
		Sk:        make([]byte, 8),
		Name:      "boom",
		Schema:    "Qm...",
		Initiator: "123",
		Type:      pb.Thread_OPEN,
		Members:   []string{},
		Sharing:   pb.Thread_SHARED,
		State:     pb.Thread_LOADED,
	}); err != nil {
		t.Error(err)
	}
	th := threadStore.Get("Qmabc")
	if th == nil {
		t.Error("could not get thread")
	}
}

func TestThreadDB_List(t *testing.T) {
	setupThreadDB()
	if err := threadStore.Add(&pb.Thread{
		Id:        "Qm123",
		Key:       ksuid.New().String(),
		Sk:        make([]byte, 8),
		Name:      "boom",
		Schema:    "Qm...",
		Initiator: "123",
		Type:      pb.Thread_PRIVATE,
		Members:   []string{},
		Sharing:   pb.Thread_NOT_SHARED,
		State:     pb.Thread_LOADED,
	}); err != nil {
		t.Error(err)
	}
	if err := threadStore.Add(&pb.Thread{
		Id:      "Qm456",
		Key:     ksuid.New().String(),
		Sk:      make([]byte, 8),
		Name:    "boom",
		Schema:  "Qm...",
		Type:    pb.Thread_PRIVATE,
		Members: []string{},
		Sharing: pb.Thread_NOT_SHARED,
		State:   pb.Thread_LOADED,
	}); err != nil {
		t.Error(err)
	}
	all := threadStore.List()
	if len(all.Items) != 2 {
		t.Error("returned incorrect number of threads")
		return
	}
}

func TestThreadDB_Count(t *testing.T) {
	setupThreadDB()
	err := threadStore.Add(&pb.Thread{
		Id:        "Qm123count",
		Key:       ksuid.New().String(),
		Sk:        make([]byte, 8),
		Name:      "boom",
		Schema:    "Qm...",
		Initiator: "123",
		Type:      pb.Thread_PRIVATE,
		Members:   []string{},
		Sharing:   pb.Thread_NOT_SHARED,
		State:     pb.Thread_LOADED,
	})
	if err != nil {
		t.Error(err)
	}
	cnt := threadStore.Count()
	if cnt != 1 {
		t.Error("returned incorrect count of threads")
		return
	}
}

func TestThreadDB_UpdateHead(t *testing.T) {
	setupThreadDB()
	err := threadStore.Add(&pb.Thread{
		Id:        "Qmabc",
		Key:       ksuid.New().String(),
		Sk:        make([]byte, 8),
		Name:      "boom",
		Schema:    "Qm...",
		Initiator: "123",
		Type:      pb.Thread_PRIVATE,
		Members:   []string{},
		Sharing:   pb.Thread_NOT_SHARED,
		State:     pb.Thread_LOADED,
	})
	if err != nil {
		t.Error(err)
	}
	if err := threadStore.UpdateHead("Qmabc", "12345"); err != nil {
		t.Error(err)
	}
	th := threadStore.Get("Qmabc")
	if th == nil {
		t.Error("could not get thread")
	}
	if th.Head != "12345" {
		t.Error("update head failed")
	}
}

func TestThreadDB_UpdateName(t *testing.T) {
	if err := threadStore.UpdateName("Qmabc", "boom2"); err != nil {
		t.Error(err)
	}
	th := threadStore.Get("Qmabc")
	if th == nil {
		t.Error("could not get thread")
	}
	if th.Name != "boom2" {
		t.Error("update name failed")
	}
}

func TestThreadDB_Delete(t *testing.T) {
	setupThreadDB()
	if err := threadStore.Add(&pb.Thread{
		Id:        "Qm789",
		Key:       ksuid.New().String(),
		Sk:        make([]byte, 8),
		Name:      "boom",
		Schema:    "Qm...",
		Initiator: "123",
		Type:      pb.Thread_PRIVATE,
		Members:   []string{},
		Sharing:   pb.Thread_NOT_SHARED,
		State:     pb.Thread_LOADED,
	}); err != nil {
		t.Error(err)
	}
	all := threadStore.List()
	if len(all.Items) == 0 {
		t.Error("returned incorrect number of threads")
		return
	}
	if err := threadStore.Delete(all.Items[0].Id); err != nil {
		t.Error(err)
	}
	stmt, err := threadStore.PrepareQuery("select id from threads where id=?")
	if err != nil {
		t.Error(err)
	}
	defer stmt.Close()
	var id string
	if err := stmt.QueryRow(all.Items[0].Id).Scan(&id); err == nil {
		t.Error("Delete failed")
	}
}
