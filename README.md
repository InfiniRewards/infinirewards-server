# InfiniRewards Server

InfiniRewards is a Web3 Loyalty and Rewards Platform that enables merchants to create and manage digital loyalty programs using blockchain technology.

## Overview

The InfiniRewards server provides a robust REST API for managing rewards and collectibles on the Starknet blockchain. It handles user authentication, merchant management, and blockchain interactions for minting and managing digital assets.

## Features

- **Authentication**
  - Phone number verification via OTP
  - API key management for secure access
  - JWT-based authentication with refresh tokens

- **User Management**
  - User registration and profile management
  - Starknet wallet integration
  - API key creation and management

- **Merchant Features**
  - Merchant account creation and management
  - Points contract deployment and management
  - Collectible contract creation and management

- **Digital Assets**
  - Points token minting and burning
  - Collectible NFT minting and management
  - Token balance checking and transfers

## Technical Stack

- **Backend**: Go (Golang)
- **Blockchain**: Starknet
- **Message Broker**: NATS with JetStream
- **Authentication**: JWT + OTP
- **Documentation**: Swagger/OpenAPI

## API Documentation

The API is fully documented using Swagger/OpenAPI. When running in non-production environments, you can access the interactive API documentation at:
