package handlers

import "time"

// Tuple to store key, time pairs
type Pair struct {
	key  string
	time int
}

func beginRecord(client *Client) {
	client.recording = true
	var notes []Pair
	for client.recording {
		payloadStruct := <-client.recordNotes
		message := payloadStruct.Message
		time := payloadStruct.Time
		notes = append(notes, Pair{key: message, time: time})
	}

	BroadcastSocketEventToAllClient(client.pool, SocketEventStruct{
		EventName: "keyboardPress",
		EventPayload: EventPayloadStruct{
			User:    client.username,
			Message: notes[0].key,
			Time:    notes[0].time,
		},
	})
	for i := 1; i < len(notes); i++ {
		delay := notes[i].time - notes[i-1].time
		time.Sleep(time.Duration(delay) * time.Millisecond)
		BroadcastSocketEventToAllClient(client.pool, SocketEventStruct{
			EventName: "keyboardPress",
			EventPayload: EventPayloadStruct{
				User:    client.username,
				Message: notes[i].key,
				Time:    notes[i].time,
			},
		})
	}
}
