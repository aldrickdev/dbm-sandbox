package providers

const (
	postgres  = "Postgres"
	mysql     = "MySQL"
	sqlserver = "SQL Server"
)

// A DBMS is a struct that contains the details for a particular database
// management system that we support.
type DBMS struct {
	// Name is the name of the DBMS and the value that will used as an option 
  // when prompting the user for what DBMS they would like to use.
	Name     string

  // versions holds all of the versions that we support for the particular 
  // DBMS.
	versions []string
}

// Helper function to create DBMS's 
func newDBMS(name string, versions []string) DBMS {
	return DBMS{
		Name:     name,
		versions: versions,
	}
}

// Returns the concrete DBMS implementation based on the provided input
func GetDBMS(DBMSName string) DBMS {
	switch DBMSName {
	case postgres:
		return PostgresDBMS()

	case mysql:
		return MySQLDBMS()

	case sqlserver:
		return SQLServerDBMS()

	default:
		return PostgresDBMS()
	}

}

// Returns the Postgres DBMS concrete implementation
// 
// Available Versions: https://hub.docker.com/_/postgres/tags
//
// Supported Versions: https://docs.datadoghq.com/database_monitoring/setup_postgres/selfhosted/?tab=postgres15
func PostgresDBMS() DBMS {
	return newDBMS(postgres, []string{"16", "15", "14", "13", "12"})
}

// Returns the MySQL DBMS concrete implementation
// 
// Available Versions: https://hub.docker.com/_/mysql/tags
//
// Supported Versions: https://docs.datadoghq.com/database_monitoring/setup_mysql/selfhosted/?tab=mysql57
func MySQLDBMS() DBMS {
	return newDBMS(mysql, []string{"8.0.37", "8.0.36", "8.0.35", "8.0.34", "8.0.33"})
}

// Returns the MySQL DBMS concrete implementation
// 
// Available Versions: https://hub.docker.com/_/microsoft-mssql-server
// 
// Supported Versions: https://docs.datadoghq.com/database_monitoring/setup_sql_server/selfhosted/?tab=sqlserver2014
func SQLServerDBMS() DBMS {
	return newDBMS(sqlserver, []string{"2022-latest", "2019-latest", "2017-latest"})
}
