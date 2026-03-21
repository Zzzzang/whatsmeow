package whatsmeow

import (
	"context"
	"encoding/hex"

	waBinary "go.mau.fi/whatsmeow/binary"
	"go.mau.fi/whatsmeow/types"
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

// 是否聊天的参与者都在拿到的设备列表里
func isAllParticipantsInAllDevices(participants []types.JID, devices []types.JID) bool {
	numParticipants := len(participants)
	numParticipantsInAllDevices := 0
	for _, participant := range participants {
		u := participant.User
		for _, device := range devices {
			if u == device.User {
				numParticipantsInAllDevices++
				break
			}
		}
	}
	if numParticipantsInAllDevices == numParticipants {
		return true
	}
	return false
}

func (cli *Client) SendSignCredential(ctx context.Context) error {

	signBytes, _ := hex.DecodeString("deb1506132575dec17a4c906151f32996594fd8766880d19005022f98cbabe40")

	return cli.sendNode(ctx, waBinary.Node{
		Tag: "iq",
		Attrs: waBinary.Attrs{
			"id":    cli.GenerateMessageID(),
			"xmlns": "privatestats",
			"to":    types.ServerJID,
			"type":  "get",
		},
		Content: []waBinary.Node{{
			Tag: "sign_credential",
			Attrs: waBinary.Attrs{
				"version": "1",
			},
			Content: []waBinary.Node{{
				Tag:     "blinded_credential",
				Content: signBytes, // random.Bytes(32),
			}},
		}},
	})
}
