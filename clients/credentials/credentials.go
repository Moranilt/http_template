package credentials

import "fmt"

type SourceStringer interface {
	SourceString() string
}

type DB struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Host     string `mapstructure:"host"`
	SSLMode  string `mapstructure:"sslmode"`
}

func (d *DB) SourceString(production bool) string {
	if !production {
		return fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s sslmode=disable",
			d.Username, d.Password, d.DBName, d.Host,
		)
	}
	return fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s",
		d.Username, d.Password, d.DBName, d.Host,
	)
}

type RabbitMQ struct {
	Host     string `mapstructure:"host"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func (r *RabbitMQ) SourceString() string {
	return fmt.Sprintf("amqp://%s:%s@%s/", r.Username, r.Password, r.Host)
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
}
