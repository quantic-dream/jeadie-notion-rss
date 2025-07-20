package main

import (
	"fmt"
	"log"
)

func main() {
	nDao, err := ConstructNotionDaoFromEnv()
	if err != nil {
		panic(fmt.Errorf("configuration error: %w", err))
	}

	// 🔄 Limpiar la base de datos antes de cargar nuevos feeds
	err = nDao.CleanUnstarredContentDatabase()
	if err != nil {
		log.Fatalf("Error cleaning content database: %v", err)
	}

	// 🧾 Ejecutar tareas RSS
	tasks := GetAllTasks()
	errs := make([]error, len(tasks))
	for i, t := range tasks {
		errs[i] = t.Run(nDao)
	}

	PanicOnErrors(errs)
}
