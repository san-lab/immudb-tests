package datastructs

import (
	"github.com/codenotary/immudb/pkg/client"
)

var THIS_BANK NameAddress
var COUNTERPART_BANKS = make(map[string]string) // First entry is THIS_BANK

type NameAddress struct {
	Name    string
	Address string
}

var STATE_CLIENT client.ImmuClient
var MSGS_CLIENT client.ImmuClient
