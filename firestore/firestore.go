package firestore

import (
	"context"
	"errors"

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

type DocumentSnapshot interface {
	DataTo(interface{}) error
	Exists() bool
}

type Client interface {
	Close() error
	Doc(string) DocumentRef
	GetAll(context.Context, ...DocumentRef) ([]DocumentSnapshot, error)
}

// delegate implementation
var (
	MergeAll                = mergeSetOption{firestore.MergeAll}
	ErrBlankFirestoreClient = errors.New("err_blankclient")
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

type client struct {
	*firestore.Client
}

func (c *client) Doc(path string) DocumentRef {
	fref := c.Client.Doc(path)
	return &documentRef{fref}
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

func newFirestoreClient(projectID string) (*firestore.Client, error) {
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
