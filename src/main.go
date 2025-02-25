package main

import (
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/sergi/go-diff/diffmatchpatch"
)


func main() {
	app := pocketbase.New()

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// serves static files from the provided public dir (if exists)
		se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))

        return se.Next()
	})
	app.OnRecordCreateRequest("note", "task").BindFunc(func(e *core.RecordRequestEvent) error {
		note := e.Record
		note.Set("user", e.Auth.Id)
		return e.Next()
	})

	app.OnRecordAfterCreateSuccess("task").BindFunc(func(e *core.RecordEvent) error {
		task := e.Record

		collection, err := app.FindCollectionByNameOrId("taskHistory")
		if err != nil {
			return err
		}

		record := core.NewRecord(collection)
		record.Set("task", task.Id)
		record.Set("startDate", task.Get("startDate"))
		record.Set("endDate", task.Get("endDate"))
		record.Set("isAllDay", task.Get("isAllDay"))
		record.Set("status", task.Get("status"))
		record.Set("name", task.Get("name"))
		record.Set("description", task.Get("description"))
		record.Set("user", task.Get("user"))

		if err := app.Save(record); err != nil {
			return err
		}

		return e.Next()
	})
	app.OnRecordAfterUpdateSuccess("task").BindFunc(func(e *core.RecordEvent) error {
		task := e.Record

		collection, err := app.FindCollectionByNameOrId("taskHistory")
		if err != nil {
			return err
		}

		record := core.NewRecord(collection)
		record.Set("task", task.Id)
		record.Set("startDate", task.Get("startDate"))
		record.Set("endDate", task.Get("endDate"))
		record.Set("isAllDay", task.Get("isAllDay"))
		record.Set("status", task.Get("status"))
		record.Set("name", task.Get("name"))
		record.Set("description", task.Get("description"))
		record.Set("user", task.Get("user"))

		if err := app.Save(record); err != nil {
			return err
		}

		return e.Next()
	})

	app.OnRecordUpdateRequest("note").BindFunc(func(e *core.RecordRequestEvent) error {
		note := e.Record
		oldNote, err := app.FindRecordById("note", note.GetString("id"))
		if err != nil {
			return err
		}

		dmp := diffmatchpatch.New()
		oldContent := oldNote.Get("content").(string)
		newContent := note.Get("content").(string)

		// Patches are made in reverse order,
		// so that the history can be computed from the current content
		patches := dmp.PatchToText(dmp.PatchMake(newContent, oldContent))

		newTitle := note.Get("title")

		if patches != "" || newTitle != "" {
			collection, err := app.FindCollectionByNameOrId("noteChange")
			if err != nil {
				return err
			}

			record := core.NewRecord(collection)
			record.Set("note", note.Id)
			record.Set("user", note.Get("user"))
			record.Set("contentPatch", patches)
			record.Set("title", newTitle)

			if err := app.Save(record); err != nil {
				return err
			}
		}

		return e.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

