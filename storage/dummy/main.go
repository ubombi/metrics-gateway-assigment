package dummy

import (
	"fmt"
	"io"
	"os"

	"github.com/ubombi/timeseries/api"
)

type Storage struct {
	Output io.Writer
}

var Stdout = &Storage{
	Output: os.Stdout,
}

func (s *Storage) Store(e api.Event) error {
	fmt.Println(e)
	return nil
}
