package main

import (
  "fmt"
  "github.com/ilyakaznacheev/cleanenv"

  "context"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

type ConfigStruct struct {
  Mongo struct {
    Port        string `yaml:"port"`
    Host        string `yaml:"host"`
    User        string `yaml:"dbuser"`
    Password    string `yaml:"dbpassword"`
    Database    string `yaml:"dbname"`
    Extra       string `yaml:"extra"`
    URI         string `yaml:"uri"`
  } `yaml:"mongo"`

  Sources     []string `yaml:"sources"`
  DefaultSource string `yaml:"default_source"`
}
var cfg ConfigStruct

func get_uri() string {
  err := cleanenv.ReadConfig("config.yml", &cfg)
  if err != nil {
    panic(err)
  }
  return cfg.Mongo.URI
}

var collection *mongo.Collection
var collections = make(map[string]*mongo.Collection)
var ctx = context.TODO()

func db_init(uri string) {
  clientOptions := options.Client().ApplyURI(uri)

  client, err := mongo.Connect(ctx, clientOptions)
  if err != nil {
    fmt.Println("error 1")
    panic(err)
    // log.Fatal(err)
  }

  err = client.Ping(ctx, nil)
  if err != nil {
    fmt.Println("error 2")
    panic(err)
    // log.Fatal(err)
  }

  fmt.Println("here")

  collection = client.Database("tracing").Collection("logs")
  for _, source := range(cfg.Sources) {
    collections[source] = client.Database("tracing").Collection(source)
    fmt.Println(source, "loaded")
  }
}

func main() {
  get_uri()
  db_init(cfg.Mongo.URI)
}
