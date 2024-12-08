package infinirewards

import (
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/fxamacker/cbor/v2"
)

// Metadata represents the metadata structure for points and collectibles
type Metadata struct {
	Metadata string                 `cbor:"description"`
	Extra    map[string]interface{} `cbor:"extra,omitempty"`
}

// EncodeCBORMetadata encodes metadata into CBOR format
func EncodeCBORMetadata(description string, extra map[string]interface{}) ([]byte, error) {
	metadata := Metadata{
		Metadata: description,
		Extra:    extra,
	}
	return cbor.Marshal(metadata)
}

// DecodeCBORMetadata decodes CBOR data into metadata
func DecodeCBORMetadata(data []byte) (*Metadata, error) {
	var metadata Metadata
	err := cbor.Unmarshal(data, &metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to decode CBOR metadata: %w", err)
	}
	return &metadata, nil
}

// StringToByteArrFeltWithCBOR converts a metadata object to a byte array felt representation
func StringToByteArrFeltWithCBOR(description string, extra map[string]interface{}) ([]*felt.Felt, error) {
	cborData, err := EncodeCBORMetadata(description, extra)
	if err != nil {
		return nil, fmt.Errorf("failed to encode CBOR metadata: %w", err)
	}
	return utils.StringToByteArrFelt(string(cborData))
}

// ByteArrFeltToCBORMetadata converts a byte array felt to metadata
func ByteArrFeltToCBORMetadata(felts []*felt.Felt) (*Metadata, error) {
	str, err := utils.ByteArrFeltToString(felts)
	if err != nil {
		return nil, fmt.Errorf("failed to convert felt array to string: %w", err)
	}
	return DecodeCBORMetadata([]byte(str))
}
