package strings

import (
	"fmt"
	"math"
	"strings"
)

// ConvertToCurrency will convert a float into human-readable currency string. The delimiters are
// optional but by will default will convert 123567.89 to format of 123,456.69
func ConvertToCurrency(value float64, delimiters ...string) string {
	var dollarDelimiter string = ","
	var centDelimiter string = "."
	if len(delimiters) > 0 {
		dollarDelimiter = delimiters[0]
	}
	if len(delimiters) == 2 {
		centDelimiter = delimiters[1]
	}

	d, c := math.Modf(value)
	dollars := fmt.Sprintf("%f", d)
	cents := fmt.Sprintf("%f", c)
	dollars = dollars[0:strings.Index(dollars, ".")] // get string before `.`
	cents = cents[strings.Index(cents, ".")+1 : 4]   // get string after `.`
	str := dollars

	var newValue string
	for {
		if len(str) > 3 {
			newValue = str[len(str)-3:] + newValue
			str = str[:len(str)-3]
			if len(str) > 0 {
				newValue = dollarDelimiter + newValue
			}
		} else {
			newValue = str + newValue
			break
		}
	}

	v := fmt.Sprintf("%s%s%s", newValue, centDelimiter, cents)

	return v
}
