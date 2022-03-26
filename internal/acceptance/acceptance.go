package acceptance

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pjfebbraro/terraform-provider-azureadb2cief/internal/provider"
	"os"
	"sync"
)

var AzureADB2CProvider *schema.Provider
var once sync.Once

func init() {
	if os.Getenv("TF_ACC") == "" {
		return
	}
	EnsureProvidersAreInitialised()
}

func EnsureProvidersAreInitialised() {
	once.Do(func() {
		AzureADB2CProvider = provider.New("dev")()
	})
}
