package whatsmeow

import (
	"context"
	waBinary "go.mau.fi/whatsmeow/binary"
)

func (cli *Client) sendIQXmppPing(query *infoQuery) (<-chan *waBinary.Node, []byte, error) {
	if cli == nil {
		return nil, nil, ErrClientIsNil
	}
	waiter := cli.waitResponse(query.ID)
	attrs := waBinary.Attrs{
		"type": string(query.Type),
	}
	if !query.To.IsEmpty() {
		attrs["to"] = query.To
	}
	data, err := cli.sendNodeAndGetData(
		context.Background(),
		waBinary.Node{
			Tag:     "iq",
			Attrs:   attrs,
			Content: query.Content,
		})
	if err != nil {
		cli.cancelResponse(query.ID, waiter)
		return nil, data, err
	}
	return waiter, data, nil
}
