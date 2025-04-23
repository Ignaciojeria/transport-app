package normalization

/*
import (
	"context"
	"transport-app/app/domain"
	"transport-app/app/shared/infrastructure/gemini"
	"transport-app/app/shared/sharedcontext"
	"transport-app/app/usecase/normalization/chile"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
	"github.com/biter777/countries"
)

type NormalizeAddressInfo func(ctx context.Context, raw domain.AddressInfo) (domain.AddressInfo, error)

func init() {
	ioc.Registry(
		NewNormalizeAddressInfo,
		gemini.NewGemini2Dot0FlashModelWrapper,
		chile.NewSingleInputPrompt)
}
func NewNormalizeAddressInfo(
	model gemini.Gemini2Dot0FlashModelWrapper,
	retrieveChileNormalizationPrompt chile.SingleInputPrompt,
) NormalizeAddressInfo {
	return func(ctx context.Context, ai domain.AddressInfo) (domain.AddressInfo, error) {
		if sharedcontext.TenantCountryFromContext(ctx) == countries.CL.Alpha2() {
			prompt := retrieveChileNormalizationPrompt(ctx, ai)
			return model.StartChat(ctx, prompt)
		}
		return ai, nil
	}
}
*/
