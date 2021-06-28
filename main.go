package main

import (
  /* debug */
  // "fmt"
  "log"

  /* config */
  "github.com/ilyakaznacheev/cleanenv"

  /* mongo */
  "context"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"

  /* webserver */
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

func load_cfg() {
  cleanenv.ReadConfig("config.yml", &cfg)
}

var collection *mongo.Collection
var collections = make(map[string]*mongo.Collection)
var ctx = context.TODO()

func db_init(uri string) {
  clientOptions := options.Client().ApplyURI(uri)

  client, err := mongo.Connect(ctx, clientOptions)
  if err != nil {
    log.Fatal(err)
  }

  err = client.Ping(ctx, nil)
  if err != nil {
    log.Fatal(err)
  }

  collection = client.Database("tracing").Collection("logs")
  for _, source := range(cfg.Sources) {
    collections[source] = client.Database("tracing").Collection(source)
  }
}

func main() {
  load_cfg()
  db_init(cfg.Mongo.URI)
}
