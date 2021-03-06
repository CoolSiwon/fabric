/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package multichannel

import (
	"testing"

	newchannelconfig "github.com/hyperledger/fabric/common/channelconfig"
	"github.com/hyperledger/fabric/common/crypto"
	"github.com/hyperledger/fabric/common/ledger/blockledger"
	mockconfigtx "github.com/hyperledger/fabric/common/mocks/configtx"
	"github.com/hyperledger/fabric/common/tools/configtxgen/configtxgentest"
	"github.com/hyperledger/fabric/common/tools/configtxgen/encoder"
	genesisconfig "github.com/hyperledger/fabric/common/tools/configtxgen/localconfig"
	cb "github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protoutil"
	"github.com/stretchr/testify/assert"
)

type mockBlockWriterSupport struct {
	*mockconfigtx.Validator
	crypto.LocalSigner
	blockledger.ReadWriter
}

func (mbws mockBlockWriterSupport) Update(bundle *newchannelconfig.Bundle) {
	mbws.Validator.SequenceVal++
}

func (mbws mockBlockWriterSupport) CreateBundle(channelID string, config *cb.Config) (*newchannelconfig.Bundle, error) {
	return nil, nil
}

func TestCreateBlock(t *testing.T) {
	seedBlock := protoutil.NewBlock(7, []byte("lasthash"))
	seedBlock.Data.Data = [][]byte{[]byte("somebytes")}

	bw := &BlockWriter{lastBlock: seedBlock}
	block := bw.CreateNextBlock([]*cb.Envelope{
		{Payload: []byte("some other bytes")},
	})

	assert.Equal(t, seedBlock.Header.Number+1, block.Header.Number)
	assert.Equal(t, protoutil.BlockDataHash(block.Data), block.Header.DataHash)
	assert.Equal(t, protoutil.BlockHeaderHash(seedBlock.Header), block.Header.PreviousHash)
}

func TestBlockSignature(t *testing.T) {
	bw := &BlockWriter{
		support: &mockBlockWriterSupport{
			LocalSigner: mockCrypto(),
		},
	}

	block := protoutil.NewBlock(7, []byte("foo"))
	bw.addBlockSignature(block)

	md := protoutil.GetMetadataFromBlockOrPanic(block, cb.BlockMetadataIndex_SIGNATURES)
	assert.Nil(t, md.Value, "Value is empty in this case")
	assert.NotNil(t, md.Signatures, "Should have signature")
}

func TestBlockLastConfig(t *testing.T) {
	lastConfigSeq := uint64(6)
	newConfigSeq := lastConfigSeq + 1
	newBlockNum := uint64(9)

	bw := &BlockWriter{
		support: &mockBlockWriterSupport{
			LocalSigner: mockCrypto(),
			Validator: &mockconfigtx.Validator{
				SequenceVal: newConfigSeq,
			},
		},
		lastConfigSeq: lastConfigSeq,
	}

	block := protoutil.NewBlock(newBlockNum, []byte("foo"))
	bw.addLastConfigSignature(block)

	assert.Equal(t, newBlockNum, bw.lastConfigBlockNum)
	assert.Equal(t, newConfigSeq, bw.lastConfigSeq)

	md := protoutil.GetMetadataFromBlockOrPanic(block, cb.BlockMetadataIndex_LAST_CONFIG)
	assert.NotNil(t, md.Value, "Value not be empty in this case")
	assert.NotNil(t, md.Signatures, "Should have signature")

	lc := protoutil.GetLastConfigIndexFromBlockOrPanic(block)
	assert.Equal(t, newBlockNum, lc)
}

func TestWriteConfigBlock(t *testing.T) {
	// TODO, use assert.PanicsWithValue once available
	t.Run("EmptyBlock", func(t *testing.T) {
		assert.Panics(t, func() { (&BlockWriter{}).WriteConfigBlock(&cb.Block{}, nil) })
	})
	t.Run("BadPayload", func(t *testing.T) {
		assert.Panics(t, func() {
			(&BlockWriter{}).WriteConfigBlock(&cb.Block{
				Data: &cb.BlockData{
					Data: [][]byte{
						protoutil.MarshalOrPanic(&cb.Envelope{Payload: []byte("bad")}),
					},
				},
			}, nil)
		})
	})
	t.Run("MissingHeader", func(t *testing.T) {
		assert.Panics(t, func() {
			(&BlockWriter{}).WriteConfigBlock(&cb.Block{
				Data: &cb.BlockData{
					Data: [][]byte{
						protoutil.MarshalOrPanic(&cb.Envelope{
							Payload: protoutil.MarshalOrPanic(&cb.Payload{}),
						}),
					},
				},
			}, nil)
		})
	})
	t.Run("BadChannelHeader", func(t *testing.T) {
		assert.Panics(t, func() {
			(&BlockWriter{}).WriteConfigBlock(&cb.Block{
				Data: &cb.BlockData{
					Data: [][]byte{
						protoutil.MarshalOrPanic(&cb.Envelope{
							Payload: protoutil.MarshalOrPanic(&cb.Payload{
								Header: &cb.Header{
									ChannelHeader: []byte("bad"),
								},
							}),
						}),
					},
				},
			}, nil)
		})
	})
	t.Run("BadChannelHeaderType", func(t *testing.T) {
		assert.Panics(t, func() {
			(&BlockWriter{}).WriteConfigBlock(&cb.Block{
				Data: &cb.BlockData{
					Data: [][]byte{
						protoutil.MarshalOrPanic(&cb.Envelope{
							Payload: protoutil.MarshalOrPanic(&cb.Payload{
								Header: &cb.Header{
									ChannelHeader: protoutil.MarshalOrPanic(&cb.ChannelHeader{}),
								},
							}),
						}),
					},
				},
			}, nil)
		})
	})
}

