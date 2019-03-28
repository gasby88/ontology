/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

package utils

import (
	"bytes"
	"fmt"
	ctypes "github.com/ontio/ontology/core/types"
	"io"
	"math/big"

	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/common/serialization"
	"github.com/ontio/ontology/vm/neovm/types"
)

func SerializeShardId(w io.Writer, id ctypes.ShardID) error {
	return WriteVarUint(w, id.ToUint64())
}

func DeserializeShardId(r io.Reader) (ctypes.ShardID, error) {
	id, err := ReadVarUint(r)
	if err != nil {
		return ctypes.ShardID{}, err
	}
	shardId, err := ctypes.NewShardID(id)
	if err != nil {
		return ctypes.ShardID{}, fmt.Errorf("generate shard id failed, err: %s", err)
	}
	return shardId, nil
}

func WriteVarUint(w io.Writer, value uint64) error {
	if err := serialization.WriteVarBytes(w, types.BigIntToBytes(big.NewInt(int64(value)))); err != nil {
		return fmt.Errorf("serialize value error:%v", err)
	}
	return nil
}

func ReadVarUint(r io.Reader) (uint64, error) {
	value, err := serialization.ReadVarBytes(r)
	if err != nil {
		return 0, fmt.Errorf("deserialize value error:%v", err)
	}
	v := types.BigIntFromBytes(value)
	if v.Cmp(big.NewInt(0)) < 0 {
		return 0, fmt.Errorf("%s", "value should not be a negative number.")
	}
	return v.Uint64(), nil
}

func WriteAddress(w io.Writer, address common.Address) error {
	if err := serialization.WriteVarBytes(w, address[:]); err != nil {
		return fmt.Errorf("serialize value error:%v", err)
	}
	return nil
}

func ReadAddress(r io.Reader) (common.Address, error) {
	from, err := serialization.ReadVarBytes(r)
	if err != nil {
		return common.Address{}, fmt.Errorf("[State] deserialize from error:%v", err)
	}
	return common.AddressParseFromBytes(from)
}

func EncodeAddress(sink *common.ZeroCopySink, addr common.Address) (size uint64) {
	return sink.WriteVarBytes(addr[:])
}

func EncodeVarUint(sink *common.ZeroCopySink, value uint64) (size uint64) {
	return sink.WriteVarBytes(types.BigIntToBytes(big.NewInt(int64(value))))
}

func DecodeVarUint(source *common.ZeroCopySource) (uint64, error) {
	value, _, irregular, eof := source.NextVarBytes()
	if eof {
		return 0, io.ErrUnexpectedEOF
	}
	if irregular {
		return 0, common.ErrIrregularData
	}
	v := types.BigIntFromBytes(value)
	if v.Cmp(big.NewInt(0)) < 0 {
		return 0, fmt.Errorf("%s", "value should not be a negative number.")
	}
	return v.Uint64(), nil
}

func DecodeAddress(source *common.ZeroCopySource) (common.Address, error) {
	from, _, irregular, eof := source.NextVarBytes()
	if eof {
		return common.Address{}, io.ErrUnexpectedEOF
	}
	if irregular {
		return common.Address{}, common.ErrIrregularData
	}

	return common.AddressParseFromBytes(from)
}

func GetUint32Bytes(num uint32) ([]byte, error) {
	bf := new(bytes.Buffer)
	if err := serialization.WriteUint32(bf, num); err != nil {
		return nil, fmt.Errorf("serialization.WriteUint32, serialize uint32 error: %v", err)
	}
	return bf.Bytes(), nil
}

func GetBytesUint32(b []byte) (uint32, error) {
	num, err := serialization.ReadUint32(bytes.NewBuffer(b))
	if err != nil {
		return 0, fmt.Errorf("serialization.ReadUint32, deserialize uint32 error: %v", err)
	}
	return num, nil
}

func GetUint64Bytes(num uint64) ([]byte, error) {
	bf := new(bytes.Buffer)
	if err := serialization.WriteUint64(bf, num); err != nil {
		return nil, fmt.Errorf("serialization.WriteUint64, serialize uint64 error: %v", err)
	}
	return bf.Bytes(), nil
}

func GetBytesUint64(b []byte) (uint64, error) {
	num, err := serialization.ReadUint64(bytes.NewBuffer(b))
	if err != nil {
		return 0, fmt.Errorf("serialization.ReadUint64, deserialize uint64 error: %v", err)
	}
	return num, nil
}
