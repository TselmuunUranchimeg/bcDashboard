# Back-end

# Project structure

| Directories | Description |
| :---------- | :---------- |
| `middlewares` | Middlewares for the server |
| `routes` | collection of `http.HandlerFunc` for routes |
| `services` | Core logic of the server |
| `tests` | Unit tests for service functions. Few tests are not included here as they are tests for private functions. |
| `token` | Golang code of the smart contracts's compiled code |

## Entity-relationship diagram (ERD)

![image-info](image/ERD%20image.png)

### Tables

| Table | Description |
| :---- | :---------- |
| `users` | Registered users. `username` was chosen to make sure of no duplicates as the authentication is rather crude and simple |
| `wallets` | Wallet addresses created by users. Private key is not stored as that defeats the purpose of decentralization. |
| `transactions` | Transfers of Ethereum and ERC20 tokens (only those from registered contracts) to registered wallet addresses |
| `contracts` | Contract information registered by the admin |
| `networks` | Rows of blockchain network URLs for the server to connect to |
| `check` | Record of last checked block number along with the `TIMESTAMP` (with timezone). |

## Background task

Since Golang doesn't have a built-in scheduler, the `BackgroundTask` is basically a while loop that will run after provided intervals. 

```
func BackgroundTask(db *sql.DB, hm *HashMap, duration int) {
	for {
		if err := backgroundClientTask(db, hm); err != nil {
			fmt.Printf("%v\n", err.Error())
			if hm.Client == nil {
				i := 1
				for hm.Client != nil {
					time.Sleep(time.Duration(i) * time.Second)
					hm.Client, err = ethclient.Dial(hm.Url)
					if err != nil {
						i += 1
					}
				}
			}
		}
		time.Sleep(time.Second * time.Duration(duration))
	}
}
```

As you might have noticed, there is a `HashMap` struct as one of the function parameters. This struct stores all of the important information about the network and will be passed to both the `BackgroundTask` function and `http.HandlerFunc` functions. The reason for this choice is that by keeping the connection alive and shared, there will be less overhead compared to starting another connection in handler functions while still keeping one in the background. 

```
type HashMap struct {
	Mu               sync.RWMutex // To avoid race conditions and other issues
	Value, Contracts map[string]bool
	Id               int
	Client           *ethclient.Client
	Url              string
}
```

Keep in mind that depending on the network, there is the risk of losing the connection. In the `BackgroundTask` function, it is still possible to reconnect, however; there will be a period, where requests will be met with errors. 

```
...
if hm.Client == nil {
    i := 1
    for hm.Client != nil {
        time.Sleep(time.Duration(i) * time.Second)
        hm.Client, err = ethclient.Dial(hm.Url)
        if err != nil {
            i += 1
        }
    }
}
...
```

## Authentication

As the main purpose of this project was to get familiar with how to effectively work with blockchains using Golang, authentication was not the main concern, Hence, only a simple JWT authentication was used. (Feel free to change it however you want.)

>When testing the authentication with tools, such as Postman, Thunder Client, remember to keep the `Secure` property as `false`. Otherwise, the cookies will not be saved for subsequent requests. 


## Processing transactions

As the [go-ethereum](https://github.com/ethereum/go-ethereum) package doesn't provide any direct way to get the sender's address, the following method was used to get the address. 
```
from, err := types.Sender(types.NewLondonSigner(tx.ChainId()), tx)
if err != nil {
    return err
}
```

When a transaction is a transfer from smart contracts, the `types.Transaction.To()` returns the address of the smart contract and not the actual recipient. The only method I cound find was to check the bytecode of the address. If the address is truly a contract address, its length will be more than 0. 

```
bytecode, err := client.CodeAt(context.Background(), address, blockNumber)
if err != nil {
    ...
}

if len(bytecode) > 0 {
    // Contract address
}
```