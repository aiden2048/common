package config
import(
	"github.com/aiden2048/common/rootConfig"
)
func LoadConfig() error {
	//加载配置
	rootConfig.LoadAppInfoList(0) //app game 配置
	return nil
}
