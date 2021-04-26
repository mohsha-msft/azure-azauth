package azauth

import (
	"github.com/JeffreyRichter/enum/enum"
	"reflect"
)

//var EAuthType = AuthType(0)

type AuthType uint8

func (AuthType) SharedKey() AuthType { return AuthType(0) }
func (AuthType) SPN() AuthType       { return AuthType(1) }
func (AuthType) MSI() AuthType       { return AuthType(2) }

func (o AuthType) String() string {
	return enum.StringInt(o, reflect.TypeOf(o))
}

var EClientType = ClientType(0)

type ClientType uint8

func (ClientType) ServiceClient() ClientType   { return ClientType(0) }
func (ClientType) ContainerClient() ClientType { return ClientType(1) }
func (ClientType) BlobClient() ClientType      { return ClientType(2) }

func (ct ClientType) String() string {
	return enum.StringInt(ct, reflect.TypeOf(ct))
}
