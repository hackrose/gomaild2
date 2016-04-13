package helo

import (
	. "github.com/trapped/gomaild2/smtp/structs"
)

func identify(c *Client, cmd Command) Reply {
	valid := RxDomain.MatchString(cmd.Args)
	//TODO: start blacklist check
	if valid {
		c.Data["identifier"] = interface{}(cmd.Args)
		c.State = Identified
		return Reply{
			Result:  Ok,
			Message: "domain validated, welcome " + cmd.Args,
		}
	} else {
		return Reply{
			Result:  ArgumentSyntaxError,
			Message: "invalid domain",
		}
	}
}

func alreadyidentified(c *Client, cmd Command) Reply {
	return Reply{
		Result:  BadSequence,
		Message: "already identified",
	}
}

func Process(c *Client, cmd Command) Reply {
	if c.State >= Identified {
		return alreadyidentified(c, cmd)
	} else {
		return identify(c, cmd)
	}
}
