package main

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/urfave/cli/v2"
)

var longReasons = map[int]string{
	100: `The initial part of a request has been received and has not yet been rejected by
the server. The server intends to send a final response after the request has
been fully received and acted upon.`,
	101: `The server understands and is willing to comply with the client's request, via
the Upgrade header field1, for a change in the application protocol being used
on this connection.`,
	102: `An interim response used to inform the client that the server has accepted the
complete request, but has not yet completed it.`,
	200: `The request has succeeded.`,
	201: `The request has been fulfilled and has resulted in one or more new resources
being created.`,
	202: `The request has been accepted for processing, but the processing has not been
completed. The request might or might not eventually be acted upon, as it might
be disallowed when processing actually takes place.`,
	203: `The request was successful but the enclosed payload has been modified from that
of the origin server's 200 OK response by a transforming proxy1.`,
	204: `The server has successfully fulfilled the request and that there is no
additional content to send in the response payload body.`,
	205: `The server has fulfilled the request and desires that the user agent reset the
"document view", which caused the request to be sent, to its original state as
received from the origin server.`,
	206: `The server is successfully fulfilling a range request for the target resource by
transferring one or more parts of the selected representation that correspond to
the satisfiable ranges found in the request's Range header field1.`,
	207: `A Multi-Status response conveys information about multiple resources in
situations where multiple status codes might be appropriate.`,
	208: `Used inside a DAV: propstat response element to avoid enumerating the internal
members of multiple bindings to the same collection repeatedly.`,
	226: `The server has fulfilled a GET request for the resource, and the response is a
representation of the result of one or more instance-manipulations applied to
the current instance.`,
	300: `The target resource has more than one representation, each with its own more
specific identifier, and information about the alternatives is being provided so
that the user (or user agent) can select a preferred representation by
redirecting its request to one or more of those identifiers.`,
	301: `The target resource has been assigned a new permanent URI and any future
references to this resource ought to use one of the enclosed URIs.`,
	302: `The target resource resides temporarily under a different URI. Since the
redirection might be altered on occasion, the client ought to continue to use
the effective request URI for future requests.`,
	303: `The server is redirecting the user agent to a different resource, as indicated
by a URI in the Location header field, which is intended to provide an indirect
response to the original request.`,
	304: `A conditional GET or HEAD request has been received and would have resulted in a
200 OK response if it were not for the fact that the condition evaluated to
false.`,
	305: `Defined in a previous version of this specification and is now deprecated, due
to security concerns regarding in-band configuration of a proxy.`,
	307: `The target resource resides temporarily under a different URI and the user agent
MUST NOT change the request method if it performs an automatic redirection to
that URI.`,
	308: `The target resource has been assigned a new permanent URI and any future
references to this resource ought to use one of the enclosed URIs.`,
	400: `The server cannot or will not process the request due to something that is
perceived to be a client error (e.g., malformed request syntax, invalid request
message framing, or deceptive request routing).`,
	401: `The request has not been applied because it lacks valid authentication
credentials for the target resource.`,
	402: `Reserved for future use.`,
	403: `The server understood the request but refuses to authorize it.`,
	404: `The origin server did not find a current representation for the target resource
or is not willing to disclose that one exists.`,
	405: `The method received in the request-line is known by the origin server but not
supported by the target resource.`,
	406: `The target resource does not have a current representation that would be
acceptable to the user agent, according to the proactive negotiation header
fields received in the request1, and the server is unwilling to supply a default
representation.`,
	407: `Similar to 401 Unauthorized, but it indicates that the client needs to
authenticate itself in order to use a proxy.`,
	408: `The server did not receive a complete request message within the time that it
was prepared to wait.`,
	409: `The request could not be completed due to a conflict with the current state of
the target resource. This code is used in situations where the user might be
able to resolve the conflict and resubmit the request.`,
	410: `The target resource is no longer available at the origin server and that this
condition is likely to be permanent.`,
	411: `The server refuses to accept the request without a defined Content-Length1.`,
	412: `One or more conditions given in the request header fields evaluated to false
when tested on the server.`,
	413: `The server is refusing to process a request because the request payload is
larger than the server is willing or able to process.`,
	414: `The server is refusing to service the request because the request-target1 is
longer than the server is willing to interpret.`,
	415: `The origin server is refusing to service the request because the payload is in a
format not supported by this method on the target resource.`,
	416: `None of the ranges in the request's Range header field1 overlap the current
extent of the selected resource or that the set of ranges requested has been
rejected due to invalid ranges or an excessive request of small or overlapping
ranges.`,
	417: `The expectation given in the request's Expect header field1 could not be met by
at least one of the inbound servers.`,
	422: `The server understands the content type of the request entity (hence a 415
Unsupported Media Type status code is inappropriate), and the syntax of the
request entity is correct (thus a 400 Bad Request status code is inappropriate)
but was unable to process the contained instructions.`,
	423: `The source or destination resource of a method is locked.`,
	424: `The method could not be performed on the resource because the requested action
depended on another action and that action failed.`,
	426: `The server refuses to perform the request using the current protocol but might
be willing to do so after the client upgrades to a different protocol.`,
	428: `The origin server requires the request to be conditional.`,
	429: `The user has sent too many requests in a given amount of time ("rate limiting").`,
	431: `The server is unwilling to process the request because its header fields are too
large. The request MAY be resubmitted after reducing the size of the request
header fields.`,
	500: `The server encountered an unexpected condition that prevented it from fulfilling
the request.`,
	501: `The server does not support the functionality required to fulfill the request.`,
	502: `The server, while acting as a gateway or proxy, received an invalid response
from an inbound server it accessed while attempting to fulfill the request.`,
	503: `The server is currently unable to handle the request due to a temporary overload
or scheduled maintenance, which will likely be alleviated after some delay.`,
	504: `The server, while acting as a gateway or proxy, did not receive a timely
response from an upstream server it needed to access in order to complete the
request.`,
	505: `The server does not support, or refuses to support, the major version of HTTP
that was used in the request message.`,
	506: `The server has an internal configuration error: the chosen variant resource is
configured to engage in transparent content negotiation itself, and is therefore
not a proper end point in the negotiation process.`,
	507: `The method could not be performed on the resource because the server is unable
to store the representation needed to successfully complete the request.`,
	508: `The server terminated an operation because it encountered an infinite loop while
processing a request with "Depth: infinity". This status indicates that the
entire operation failed.`,
	510: `The policy for accessing the resource has not been met in the request. The
server should send back all the information necessary for the client to issue an
extended request.`,
	511: `The client needs to authenticate to gain network access.`,
}

