package engine

import (
	"GoFlow/routines"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"sync"
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

func transformVariablesInJSONPayload(variables map[string]interface{}) (string, error) {

	fmt.Printf("Variables: %#v", variables)

	for k, v := range variables {

		switch t := v.(type) {
		case string:
			dt, err := time.Parse("2006-01-02T15:04:05.000", v.(string))
			if err == nil {
				variables[k] = dt
				fmt.Println(k, dt, "(datetime)")
			} else {
				variables[k] = v
				fmt.Println(k, v, "(string)")
			}
		case float64:
			variables[k] = v.(float64)
			fmt.Println(k, v, "(float64)")
		case int:
			variables[k] = v.(int)
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

	j, err := json.Marshal(variables)

	return string(j), err
}

func doProcessEvent(e ProcessEvent, variables *map[string]interface{}, engine *ProcessEngine, waitGroup *sync.WaitGroup) error {
	fmt.Printf("\n\nProcessing event: %#v", e)

	r := routines.New(engine.logger, variables)
	params := make([]interface{}, len(e.Parameters))

	for i, p := range e.Parameters {
		params[i] = p.Value
	}

	switch e.EventType {
	case "function":
		_, err := r.CallFunc(e.Name, params)

		if err != nil {
			return err
		}
	case "validator":
		_, err := r.CallFunc(e.Name, params)

		if err != nil {
			return err
		}
	}
	waitGroup.Done()
	return nil
}

func (engine *ProcessEngine) NewInstance(payload CreateProcessPayload) (instanceNumber int, err error) {

	tx, err := engine.db.Beginx()
	if err != nil {
		engine.logger.Fatalf("Error beginning transaction: %v", err)
		return 0, err
	}

	pmap, err := engine.GetProcessMap(payload.ProcessName)
	if err != nil {
		tx.Rollback()
		engine.logger.Fatalf("Error retrieving map from database: %v", err)
	}
	fmt.Printf("Map: %#v", pmap)

	engine.logger.Printf("Map: %s", pmap)

	var startNode ProcessNode
	var found bool

	for _, n := range pmap.Nodes {
		if n.NodeType == "start" {
			startNode = n
			found = true
			break
		}
	}

	if found == false {
		tx.Rollback()
		engine.logger.Fatalf("Error retrieving start node")
	}

	engine.logger.Printf("Start node found! %#v", startNode)

	// Processing pre and post events
	pre := startNode.Events.Pre
	post := startNode.Events.Post

	var eventsWG sync.WaitGroup
	eventsWG.Add(len(pre))
	errs := make([]error, 0, len(pre)+len(post))

	for _, ev := range pre {
		go func() {
			err = doProcessEvent(ev, &payload.Variables, engine, &eventsWG)
			if err != nil {
				_ = append(errs, err)
			}
		}()
	}
	eventsWG.Wait()

	if len(errs) > 0 {
		tx.Rollback()
		engine.logger.Fatalf("Errors occured during pre-events: %#v", errs)
		return
	}

	eventsWG.Add(len(post))
	for _, ev := range post {
		go func() {
			err = doProcessEvent(ev, &payload.Variables, engine, &eventsWG)
			if err != nil {
				errs[len(errs)] = err
			}
		}()
	}
	eventsWG.Wait()

	if len(errs) > 0 {
		engine.logger.Fatalf("Errors occured during post-events: %#v", errs)
		return
	}

	j, err := transformVariablesInJSONPayload(payload.Variables)

	fmt.Printf("JSON Variables: %s", string(j))

	tx.Commit()

	return 0, nil
}
