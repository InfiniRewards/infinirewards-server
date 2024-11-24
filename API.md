# InfiniRewards API
API for InfiniRewards - A Web3 Loyalty and Rewards Platform

## Version: 1.0

**Contact information:**  
API Support  
http://www.infinirewards.io/support  
support@infinirewards.io  

**License:** [Apache 2.0](http://www.apache.org/licenses/LICENSE-2.0.html)

### Security
**BearerAuth**  

| apiKey | *API Key* |
| ------ | --------- |
| Description | Enter your bearer token in the format **Bearer <token>** |
| Name | Authorization |
| In | header |

**Schemes:** http, https

---
### /auth/authenticate

#### POST
##### Summary

Authenticate user

##### Description

Authenticate user using OTP or API key

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| request | body | Authentication Request | Yes | [models.AuthenticateRequest](#modelsauthenticaterequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Authentication successful | [models.AuthenticateResponse](#modelsauthenticateresponse) |
| 400 | Invalid request format or validation failed | [models.ErrorResponse](#modelserrorresponse) |
| 401 | Authentication failed | [models.ErrorResponse](#modelserrorresponse) |
| 429 | Too many attempts | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

### /auth/refresh-token

#### POST
##### Summary

Refresh token

##### Description

Refresh an existing authentication token

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| request | body | Token Refresh Request | Yes | [models.RefreshTokenRequest](#modelsrefreshtokenrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Token refreshed successfully | [models.RefreshTokenResponse](#modelsrefreshtokenresponse) |
| 400 | Invalid request format | [models.ErrorResponse](#modelserrorresponse) |
| 401 | Invalid or expired token | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /auth/request-otp

#### POST
##### Summary

Request OTP

##### Description

Request a one-time password for authentication

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| request | body | OTP Request | Yes | [models.RequestOTPRequest](#modelsrequestotprequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OTP sent successfully | [models.RequestOTPResponse](#modelsrequestotpresponse) |
| 400 | Invalid request format or validation failed | [models.ErrorResponse](#modelserrorresponse) |
| 429 | Too many requests | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

---
### /collectibles/{address}

#### GET
##### Summary

Get collectible details

##### Description

Get detailed information about a collectible contract

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| address | path | Contract address | Yes | string (hex) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Collectible details retrieved successfully | [models.GetCollectibleDetailsResponse](#modelsgetcollectibledetailsresponse) |
| 400 | Invalid contract address | [models.ErrorResponse](#modelserrorresponse) |
| 404 | Contract not found | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

### /collectibles/{address}/balance/{tokenId}

#### GET
##### Summary

Get collectible balance

##### Description

Get balance of collectible tokens for a specific token ID

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| address | path | Contract address | Yes | string (hex) |
| tokenId | path | Token ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Balance retrieved successfully | [models.GetCollectibleBalanceResponse](#modelsgetcollectiblebalanceresponse) |
| 400 | Invalid request parameters | [models.ErrorResponse](#modelserrorresponse) |
| 401 | Missing or invalid authentication token | [models.ErrorResponse](#modelserrorresponse) |
| 404 | Contract or token not found | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /collectibles/{address}/purchase

#### POST
##### Summary

Purchase collectible

##### Description

Purchase a collectible token using points

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| address | path | Contract address | Yes | string (hex) |
| request | body | Purchase details | Yes | [models.PurchaseCollectibleRequest](#modelspurchasecollectiblerequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Purchase successful | [models.PurchaseCollectibleResponse](#modelspurchasecollectibleresponse) |
| 400 | Invalid request format or validation failed | [models.ErrorResponse](#modelserrorresponse) |
| 401 | Missing or invalid authentication token | [models.ErrorResponse](#modelserrorresponse) |
| 402 | Insufficient points balance | [models.ErrorResponse](#modelserrorresponse) |
| 404 | Collectible not found | [models.ErrorResponse](#modelserrorresponse) |
| 410 | Collectible expired | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /collectibles/{address}/redeem

#### POST
##### Summary

Redeem collectible

##### Description

Redeem a collectible token

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| address | path | Contract address | Yes | string (hex) |
| request | body | Redemption details | Yes | [models.RedeemCollectibleRequest](#modelsredeemcollectiblerequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Redemption successful | [models.RedeemCollectibleResponse](#modelsredeemcollectibleresponse) |
| 400 | Invalid request format or validation failed | [models.ErrorResponse](#modelserrorresponse) |
| 401 | Missing or invalid authentication token | [models.ErrorResponse](#modelserrorresponse) |
| 403 | Not authorized to redeem | [models.ErrorResponse](#modelserrorresponse) |
| 404 | Collectible not found | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /collectibles/{address}/token-data/{tokenId}

#### GET
##### Summary

Get token data

##### Description

Get metadata for a collectible token

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| address | path | Contract address | Yes | string (hex) |
| tokenId | path | Token ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Token data retrieved successfully | [models.GetTokenDataResponse](#modelsgettokendataresponse) |
| 400 | Invalid request parameters | [models.ErrorResponse](#modelserrorresponse) |
| 404 | Token data not found | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

#### PUT
##### Summary

Set token data

##### Description

Set metadata for a collectible token

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| address | path | Contract address | Yes | string (hex) |
| tokenId | path | Token ID | Yes | integer |
| request | body | Token data | Yes | [models.SetTokenDataRequest](#modelssettokendatarequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Token data updated successfully | [models.SetTokenDataResponse](#modelssettokendataresponse) |
| 400 | Invalid request format or validation failed | [models.ErrorResponse](#modelserrorresponse) |
| 401 | Missing or invalid authentication token | [models.ErrorResponse](#modelserrorresponse) |
| 403 | Not authorized to set token data | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /collectibles/{address}/uri/{tokenId}

#### GET
##### Summary

Get collectible URI

##### Description

Get the URI for a specific collectible token's metadata

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| address | path | Contract address | Yes | string (hex) |
| tokenId | path | Token ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | URI retrieved successfully | [models.GetCollectibleURIResponse](#modelsgetcollectibleuriresponse) |
| 400 | Invalid request parameters | [models.ErrorResponse](#modelserrorresponse) |
| 404 | Token not found | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

### /collectibles/{address}/valid/{tokenId}

#### GET
##### Summary

Check collectible validity

##### Description

Check if a collectible token is valid (not expired)

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| address | path | Contract address | Yes | string (hex) |
| tokenId | path | Token ID | Yes | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Validity status retrieved | [models.IsCollectibleValidResponse](#modelsiscollectiblevalidresponse) |
| 400 | Invalid request parameters | [models.ErrorResponse](#modelserrorresponse) |
| 404 | Token not found | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

### /merchant/collectibles/mint

#### POST
##### Summary

Mint collectible tokens

##### Description

Mint new collectible tokens for a specified recipient

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| request | body | Mint Request | Yes | [models.MintCollectibleRequest](#modelsmintcollectiblerequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Mint successful | [models.MintCollectibleResponse](#modelsmintcollectibleresponse) |
| 400 | Invalid request format or validation failed | [models.ErrorResponse](#modelserrorresponse) |
| 401 | Missing or invalid authentication token | [models.ErrorResponse](#modelserrorresponse) |
| 403 | Insufficient permissions to mint | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

---
### /merchant

#### GET
##### Summary

Get merchant details

##### Description

Get details of a merchant

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Merchant details | [models.Merchant](#modelsmerchant) |
| 401 | Missing or invalid authentication token | [models.ErrorResponse](#modelserrorresponse) |
| 403 | Not a merchant account | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### POST
##### Summary

Create new merchant

##### Description

Create a new merchant account with initial points contract

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| request | body | Merchant creation request | Yes | [models.CreateMerchantRequest](#modelscreatemerchantrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Merchant created successfully | [models.CreateMerchantResponse](#modelscreatemerchantresponse) |
| 400 | Invalid request format or validation failed | [models.ErrorResponse](#modelserrorresponse) |
| 409 | Merchant already exists | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

### /merchant/collectible-contracts

#### GET
##### Summary

Get merchant's collectible contracts

##### Description

Get all collectible contracts associated with a merchant

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of collectible contracts | [models.GetCollectibleContractsResponse](#modelsgetcollectiblecontractsresponse) |
| 401 | Missing or invalid authentication token | [models.ErrorResponse](#modelserrorresponse) |
| 403 | Not a merchant account | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /merchant/collectible/upgrade

#### POST
##### Summary

Upgrade collectible contract

##### Description

Upgrade the collectible contract

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| request | body | Upgrade Collectible Contract Request | Yes | [models.UpgradeCollectibleContractRequest](#modelsupgradecollectiblecontractrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.UpgradeCollectibleContractResponse](#modelsupgradecollectiblecontractresponse) |
| 400 | Bad Request | string |
| 401 | Unauthorized | string |
| 500 | Internal Server Error | string |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /merchant/collectibles

#### POST
##### Summary

Create collectible contract

##### Description

Create a new collectible contract for a merchant

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| request | body | Collectible creation request | Yes | [models.CreateCollectibleRequest](#modelscreatecollectiblerequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Created collectible details | [models.CreateCollectibleResponse](#modelscreatecollectibleresponse) |
| 400 | Invalid request format or validation failed | [models.ErrorResponse](#modelserrorresponse) |
| 401 | Missing or invalid authentication token | [models.ErrorResponse](#modelserrorresponse) |
| 403 | Not authorized to create collectibles | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /merchant/points

#### POST
##### Summary

Create points contract

##### Description

Create an additional points contract for a merchant

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| request | body | Points contract creation request | Yes | [models.CreatePointsContractRequest](#modelscreatepointscontractrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | Created points contract details | [models.CreatePointsContractResponse](#modelscreatepointscontractresponse) |
| 400 | Invalid request format or validation failed | [models.ErrorResponse](#modelserrorresponse) |
| 401 | Missing or invalid authentication token | [models.ErrorResponse](#modelserrorresponse) |
| 403 | Not authorized to create points contracts | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /merchant/points-contracts

#### GET
##### Summary

Get merchant's points contracts

##### Description

Get all points contracts associated with a merchant

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of points contracts | [models.GetPointsContractsResponse](#modelsgetpointscontractsresponse) |
| 401 | Missing or invalid authentication token | [models.ErrorResponse](#modelserrorresponse) |
| 403 | Not a merchant account | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /merchant/points/upgrade

#### POST
##### Summary

Upgrade points contract

##### Description

Upgrade the points contract

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| request | body | Upgrade Points Contract Request | Yes | [models.UpgradePointsContractRequest](#modelsupgradepointscontractrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.UpgradePointsContractResponse](#modelsupgradepointscontractresponse) |
| 400 | Bad Request | string |
| 401 | Unauthorized | string |
| 500 | Internal Server Error | string |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /merchant/upgrade

#### POST
##### Summary

Upgrade merchant contract

##### Description

Upgrade the merchant contract

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| request | body | Upgrade Merchant Contract Request | Yes | [models.UpgradeMerchantContractRequest](#modelsupgrademerchantcontractrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.UpgradeMerchantContractResponse](#modelsupgrademerchantcontractresponse) |
| 400 | Bad Request | string |
| 401 | Unauthorized | string |
| 500 | Internal Server Error | string |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

---
### /points/burn

#### POST
##### Summary

Burn points tokens

##### Description

Burn points tokens from a merchant's account

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| request | body | Burn Request | Yes | [models.BurnPointsRequest](#modelsburnpointsrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Points burned successfully | [models.BurnPointsResponse](#modelsburnpointsresponse) |
| 400 | Invalid request format or validation failed | [models.ErrorResponse](#modelserrorresponse) |
| 401 | Missing or invalid authentication token | [models.ErrorResponse](#modelserrorresponse) |
| 403 | Not authorized to burn points | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /points/mint

#### POST
##### Summary

Mint points tokens

##### Description

Mint new points tokens for a specified recipient

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| request | body | Mint Request | Yes | [models.MintPointsRequest](#modelsmintpointsrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Points minted successfully | [models.MintPointsResponse](#modelsmintpointsresponse) |
| 400 | Invalid request format or validation failed | [models.ErrorResponse](#modelserrorresponse) |
| 401 | Missing or invalid authentication token | [models.ErrorResponse](#modelserrorresponse) |
| 403 | Not authorized to mint points | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /points/transfer

#### POST
##### Summary

Transfer points between accounts

##### Description

Transfer points from one account to another

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| request | body | Transfer Request | Yes | [models.TransferPointsRequest](#modelstransferpointsrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Points transferred successfully | [models.TransferPointsResponse](#modelstransferpointsresponse) |
| 400 | Invalid request format or validation failed | [models.ErrorResponse](#modelserrorresponse) |
| 401 | Missing or invalid authentication token | [models.ErrorResponse](#modelserrorresponse) |
| 403 | Not authorized to transfer points | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /points/{address}/balance

#### GET
##### Summary

Get points balance

##### Description

Get the points balance for a specific account

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| address | path | Points contract address | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | Balance retrieved | [models.GetPointsBalanceResponse](#modelsgetpointsbalanceresponse) |
| 400 | Invalid request format | [models.ErrorResponse](#modelserrorresponse) |
| 401 | Unauthorized | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

---
### /user

#### GET
##### Summary

Get user details

##### Description

Get authenticated user details

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | User details retrieved successfully | [models.User](#modelsuser) |
| 401 | Authentication error | [models.ErrorResponse](#modelserrorresponse) |
| 404 | User not found | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### PUT
##### Summary

Update user

##### Description

Update authenticated user details

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| request | body | User Update Request | Yes | [models.UpdateUserRequest](#modelsupdateuserrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | User updated successfully | [models.User](#modelsuser) |
| 400 | Invalid request format | [models.ErrorResponse](#modelserrorresponse) |
| 401 | Unauthorized access | [models.ErrorResponse](#modelserrorresponse) |
| 404 | User not found | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### POST
##### Summary

Create user

##### Description

Create a new user

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| request | body | User Creation Request | Yes | [models.CreateUserRequest](#modelscreateuserrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | User created successfully | [models.User](#modelsuser) |
| 400 | Invalid request format | [models.ErrorResponse](#modelserrorresponse) |
| 409 | User already exists | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

#### DELETE
##### Summary

Delete user

##### Description

Delete authenticated user

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | User deleted successfully | [models.MessageResponse](#modelsmessageresponse) |
| 401 | Unauthorized access | [models.ErrorResponse](#modelserrorresponse) |
| 404 | User not found | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /user/upgrade

#### POST
##### Summary

Upgrade User Contract

##### Description

Upgrade a user contract

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| request | body | Upgrade User Contract Request | Yes | [models.UpgradeUserContractRequest](#modelsupgradeusercontractrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [models.UpgradeUserContractResponse](#modelsupgradeusercontractresponse) |
| 400 | Bad Request | string |
| 401 | Unauthorized | string |
| 500 | Internal Server Error | string |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

---
### /user/api-keys

#### GET
##### Summary

List API keys

##### Description

List all API keys for authenticated user

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | List of API keys | [ [models.APIKey](#modelsapikey) ] |
| 401 | Unauthorized access | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

#### POST
##### Summary

Create API key

##### Description

Create a new API key for authenticated user

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| request | body | API Key Creation Request | Yes | [models.CreateAPIKeyRequest](#modelscreateapikeyrequest) |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 201 | API key created successfully | [models.APIKey](#modelsapikey) |
| 400 | Invalid request format | [models.ErrorResponse](#modelserrorresponse) |
| 401 | Unauthorized access | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

### /user/api-keys/{keyId}

#### DELETE
##### Summary

Delete API key

##### Description

Delete an API key for authenticated user

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ------ |
| keyId | path | API Key ID | Yes | string |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | API key deleted successfully | [models.MessageResponse](#modelsmessageresponse) |
| 400 | Invalid request format | [models.ErrorResponse](#modelserrorresponse) |
| 401 | Unauthorized access | [models.ErrorResponse](#modelserrorresponse) |
| 500 | Internal server error | [models.ErrorResponse](#modelserrorresponse) |

##### Security

| Security Schema | Scopes |
| --------------- | ------ |
| BearerAuth |  |

---
### Models

#### models.APIKey

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| createdAt | string |  | No |
| id | string |  | No |
| name | string |  | No |
| secret | string |  | No |
| updatedAt | string |  | No |
| userId | string |  | No |

#### models.AuthenticateRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| device | string | Device is a unique device identifier example: device_123 | Yes |
| id | string | ID is either a verification ID (for OTP) or API key ID example: 01HNAJ6640M9JRRJFQSZZVE3HH | Yes |
| method | string | Method must be either "otp" or "secret" example: otp<br>*Enum:* `"otp"`, `"secret"` | Yes |
| signature | string | Signature is a unique device identifier example: device_signature_123 | Yes |
| token | string | Token is either an OTP code or API key secret example: 123456 | Yes |

#### models.AuthenticateResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| token | [models.Token](#modelstoken) |  | No |

#### models.BurnPointsRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| amount | string | Amount of points to burn example: 50 | Yes |
| pointsContract | string | PointsContract is the address of the points contract example: 0x1234567890abcdef1234567890abcdef12345678 | Yes |

#### models.BurnPointsResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| transactionHash | string | TransactionHash is the hash of the burn transaction example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890 | No |

#### models.CollectibleContractInfo

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| address | string | Address of the contract example: 0x1234567890abcdef1234567890abcdef12345678 | No |
| description | string | Description of the collection example: Limited collectibles | No |
| name | string | Name of the collectible collection example: Special Edition | No |
| pointsContract | string | PointsContract is the address of the points contract used for purchases example: 0x1234567890abcdef1234567890abcdef12345678 | No |
| tokenDescriptions | [ string ] | TokenDescriptions lists descriptions for each token example: ["Gold","Silver","Bronze"] | No |
| tokenExpiries | [ integer ] | TokenExpiries lists expiry timestamps for each token example: [1735689600,1735689600,1735689600] | No |
| tokenIds | [ string ] | TokenIDs lists all token IDs in the collection example: ["1","2","3"] | No |
| tokenPrices | [ string ] | TokenPrices lists prices for each token example: ["100","200","300"] | No |
| tokenSupplies | [ string ] | TokenSupplies lists supplies for each token example: [100,200,300] | No |

#### models.CreateAPIKeyRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| name | string |  | No |

#### models.CreateCollectibleRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| description | string | Description of the collectible collection example: Limited edition collectibles | Yes |
| name | string | Name of the collectible contract example: Special Edition | Yes |

#### models.CreateCollectibleResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| address | string | Address is the deployed contract address example: 0x1234567890abcdef1234567890abcdef12345678 | No |
| transactionHash | string | TransactionHash is the hash of the creation transaction example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890 | No |

#### models.CreateMerchantRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| decimals | integer | Decimals for the points token example: 18 | Yes |
| name | string | Name of the merchant example: My Store | Yes |
| symbol | string | Symbol for the points token example: PTS | Yes |

#### models.CreateMerchantResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| merchantAddress | string | MerchantAddress is the deployed merchant contract address example: 0x1234567890abcdef1234567890abcdef12345678 | No |
| pointsAddress | string | PointsAddress is the deployed points contract address example: 0x9876543210abcdef1234567890abcdef12345678 | No |
| transactionHash | string | TransactionHash is the hash of the creation transaction example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890 | No |

#### models.CreatePointsContractRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| decimals | string | Decimals specifies the number of decimal places example: 18 | Yes |
| description | string | Description of the points system example: Premium tier loyalty points | Yes |
| name | string | Name of the points token example: Premium Points | Yes |
| symbol | string | Symbol for the points token (3-4 characters) example: PPT | Yes |

#### models.CreatePointsContractResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| address | string | Address is the deployed contract address example: 0x1234567890abcdef1234567890abcdef12345678 | No |
| transactionHash | string | TransactionHash is the hash of the creation transaction example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890 | No |

#### models.CreateUserRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| avatar | string | Avatar is the user's avatar | No |
| email | string | Email is the user's email example: john.doe@example.com | No |
| name | string | Name is the user's name example: John Doe | No |

#### models.ErrorResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| code | string |  | No |
| details |  |  | No |
| message | string |  | No |

#### models.GetCollectibleBalanceResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| balance | string | Balance is the number of tokens owned example: 100 | No |

#### models.GetCollectibleContractsResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| contracts | [ [models.CollectibleContractInfo](#modelscollectiblecontractinfo) ] | Contracts is the list of collectible contracts | No |

#### models.GetCollectibleDetailsResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| address | string | Address of the collectible collection example: 0x1234567890abcdef1234567890abcdef12345678 | No |
| description | string | Description of the collectible collection example: Special Edition Collectibles | No |
| name | string | Name of the collectible collection example: Special Edition Collectibles | No |
| pointsContract | string | PointsContract is the address of the points contract example: 0x1234567890abcdef1234567890abcdef12345678 | No |
| tokenBalances | [ string ] | TokenBalances lists balances for each token example: ["10","20","30"] | No |
| tokenDescriptions | [ string ] | TokenDescriptions lists descriptions for each token example: ["Gold","Silver","Bronze"] | No |
| tokenExpiries | [ integer ] | TokenExpiries lists expiry timestamps for each token example: [1735689600,1735689600,1735689600] | No |
| tokenIds | [ string ] | TokenIDs lists all token IDs in the collection example: ["1","2","3"] | No |
| tokenPrices | [ string ] | TokenPrices lists prices for each token example: ["100","200","300"] | No |
| tokenSupplies | [ integer ] | TokenSupplies lists supplies for each token example: [100,200,300] | No |

#### models.GetCollectibleURIResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| uri | string | URI of the collectible metadata example: https://example.com/metadata/1 | No |

#### models.GetPointsBalanceResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| balance | string | Balance is the number of points owned example: 1000 | No |
| decimals | integer | Decimals of the points token example: 18 | No |
| description | string | Description of the points token example: Premium points for our most loyal customers | No |
| name | string | Name of the points token example: Premium Points | No |
| symbol | string | Symbol of the points token example: PPM | No |

#### models.GetPointsContractsResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| contracts | [ [models.PointsContractInfo](#modelspointscontractinfo) ] | Contracts is the list of points contracts | No |

#### models.GetTokenDataResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| description | string | Description is the token's metadata description example: Limited edition collectible | No |
| expiry | integer | Expiry is the Unix timestamp when the token expires example: 1735689600 | No |
| pointsContract | string | PointsContract is the address of the points contract used for purchases example: 0x1234567890abcdef1234567890abcdef12345678 | No |
| price | string | Price is the cost in points to purchase the token example: 100 | No |

#### models.IsCollectibleValidResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| isValid | boolean | IsValid indicates if the collectible is still valid example: true | No |

#### models.Merchant

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| address | string |  | No |
| createdAt | string |  | No |
| decimals | integer |  | No |
| id | string |  | No |
| name | string |  | No |
| symbol | string |  | No |
| updatedAt | string |  | No |

#### models.MessageResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| message | string | *Example:* `"message"` | No |

#### models.MintCollectibleRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| amount | string | Amount is the number of tokens to mint example: 1 | Yes |
| collectibleAddress | string | CollectibleAddress is the contract address of the collectible example: 0x1234567890abcdef1234567890abcdef12345678 | Yes |
| to | string | To is the recipient address example: 0x9876543210abcdef1234567890abcdef12345678 | Yes |
| tokenId | string | TokenId is the ID of the token to mint example: 1 | Yes |

#### models.MintCollectibleResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| transactionHash | string | TransactionHash is the hash of the mint transaction example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890 | No |

#### models.MintPointsRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| amount | string | Amount of points to mint (in smallest unit) example: 100 | Yes |
| pointsContract | string | PointsContract is the address of the points contract example: 0x1234567890abcdef1234567890abcdef12345678 | Yes |
| recipient | string | Recipient is the address receiving the points example: 0x9876543210abcdef1234567890abcdef12345678 | Yes |

#### models.MintPointsResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| transactionHash | string | TransactionHash is the hash of the mint transaction example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890 | No |

#### models.PointsContractInfo

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| address | string | Address of the contract example: 0x1234567890abcdef1234567890abcdef12345678 | No |
| decimals | integer | Decimals places for the token example: 18 | No |
| description | string | Description of the points token example: Loyalty points for Store XYZ | No |
| name | string | Name of the points token example: Store Points | No |
| symbol | string | Symbol of the points token example: SP | No |
| totalSupply | integer | TotalSupply of the points token example: 1000000 | No |

#### models.PurchaseCollectibleRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| amount | string | Amount is the number of tokens to purchase example: 1 | Yes |
| tokenId | string | TokenId is the ID of the token to purchase example: 1 | Yes |
| user | string | User is the address purchasing the collectible example: 0x1234567890abcdef1234567890abcdef12345678 | Yes |

#### models.PurchaseCollectibleResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| transactionHash | string | TransactionHash is the hash of the purchase transaction example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890 | No |

#### models.RedeemCollectibleRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| amount | string | Amount is the number of tokens to redeem example: 1 | Yes |
| tokenId | string | TokenId is the ID of the token to redeem example: 1 | Yes |
| user | string | User is the address redeeming the collectible example: 0x1234567890abcdef1234567890abcdef12345678 | Yes |

#### models.RedeemCollectibleResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| transactionHash | string | TransactionHash is the hash of the redemption transaction example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890 | No |

#### models.RefreshTokenRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| refreshToken | string |  | No |

#### models.RefreshTokenResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| creds | string |  | No |
| token | [models.Token](#modelstoken) |  | No |

#### models.RequestOTPRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| phoneNumber | string | PhoneNumber must be in E.164 format (e.g. +60123456789) example: +60123456789 | Yes |

#### models.RequestOTPResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| id | string | *Example:* `"phoneNumberVerification:01HNAJ6640M9JRRJFQSZZVE3HH"` | No |
| message | string | *Example:* `"message"` | No |

#### models.SetTokenDataRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| description | string | Description is the token's metadata description example: Limited edition collectible | Yes |
| expiry | integer | Expiry is the Unix timestamp when the token expires example: 1735689600 | Yes |
| pointsContract | string | PointsContract is the address of the points contract used for purchases example: 0x1234567890abcdef1234567890abcdef12345678 | Yes |
| price | string | Price is the cost in points to purchase the token example: 100 | Yes |

#### models.SetTokenDataResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| transactionHash | string | TransactionHash is the hash of the set data transaction example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890 | No |

#### models.Token

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| accessToken | string |  | No |
| accessTokenExpiry | string | RefreshToken       string    `json:"refreshToken"` (use ID as refresh token) | No |
| createdAt | string |  | No |
| device | string |  | No |
| id | string |  | No |
| refreshTokenExpiry | string |  | No |
| service | string |  | No |
| user | string |  | No |

#### models.TransferPointsRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| amount | string | Amount of points to transfer example: 25 | Yes |
| pointsContract | string | PointsContract is the address of the points contract example: 0x1234567890abcdef1234567890abcdef12345678 | Yes |
| to | string | To is the recipient's address example: 0x9876543210abcdef1234567890abcdef12345678 | Yes |

#### models.TransferPointsResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| transactionHash | string | TransactionHash is the hash of the transfer transaction example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890 | No |

#### models.UpdateUserRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| avatar | string | Avatar is the user's avatar | No |
| email | string | Email is the user's email example: john.doe@example.com | No |
| name | string | Name is the user's name example: John Doe | No |

#### models.UpgradeCollectibleContractRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| collectibleAddress | string | CollectibleAddress is the address of the collectible contract to upgrade example: 0x1234567890abcdef1234567890abcdef12345678 | No |
| newClassHash | string | NewClassHash is the class hash of the new implementation contract example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890 | No |

#### models.UpgradeCollectibleContractResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| transactionHash | string | TransactionHash is the hash of the upgrade transaction example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890 | No |

#### models.UpgradeMerchantContractRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| newClassHash | string | NewClassHash is the class hash of the new implementation contract example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890 | No |

#### models.UpgradeMerchantContractResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| transactionHash | string | TransactionHash is the hash of the upgrade transaction example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890 | No |

#### models.UpgradePointsContractRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| newClassHash | string | NewClassHash is the class hash of the new implementation contract example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890 | No |
| pointsContract | string | PointsContract is the address of the points contract to upgrade example: 0x1234567890abcdef1234567890abcdef12345678 | No |

#### models.UpgradePointsContractResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| transactionHash | string | TransactionHash is the hash of the upgrade transaction example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890 | No |

#### models.UpgradeUserContractRequest

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| newClassHash | string | NewClassHash is the class hash of the new implementation contract example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890 | No |

#### models.UpgradeUserContractResponse

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| transactionHash | string | TransactionHash is the hash of the upgrade transaction example: 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890 | No |

#### models.User

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| accountAddress | string | AccountAddress is the user's StarkNet account address | No |
| avatar | string | Avatar is the user's avatar | No |
| createdAt | string | CreatedAt is the time the user was created | No |
| email | string | Email is the user's email example: john.doe@example.com | No |
| id | string | ID is the user's ID example: 0x1234567890abcdef1234567890abcdef12345678 | No |
| name | string | Name is the user's name example: John Doe | No |
| phoneNumber | string | PhoneNumber is the user's phone number example: 1234567890 | No |
| privateKey | string | PrivateKey is the user's StarkNet private key | No |
| publicKey | string | PublicKey is the user's StarkNet public key | No |
| updatedAt | string | UpdatedAt is the time the user was last updated | No |
