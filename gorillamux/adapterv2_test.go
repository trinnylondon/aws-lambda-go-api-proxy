package gorillamux_test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gorilla/mux"
	"github.com/trinnylondon/aws-lambda-go-api-proxy/gorillamux"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GorillaMuxAdapterV2 tests", func() {
	Context("Simple ping request", func() {
		It("Proxies the event correctly", func() {
			homeHandler := func(w http.ResponseWriter, req *http.Request) {
				w.Header().Add("unfortunately-required-header", "")
				fmt.Fprintf(w, "Home Page")
			}

			productsHandler := func(w http.ResponseWriter, req *http.Request) {
				w.Header().Add("unfortunately-required-header", "")
				fmt.Fprintf(w, "Products Page")
			}

			r := mux.NewRouter()
			r.HandleFunc("/", homeHandler)
			r.HandleFunc("/products", productsHandler)

			adapter := gorillamux.NewV2(r)

			homePageReq := events.APIGatewayV2HTTPRequest{
				RequestContext: events.APIGatewayV2HTTPRequestContext{
					HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
						Method: http.MethodGet,
						Path:   "/",
					},
				},
			}

			homePageResp, homePageReqErr := adapter.ProxyWithContext(context.Background(), homePageReq)

			Expect(homePageReqErr).To(BeNil())
			Expect(homePageResp.StatusCode).To(Equal(200))
			Expect(homePageResp.Body).To(Equal("Home Page"))

			productsPageReq := events.APIGatewayV2HTTPRequest{
				RequestContext: events.APIGatewayV2HTTPRequestContext{
					HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
						Method: http.MethodGet,
						Path:   "/products",
					},
				},
			}

			productsPageResp, productsPageReqErr := adapter.Proxy(productsPageReq)

			Expect(productsPageReqErr).To(BeNil())
			Expect(productsPageResp.StatusCode).To(Equal(200))
			Expect(productsPageResp.Body).To(Equal("Products Page"))
		})
	})
})
