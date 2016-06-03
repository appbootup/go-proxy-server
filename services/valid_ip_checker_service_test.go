package services_test

import (
	"github.com/GetSimpl/proxy-server/db"
	"github.com/GetSimpl/proxy-server/services"
	"github.com/GetSimpl/go-simpl/logger"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net"
)

func clearRedis(){
	redisClient := db.Get()
	_, err :=  redisClient.Del(services.CLIENT_IP_ADDRESSESS).Result()
	if err != nil {
		panic(err)
		logger.E(err)
	}
}

var _ = BeforeSuite(func(){
	db.Init()
	logger.Init(1)
})

var _ = AfterSuite(func(){
	clearRedis()
	db.CloseConnection()
})

var _ = Describe("ValidIpCheckerService", func() {
	Context("When a set of valid Ips are provided", func(){
		BeforeEach(func(){
			clearRedis()
			validIpAddress := []string{"127.0.0.1"}
			dbClient := db.Get()
			_, err := dbClient.LPush(services.CLIENT_IP_ADDRESSESS, validIpAddress...).Result()
			if (err != nil) {
				logger.E("Unable to push seeds into the Redis")
				panic(err)
			}

		})


		Context(" And when a request originates from a valid IP", func(){
			It("Should not return an error", func(){
				rulesEngineChecker := services.NewRulesEngineChecker()
				req, _ := http.NewRequest("GET", "http://127.0.0.1:8000/api/v1/simpl_buy/approved", nil)
				req.RemoteAddr = "127.0.0.1:8000"
				_, _, splitErr := net.SplitHostPort(req.RemoteAddr)

				if splitErr != nil {
					logger.E(splitErr)
					panic(splitErr)
				}

				checker := services.IsIPValid(req, rulesEngineChecker)
				isSuccess, _ := checker.Success()
				Expect(isSuccess).To(BeTrue())
			})
		})

		Context("And when a request originates from an Invalid IP", func(){
			It("Should return an error", func() {
				rulesEngineChecker := services.NewRulesEngineChecker()
				req, _ := http.NewRequest("GET", "http://127.0.0.1:9000/api/v1/simpl_buy/approved", nil)
				req.RemoteAddr = "127.0.0.2:8000"
				_, _, splitErr := net.SplitHostPort(req.RemoteAddr)

				if splitErr != nil {
					logger.E(splitErr)
					panic(splitErr)
				}

				checker := services.IsIPValid(req, rulesEngineChecker)
				isSuccess, errors := checker.Success()
				Expect(isSuccess).To(BeFalse())
				Expect(errors[0]).To(Equal("[Invalid IP Address] Host: 127.0.0.2:8000 trying to access: 127.0.0.1:9000"))
			})
		})
	})
	Context("When a valid set of Ips are not provided", func() {
		BeforeEach(func(){
			clearRedis()
			validClientIps, _ := services.GetValidClientIps()
			Expect(len(validClientIps)).To(Equal(0))

		})

		Context("And when a request originates from a random IP", func(){
			It("Should be allowed to reach the application server", func(){
				rulesEngineChecker := services.NewRulesEngineChecker()
				req, _ := http.NewRequest("GET", "http://127.0.0.1:9000/api/v1/simpl_buy/approved", nil)
				req.RemoteAddr = "127.0.0.2:8000"
				_, _, splitErr := net.SplitHostPort(req.RemoteAddr)

				if splitErr != nil {
					logger.E(splitErr)
					panic(splitErr)
				}

				checker := services.IsIPValid(req, rulesEngineChecker)
				isSuccess, errors := checker.Success()
				Expect(isSuccess).To(BeTrue())
				Expect(len(errors)).To(Equal(0))
			})
		})
	})
})
