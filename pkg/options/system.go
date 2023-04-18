package options

type System struct {
	Env             string `mapstructure:"env" json:"env" yaml:"env"`                                  // 环境值
	Addr            int    `mapstructure:"addr" json:"addr" yaml:"addr"`                               // 端口值
	DbType          string `mapstructure:"db-type" json:"db-type" yaml:"db-type"`                      // 数据库类型:mysql(默认)|sqlite|sqlserver|postgresql
	UseRedis        bool   `mapstructure:"use-redis" json:"use-redis" yaml:"use-redis"`                // 使用redis
	UseMultipoint   bool   `mapstructure:"use-multipoint" json:"use-multipoint" yaml:"use-multipoint"` // 多点登录拦截
	LimitCountIP    int64  `mapstructure:"iplimit-count" json:"iplimit-count" yaml:"iplimit-count"`
	LimitIntervalIP int64  `mapstructure:"iplimit-interval" json:"iplimit-interval" yaml:"iplimit-interval"`
	RouterPrefix    string `mapstructure:"router-prefix" json:"router-prefix" yaml:"router-prefix"`
}
