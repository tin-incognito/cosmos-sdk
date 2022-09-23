package mlsag

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/privacy/common"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/operation"
)

// Ring is the struct for a MLSAG ring. It is a matrix of public keys.
// One of these rows is the actual coin transfer, others consist of "decoy" public keys from the chain data.
type Ring struct {
	keys [][]*operation.Point
}

// GetKeys returns the public keys in this Ring as a 2-dimensional array
func (ring Ring) GetKeys() [][]*operation.Point {
	return ring.keys
}

func NewRing(keys [][]*operation.Point) *Ring {
	return &Ring{keys}
}

// ToBytes encodes this Ring's dimensions & contents into a byte array.
func (ring Ring) ToBytes() ([]byte, error) {
	k := ring.keys
	if len(k) == 0 {
		return nil, errors.New("RingToBytes: Ring is empty")
	}
	// Make sure that the ring size is a rectangle row*column
	for i := 1; i < len(k); i++ {
		if len(k[i]) != len(k[0]) {
			return nil, errors.New("RingToBytes: Ring is not a proper rectangle row*column")
		}
	}
	n := len(k)
	m := len(k[0])
	if n > 255 || m > 255 {
		return nil, errors.New("RingToBytes: Ring size is too large")
	}
	b := make([]byte, 3)
	b[0] = MlsagPrefix
	b[1] = byte(n)
	b[2] = byte(m)

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			b = append(b, k[i][j].ToBytesS()...)
		}
	}

	return b, nil
}

// FromBytes decodes a byte array into a Ring
func (ring *Ring) FromBytes(b []byte) (*Ring, error) {
	if len(b) < 3 {
		return nil, errors.New("RingFromBytes: byte length is too short")
	}
	if b[0] != MlsagPrefix {
		return nil, errors.New("RingFromBytes: byte[0] is not MlsagPrefix")
	}
	n := int(b[1])
	m := int(b[2])
	// fmt.Println(b[3 : 3+operation.Ed25519KeySize])

	if len(b) != operation.Ed25519KeySize*n*m+3 {
		return nil, errors.New("RingFromBytes: byte length is not correct")
	}
	offset := 3
	key := make([][]*operation.Point, 0)
	for i := 0; i < n; i++ {
		curRow := make([]*operation.Point, m)
		for j := 0; j < m; j++ {
			currentByte := b[offset : offset+operation.Ed25519KeySize]
			offset += operation.Ed25519KeySize
			currentPoint, err := new(operation.Point).FromBytesS(currentByte)
			if err != nil {
				return nil, errors.New("RingFromBytes: byte contains incorrect point")
			}
			curRow[j] = currentPoint
		}
		key = append(key, curRow)
	}
	ring = NewRing(key)
	return ring, nil
}

func createFakePublicKeyArray(length int) []*operation.Point {
	K := make([]*operation.Point, length)
	for i := 0; i < length; i++ {
		K[i] = operation.RandomPoint()
	}
	return K
}

// NewRandomRing creates a random ring with dimension: (numFake; len(privateKeys)) where we generate fake public keys inside
func NewRandomRing(privateKeys []*operation.Scalar, numFake, pi int) (K *Ring) {
	m := len(privateKeys)

	K = new(Ring)
	K.keys = make([][]*operation.Point, numFake)
	for i := 0; i < numFake; i++ {
		if i != pi {
			K.keys[i] = createFakePublicKeyArray(m)
		} else {
			K.keys[pi] = make([]*operation.Point, m)
			for j := 0; j < m; j++ {
				K.keys[i][j] = parsePublicKey(privateKeys[j], j == m-1)
			}
		}
	}
	return
}

// Mlsag struct contains the data needed to sign a Ring.
//
// pi : the row index of the "real" coins
type Mlsag struct {
	R           *Ring
	pi          int
	keyImages   []*operation.Point
	privateKeys []*operation.Scalar
}

