// Code generated by gotdgen, DO NOT EDIT.

package mt

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

// RPCDropAnswerRequest represents TL type `rpc_drop_answer#58e4a740`.
type RPCDropAnswerRequest struct {
	// ReqMsgID field of RPCDropAnswerRequest.
	ReqMsgID int64
}

// RPCDropAnswerRequestTypeID is TL type id of RPCDropAnswerRequest.
const RPCDropAnswerRequestTypeID = 0x58e4a740

// Ensuring interfaces in compile-time for RPCDropAnswerRequest.
var (
	_ bin.Encoder     = &RPCDropAnswerRequest{}
	_ bin.Decoder     = &RPCDropAnswerRequest{}
	_ bin.BareEncoder = &RPCDropAnswerRequest{}
	_ bin.BareDecoder = &RPCDropAnswerRequest{}
)

func (r *RPCDropAnswerRequest) Zero() bool {
	if r == nil {
		return true
	}
	if !(r.ReqMsgID == 0) {
		return false
	}

	return true
}

// String implements fmt.Stringer.
func (r *RPCDropAnswerRequest) String() string {
	if r == nil {
		return "RPCDropAnswerRequest(nil)"
	}
	type Alias RPCDropAnswerRequest
	return fmt.Sprintf("RPCDropAnswerRequest%+v", Alias(*r))
}

// TypeID returns type id in TL schema.
//
// See https://core.telegram.org/mtproto/TL-tl#remarks.
func (*RPCDropAnswerRequest) TypeID() uint32 {
	return RPCDropAnswerRequestTypeID
}

// TypeName returns name of type in TL schema.
func (*RPCDropAnswerRequest) TypeName() string {
	return "rpc_drop_answer"
}

// TypeInfo returns info about TL type.
func (r *RPCDropAnswerRequest) TypeInfo() tdp.Type {
	typ := tdp.Type{
		Name: "rpc_drop_answer",
		ID:   RPCDropAnswerRequestTypeID,
	}
	if r == nil {
		typ.Null = true
		return typ
	}
	typ.Fields = []tdp.Field{
		{
			Name:       "ReqMsgID",
			SchemaName: "req_msg_id",
		},
	}
	return typ
}

// Encode implements bin.Encoder.
func (r *RPCDropAnswerRequest) Encode(b *bin.Buffer) error {
	if r == nil {
		return fmt.Errorf("can't encode rpc_drop_answer#58e4a740 as nil")
	}
	b.PutID(RPCDropAnswerRequestTypeID)
	return r.EncodeBare(b)
}

// EncodeBare implements bin.BareEncoder.
func (r *RPCDropAnswerRequest) EncodeBare(b *bin.Buffer) error {
	if r == nil {
		return fmt.Errorf("can't encode rpc_drop_answer#58e4a740 as nil")
	}
	b.PutLong(r.ReqMsgID)
	return nil
}

// Decode implements bin.Decoder.
func (r *RPCDropAnswerRequest) Decode(b *bin.Buffer) error {
	if r == nil {
		return fmt.Errorf("can't decode rpc_drop_answer#58e4a740 to nil")
	}
	if err := b.ConsumeID(RPCDropAnswerRequestTypeID); err != nil {
		return fmt.Errorf("unable to decode rpc_drop_answer#58e4a740: %w", err)
	}
	return r.DecodeBare(b)
}

// DecodeBare implements bin.BareDecoder.
func (r *RPCDropAnswerRequest) DecodeBare(b *bin.Buffer) error {
	if r == nil {
		return fmt.Errorf("can't decode rpc_drop_answer#58e4a740 to nil")
	}
	{
		value, err := b.Long()
		if err != nil {
			return fmt.Errorf("unable to decode rpc_drop_answer#58e4a740: field req_msg_id: %w", err)
		}
		r.ReqMsgID = value
	}
	return nil
}

// GetReqMsgID returns value of ReqMsgID field.
func (r *RPCDropAnswerRequest) GetReqMsgID() (value int64) {
	if r == nil {
		return
	}
	return r.ReqMsgID
}