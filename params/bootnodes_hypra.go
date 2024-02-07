package params

// HypraBootnodes are the enode URLs of the P2P bootstrap nodes running
// reliably and availably on the MintMe network.
var HypraBootnodes = []string{
	// Foundation Bootnodes
	"enode://b3d8c6ad187f54860bd288e8e343c5cb05db023b3a74a4cd9d85cae3e2677074f92b3afecfd2bb445f9cba151848d3294abff9bedcee5d437ff161300f5144e9@77.100.75.201:30303", // Dev
	"enode://301c2d2d622fe3d49f9a9d5a294a1a65ce0f686a10b5b6ea2e965533b7e84ecea25f1f2eec78e6fa948ca129ec5f9a8fe731d9641df0163e4847ded09dbfd1e4@54.36.108.60:30303",  // Explorer

	// Communtiy Bootnodes
	"enode://959f6378ee6162f57977e5e6ab8dd56bd8ef5d1bc2a1bb01c6b41cfc2d07ea490d4c939c7625f13886c684b221a9c3e710e4a66a718a3231c40d2536c344df9d@27.254.39.27:30308",
	"enode://e82bf286f09a7b86f5528a0e7c29928a8bb0bf9416d9678a91da9e2729480700a71777490ed115cad82b9f75268fc1f9a0d9483bb65acd6665708778c2d035f5@178.234.84.24:30303?discport=1337",
	"enode://fe072785d5044f22b393df8a364dcc92d927b9f88aff14bff2484db20caa8350a07df3b9b1f0fb0b222304f426ab887ad9829bff6948aba84e3b5f1776b8dd52@195.201.122.219:30303",
}

// Once Hypra network has DNS discovery set up,
// this value can be configured.
// var HypraDNSNetwork = "enrtree://AJE62Q4DUX4QMMXEHCSSCSC65TDHZYSMONSD64P3WULVLSF6MRQ3K@example.network"