func NewMlsag(privateKeys []*operation.Scalar, R *Ring, pi int) *Mlsag {
	return &Mlsag{
		R,
		pi,
		ParseKeyImages(privateKeys),
		privateKeys,
	}
}

// Parse public key from private key
func parsePublicKey(privateKey *operation.Scalar, isLast bool) *operation.Point {
	// isLast will commit to random base G
	if isLast {
		return new(operation.Point).ScalarMult(
			operation.PedCom.G[operation.PedersenRandomnessIndex],
			privateKey,
		)
	}
	return new(operation.Point).ScalarMultBase(privateKey)
}

func ParseKeyImages(privateKeys []*operation.Scalar) []*operation.Point {
	m := len(privateKeys)

	result := make([]*operation.Point, m)
	for i := 0; i < m; i++ {
		publicKey := parsePublicKey(privateKeys[i], i == m-1)
		hashPoint := operation.HashToPoint(publicKey.ToBytesS())
		result[i] = new(operation.Point).ScalarMult(hashPoint, privateKeys[i])
	}
	return result
}

func (mlsag *Mlsag) createRandomChallenges() (alpha []*operation.Scalar, r [][]*operation.Scalar) {
	m := len(mlsag.privateKeys)
	n := len(mlsag.R.keys)

	alpha = make([]*operation.Scalar, m)
	for i := 0; i < m; i++ {
		alpha[i] = operation.RandomScalar()
	}
	r = make([][]*operation.Scalar, n)
	for i := 0; i < n; i++ {
		r[i] = make([]*operation.Scalar, m)
		if i == mlsag.pi {
			continue
		}
		for j := 0; j < m; j++ {
			r[i][j] = operation.RandomScalar()
		}
	}
	return
}

func calculateFirstC(digest [common.HashSize]byte, alpha []*operation.Scalar, K []*operation.Point) (*operation.Scalar, error) {
	if len(alpha) != len(K) {
		return nil, errors.New("Error in MLSAG: Calculating first C must have length of alpha be the same with length of ring R")
	}
	var b []byte
	b = append(b, digest[:]...)

	// Process columns before the last
	for i := 0; i < len(K)-1; i++ {
		alphaG := new(operation.Point).ScalarMultBase(alpha[i])

		H := operation.HashToPoint(K[i].ToBytesS())
		alphaH := new(operation.Point).ScalarMult(H, alpha[i])

		b = append(b, alphaG.ToBytesS()...)
		b = append(b, alphaH.ToBytesS()...)
	}

	// Process last column
	alphaG := new(operation.Point).ScalarMult(
		operation.PedCom.G[operation.PedersenRandomnessIndex],
		alpha[len(K)-1],
	)
	b = append(b, alphaG.ToBytesS()...)

	return operation.HashToScalar(b), nil
}

func calculateNextC(digest [common.HashSize]byte, r []*operation.Scalar, c *operation.Scalar, K []*operation.Point, keyImages []*operation.Point) (*operation.Scalar, error) {
	if len(r) != len(K) || len(r) != len(keyImages) {
		return nil, errors.New("Error in MLSAG: Calculating next C must have length of r be the same with length of ring R and same with length of keyImages")
	}
	var b []byte
	b = append(b, digest[:]...)

	// Below is the mathematics within the Monero paper:
	// If you are reviewing my code, please refer to paper
	// rG: r*G
	// cK: c*R
	// rG_cK: rG + cK
	//
	// HK: H_p(K_i)
	// rHK: r_i*H_p(K_i)
	// cKI: c*R~ (KI as keyImage)
	// rHK_cKI: rHK + cKI

	// Process columns before the last
	for i := 0; i < len(K)-1; i++ {
		rG := new(operation.Point).ScalarMultBase(r[i])
		cK := new(operation.Point).ScalarMult(K[i], c)
		rG_cK := new(operation.Point).Add(rG, cK)

		HK := operation.HashToPoint(K[i].ToBytesS())
		rHK := new(operation.Point).ScalarMult(HK, r[i])
		cKI := new(operation.Point).ScalarMult(keyImages[i], c)
		rHK_cKI := new(operation.Point).Add(rHK, cKI)

		b = append(b, rG_cK.ToBytesS()...)
		b = append(b, rHK_cKI.ToBytesS()...)
	}

	// Process last column
	rG := new(operation.Point).ScalarMult(
		operation.PedCom.G[operation.PedersenRandomnessIndex],
		r[len(K)-1],
	)
	cK := new(operation.Point).ScalarMult(K[len(K)-1], c)
	rG_cK := new(operation.Point).Add(rG, cK)
	b = append(b, rG_cK.ToBytesS()...)

	return operation.HashToScalar(b), nil
}

