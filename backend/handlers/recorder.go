package handlers

import "time"

/*
* @function beginRecord
* @description
* While client's recording field is true, records notes that are sent into the backend in a notes array
* This array is saved to the client struct for later playback.

* @exported: false
* @return N/A
 */
func beginRecord(client *Client) {
	client.recordedNotes = &[]SocketEventStruct{}
	client.recording = true
	for client.recording {
		payloadStruct := <-client.recordNotes
		*(client.recordedNotes) = append(*(client.recordedNotes), payloadStruct)
	}
}

/*
* @function playRecording
* @description
* Play the recorded notes saved in recordNotes channel. Sleep for the amount of time between notes to
* simulate pauses. Passes in all of these notes as they play into the broadcastsocketevent function so
* they can be executed.

* @exported: false
* @return N/A
 */
func playRecording(client *Client) {
	notes := *client.recordedNotes
	BroadcastSocketEventToAllClient(client, true, SocketEventStruct{
		EventName: "keyboardPress",
		EventPayload: EventPayloadStruct{
			User:    notes[0].EventName,
			Message: notes[0].EventPayload.Message,
			Time:    notes[0].EventPayload.Time,
		},
	})
	for i := 1; i < len(notes); i++ {
		delay := notes[i].EventPayload.Time - notes[i-1].EventPayload.Time
		time.Sleep(time.Duration(delay) * time.Millisecond)
		BroadcastSocketEventToAllClient(client, true, SocketEventStruct{
			EventName: "keyboardPress",
			EventPayload: EventPayloadStruct{
				User:    notes[i].EventName,
				Message: notes[i].EventPayload.Message,
				Time:    notes[i].EventPayload.Time,
			},
		})
	}
}
