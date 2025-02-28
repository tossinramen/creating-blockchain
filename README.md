# Library Checkout Blockchain

This Go application implements a simple blockchain for tracking library book checkouts. Each checkout transaction is stored as a block in the blockchain, ensuring data integrity and security.

## Features
- Implements a blockchain to store book checkout records.
- Supports adding new books and tracking checkouts.
- Ensures valid block chaining with cryptographic hashing.
- Provides an API to interact with the blockchain.

## API Endpoints
- `GET /` - Retrieve the blockchain.
- `POST /` - Add a new book checkout record.
- `POST /new` - Register a new book.

