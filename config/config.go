package config
import(
	"encoding/json"
	"os"
	"log"
)

type DBConfig struct{
	
	User string `json:"user"`
	Password string `json:"password"`
	Host string `json:"host"`
	Name string `json:"name"`
	
	
}

type Config struct{
	DB DBConfig `json:"db"`
	AuthToken string `json:"auth_token"`
}
type ConfigTest struct{
	DB DBConfig `json:"db"`
}

func LoadConfig(filename string) (*Config, error){
	file, err:=os.Open(filename)
	if err!=nil{
		return nil, err
	}
	defer file.Close()
	var config Config
	decoder:=json.NewDecoder(file)
	err=decoder.Decode(&config)
	if err!=nil{
		return nil,err
	}
	return &config, nil
}
func LoadTestConfig(filename string) *ConfigTest{
	file, err:=os.Open(filename)
	if err!=nil{
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()
	
	decoder:=json.NewDecoder(file)
	config:=&ConfigTest{}
	err=decoder.Decode(config)
	if err!=nil{
		log.Fatalf("Failed to decode config file: %v", err)
	}
	return config
}

	

