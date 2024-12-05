package plugins

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

type RedisPlugin struct {
	client *redis.Client
	app    *pocketbase.PocketBase
}

type event struct {
	Event      string      `json:"event"`
	Collection string      `json:"collection"`
	Record     interface{} `json:"record"`
	Source     string      `json:"source"`
}

func NewRedisPlugin(app *pocketbase.PocketBase, redisURL string) *RedisPlugin {
	client := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		app.Logger().Error("Failed to connect to Redis:", err)
		return nil
	}

	return &RedisPlugin{
		client: client,
		app:    app,
	}
}

func (rp *RedisPlugin) Register(app *pocketbase.PocketBase) {
	app.OnRecordCreate().BindFunc(func(e *core.RecordEvent) error {
		if err := e.Next(); err != nil {
			return err
		}
		source := e.Record.GetString("eventSource")
		if source == "" {
			source = "pocketbase"
		}
		return rp.publishEvent("create", e.Record.Collection().Name, e.Record, source)
	})

	app.OnRecordUpdate().BindFunc(func(e *core.RecordEvent) error {
		if err := e.Next(); err != nil {
			return err
		}
		source := e.Record.GetString("eventSource")
		if source == "" {
			source = "pocketbase"
		}
		return rp.publishEvent("update", e.Record.Collection().Name, e.Record, source)
	})

	app.OnRecordDelete().BindFunc(func(e *core.RecordEvent) error {
		if err := e.Next(); err != nil {
			return err
		}
		source := e.Record.GetString("eventSource")
		if source == "" {
			source = "pocketbase"
		}
		return rp.publishEvent("delete", e.Record.Collection().Name, e.Record, source)
	})
}

func (rp *RedisPlugin) publishEvent(eventType, collection string, record interface{}, source string) error {
	e := event{
		Event:      eventType,
		Collection: collection,
		Record:     record,
		Source:     source,
	}
	payload, err := json.Marshal(e)
	if err != nil {
		rp.app.Logger().Error("Failed to marshal event:", err)
		return err
	}

	rp.app.Logger().Info("Publishing Redis event:",
		"type", eventType,
		"collection", collection,
		"source", source,
	)

	ctx := context.Background()
	return rp.client.Publish(ctx, "pocketbase:events:publisher", payload).Err()
}

func (rp *RedisPlugin) ListenAndProcessEvents(app *pocketbase.PocketBase) {
	ctx := context.Background()
	sub := rp.client.Subscribe(ctx, "pocketbase:events:receiver")

	rp.app.Logger().Info("Started listening for Redis events")

	go func() {
		for msg := range sub.Channel() {
			rp.app.Logger().Debug("Received message:", msg.Payload)
			var e event
			if err := json.Unmarshal([]byte(msg.Payload), &e); err != nil {
				rp.app.Logger().Error("Failed to unmarshal event:", err)
				continue
			}

			if e.Source == "pocketbase" {
				rp.app.Logger().Debug("Skipping pocketbase sourced event")
				continue
			}

			rp.processEvent(app, e)
		}
	}()
}

func (rp *RedisPlugin) processEvent(app *pocketbase.PocketBase, e event) {
	ctx := context.WithValue(context.Background(), "eventSource", e.Source)

	switch e.Event {
	case "create":
		rp.createRecord(app, e.Collection, e.Record.(map[string]interface{}), ctx)
	case "update":
		rp.updateRecord(app, e.Collection, e.Record.(map[string]interface{}), ctx)
	case "delete":
		rp.deleteRecord(app, e.Collection, e.Record.(map[string]interface{}), ctx)
	}
}

func (rp *RedisPlugin) createRecord(app *pocketbase.PocketBase, collection string, record map[string]interface{}, ctx context.Context) {
	c, err := app.FindCollectionByNameOrId(collection)
	if err != nil {
		rp.app.Logger().Error("Collection not found:", err)
		return
	}

	newRecord := core.NewRecord(c)
	for key, value := range record {
		if key != "id" && key != "created" && key != "updated" &&
			key != "collectionId" && key != "collectionName" {
			newRecord.Set(key, value)
		}
	}

	source := ctx.Value("eventSource").(string)
	newRecord.Set("eventSource", source)

	if err := app.Save(newRecord); err != nil {
		rp.app.Logger().Error("Failed to save record:", err)
		return
	}

	rp.app.Logger().Info("Created record from Redis event:",
		"collection", collection,
		"source", source,
	)
}

func (rp *RedisPlugin) updateRecord(app *pocketbase.PocketBase, collection string, record map[string]interface{}, ctx context.Context) {
	id := record["id"].(string)
	data := record

	existingRecord, err := app.FindRecordById(collection, id)
	if err != nil {
		rp.app.Logger().Error("Record not found:", err)
		return
	}

	for key, value := range data {
		existingRecord.Set(key, value)
	}

	source := ctx.Value("eventSource").(string)
	existingRecord.Set("eventSource", source)

	if err := app.Save(existingRecord); err != nil {
		rp.app.Logger().Error("Failed to save record:", err)
		return
	}

	rp.app.Logger().Info("Updated record from Redis event:",
		"collection", collection,
		"id", id,
		"source", source,
	)
}

func (rp *RedisPlugin) deleteRecord(app *pocketbase.PocketBase, collection string, record map[string]interface{}, ctx context.Context) {
	id := record["id"].(string)

	existingRecord, err := app.FindRecordById(collection, id)
	if err != nil {
		rp.app.Logger().Error("Record not found:", err)
		return
	}

	if err := app.Delete(existingRecord); err != nil {
		rp.app.Logger().Error("Failed to delete record:", err)
		return
	}

	rp.app.Logger().Info("Deleted record from Redis event:",
		"collection", collection,
		"id", id,
		"source", ctx.Value("eventSource").(string),
	)
}
