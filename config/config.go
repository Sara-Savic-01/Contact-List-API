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

	

