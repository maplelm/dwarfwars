package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/maplelm/dwarfwars/pkg/settings"

	_ "github.com/go-sql-driver/mysql"
)

func ValidateSQLServers(servers []settings.SQLServerSettings) (err error) {
dbloop:
	for _, s := range servers {
		//connstr, err := s.ConnString()
		//if err != nil {
		//	return err
		//}
		conn, err := sql.Open(s.Driver, fmt.Sprintf("%s:%s@tcp(%s:%d)/", s.Username, s.Password, s.Addr, s.Port))
		if err != nil {
			return err
		}
		defer conn.Close()

		// Check to make sure the database exists

		log.Printf("(SQL Database Validation) Running Query: CREATE DATABASE IF NOT EXISTS '%s' ;", s.DefaultDB)
		_, err = conn.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", s.DefaultDB))
		if err != nil {
			return err
		}
		// Check to make sure all tables that are expected exist

		for tn, tc := range s.Tables {
			q := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s (", s.DefaultDB, tn)
			pk := ""
			if l := len(tc); l <= 0 {
				log.Printf("(SQL Validation) Adding %d columns to table creation query, did not create table %s", len(tc), fmt.Sprintf("%s.%s", s.DefaultDB, tn))
				continue dbloop

			}
			for i, each := range tc {
				q += fmt.Sprintf("%s %s ", each.Name, each.ColType)
				if each.AutoIncrement {
					q += "auto_increment "
				}
				if strings.Contains(each.ColType, "TEXT") || strings.Contains(each.ColType, "VARCHAR") {
					q += "CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci "
				}
				if each.DefaultValue != nil && len(*each.DefaultValue) > 0 {
					q += fmt.Sprintf(" DEFAULT %s ", *each.DefaultValue)
				}
				if each.Nullable {
					q += "NOT NULL"
				}
				//if each.PrimaryKey && !strings.Contains(each.ColType, "TEXT") && !strings.Contains(each.ColType, "BLOB") {
				if each.PrimaryKey {
					log.Printf("(SQL Validation) Primary Key For %s.%s is %s, \n\tPrimary Key: %t\n\tColumn Type: %s\n\tDefault Value: %v\n\tForiegn Key: %t\n\tAuto Increment: %t\n\tNullable: %t",
						s.DefaultDB,
						tn,
						each.Name,
						each.PrimaryKey,
						each.ColType,
						each.DefaultValue,
						each.ForiegnKey,
						each.AutoIncrement,
						each.Nullable,
					)
					pk = each.Name
				}
				if i < len(tc)-1 {
					q += ",\n"
				} else if len(pk) > 0 {
					q += fmt.Sprintf(",\nCONSTRAINT %s_%s PRIMARY KEY (%s)\n)", tn, each.Name, pk)
				} else {
					q += "\n)"
				}

			}
			q += "\nENGINE=InnoDB\nDEFAULT CHARSET=utf8mb4\nCOLLATE=utf8mb4_0900_ai_ci;"

			log.Printf("(SQL Validation) Running Query %s", q)
			_, err := conn.Exec(q)
			if err != nil {
				return err
			}
		}
	}
	return
}
