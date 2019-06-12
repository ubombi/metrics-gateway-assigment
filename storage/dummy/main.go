package dummy

import (
	"fmt"
	"io"
	"os"

	"github.com/ubombi/timeseries/storage"
)

type Storage struct {
	Output io.Writer
}

var Stdout storage.Interface = &Storage{
	Output: os.Stdout,
}

func (s *Storage) Store(e storage.Event) error {
	fmt.Println(e)
	return nil
}
