package graphql

import (
	"github.com/firmeve/firmeve/http"
	"github.com/firmeve/firmeve/kernel"
	"github.com/firmeve/firmeve/kernel/contract"
)

type Provider struct {
	kernel.BaseProvider
}

func (p Provider) Name() string {
	return `graphql`
}

func (p Provider) Register() {

}

func (p *Provider) Boot() {
	router := p.Application.Resolve(`http.router`).(*http.Router)
	router.POST(`/graphql`, func(c contract.Context) {

		// @todo
		// 可以考虑每个query对应每个schema
		// 在这里分发graphql
		// register 多个graphql
		/*
			result := graphql.Do(graphql.Params{
				// 这里类似于routes，里面注册多个schema，然后根据Query来查找
				// 或者直接注册一个顶层的，然后在顶层下面直接Query
				Schema:        schema,
				RequestString: query,
			})
			if len(result.Errors) > 0 {
				fmt.Printf("errors: %v", result.Errors)
			}

			//json.NewEncoder(w).Encode(result)
			c.RenderWith(200, render.JSON, result)
		*/
		//return result
		//
		//var schema, _ = graphql.NewSchema(
		//	graphql.SchemaConfig{
		//		Query:    queryType,
		//		Mutation: mutationType,
		//	},
		//)
		//
		//func executeQuery(query string, schema graphql.Schema) *graphql.Result {
		//
		//})
		//	if len(result.Errors) > 0 {
		//	fmt.Printf("errors: %v", result.Errors)
		//}
		//	return result
		//}
		//
		//func main() {
		//	http.HandleFunc("/product", func(w http.ResponseWriter, r *http.Request) {
		//		result := executeQuery(r.URL.Query().Get("query"), schema)
		//		json.NewEncoder(w).Encode(result)
		//	})
	})
}
