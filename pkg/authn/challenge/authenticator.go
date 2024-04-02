package challenge

import (
	"context"
	"crypto/rsa"
	"embed"
	"encoding"
	"encoding/json"

	"github.com/bacalhau-project/bacalhau/pkg/authn"
	"github.com/bacalhau-project/bacalhau/pkg/lib/policy"
	"github.com/bacalhau-project/bacalhau/pkg/system"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

//go:embed *.rego
var policies embed.FS

// The data that will be passed to the authn policy, once verification of the
// input challenge has passed.
type policyData struct {
	SigningKey jwk.Key `json:"signingKey"`
	NodeID     string  `json:"nodeId"`
	ClientID   string  `json:"clientId"`
}

// The data that the user will supply to us to try and authenticate.
type response struct {
	PhraseSignature string `json:"PhraseSignature"`

	// A base64 encoded PKCS1 public key, as generated by config.encodePublicKey
	PublicKey string `json:"PublicKey"`
}

// The data that we will send to the user to allow them to sign a challenge.
type request struct {
	InputPhrase []byte `json:"InputPhrase"`
}

type challengeAuthenticator struct {
	authnPolicy *policy.Policy
	getPhrase   encoding.BinaryMarshaler
	key         jwk.Key
	nodeID      string

	challenge policy.Query[policyData, string]
}

func NewAuthenticator(p *policy.Policy, getPhrase encoding.BinaryMarshaler, key *rsa.PrivateKey, nodeID string) authn.Authenticator {
	return challengeAuthenticator{
		authnPolicy: p,
		getPhrase:   getPhrase,
		key:         lo.Must(jwk.New(key)),
		nodeID:      nodeID,
		challenge:   policy.AddQuery[policyData, string](p, authn.PolicyTokenRule),
	}
}

func (authenticator challengeAuthenticator) Authenticate(ctx context.Context, req []byte) (authn.Authentication, error) {
	var userInput response
	err := json.Unmarshal(req, &userInput)
	if err != nil {
		return authn.Error(errors.Wrap(err, "invalid authentication data"))
	}

	inputPhrase, err := authenticator.getPhrase.MarshalBinary()
	if err != nil {
		return authn.Error(err)
	}

	err = system.Verify(inputPhrase, userInput.PhraseSignature, userInput.PublicKey)
	if err != nil {
		// Don't return an error here because this is likely a bad user request.
		return authn.Failed(err.Error()), nil
	}

	userKey, err := system.DecodePublicKey(userInput.PublicKey)
	if err != nil {
		return authn.Error(err)
	}

	data := policyData{
		SigningKey: authenticator.key,
		NodeID:     authenticator.nodeID,
		ClientID:   system.ConvertToClientID(userKey),
	}

	token, err := authenticator.challenge(ctx, data)
	if errors.Is(err, policy.ErrNoResult) {
		return authn.Failed("signature verified but user credentials rejected"), nil
	} else if err != nil {
		return authn.Error(err)
	}

	return authn.Authentication{Success: true, Token: token}, nil
}

func (challengeAuthenticator) IsInstalled(context.Context) (bool, error) {
	return true, nil
}

func (authenticator challengeAuthenticator) Requirement() authn.Requirement {
	req := request{
		InputPhrase: lo.Must(authenticator.getPhrase.MarshalBinary()),
	}

	params := json.RawMessage(lo.Must(json.Marshal(req)))
	return authn.Requirement{
		Type:   authn.MethodTypeChallenge,
		Params: &params,
	}
}

// AnonymousModePolicy grants read-only access to all namespaces and full access
// to a namespace matching the user's client ID, derived from their submitted
// public key.
var AnonymousModePolicy *policy.Policy = lo.Must(policy.FromFS(policies, "challenge_ns_anon.rego"))
