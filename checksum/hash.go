package checksum

import "hash"

type Hash8 interface {
	hash.Hash

	Sum8() uint8
}

type Hash16 interface {
	hash.Hash

	Sum16() uint16
}
