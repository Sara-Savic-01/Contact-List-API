package config
import(
	"encoding/json"
	"os"
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

var AppConfig Config

func LoadConfig(filename string) (*Config, error){
	file, err:=os.Open(filename)
	if err!=nil{
		return nil, err
	}
	defer file.Close()
	
	decoder:=json.NewDecoder(file)
	err=decoder.Decode(&AppConfig)
	if err!=nil{
		return nil,err
	}
	return &config, nil
}

	

