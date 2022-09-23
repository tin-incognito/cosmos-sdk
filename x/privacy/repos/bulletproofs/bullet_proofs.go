package bulletproofs

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/privacy/common"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/operation"
)

// AggregatedRangeWitness contains the prover's secret data (the actual values to be proven & the generated random blinders)
// needed for creating a range proof.
type AggregatedRangeWitness struct {
	values []uint64
	rands  []*operation.Scalar
}

// AggregatedRangeProof is the struct for Bulletproof.
// The statement being proven is that output coins' values are in the uint64 range.
type AggregatedRangeProof struct {
	cmsValue          []*operation.Point
	a                 *operation.Point
	s                 *operation.Point
	t1                *operation.Point
	t2                *operation.Point
	tauX              *operation.Scalar
	tHat              *operation.Scalar
	mu                *operation.Scalar
	innerProductProof *InnerProductProof
}

type bulletproofParams struct {
	g  []*operation.Point
	h  []*operation.Point
	u  *operation.Point
	cs *operation.Point

	precomps []operation.PrecomputedPoint
}

// AggParam contains global Bulletproofs parameters `g, h, u, cs`
var AggParam = newBulletproofParams(common.MaxOutputCoin)

// ValidateSanity performs sanity checks for this proof.
func (proof AggregatedRangeProof) ValidateSanity() bool {
	for i := 0; i < len(proof.cmsValue); i++ {
		if !proof.cmsValue[i].PointValid() {
			return false
		}
	}
	if !proof.a.PointValid() || !proof.s.PointValid() || !proof.t1.PointValid() || !proof.t2.PointValid() {
		return false
	}
	if !proof.tauX.ScalarValid() || !proof.tHat.ScalarValid() || !proof.mu.ScalarValid() {
		return false
	}

	return proof.innerProductProof.ValidateSanity()
}

func NewAggregatedRangeProof() *AggregatedRangeProof {
	proof := &AggregatedRangeProof{}
	proof.a = operation.NewIdentityPoint()
	proof.s = operation.NewIdentityPoint()
	proof.t1 = operation.NewIdentityPoint()
	proof.t2 = operation.NewIdentityPoint()
	proof.tauX = new(operation.Scalar)
	proof.tHat = new(operation.Scalar)
	proof.mu = new(operation.Scalar)
	proof.innerProductProof = NewInnerProductProof()
	return proof
}

// IsNil returns true if any field in this proof is nil
func (proof AggregatedRangeProof) IsNil() bool {
	if proof.a == nil {
		return true
	}
	if proof.s == nil {
		return true
	}
	if proof.t1 == nil {
		return true
	}
	if proof.t2 == nil {
		return true
	}
	if proof.tauX == nil {
		return true
	}
	if proof.tHat == nil {
		return true
	}
	if proof.mu == nil {
		return true
	}
	return proof.innerProductProof == nil
}

func (proof AggregatedRangeProof) GetCommitments() []*operation.Point { return proof.cmsValue }

func (proof *AggregatedRangeProof) SetCommitments(cmsValue []*operation.Point) {
	proof.cmsValue = cmsValue
}

// Bytes marshals the proof into a byte slice
func (proof AggregatedRangeProof) Bytes() []byte {
	var res []byte

	if proof.IsNil() {
		return []byte{}
	}

	res = append(res, byte(len(proof.cmsValue)))
	for i := 0; i < len(proof.cmsValue); i++ {
		res = append(res, proof.cmsValue[i].ToBytesS()...)
	}

	res = append(res, proof.a.ToBytesS()...)
	res = append(res, proof.s.ToBytesS()...)
	res = append(res, proof.t1.ToBytesS()...)
	res = append(res, proof.t2.ToBytesS()...)

	res = append(res, proof.tauX.ToBytesS()...)
	res = append(res, proof.tHat.ToBytesS()...)
	res = append(res, proof.mu.ToBytesS()...)
	res = append(res, proof.innerProductProof.Bytes()...)

	return res
}

