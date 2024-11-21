package routes

import (
	"infinirewards/controllers"
	"infinirewards/middleware"
	"net/http"
)

func SetInfiniRewardsRoutes(mux *http.ServeMux) {
	// Collectible endpoints
	//	@Summary		Mint Collectible
	//	@Description	Mint new collectible tokens
	//	@Tags			collectibles
	//	@Accept			json
	//	@Produce		json
	//	@Security		BearerAuth
	//	@Param			request	body		models.MintCollectibleRequest	true	"Mint Request"
	//	@Success		200		{object}	models.MintCollectibleResponse
	//	@Failure		400		{string}	string	"Bad Request"
	//	@Failure		401		{string}	string	"Unauthorized"
	//	@Failure		500		{string}	string	"Internal Server Error"
	//	@Router			/collectibles/mint [post]
	mux.HandleFunc("POST /collectibles/mint", middleware.AuthMiddleware(controllers.MintCollectibleHandler))

	//	@Summary		Get Collectible Balance
	//	@Description	Get balance of collectible tokens
	//	@Tags			collectibles
	//	@Accept			json
	//	@Produce		json
	//	@Security		BearerAuth
	//	@Param			address	path		string	true	"Contract Address"
	//	@Param			tokenId	path		string	true	"Token ID"
	//	@Success		200		{object}	models.GetCollectibleBalanceResponse
	//	@Failure		400		{string}	string	"Bad Request"
	//	@Failure		401		{string}	string	"Unauthorized"
	//	@Failure		500		{string}	string	"Internal Server Error"
	//	@Router			/collectibles/{address}/balance/{tokenId} [get]
	mux.HandleFunc("GET /collectibles/{address}/balance/{tokenId}", middleware.AuthMiddleware(controllers.GetCollectibleBalanceHandler))

	//	@Summary		Get Collectible URI
	//	@Description	Get URI for collectible token
	//	@Tags			collectibles
	//	@Accept			json
	//	@Produce		json
	//	@Param			address	path		string	true	"Contract Address"
	//	@Param			tokenId	path		string	true	"Token ID"
	//	@Success		200		{object}	models.GetCollectibleURIResponse
	//	@Failure		400		{string}	string	"Bad Request"
	//	@Failure		500		{string}	string	"Internal Server Error"
	//	@Router			/collectibles/{address}/uri/{tokenId} [get]
	mux.HandleFunc("GET /collectibles/{address}/uri/{tokenId}", controllers.GetCollectibleURIHandler)

	//	@Summary		Check Collectible Validity
	//	@Description	Check if a collectible token is valid
	//	@Tags			collectibles
	//	@Accept			json
	//	@Produce		json
	//	@Param			address	path		string	true	"Contract Address"
	//	@Param			tokenId	path		string	true	"Token ID"
	//	@Success		200		{object}	models.IsCollectibleValidResponse
	//	@Failure		400		{string}	string	"Bad Request"
	//	@Failure		500		{string}	string	"Internal Server Error"
	mux.HandleFunc("GET /collectibles/{address}/valid/{tokenId}", controllers.IsCollectibleValidHandler)

	//	@Summary		Set Token Data
	//	@Description	Set token data for collectible
	//	@Tags			collectibles
	//	@Accept			json
	//	@Produce		json
	//	@Security		BearerAuth
	//	@Param			address	path		string						true	"Contract Address"
	//	@Param			tokenId	path		string						true	"Token ID"
	//	@Param			request	body		models.SetTokenDataRequest	true	"Token Data"
	//	@Success		200		{object}	models.SetTokenDataResponse
	//	@Failure		400		{string}	string	"Bad Request"
	//	@Failure		401		{string}	string	"Unauthorized"
	//	@Failure		500		{string}	string	"Internal Server Error"
	//	@Router			/collectibles/{address}/token-data/{tokenId} [put]
	mux.HandleFunc("PUT /collectibles/{address}/token-data/{tokenId}", middleware.AuthMiddleware(controllers.SetTokenDataHandler))

	//	@Summary		Get Token Data
	//	@Description	Get token data for collectible
	//	@Tags			collectibles
	//	@Accept			json
	//	@Produce		json
	//	@Param			address	path		string	true	"Contract Address"
	//	@Param			tokenId	path		string	true	"Token ID"
	//	@Success		200		{object}	models.GetTokenDataResponse
	//	@Failure		400		{string}	string	"Bad Request"
	//	@Failure		500		{string}	string	"Internal Server Error"
	//	@Router			/collectibles/{address}/token-data/{tokenId} [get]
	mux.HandleFunc("GET /collectibles/{address}/token-data/{tokenId}", controllers.GetTokenDataHandler)

	//	@Summary		Redeem Collectible
	//	@Description	Redeem collectible tokens
	//	@Tags			collectibles
	//	@Accept			json
	//	@Produce		json
	//	@Security		BearerAuth
	//	@Param			request	body		models.RedeemCollectibleRequest	true	"Redeem Request"
	//	@Success		200		{object}	models.RedeemCollectibleResponse
	//	@Failure		400		{string}	string	"Bad Request"
	//	@Failure		401		{string}	string	"Unauthorized"
	//	@Failure		500		{string}	string	"Internal Server Error"
	//	@Router			/collectibles/{address}/redeem [post]
	mux.HandleFunc("POST /collectibles/{address}/redeem", middleware.AuthMiddleware(controllers.RedeemCollectibleHandler))

	// Points endpoints
	//	@Summary		Mint Points
	//	@Description	Mint new points tokens
	//	@Tags			points
	//	@Accept			json
	//	@Produce		json
	//	@Security		BearerAuth
	//	@Param			request	body		models.MintPointsRequest	true	"Mint Request"
	//	@Success		200		{object}	models.MintPointsResponse
	//	@Failure		400		{string}	string	"Bad Request"
	//	@Failure		401		{string}	string	"Unauthorized"
	//	@Failure		500		{string}	string	"Internal Server Error"
	//	@Router			/points/mint [post]
	mux.HandleFunc("POST /points/mint", middleware.AuthMiddleware(controllers.MintPointsHandler))

	//	@Summary		Burn Points
	//	@Description	Burn points tokens
	//	@Tags			points
	//	@Accept			json
	//	@Produce		json
	//	@Security		BearerAuth
	//	@Param			request	body		models.BurnPointsRequest	true	"Burn Request"
	//	@Success		200		{object}	models.BurnPointsResponse
	//	@Failure		400		{string}	string	"Bad Request"
	//	@Failure		401		{string}	string	"Unauthorized"
	//	@Failure		500		{string}	string	"Internal Server Error"
	//	@Router			/points/burn [post]
	mux.HandleFunc("POST /points/burn", middleware.AuthMiddleware(controllers.BurnPointsHandler))

	//	@Summary		Get Points Balance
	//	@Description	Get points balance
	//	@Tags			points
	//	@Accept			json
	//	@Produce		json
	//	@Security		BearerAuth
	//	@Param			address	path		string	true	"Contract Address"
	//	@Success		200		{object}	models.GetPointsBalanceResponse
	//	@Failure		400		{string}	string	"Bad Request"
	//	@Failure		401		{string}	string	"Unauthorized"
	//	@Failure		500		{string}	string	"Internal Server Error"
	//	@Router			/points/{address}/balance [get]
	mux.HandleFunc("GET /points/{address}/balance", middleware.AuthMiddleware(controllers.GetPointsBalanceHandler))

	//	@Summary		Transfer Points
	//	@Description	Transfer points between accounts
	//	@Tags			points
	//	@Accept			json
	//	@Produce		json
	//	@Security		BearerAuth
	//	@Param			request	body		models.TransferPointsRequest	true	"Transfer Request"
	//	@Success		200		{object}	models.TransferPointsResponse
	//	@Failure		400		{string}	string	"Bad Request"
	//	@Failure		401		{string}	string	"Unauthorized"
	//	@Failure		500		{string}	string	"Internal Server Error"
	//	@Router			/points/transfer [post]
	mux.HandleFunc("POST /points/transfer", middleware.AuthMiddleware(controllers.TransferPointsHandler))

	// Merchant endpoints
	//	@Summary		Get Points Contracts
	//	@Description	Get merchant's points contracts
	//	@Tags			merchants
	//	@Accept			json
	//	@Produce		json
	//	@Security		BearerAuth
	//	@Success		200	{object}	models.GetPointsContractsResponse
	//	@Failure		401	{string}	string	"Unauthorized"
	//	@Failure		500	{string}	string	"Internal Server Error"
	//	@Router			/merchant/points-contracts [get]
	mux.HandleFunc("GET /merchant/points-contracts", middleware.AuthMiddleware(controllers.GetPointsContractsHandler))

	//	@Summary		Get Collectible Contracts
	//	@Description	Get merchant's collectible contracts
	//	@Tags			merchants
	//	@Accept			json
	//	@Produce		json
	//	@Security		BearerAuth
	//	@Success		200	{object}	models.GetCollectibleContractsResponse
	//	@Failure		401	{string}	string	"Unauthorized"
	//	@Failure		500	{string}	string	"Internal Server Error"
	//	@Router			/merchant/collectible-contracts [get]
	mux.HandleFunc("GET /merchant/collectible-contracts", middleware.AuthMiddleware(controllers.GetCollectibleContractsHandler))

	// Factory endpoints
	//	@Summary		Create Merchant
	//	@Description	Create a new merchant account
	//	@Tags			factory
	//	@Accept			json
	//	@Produce		json
	//	@Param			request	body		models.CreateMerchantRequest	true	"Merchant Creation Request"
	//	@Success		201		{object}	models.CreateMerchantResponse
	//	@Failure		400		{string}	string	"Bad Request"
	//	@Failure		500		{string}	string	"Internal Server Error"
	//	@Router			/merchant [post]
	mux.HandleFunc("POST /merchant", middleware.AuthMiddleware(controllers.CreateMerchantHandler))

	//	@Summary		Create Collectible
	//	@Description	Create a new collectible contract
	//	@Tags			factory
	//	@Accept			json
	//	@Produce		json
	//	@Security		BearerAuth
	//	@Param			request	body		models.CreateCollectibleRequest	true	"Collectible Creation Request"
	//	@Success		201		{object}	models.CreateCollectibleResponse
	//	@Failure		400		{string}	string	"Bad Request"
	//	@Failure		401		{string}	string	"Unauthorized"
	//	@Failure		500		{string}	string	"Internal Server Error"
	//	@Router			/collectibles [post]
	mux.HandleFunc("POST /collectibles", middleware.AuthMiddleware(controllers.CreateCollectibleHandler))

	//	@Summary		Create Points Contract
	//	@Description	Create a new points contract
	//	@Tags			factory
	//	@Accept			json
	//	@Produce		json
	//	@Security		BearerAuth
	//	@Param			request	body		models.CreatePointsContractRequest	true	"Points Contract Creation Request"
	//	@Success		201		{object}	models.CreatePointsContractResponse
	//	@Failure		400		{string}	string	"Bad Request"
	//	@Failure		401		{string}	string	"Unauthorized"
	//	@Failure		500		{string}	string	"Internal Server Error"
	//	@Router			/points-contracts [post]
	mux.HandleFunc("POST /points-contracts", middleware.AuthMiddleware(controllers.CreatePointsContractHandler))
}