func TestGoodWriteConfig(t *testing.T) {
	confSys := configtxgentest.Load(genesisconfig.SampleInsecureSoloProfile)
	genesisBlockSys := encoder.New(confSys).GenesisBlock()
	_, l := newRAMLedgerAndFactory(10, genesisconfig.TestChainID, genesisBlockSys)

	bw := &BlockWriter{
		support: &mockBlockWriterSupport{
			LocalSigner: mockCrypto(),
			ReadWriter:  l,
			Validator:   &mockconfigtx.Validator{},
		},
	}

	ctx := makeConfigTx(genesisconfig.TestChainID, 1)
	block := protoutil.NewBlock(1, protoutil.BlockHeaderHash(genesisBlockSys.Header))
	block.Data.Data = [][]byte{protoutil.MarshalOrPanic(ctx)}
	consenterMetadata := []byte("foo")
	bw.WriteConfigBlock(block, consenterMetadata)

	// Wait for the commit to complete
	bw.committingBlock.Lock()
	bw.committingBlock.Unlock()

	cBlock := blockledger.GetBlock(l, block.Header.Number)
	assert.Equal(t, block.Header, cBlock.Header)
	assert.Equal(t, block.Data, cBlock.Data)

	omd := protoutil.GetMetadataFromBlockOrPanic(block, cb.BlockMetadataIndex_ORDERER)
	assert.Equal(t, consenterMetadata, omd.Value)
}

func TestRaceWriteConfig(t *testing.T) {
	confSys := configtxgentest.Load(genesisconfig.SampleInsecureSoloProfile)
	genesisBlockSys := encoder.New(confSys).GenesisBlock()
	_, l := newRAMLedgerAndFactory(10, genesisconfig.TestChainID, genesisBlockSys)

	bw := &BlockWriter{
		support: &mockBlockWriterSupport{
			LocalSigner: mockCrypto(),
			ReadWriter:  l,
			Validator:   &mockconfigtx.Validator{},
		},
	}

	ctx := makeConfigTx(genesisconfig.TestChainID, 1)
	block1 := protoutil.NewBlock(1, protoutil.BlockHeaderHash(genesisBlockSys.Header))
	block1.Data.Data = [][]byte{protoutil.MarshalOrPanic(ctx)}
	consenterMetadata1 := []byte("foo")

	ctx = makeConfigTx(genesisconfig.TestChainID, 1)
	block2 := protoutil.NewBlock(2, protoutil.BlockHeaderHash(block1.Header))
	block2.Data.Data = [][]byte{protoutil.MarshalOrPanic(ctx)}
	consenterMetadata2 := []byte("bar")

	bw.WriteConfigBlock(block1, consenterMetadata1)
	bw.WriteConfigBlock(block2, consenterMetadata2)

	// Wait for the commit to complete
	bw.committingBlock.Lock()
	bw.committingBlock.Unlock()

	cBlock := blockledger.GetBlock(l, block1.Header.Number)
	assert.Equal(t, block1.Header, cBlock.Header)
	assert.Equal(t, block1.Data, cBlock.Data)
	expectedLastConfigBlockNumber := block1.Header.Number
	testLastConfigBlockNumber(t, block1, expectedLastConfigBlockNumber)

	cBlock = blockledger.GetBlock(l, block2.Header.Number)
	assert.Equal(t, block2.Header, cBlock.Header)
	assert.Equal(t, block2.Data, cBlock.Data)
	expectedLastConfigBlockNumber = block2.Header.Number
	testLastConfigBlockNumber(t, block2, expectedLastConfigBlockNumber)

	omd := protoutil.GetMetadataFromBlockOrPanic(block1, cb.BlockMetadataIndex_ORDERER)
	assert.Equal(t, consenterMetadata1, omd.Value)
}
