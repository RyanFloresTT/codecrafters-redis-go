package rdb

type Config struct {
	Dir        string `json:"dir"`
	DbFileName string `json:"db_file_name"`
}

var RDBConfig = Config{
	Dir:        "/var/lib/redis",
	DbFileName: "dump.rdb",
}

func LoadConfig() Config {
	return RDBConfig
}
