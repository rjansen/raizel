package firestore

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/firestore"
)

type SetOption interface {
	delegate() firestore.SetOption
}

type DocumentRef interface {
	Get(context.Context) (DocumentSnapshot, error)
	Set(context.Context, interface{}, ...SetOption) error
	Delete(context.Context) error
	delegate() *firestore.DocumentRef
}

type CollectionRef interface {
	Query
	delegate() *firestore.CollectionRef
}

type DocumentSnapshot interface {
	DataTo(interface{}) error
	Exists() bool
}

type DocumentIterator interface {
	GetAll() ([]DocumentSnapshot, error)
}

type Direction = firestore.Direction

type Query interface {
	Documents(context.Context) DocumentIterator
	Where(string, string, interface{}) Query
	OrderBy(string, Direction) Query
	Offset(int) Query
	Limit(int) Query
}

type Client interface {
	Close() error
	Doc(string) DocumentRef
	Collection(string) CollectionRef
	GetAll(context.Context, ...DocumentRef) ([]DocumentSnapshot, error)
	Batch() WriteBatch
}

type WriteBatch interface {
	Set(DocumentRef, interface{}, ...SetOption) WriteBatch
	Delete(DocumentRef) WriteBatch
	Commit(context.Context) error
}

// delegate implementation
var (
	MergeAll                          = mergeSetOption{firestore.MergeAll}
	Asc                     Direction = firestore.Asc
	Desc                    Direction = firestore.Desc
	ErrBlankFirestoreClient           = errors.New("err_blankclient")
)

type mergeSetOption struct {
	firestore.SetOption
}

func (option mergeSetOption) delegate() firestore.SetOption {
	return option.SetOption
}

type documentRef struct {
	*firestore.DocumentRef
}

func (doc *documentRef) delegate() *firestore.DocumentRef {
	return doc.DocumentRef
}

func (doc *documentRef) Get(ctx context.Context) (DocumentSnapshot, error) {
	fmt.Println("firestore_get")
	return doc.DocumentRef.Get(ctx)
}

func (doc *documentRef) Set(ctx context.Context, data interface{}, opts ...SetOption) error {
	fopts := make([]firestore.SetOption, len(opts))
	for index, opt := range opts {
		fopts[index] = opt.delegate()
	}
	_, err := doc.DocumentRef.Set(ctx, data, fopts...)
	return err
}

func (doc *documentRef) Delete(ctx context.Context) error {
	_, err := doc.DocumentRef.Delete(ctx)
	return err
}

type documentIterator struct {
	*firestore.DocumentIterator
}

func (iter *documentIterator) GetAll() ([]DocumentSnapshot, error) {
	docs, err := iter.DocumentIterator.GetAll()
	if err != nil {
		return nil, err
	}
	delegates := make([]DocumentSnapshot, len(docs))
	for index, doc := range docs {
		delegates[index] = doc
	}
	return delegates, nil
}

type query struct {
	firestore.Query
}

func (q query) Where(path, op string, value interface{}) Query {
	return query{
		Query: q.Query.Where(path, op, value),
	}
}

func (q query) Documents(ctx context.Context) DocumentIterator {
	return &documentIterator{
		DocumentIterator: q.Query.Documents(ctx),
	}
}
func (q query) OrderBy(path string, direction Direction) Query {
	return query{
		Query: q.Query.OrderBy(path, direction),
	}
}

func (q query) Offset(n int) Query {
	return query{
		Query: q.Query.Offset(n),
	}
}

func (q query) Limit(n int) Query {
	return query{
		Query: q.Query.Limit(n),
	}
}

type collectionRef struct {
	query
	*firestore.CollectionRef
}

func (coll *collectionRef) delegate() *firestore.CollectionRef {
	return coll.CollectionRef
}

func (coll *collectionRef) Where(path, op string, value interface{}) Query {
	return query{
		Query: coll.CollectionRef.Where(path, op, value),
	}
}

func (coll *collectionRef) Documents(ctx context.Context) DocumentIterator {
	return &documentIterator{
		DocumentIterator: coll.CollectionRef.Documents(ctx),
	}
}

type writeBatch struct {
	*firestore.WriteBatch
}

func (w *writeBatch) Set(ref DocumentRef, data interface{}, opts ...SetOption) WriteBatch {
	fopts := make([]firestore.SetOption, len(opts))
	for index, opt := range opts {
		fopts[index] = opt.delegate()
	}
	return &writeBatch{
		WriteBatch: w.WriteBatch.Set(ref.delegate(), data, fopts...),
	}
}

func (w *writeBatch) Delete(ref DocumentRef) WriteBatch {
	return &writeBatch{
		WriteBatch: w.WriteBatch.Delete(ref.delegate()),
	}
}

func (w *writeBatch) Commit(ctx context.Context) error {
	_, err := w.WriteBatch.Commit(ctx)
	return err
}

type client struct {
	*firestore.Client
}

func (c *client) Doc(path string) DocumentRef {
	return &documentRef{
		DocumentRef: c.Client.Doc(path),
	}
}

func (c *client) Collection(path string) CollectionRef {
	return &collectionRef{
		CollectionRef: c.Client.Collection(path),
	}
}

func (c *client) GetAll(ctx context.Context, refs ...DocumentRef) ([]DocumentSnapshot, error) {
	frefs := make([]*firestore.DocumentRef, len(refs))
	for index, ref := range refs {
		frefs[index] = ref.delegate()
	}
	fdocs, err := c.Client.GetAll(ctx, frefs)
	if err != nil {
		return nil, err
	}
	docs := make([]DocumentSnapshot, len(fdocs))
	for index, fdoc := range fdocs {
		docs[index] = fdoc
	}
	return docs, nil
}

func (c *client) Batch() WriteBatch {
	return &writeBatch{
		WriteBatch: c.Client.Batch(),
	}
}

func newFirestoreClient(projectID string) (*firestore.Client, error) {
	fmt.Println("begin_firestore_client")
	defer fmt.Println("end_firestore_client")
	return firestore.NewClient(context.Background(), projectID)
}

func newClient(fclient *firestore.Client) (Client, error) {
	if fclient == nil {
		return nil, ErrBlankFirestoreClient
	}
	return &client{fclient}, nil
}

func NewClient(projectID string) Client {
	fcli, err := newFirestoreClient(projectID)
	if err != nil {
		panic(err)
	}
	cli, err := newClient(fcli)
	if err != nil {
		panic(err)
	}
	return cli
}