func printRange(n int, printLong bool) {
	keys := make([]int, 0, len(longReasons))
	for k := range longReasons {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, code := range keys {
		if code/100 == n {
			fmt.Printf("%d %s\n", code, http.StatusText(code))
			if printLong {
				fmt.Printf("%s\n\n", longReasons[code])
			}
		}
	}
}

func printCode(c *cli.Context, printLong bool) error {
	if code := c.Args().First(); code != "" {
		n, err := strconv.Atoi(code)
		if err != nil {
			return fmt.Errorf("'%s' is not a valid HTTP status code: not a number", code)
		}
		if http.StatusText(n) == "" {
			return fmt.Errorf("'%d' is not a valid HTTP status code: out of range", n)
		}
		fmt.Printf("%d %s\n", n, http.StatusText(n))
		if printLong {
			fmt.Printf("%s\n\n", longReasons[n])
		}
		return nil
	}
	return nil
}

func printCodeGroups(c *cli.Context, printLong bool) {
	if c.Bool("1") {
		printRange(1, printLong)
	}
	if c.Bool("2") {
		printRange(2, printLong)
	}
	if c.Bool("3") {
		printRange(3, printLong)
	}
	if c.Bool("4") {
		printRange(4, printLong)
	}
	if c.Bool("5") {
		printRange(5, printLong)
	}
}

func printCodes(c *cli.Context, printLong bool) error {
	if err := printCode(c, printLong); err != nil {
		return err
	}
	printCodeGroups(c, printLong)
	return nil
}

func main() {
	app := &cli.App{
		Name:  "http-status",
		Usage: "explain http status codes",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "1",
				Usage: "explain 100s status codes",
			},
			&cli.BoolFlag{
				Name:  "2",
				Usage: "explain 200s status codes",
			},
			&cli.BoolFlag{
				Name:  "3",
				Usage: "explain 300s status codes",
			},
			&cli.BoolFlag{
				Name:  "4",
				Usage: "explain 400s status codes",
			},
			&cli.BoolFlag{
				Name:  "5",
				Usage: "explain 500s status codes",
			},
			&cli.BoolFlag{
				Name:    "l",
				Usage:   "give longer explanation",
				Aliases: []string{"long"},
			},
		},
		EnableBashCompletion: true,
		Action: func(c *cli.Context) error {
			return printCodes(c, c.Bool("l"))
		},
		Version: "0.1.0",
	}
	cli.AppHelpTemplate = `{{.Name}} - {{.Usage}}

Usage: {{.HelpName}} {{if .VisibleFlags}}[options]{{end}} [status-code]
   {{range .VisibleFlags}}{{.}}
   {{end}}
`

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
