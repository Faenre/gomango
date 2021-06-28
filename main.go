package main

import (
  "fmt"
  "github.com/ilyakaznacheev/cleanenv"
)


type ConfigDatabase struct {
    // Server struct {
    //     Port string `yaml:"port"`
    //     Host string `yaml:"host"`
    // } `yaml:"server"`
    // Database struct {
    //     Username string `yaml:"user"`
    //     Password string `yaml:"pass"`
    // } `yaml:"database"`
  Mongo struct {
    Port      string `yaml:"port"`
    Host      string `yaml:"host"`
    User      string `yaml:"dbuser"`
    Password  string `yaml:"dbpassword"`
    Database  string `yaml:"dbname"`
    Extra     string `yaml:"extra"`
  } `yaml:"mongo"`
  URI string `yaml:"DB_URI"`
}


func main() {
  var cfg ConfigDatabase
  err := cleanenv.ReadConfig("config.yml", &cfg)
  if err == nil {
    fmt.Println(cfg.Mongo.Port)
    fmt.Println(cfg.URI)
  } else {
    fmt.Println(err)
  }
}

// func init() {

//   clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
//   client, err := mongo.Connect(ctx, clientOptions)

//   if err != nil {
//     log.Fatal(err)
//   }
// }
