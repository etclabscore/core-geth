package params

// ABCBootnodes are the enode URLs of the P2P bootstrap nodes running
// reliably and availably on the ABC network.
// They will be the first point of contact for nodes coming online
// to find peers.
var ABCBootnodes = []string{
	"enode://454b484b652c07c72adebfabf8bc2bd95b420b16952ef3de76a9c00ef63f07cca02a20bd2363426f9e6fe372cef96a42b0fec3c747d118f79fd5e02f2a4ebc4f@42.152.120.98:30303",
}

// Once ABC network has DNS discovery set up,
// this value can be configured.
// var ABCDNSNetwork = "enrtree://AJE62Q4DUX4QMMXEHCSSCSC65TDHZYSMONSD64P3WULVLSF6MRQ3K@example.network"
