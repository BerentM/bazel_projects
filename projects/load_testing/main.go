package main

import (
	"database/sql"
	"fmt"
	myproto "load_testing/proto"
	"log"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	_ "github.com/lib/pq"
)

func queryDb(db *sql.DB) map[string][]string {
	types := make(map[string][]string, 0)

	rows, err := db.Query(`SELECT DISTINCT v_event_type, jb_event_payload FROM t_primary_events where v_event_type IN ('BROWSING', 'CLICK')`)
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	for rows.Next() {
		var eventType string
		var payload string
		if err := rows.Scan(&eventType, &payload); err != nil {
			log.Fatal("Failed to scan row:", err)
		}
		types[eventType] = append(types[eventType], payload)
	}
	return types
}

func setPayloadJson(data string, protoObject proto.Message) ([]byte, error) {
	if err := protojson.Unmarshal([]byte(data), protoObject); err != nil {
		return nil, err
	}
	body, err := proto.Marshal(protoObject)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func main() {
	var requests [][]byte
	var keys []string

	connStr := "postgresql://gorm:gorm@localhost:5453/gorm?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	events := queryDb(db)
	for key := range events {
		keys = append(keys, key)
	}

	// Create a protobuf message with the query result
	for _, key := range keys {
		for _, event := range events[key] {
			var err error
			var payload []byte
			switch key {
			case "BROWSING":
				payload, err = setPayloadJson(event, &myproto.BrowsingEvent{})
			case "CLICK":
				payload, err = setPayloadJson(event, &myproto.ClickEvent{})
			case "DOWNLOAD":
			case "CLIPBOARD":
			case "TOP_SITES":
				// var tmpJson []byte
				// lCount := len(event.TopSitesEvent)
				// payloadJson = append(payloadJson, '[')
				// for i := 0; i < lCount; i++ {
				// 	tmpJson, err = protojson.Marshal(event.TopSitesEvent[i])
				// 	if i > 0 {
				// 		payloadJson = append(payloadJson, ',')
				// 		payloadJson = append(payloadJson, tmpJson...)
				// 	} else {
				// 		payloadJson = append(payloadJson, tmpJson...)
				// 	}
				// }
				// payloadJson = append(payloadJson, ']')
			case "BROWSING_DATA":
			case "EXTENSION_MANAGEMENTS":
			case "CONTENT_SETTINGS":
			case "SENSOR":
			case "BROWSER_LOCATION":
			case "NETWORK_ACTIVITY":
			case "TAB_OPENED":
			case "LOCATION":
			case "SOFTWARE":
			case "HARDWARE":
			case "NETWORK_INFO":
			case "SCROLL":
			case "CHANGE":
			case "INPUT":
			case "SUBMIT":
			case "TAB_ACTIVATED":
			case "TAB_CLOSED":
			case "LOGIN":
			case "PARTNER":
			case "TELEMETRY_INFO":
			case "SCREENSHOT":
			case "UNKNOWN":
			default:
			}

			if err != nil {
				log.Fatal(err)
			}
			requests = append(requests, payload)
		}
	}

	for i, request := range requests {
		// TODO: Do something with that request
		fmt.Println(i, request)
	}
}
