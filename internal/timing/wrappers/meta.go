package timingwrappers

//go:generate go run ./generate/generate.go -in gitlab.cmpayments.local/creditcard/authorization/internal/authorization/app.Repository -out AuthorizationRepository
//go:generate go run ./generate/generate.go -in gitlab.cmpayments.local/creditcard/authorization/internal/authorization/app.Tokenizer -out AuthorizationTokenizer

//go:generate go run ./generate/generate.go -in gitlab.cmpayments.local/creditcard/authorization/internal/refund/app.Repository -out RefundRepository
//go:generate go run ./generate/generate.go -in gitlab.cmpayments.local/creditcard/authorization/internal/capture/app.CaptureRepository -out CaptureRepository
//go:generate go run ./generate/generate.go -in gitlab.cmpayments.local/creditcard/authorization/internal/reversal/app.ReversalRepository -out ReversalRepository

//go:generate go run ./generate/generate.go -in gitlab.cmpayments.local/creditcard/platform/events/pubsub.Publisher -out Publisher
