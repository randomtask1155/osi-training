package colors 

import(
  "fmt"
)

var (
  
  // Colors is a list of bash color codes
  Colors = []string{}
)

// ColorizeError conver error message into red string
func ColorizeError(err error) string {
  return fmt.Sprintf("\033[91m%s\033[0m", err)
}

// ColorizeText takes string and wraps it in color terminators
func ColorizeText(s string) string {
  return fmt.Sprintf("\033[94m%s\033[0m", s)
}