package stub

func NewWebhookVerifierStub() *WebhookVerifierStub { return &WebhookVerifierStub{} }

type WebhookVerifierStub struct{}

func (w *WebhookVerifierStub) VerifySignature(headers map[string]string, body []byte) (bool, error) {
	return true, nil
}
