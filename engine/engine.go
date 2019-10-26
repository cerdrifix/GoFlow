package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

type ProcessEngine struct {
	logger  *log.Logger
	db      *sqlx.DB
	context context.Context
}

func New(logger *log.Logger, db *sqlx.DB, context context.Context) *ProcessEngine {
	return &ProcessEngine{
		logger:  logger,
		db:      db,
		context: context,
	}
}

func (engine *ProcessEngine) GetProcessMap(name string) (processMap ProcessMap, err error) {

	q := fmt.Sprintf("SELECT * FROM public.fn_maps_getlatestbyname('%s')", name)

	rows, err := engine.db.QueryxContext(engine.context, q)
	if err != nil {
		engine.logger.Printf("Unable to query db")
		return processMap, err
	}
	for rows.Next() {
		var name string
		var version int
		var data string

		err = rows.Scan(&name, &version, &data)
		if err != nil {
			engine.logger.Printf("Error during row scan")
			return processMap, err
		}

		engine.logger.Printf("\nMap name: %s\n    version: %d\n    data: %s", name, version, data)

		err = json.Unmarshal([]byte(data), &processMap)
		if err != nil {
			engine.logger.Printf("Error during unmarshaling \njson: %s\nerror: %v", data, err)
			return processMap, err
		}
	}

	engine.logger.Printf("Json parsed: %#v", processMap)
	return processMap, nil
}

func (engine *ProcessEngine) NewInstance(payload CreateProcessPayload) (instanceNumber int, err error) {

	m, err := engine.GetProcessMap(payload.ProcessName)
	if err != nil {
		engine.logger.Fatalf("Error retrieving map from database: %v", err)
	}

	fmt.Println(m)

	fmt.Printf("Variables: %#v", payload.Variables)

	for k, v := range payload.Variables {

		switch t := v.(type) {
		case string:
			dt, err := time.Parse("2006-01-02T15:04:05.000", v.(string))
			if err == nil {
				payload.Variables[k] = dt
				fmt.Println(k, dt, "(datetime)")
			} else {
				payload.Variables[k] = v
				fmt.Println(k, v, "(string)")
			}
		case float64:
			payload.Variables[k] = v.(float64)
			fmt.Println(k, v, "(float64)")
		case int:
			payload.Variables[k] = v.(int)
			fmt.Println(k, v, "(int)")
		case []interface{}:
			fmt.Println(k, "(array):")
			for i, u := range t {
				fmt.Println("    ", i, u)
			}
		default:
			fmt.Println(k, v, "(unknown)")
		}
	}

	return 0, nil
}
