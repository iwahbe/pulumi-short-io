package short

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
)

const VERSION = "0.1.0"

func Provider() p.Provider {
	return infer.NewProvider().WithResources(
		infer.Resource[*Link, LinkArgs, LinkState](),
	).
		WithConfig(infer.Config[*Config]()).
		WithModuleMap(map[tokens.ModuleName]tokens.ModuleName{
			"pulumi-short-io": "index",
		}).
		WithLanguageMap(map[string]any{
			"nodejs": map[string]bool{
				"respectSchemaVersion": true,
			},
		})
}

type Config struct {
	Token   string `pulumi:"token"`
	Version string `pulumi:"version"`
}

func (c *Config) Annotate(a infer.Annotator) {
	a.Describe(&c.Token, "The authentication token for short.io")
	a.SetDefault(&c.Token, nil, "SHORT_IO_TOKEN")
}

type Link struct{}
type LinkArgs struct {
	Domain string `pulumi:"domain"`
	Short  string `pulumi:"short"`
	Long   string `pulumi:"long"`
}
type LinkState struct {
	LinkArgs

	IdString string `pulumi:"idString"`
}

func (*Link) Create(ctx p.Context, name string, input LinkArgs, preview bool) (string, LinkState, error) {
	url := "https://api.short.io/links"

	payload := strings.NewReader(
		fmt.Sprintf("{\"allowDuplicates\":false, \"domain\": %q, \"originalURL\": %q, \"path\": %q}",
			input.Domain, input.Long, input.Short,
		))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return "", LinkState{}, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", infer.GetConfig[*Config](ctx).Token)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", LinkState{}, err
	}

	result := map[string]any{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", LinkState{}, err
	}

	return name, LinkState{
		LinkArgs: input,
		IdString: result["idString"].(string),
	}, nil
}
