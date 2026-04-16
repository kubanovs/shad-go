//go:build !solution

package hotelbusiness

import "sort"

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type GuestEvent struct {
	Date      int
	EventType int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func sortGuestEvent(guestEvents []GuestEvent, i, j int) bool {
	if guestEvents[i].Date < guestEvents[j].Date {
		return true
	} else if guestEvents[i].Date > guestEvents[j].Date {
		return false
	} else {
		return guestEvents[i].EventType < guestEvents[j].EventType
	}
}

func ComputeLoad(guests []Guest) []Load {
	if len(guests) == 0 {
		return []Load{}
	}

	guestEvents := []GuestEvent{}

	for _, guest := range guests {
		guestEvents = append(guestEvents, GuestEvent{guest.CheckInDate, 1})
		guestEvents = append(guestEvents, GuestEvent{guest.CheckOutDate, -1})
	}

	sort.Slice(guestEvents, func(i, j int) bool {
		return sortGuestEvent(guestEvents, i, j)
	})

	var loads []Load
	var curGuests int

	for i, guestEvent := range guestEvents {
		if i != 0 && (len(loads) == 0 || curGuests != loads[len(loads)-1].GuestCount) && guestEvents[i-1].Date != guestEvents[i].Date {
			loads = append(loads, Load{guestEvents[i-1].Date, curGuests})
		}

		switch guestEvent.EventType {
		case -1:
			curGuests -= 1
		case 1:
			curGuests += 1
		}
	}

	loads = append(loads, Load{guestEvents[len(guestEvents)-1].Date, curGuests})

	return loads
}
