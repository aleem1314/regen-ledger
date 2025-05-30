package v1

import (
	"strconv"
	"strings"
	"testing"

	"github.com/cosmos/gogoproto/jsonpb"
	"github.com/regen-network/gocuke"
	"github.com/stretchr/testify/require"
)

type batchIssuance struct {
	t        gocuke.TestingT
	issuance *BatchIssuance
	err      error
}

func TestBatchIssuance(t *testing.T) {
	gocuke.NewRunner(t, &batchIssuance{}).Path("./features/types_batch_issuance.feature").Run()
}

func (s *batchIssuance) Before(t gocuke.TestingT) {
	s.t = t
}

func (s *batchIssuance) TheBatchIssuance(a gocuke.DocString) {
	s.issuance = &BatchIssuance{}
	err := jsonpb.UnmarshalString(a.Content, s.issuance)
	require.NoError(s.t, err)
}

func (s *batchIssuance) RetirementReasonWithLength(a string) {
	length, err := strconv.ParseInt(a, 10, 64)
	require.NoError(s.t, err)

	s.issuance.RetirementReason = strings.Repeat("x", int(length))
}

func (s *batchIssuance) TheBatchIssuanceIsValidated() {
	s.err = s.issuance.Validate()
}

func (s *batchIssuance) ExpectTheError(a string) {
	require.EqualError(s.t, s.err, a)
}

func (s *batchIssuance) ExpectNoError() {
	require.NoError(s.t, s.err)
}
