package main

import (
	"fmt"
	"os"
	"regexp"
	"bufio"
)

func main() {

/*
 * from  https://github.com/StefanSchroeder/Golang-Regex-Tutorial/blob/master/01-chapter3.markdown
 *
 * Interestingly the RFC 2822 which defines the format of
 * email-addresses is pretty permissive. That makes it hard to come up
 * with a simple regular expression that matches a valid email
 * address. In most cases though your application can make some
 * assumptions about addresses and I found this one sufficient for all
 * practical purposes:
 *
 * (\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3})
 *
 * It must start with a character of the \w class. Then we can have
 * any number of characters including the hyphen, the '.' and the
 * underscore. We want the last character before the @ to be a
 * 'regular' character again. We repeat the same pattern for the
 * domain, only that the suffix (part behind the last dot) can be only
 * 2 or 3 characters. This will cover most cases. If you come across
 * an email address that does not match this regexp it has probably
 * deliberately been setup to annoy you and you can therefore ignore
 * it.
*/

/*
 *  According to https://golang.org/pkg/regexp/syntax/
 *
 *  \w             word characters (== [0-9A-Za-z_])
 *
 */

	// regex, err := regexp.Compile("\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3}")
	regex, err := regexp.Compile("[0-9A-Za-z_][-.0-9A-Za-z_]*[0-9A-Za-z_]@[0-9A-Za-z_][-.0-9A-Za-z_]*[0-9A-Za-z_][.][0-9A-Za-z_]{2,3}")
	if err != nil {
		fmt.Println("regular expression error", err.Error());
		return
	}

	var text string

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("email addr: ")
		text, _ = reader.ReadString('\n')
		text = text[:len(text)-1] // chop off '\n'
		if regex.MatchString(text) {
			fmt.Println("ok")
		} else {
			fmt.Println("not ok")
		}
	}
}