// SetBytes unmarshals the proof from a byte slice
func (proof *AggregatedRangeProof) SetBytes(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}

	lenValues := int(bytes[0])
	offset := 1
	var err error

	proof.cmsValue = make([]*operation.Point, lenValues)
	for i := 0; i < lenValues; i++ {
		if offset+operation.Ed25519KeySize > len(bytes) {
			return fmt.Errorf("range-proof byte unmarshaling failed")
		}
		proof.cmsValue[i], err = new(operation.Point).FromBytesS(bytes[offset : offset+operation.Ed25519KeySize])
		if err != nil {
			return err
		}
		offset += operation.Ed25519KeySize
	}

	if offset+operation.Ed25519KeySize > len(bytes) {
		return fmt.Errorf("range-proof byte unmarshaling failed")
	}
	proof.a, err = new(operation.Point).FromBytesS(bytes[offset : offset+operation.Ed25519KeySize])
	if err != nil {
		return err
	}
	offset += operation.Ed25519KeySize

	if offset+operation.Ed25519KeySize > len(bytes) {
		return fmt.Errorf("range-proof byte unmarshaling failed")
	}
	proof.s, err = new(operation.Point).FromBytesS(bytes[offset : offset+operation.Ed25519KeySize])
	if err != nil {
		return err
	}
	offset += operation.Ed25519KeySize

	if offset+operation.Ed25519KeySize > len(bytes) {
		return fmt.Errorf("range-proof byte unmarshaling failed")
	}
	proof.t1, err = new(operation.Point).FromBytesS(bytes[offset : offset+operation.Ed25519KeySize])
	if err != nil {
		return err
	}
	offset += operation.Ed25519KeySize

	if offset+operation.Ed25519KeySize > len(bytes) {
		return fmt.Errorf("range-proof byte unmarshaling failed")
	}
	proof.t2, err = new(operation.Point).FromBytesS(bytes[offset : offset+operation.Ed25519KeySize])
	if err != nil {
		return err
	}
	offset += operation.Ed25519KeySize

	if offset+operation.Ed25519KeySize > len(bytes) {
		return fmt.Errorf("range-proof byte unmarshaling failed")
	}
	proof.tauX = new(operation.Scalar).FromBytesS(bytes[offset : offset+operation.Ed25519KeySize])
	offset += operation.Ed25519KeySize

	if offset+operation.Ed25519KeySize > len(bytes) {
		return fmt.Errorf("range-proof byte unmarshaling failed")
	}
	proof.tHat = new(operation.Scalar).FromBytesS(bytes[offset : offset+operation.Ed25519KeySize])
	offset += operation.Ed25519KeySize

	if offset+operation.Ed25519KeySize > len(bytes) {
		return fmt.Errorf("range-proof byte unmarshaling failed")
	}
	proof.mu = new(operation.Scalar).FromBytesS(bytes[offset : offset+operation.Ed25519KeySize])
	offset += operation.Ed25519KeySize

	if offset >= len(bytes) {
		return fmt.Errorf("range-proof byte unmarshaling failed")
	}

	proof.innerProductProof = new(InnerProductProof)
	return proof.innerProductProof.SetBytes(bytes[offset:])
}

// Set sets the values of both `wit`'s members
func (wit *AggregatedRangeWitness) Set(values []uint64, rands []*operation.Scalar) {
	numValue := len(values)
	wit.values = make([]uint64, numValue)
	wit.rands = make([]*operation.Scalar, numValue)

	for i := range values {
		wit.values[i] = values[i]
		wit.rands[i] = new(operation.Scalar).Set(rands[i])
	}
}

