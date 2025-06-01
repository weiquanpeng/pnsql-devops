package global

type Postgresql struct {
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Dbname   string `mapstructure:"dbname" json:"dbname" yaml:"dbname"`
}
