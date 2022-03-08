package members

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/gotd/td/internal/testutil"
	"github.com/gotd/td/tg"
)

func TestChannelMembers_Count(t *testing.T) {
	a := require.New(t)
	ctx := context.Background()
	mock, m := testManager(t)

	ch := m.Channel(getTestChannel())
	members, err := Channel(ctx, ch)
	a.NoError(err)

	mock.ExpectCall(&tg.ChannelsGetParticipantsRequest{
		Channel: ch.InputChannel(),
		Filter:  &tg.ChannelParticipantsRecent{},
		Offset:  0,
		Limit:   1,
	}).ThenErr(testutil.TestError())
	_, err = members.Count(ctx)
	a.Error(err)

	mock.ExpectCall(&tg.ChannelsGetParticipantsRequest{
		Channel: ch.InputChannel(),
		Filter:  &tg.ChannelParticipantsRecent{},
		Offset:  0,
		Limit:   1,
	}).ThenResult(&tg.ChannelsChannelParticipantsNotModified{})
	_, err = members.Count(ctx)
	a.Error(err)

	mock.ExpectCall(&tg.ChannelsGetParticipantsRequest{
		Channel: ch.InputChannel(),
		Filter:  &tg.ChannelParticipantsRecent{},
		Offset:  0,
		Limit:   1,
	}).ThenResult(&tg.ChannelsChannelParticipants{
		Count: 10,
	})
	count, err := members.Count(ctx)
	a.NoError(err)
	a.Equal(10, count)
}

func TestChannelMembers_ForEach(t *testing.T) {
	a := require.New(t)
	ctx := context.Background()
	mock, m := testManager(t)

	now := time.Now()
	date := int(now.Unix())

	rawCh := getTestChannel()
	rawCh.Date = date
	ch := m.Channel(rawCh)
	members, err := Channel(ctx, ch)
	a.NoError(err)

	mock.ExpectCall(&tg.ChannelsGetFullChannelRequest{
		Channel: ch.InputChannel(),
	}).ThenResult(&tg.MessagesChatFull{
		FullChat: getTestChannelFull(),
	})
	mock.ExpectCall(&tg.ChannelsGetParticipantsRequest{
		Channel: ch.InputChannel(),
		Filter:  &tg.ChannelParticipantsRecent{},
		Offset:  0,
		Limit:   100,
	}).ThenResult(&tg.ChannelsChannelParticipants{
		Count: 10,
		Participants: []tg.ChannelParticipantClass{
			&tg.ChannelParticipant{
				UserID: 10,
				Date:   date,
			},
			&tg.ChannelParticipantSelf{
				UserID:    10,
				InviterID: 11,
				Date:      date,
			},
			&tg.ChannelParticipantCreator{
				UserID: 10,
				Rank:   "rank",
			},
			&tg.ChannelParticipantAdmin{
				UserID:    10,
				InviterID: 11,
				Date:      date,
				Rank:      "rank",
			},
			&tg.ChannelParticipantBanned{
				Peer: &tg.PeerUser{UserID: 10},
				Date: date,
			},
			&tg.ChannelParticipantLeft{
				Peer: &tg.PeerUser{UserID: 10},
			},
		},
		Users: []tg.UserClass{
			&tg.User{
				ID:         10,
				AccessHash: 10,
			},
			&tg.User{
				ID:         11,
				AccessHash: 10,
			},
		},
	}).ExpectCall(&tg.ChannelsGetParticipantsRequest{
		Channel: ch.InputChannel(),
		Filter:  &tg.ChannelParticipantsRecent{},
		Offset:  100,
		Limit:   100,
	}).ThenResult(&tg.ChannelsChannelParticipants{
		Count: 10,
	})

	expected := []struct {
		Status      Status
		JoinDate    time.Time
		JoinDateSet bool
		Rank        string
		RankSet     bool
		InviterID   int64
	}{
		{Status: Plain, JoinDate: now, JoinDateSet: true},
		{Status: Plain, JoinDate: now, JoinDateSet: true, InviterID: 11},
		{Status: Creator, Rank: "rank", RankSet: true, JoinDate: now, JoinDateSet: true},
		{Status: Admin, Rank: "rank", RankSet: true, JoinDate: now, JoinDateSet: true, InviterID: 11},
		{Status: Banned, JoinDate: now, JoinDateSet: true},
		{Status: Left},
	}

	i := 0
	a.NoError(members.ForEach(ctx, func(m Member) error {
		p := m.(ChannelMember)
		e := expected[i]

		a.Equal(e.Status, p.Status(), i)
		a.Equal(int64(10), p.User().ID())
		if join, ok := p.JoinDate(); e.JoinDateSet {
			a.True(ok, i)
			a.Equal(e.JoinDate.Unix(), join.Unix(), i)
		} else {
			a.False(ok, i)
		}

		if rank, ok := p.Rank(); e.RankSet {
			a.True(ok, i)
			a.Equal(e.Rank, rank, i)
		} else {
			a.False(ok, i)
		}

		if inviter, ok := p.InvitedBy(); e.InviterID != 0 {
			a.True(ok, i)
			a.Equal(e.InviterID, inviter.ID())
		} else {
			a.False(ok, i)
		}

		i++
		return nil
	}))
}
