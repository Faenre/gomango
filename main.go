package main

import (
  /* debug */
  "fmt"
  "log"

  /* config */
  "github.com/ilyakaznacheev/cleanenv"

  /* mongo */
  "context"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"

  /* webserver */
  "net/http"
  "io/ioutil"
)

func output(contents ...string) {
  if (debug) { fmt.Println(contents) }
}

type TraceData struct {
  Headers map[string][]string
  Collection string
  Content string
}

/* load configs from config.yml */
type ConfigStruct struct {
  Mongo struct {
    // Port        string `yaml:"port"`
    // Host        string `yaml:"host"`
    // User        string `yaml:"dbuser"`
    // Password    string `yaml:"dbpassword"`
    // Database    string `yaml:"dbname"`
    // Extra       string `yaml:"extra"`
    URI         string  `yaml:"uri"`
  } `yaml:"mongo"`

  SourceHeader  string  `yaml:"source_header"`
  Sources     []string  `yaml:"sources"`
  DefaultSource string  `yaml:"default_source"`

  Debug         bool    `yaml:"debug"`
  HTTPPort      string  `yaml:"http_port"`
  HTTPSPort     string  `yaml:"https_port"`

  FullCert      string  `yaml:"fullcert"`
  PrivateKey    string  `yaml:"privatekey"`
}
var debug bool = false;
var cfg ConfigStruct
func load_cfg() {
  cleanenv.ReadConfig("config.yml", &cfg)
  debug = cfg.Debug
}

/* connect to db */
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

func post_to_db(collection *mongo.Collection, td TraceData) {
  if _, err := collection.InsertOne(context.TODO(), td); err != nil {
    output("Error writing!")
  } else {
    output("Write success, into:", td.Collection)
  }
}

/* web server */
func get_source_from_headers(headers map[string][]string) string {
  header := headers[cfg.SourceHeader]
  if headers == nil { return cfg.DefaultSource }
  if len(header) == 0 { return cfg.DefaultSource }

  source := header[0]
  _, ok := collections[source]
  if !ok { return cfg.DefaultSource }

  return source
}

func form_handler(w http.ResponseWriter, r *http.Request) {
  if err := r.ParseForm(); err != nil { return }

  enableCors(&w)

  headers := headers_to_map(r.Header)

  collectionName := get_source_from_headers(headers)
  collection := collections[collectionName]

  buf, _ := ioutil.ReadAll(r.Body)
  content := fmt.Sprintf("%q", buf[:])

  td := TraceData{headers, collectionName, content}

  post_to_db(collection, td)
}

func enableCors(w *http.ResponseWriter) {
  allowedHeaders := "*" // "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token, Tracer-Source"
  (*w).Header().Set("Access-Control-Allow-Origin", "*")
  (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
  (*w).Header().Set("Access-Control-Allow-Headers", allowedHeaders)
  (*w).Header().Set("Access-Control-Expose-Headers", "Authorization")
}

func headers_to_map(headers http.Header) map[string][]string {
  m := make(map[string][]string)
  for h, va := range(headers) {
    for _, v := range(va) { m[h] = append(m[h], v) }
  }
  return m
}

/* orchestrate */
func main() {
  load_cfg()

  output("Connecting to mongoDB...")
  db_init(cfg.Mongo.URI)

  output("Registering form func...")
  http.HandleFunc("/tracelog", form_handler)

  // output("Now listening on:", cfg.HTTPPort)
  // if err := http.ListenAndServe(cfg.HTTPPort, nil); err != nil {
  //   log.Fatal(err)
  // }
  output("Now listening on:", cfg.HTTPSPort)
  if err := http.ListenAndServeTLS(cfg.HTTPSPort, cfg.FullCert, cfg.PrivateKey, nil); err != nil {
    log.Fatal(err)
  }
}
