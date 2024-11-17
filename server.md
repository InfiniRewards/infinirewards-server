# Services

## Auth
Authentication service.
### Auth Steps
1. Connect to InfiniRewards NATS `public` service.
2. Get a temporary JWT token from the `public` service. (JWT with permissions only to use the `auth` service).
3. Use the JWT token and user input to authenticate with the `auth` service.
4. If authentication is successful, the `auth` service will return another JWT token. (JWT with permissions to access 'infinirewards' service).
5. Use the JWT token to access `infinirewards` services.
### Endpoints
#### Phone OTP Auth
* getOtp
  * User enters phone number and request OTP. Server generates a TOTP secret and 6 digit OTP is sent to the user's phone number.
* verifyOtp
  * User calls `auth` service with phone number and OTP to authenticate.

#### Private Key Auth
Private key authentication service.
* createKey
  * Admin generate keys on dashboard. (key and secret).
* loginWithKey
  * User can use the key and secret to authenticate with `auth` service.

## InfiniRewards
InfiniRewards service.
### Endpoints
#### Account Management
* User can view their account points balance.
* User can view their account activity history.
* User can view their account vouchers.
* User can view their account subscribed memberships.
* User can manage their account.

#### Points Management
* mintPoints
  * call contract onchain with user's private key stored in db.
  * merchant can call this api to mint points to user's account.
* purchaseVouchers
  * call contract onchain with user's private key stored in db.
  * user can call this api to purchase vouchers issued by merchant.
* getPointsBalance
  * call contract onchain with user's private key stored in db.
  * user can call this api to get their points balance.


#### Voucher Management
* mintVoucher
  * call contract onchain with user's private key stored in db.
  * merchant can call this api to mint vouchers to user's account or to shop account to be purchased by user.
* transferVoucher
  * call contract onchain with user's private key stored in db.
  * user can call this api to transfer vouchers to another user.
* getVoucher
  * call contract onchain with user's private key stored in db.
  * user can call this api to get the voucher details.
* getVoucherBalance
  * call contract onchain with user's private key stored in db.
  * user can call this api to get the voucher balance.
* redeemVouchers
  * call contract onchain with user's private key stored in db. 
  * user can call this api to redeem vouchers.
* getVoucherRedemptionHistory
  * call contract onchain with user's private key stored in db.
  * user can call this api to get the voucher redemption history.

#### Membership Management
* getMembership
  * call contract onchain with user's private key stored in db.
  * user can call this api to get the membership details.
* getActiveMenberships
  * call contract onchain with user's private key stored in db.
  * user can call this api to get the active memberships.
* getMembershipHistory
  * call contract onchain with user's private key stored in db.
  * user can call this api to get the membership history.
* joinMembership
  * call contract onchain with user's private key stored in db.
  * user can call this api to join a membership.

### Public
Public services are services that are not part of the core InfiniRewards system. They are used to provide information about the InfiniRewards system to the public.
### Endpoints
#### Auth
* auth
  * Get temporary JWT token with limited permissions to access `auth` service.


# Data Models 

## User
Stored in SurrealDB
* id : string
* phoneNumber : string
* email : string
* name : string
* avatar : string
* createdAt : Date
* updatedAt : Date
* accountAddress : string
* privateKey : string
* address : string
* apiKeys : APIKey[]

## APIKey
Stored in SurrealDB
* id : string
* userId : string
* secret : string
* createdAt : Date
* updatedAt : Date

## Voucher (InfiniRewardsCollectible)
Stored onchain
* token_id : u256 (unique identifier for the voucher)
* price : u256 (price in points required to redeem the voucher)
* expiry : u64 (timestamp when the voucher expires)
* balance : u256 (number of vouchers available, tracked by ERC1155 balance)
* uri : string (metadata URI containing title, description, image)

## Membership (InfiniRewardsCollectible)
Stored onchain
* token_id : u256 (unique identifier for the membership)
* price : u256 (price in points required to join the membership)
* expiry : u64 (timestamp when the membership expires)
* balance : u256 (number of memberships available, tracked by ERC1155 balance)
* uri : string (metadata URI containing name, description, image)

## Points (InfiniRewardsPoints)
Stored onchain
* name : ByteArray
* symbol : ByteArray
* decimals : u8
* balance : u256 (tracked by ERC20 balance for each user)