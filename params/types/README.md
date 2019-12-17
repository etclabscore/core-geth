`type Genesis struct` and `type GenesisAccount struct` have custom JSON un/marshaling methods.

These methods are defined in the files prefixed with `gen_`, and were generated
initially with the following commands:

```go
gencodec -type Genesis -field-override genesisSpecMarshaling -out gen_genesis.go
gencodec -type GenesisAccount -field-override genesisAccountMarshaling -out gen_genesis_account.go
```

The source for the `gencodec` command is here: https://github.com/fjl/gencodec. 

While these were once unmodified and thus "pure" generated files, that is no longer the case.
MultiGeth's needs for custom JSON un/marshaling have outgrown the capabilities
of the tool.

This is a comparative `git diff` which could be used to modify the original
generated files with a command such as `git apply`.


```txt
diff --git a/params/types/gen_genesis.go b/params/types/gen_genesis.go
index 392b710fd..573d8fba7 100644
--- a/params/types/gen_genesis.go
+++ b/params/types/gen_genesis.go
@@ -11,6 +11,7 @@ import (
 	"github.com/ethereum/go-ethereum/common/hexutil"
 	"github.com/ethereum/go-ethereum/common/math"
 	common0 "github.com/ethereum/go-ethereum/params/types/common"
+	"github.com/ethereum/go-ethereum/params/types/goethereum"
 )
 
 var _ = (*genesisSpecMarshaling)(nil)
@@ -69,9 +70,20 @@ func (g *Genesis) UnmarshalJSON(input []byte) error {
 		ParentHash *common.Hash                                `json:"parentHash"`
 	}
 	var dec Genesis
-	if err := json.Unmarshal(input, &dec); err != nil {
-		return err
+
+	// Note that this logic is importantly relate to the logic in params/convert/json.go, for ChainConfigurator
+	// unmarshaling.
+	dec.Config = &MultiGethChainConfig{}
+	if err := json.Unmarshal(input, &dec); err != nil || common0.IsValid(dec.Config, nil) != nil {
+		dec.Config = &goethereum.ChainConfig{}
+		if err := json.Unmarshal(input, &dec); err != nil {
+			return err
+		}
+		if err := common0.IsValid(dec.Config, nil); err !=nil {
+			return err
+		}
 	}
+
 	if dec.Config != nil {
 		g.Config = dec.Config
 	}
```