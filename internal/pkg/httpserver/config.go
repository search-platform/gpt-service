package httpserver

type Config struct {
	Port    int  `env:"PORT"`
	Prefork bool `env:"HTTP_PREFORK"`

	AllowOrigins string `env:"CORS_ALLOW_ORIGINS"`

	UseProtobuf   bool `env:"HTTP_PROTOBUF" envDefault:"true"`
	UseProtoNames bool `env:"HTTP_PROTOBUF_USE_PROTONAME"`
}