func (mlsag *Mlsag) calculateC(message [common.HashSize]byte, alpha []*operation.Scalar, r [][]*operation.Scalar) ([]*operation.Scalar, error) {
	m := len(mlsag.privateKeys)
	n := len(mlsag.R.keys)

	c := make([]*operation.Scalar, n)
	firstC, err := calculateFirstC(
		message,
		alpha,
		mlsag.R.keys[mlsag.pi],
	)
	if err != nil {
		return nil, err
	}

	c[(mlsag.pi+1)%n] = firstC
	for i := (mlsag.pi + 1) % n; i != mlsag.pi; i = (i + 1) % n {
		next := (i + 1) % n
		nextC, err := calculateNextC(message, r[i], c[i], mlsag.R.keys[i], mlsag.keyImages)
		if err != nil {
			return nil, err
		}

		c[next] = nextC
	}

	for i := 0; i < m; i++ {
		ck := new(operation.Scalar).Mul(c[mlsag.pi], mlsag.privateKeys[i])
		r[mlsag.pi][i] = new(operation.Scalar).Sub(alpha[i], ck)
	}

	return c, nil
}

// check l*KI = 0 by checking KI is a valid point
func verifyKeyImages(keyImages []*operation.Point) bool {
	for _, ki := range keyImages {
		if !ki.PointValid() {
			return false
		}
	}
	return true
}

func verifyRing(sig *Sig, R *Ring, message [common.HashSize]byte) (bool, error) {
	c := *sig.c
	cBefore := *sig.c
	if len(R.keys) != len(sig.r) {
		return false, errors.New("MLSAG Error : Malformed Ring")
	}
	for i := 0; i < len(sig.r); i++ {
		nextC, err := calculateNextC(
			message,
			sig.r[i], &c,
			R.keys[i],
			sig.keyImages,
		)
		if err != nil {
			return false, err
		}
		c = *nextC
	}
	return bytes.Equal(c.ToBytesS(), cBefore.ToBytesS()), nil
}

// Verify does verification of a ring signature on a message.
func Verify(sig *Sig, K *Ring, message []byte) (bool, error) {
	if len(message) != common.HashSize {
		return false, fmt.Errorf("Cannot mlsag verify the message because its length is not 32, maybe it has not been hashed")
	}
	message32byte := [32]byte{}
	copy(message32byte[:], message)
	b1 := verifyKeyImages(sig.keyImages)
	if !b1 {
		return false, fmt.Errorf("verify key-images %v failed", sig.keyImages)
	}
	b2, err := verifyRing(sig, K, message32byte)
	return b2, err
}

// Sign uses the private key in this Mlsag to sign a message (which is the transaction hash).
// It returns a Sig.
func (mlsag *Mlsag) Sign(message []byte) (*Sig, error) {
	if len(message) != common.HashSize {
		return nil, errors.New("Cannot mlsag sign the message because its length is not 32, maybe it has not been hashed")
	}
	message32byte := [32]byte{}
	copy(message32byte[:], message)

	alpha, r := mlsag.createRandomChallenges()          // step 2 in paper
	c, err := mlsag.calculateC(message32byte, alpha, r) // step 3 and 4 in paper

	if err != nil {
		return nil, err
	}
	return &Sig{
		c[0], mlsag.keyImages, r,
	}, nil
}
