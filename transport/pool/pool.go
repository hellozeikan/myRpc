package pool
type Pool interface {
	Get( address string) (*Conn, error)
}
func GetPool(){}

func Get(){}

func Pull(){}