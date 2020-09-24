package dal

import "os"

var databaseName = os.Getenv("SWITCH_DB_NAME")
var databaseUserName = os.Getenv("SWITCH_DB_USER")
var databasePassword = os.Getenv("SWITCH_DB_PASSWORD")
var databaseHost = os.Getenv("SWITCH_DB_HOST")
var databasePort = os.Getenv("SWITCH_DB_PORT")