func (wit AggregatedRangeWitness) Prove() (*AggregatedRangeProof, error) {
	proof := new(AggregatedRangeProof)
	numValue := len(wit.values)
	if numValue > common.MaxOutputCoin {
		return nil, fmt.Errorf("output count exceeds MaxOutputCoin")
	}
	numValuePad := roundUpPowTwo(numValue)
	maxExp := common.MaxExp
	N := maxExp * numValuePad

	aggParam := setAggregateParams(N)

	values := make([]uint64, numValuePad)
	rands := make([]*operation.Scalar, numValuePad)
	for i := range wit.values {
		values[i] = wit.values[i]
		rands[i] = new(operation.Scalar).Set(wit.rands[i])
	}
	for i := numValue; i < numValuePad; i++ {
		values[i] = uint64(0)
		rands[i] = new(operation.Scalar).FromUint64(0)
	}

	// Pedersen commitments: V = g^v * h^r
	proof.cmsValue = make([]*operation.Point, numValue)
	for i := 0; i < numValue; i++ {
		proof.cmsValue[i] = operation.PedCom.CommitAtIndex(new(operation.Scalar).FromUint64(values[i]), rands[i], operation.PedersenValueIndex)
	}
	// Convert values to binary array
	aL := make([]*operation.Scalar, N)
	aR := make([]*operation.Scalar, N)
	sL := make([]*operation.Scalar, N)
	sR := make([]*operation.Scalar, N)

	for i, value := range values {
		tmp := ConvertUint64ToBinary(value, maxExp)
		for j := 0; j < maxExp; j++ {
			aL[i*maxExp+j] = tmp[j]
			aR[i*maxExp+j] = new(operation.Scalar).Sub(tmp[j], new(operation.Scalar).FromUint64(1))
			sL[i*maxExp+j] = operation.RandomScalar()
			sR[i*maxExp+j] = operation.RandomScalar()
		}
	}
	// LINE 40-50
	// Commitment to aL, aR:  A = h^alpha * G^aL * H^aR
	// Commitment to sL, sR : S = h^rho   * G^sL * H^sR
	var alpha, rho *operation.Scalar
	alpha = operation.RandomScalar()
	rho = operation.RandomScalar()
	mbuilder := operation.NewMultBuilder(false)
	_, err := encodeVectors(aL, aR, aggParam.g, aggParam.h, mbuilder)
	if err != nil {
		return nil, err
	}
	mbuilder.AppendSingle(alpha, operation.HBase)
	proof.a = mbuilder.Eval() // evaluate & clear builder
	_, err = encodeVectors(sL, sR, aggParam.g, aggParam.h, mbuilder)
	if err != nil {
		return nil, err
	}
	mbuilder.AppendSingle(rho, operation.HBase)
	proof.s = mbuilder.Eval()

	// challenge y, z
	y := generateChallenge(aggParam.cs.ToBytesS(), []*operation.Point{proof.a, proof.s})
	z := generateChallenge(y.ToBytesS(), []*operation.Point{proof.a, proof.s})

	// LINE 51-54
	twoNumber := new(operation.Scalar).FromUint64(2)
	twoVectorN := powerVector(twoNumber, maxExp)

	// HPrime = H^(y^(1-i))
	HPrime := computeHPrime(y, N, aggParam.h)

	// l(X) = (aL -z*1^n) + sL*X; r(X) = y^n hada (aR +z*1^n + sR*X) + z^2 * 2^n
	yVector := powerVector(y, N)
	hadaProduct, err := hadamardProduct(yVector, vectorAddScalar(aR, z))
	if err != nil {
		return nil, err
	}
	vectorSum := make([]*operation.Scalar, N)
	zTmp := new(operation.Scalar).Set(z)
	for j := 0; j < numValuePad; j++ {
		zTmp.Mul(zTmp, z)
		for i := 0; i < maxExp; i++ {
			vectorSum[j*maxExp+i] = new(operation.Scalar).Mul(twoVectorN[i], zTmp)
		}
	}
	zNeg := new(operation.Scalar).Sub(new(operation.Scalar).FromUint64(0), z)
	l0 := vectorAddScalar(aL, zNeg)
	l1 := sL
	var r0, r1 []*operation.Scalar
	if r0, err = vectorAdd(hadaProduct, vectorSum); err != nil {
		return nil, err
	} else if r1, err = hadamardProduct(yVector, sR); err != nil {
		return nil, err
	}

	// t(X) = <l(X), r(X)> = t0 + t1*X + t2*X^2
	// t1 = <l1, ro> + <l0, r1>, t2 = <l1, r1>
	var t1, t2 *operation.Scalar
	if ip3, err := innerProduct(l1, r0); err != nil {
		return nil, err
	} else if ip4, err := innerProduct(l0, r1); err != nil {
		return nil, err
	} else {
		t1 = new(operation.Scalar).Add(ip3, ip4)
		if t2, err = innerProduct(l1, r1); err != nil {
			return nil, err
		}
	}

	// commitment to t1, t2
	tau1 := operation.RandomScalar()
	tau2 := operation.RandomScalar()
	proof.t1 = operation.PedCom.CommitAtIndex(t1, tau1, operation.PedersenValueIndex)
	proof.t2 = operation.PedCom.CommitAtIndex(t2, tau2, operation.PedersenValueIndex)

	x := generateChallenge(z.ToBytesS(), []*operation.Point{proof.t1, proof.t2})
	xSquare := new(operation.Scalar).Mul(x, x)

	// lVector = aL - z*1^n + sL*x
	// rVector = y^n hada (aR +z*1^n + sR*x) + z^2*2^n
	// tHat = <lVector, rVector>
	lVector, err := vectorAdd(vectorAddScalar(aL, zNeg), vectorMulScalar(sL, x))
	if err != nil {
		return nil, err
	}
	tmpVector, err := vectorAdd(vectorAddScalar(aR, z), vectorMulScalar(sR, x))
	if err != nil {
		return nil, err
	}
	rVector, err := hadamardProduct(yVector, tmpVector)
	if err != nil {
		return nil, err
	}
	rVector, err = vectorAdd(rVector, vectorSum)
	if err != nil {
		return nil, err
	}
	proof.tHat, err = innerProduct(lVector, rVector)
	if err != nil {
		return nil, err
	}

	// blinding value for tHat (agg): tauX = tau2 * x^2 + tau1 * x + <z^(1+m), rands>
	proof.tauX = new(operation.Scalar).Mul(tau2, xSquare)
	proof.tauX.Add(proof.tauX, new(operation.Scalar).Mul(tau1, x))
	zTmp = new(operation.Scalar).Set(z)
	tmpBN := new(operation.Scalar)
	for j := 0; j < numValuePad; j++ {
		zTmp.Mul(zTmp, z)
		proof.tauX.Add(proof.tauX, tmpBN.Mul(zTmp, rands[j]))
	}

	// alpha, rho blind A, S
	// mu = alpha + rho * x
	proof.mu = new(operation.Scalar).Add(alpha, new(operation.Scalar).Mul(rho, x))

	// instead of sending left vector and right vector, we use inner sum argument to reduce proof size from 2*n to 2(log2(n)) + 2
	innerProductWit := new(InnerProductWitness)
	innerProductWit.a = lVector
	innerProductWit.b = rVector
	// u' = u^x
	uPrime := new(operation.Point).ScalarMult(aggParam.u, operation.HashToScalar(x.ToBytesS()))

	// for inner-product proof, compute P = g^l * h^r * u ^ (l*r); subsitute h <- h', u <- u'
	_, err = encodeVectors(lVector, rVector, aggParam.g, HPrime, mbuilder)
	if err != nil {
		return nil, err
	}
	mbuilder.AppendSingle(proof.tHat, uPrime)
	innerProductWit.p = mbuilder.Eval()

	proof.innerProductProof, err = innerProductWit.Prove(aggParam.g, HPrime, uPrime, x.ToBytesS())
	if err != nil {
		return nil, err
	}

	return proof, nil
}
