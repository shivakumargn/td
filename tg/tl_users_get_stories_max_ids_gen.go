// Code generated by gotdgen, DO NOT EDIT.

package tg

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"go.uber.org/multierr"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/tdjson"
	"github.com/gotd/td/tdp"
	"github.com/gotd/td/tgerr"
)

// No-op definition for keeping imports.
var (
	_ = bin.Buffer{}
	_ = context.Background()
	_ = fmt.Stringer(nil)
	_ = strings.Builder{}
	_ = errors.Is
	_ = multierr.AppendInto
	_ = sort.Ints
	_ = tdp.Format
	_ = tgerr.Error{}
	_ = tdjson.Encoder{}
)

// UsersGetStoriesMaxIDsRequest represents TL type `users.getStoriesMaxIDs#ca1cb9ab`.
//
// See https://core.telegram.org/method/users.getStoriesMaxIDs for reference.
type UsersGetStoriesMaxIDsRequest struct {
	// ID field of UsersGetStoriesMaxIDsRequest.
	ID []InputUserClass
}

// UsersGetStoriesMaxIDsRequestTypeID is TL type id of UsersGetStoriesMaxIDsRequest.
const UsersGetStoriesMaxIDsRequestTypeID = 0xca1cb9ab

// Ensuring interfaces in compile-time for UsersGetStoriesMaxIDsRequest.
var (
	_ bin.Encoder     = &UsersGetStoriesMaxIDsRequest{}
	_ bin.Decoder     = &UsersGetStoriesMaxIDsRequest{}
	_ bin.BareEncoder = &UsersGetStoriesMaxIDsRequest{}
	_ bin.BareDecoder = &UsersGetStoriesMaxIDsRequest{}
)

func (g *UsersGetStoriesMaxIDsRequest) Zero() bool {
	if g == nil {
		return true
	}
	if !(g.ID == nil) {
		return false
	}

	return true
}

// String implements fmt.Stringer.
func (g *UsersGetStoriesMaxIDsRequest) String() string {
	if g == nil {
		return "UsersGetStoriesMaxIDsRequest(nil)"
	}
	type Alias UsersGetStoriesMaxIDsRequest
	return fmt.Sprintf("UsersGetStoriesMaxIDsRequest%+v", Alias(*g))
}

// FillFrom fills UsersGetStoriesMaxIDsRequest from given interface.
func (g *UsersGetStoriesMaxIDsRequest) FillFrom(from interface {
	GetID() (value []InputUserClass)
}) {
	g.ID = from.GetID()
}

// TypeID returns type id in TL schema.
//
// See https://core.telegram.org/mtproto/TL-tl#remarks.
func (*UsersGetStoriesMaxIDsRequest) TypeID() uint32 {
	return UsersGetStoriesMaxIDsRequestTypeID
}

// TypeName returns name of type in TL schema.
func (*UsersGetStoriesMaxIDsRequest) TypeName() string {
	return "users.getStoriesMaxIDs"
}

// TypeInfo returns info about TL type.
func (g *UsersGetStoriesMaxIDsRequest) TypeInfo() tdp.Type {
	typ := tdp.Type{
		Name: "users.getStoriesMaxIDs",
		ID:   UsersGetStoriesMaxIDsRequestTypeID,
	}
	if g == nil {
		typ.Null = true
		return typ
	}
	typ.Fields = []tdp.Field{
		{
			Name:       "ID",
			SchemaName: "id",
		},
	}
	return typ
}

// Encode implements bin.Encoder.
func (g *UsersGetStoriesMaxIDsRequest) Encode(b *bin.Buffer) error {
	if g == nil {
		return fmt.Errorf("can't encode users.getStoriesMaxIDs#ca1cb9ab as nil")
	}
	b.PutID(UsersGetStoriesMaxIDsRequestTypeID)
	return g.EncodeBare(b)
}

// EncodeBare implements bin.BareEncoder.
func (g *UsersGetStoriesMaxIDsRequest) EncodeBare(b *bin.Buffer) error {
	if g == nil {
		return fmt.Errorf("can't encode users.getStoriesMaxIDs#ca1cb9ab as nil")
	}
	b.PutVectorHeader(len(g.ID))
	for idx, v := range g.ID {
		if v == nil {
			return fmt.Errorf("unable to encode users.getStoriesMaxIDs#ca1cb9ab: field id element with index %d is nil", idx)
		}
		if err := v.Encode(b); err != nil {
			return fmt.Errorf("unable to encode users.getStoriesMaxIDs#ca1cb9ab: field id element with index %d: %w", idx, err)
		}
	}
	return nil
}

// Decode implements bin.Decoder.
func (g *UsersGetStoriesMaxIDsRequest) Decode(b *bin.Buffer) error {
	if g == nil {
		return fmt.Errorf("can't decode users.getStoriesMaxIDs#ca1cb9ab to nil")
	}
	if err := b.ConsumeID(UsersGetStoriesMaxIDsRequestTypeID); err != nil {
		return fmt.Errorf("unable to decode users.getStoriesMaxIDs#ca1cb9ab: %w", err)
	}
	return g.DecodeBare(b)
}

// DecodeBare implements bin.BareDecoder.
func (g *UsersGetStoriesMaxIDsRequest) DecodeBare(b *bin.Buffer) error {
	if g == nil {
		return fmt.Errorf("can't decode users.getStoriesMaxIDs#ca1cb9ab to nil")
	}
	{
		headerLen, err := b.VectorHeader()
		if err != nil {
			return fmt.Errorf("unable to decode users.getStoriesMaxIDs#ca1cb9ab: field id: %w", err)
		}

		if headerLen > 0 {
			g.ID = make([]InputUserClass, 0, headerLen%bin.PreallocateLimit)
		}
		for idx := 0; idx < headerLen; idx++ {
			value, err := DecodeInputUser(b)
			if err != nil {
				return fmt.Errorf("unable to decode users.getStoriesMaxIDs#ca1cb9ab: field id: %w", err)
			}
			g.ID = append(g.ID, value)
		}
	}
	return nil
}

// GetID returns value of ID field.
func (g *UsersGetStoriesMaxIDsRequest) GetID() (value []InputUserClass) {
	if g == nil {
		return
	}
	return g.ID
}

// MapID returns field ID wrapped in InputUserClassArray helper.
func (g *UsersGetStoriesMaxIDsRequest) MapID() (value InputUserClassArray) {
	return InputUserClassArray(g.ID)
}

// UsersGetStoriesMaxIDs invokes method users.getStoriesMaxIDs#ca1cb9ab returning error if any.
//
// See https://core.telegram.org/method/users.getStoriesMaxIDs for reference.
func (c *Client) UsersGetStoriesMaxIDs(ctx context.Context, id []InputUserClass) ([]int, error) {
	var result IntVector

	request := &UsersGetStoriesMaxIDsRequest{
		ID: id,
	}
	if err := c.rpc.Invoke(ctx, request, &result); err != nil {
		return nil, err
	}
	return []int(result.Elems), nil
}